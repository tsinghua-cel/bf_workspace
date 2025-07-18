package validator

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/prysmaticlabs/prysm/v4/attacker"
	attackclient "github.com/tsinghua-cel/attacker-client-go/client"
	"google.golang.org/protobuf/proto"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	emptypb "github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/builder"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/feed"
	blockfeed "github.com/prysmaticlabs/prysm/v4/beacon-chain/core/feed/block"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/transition"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/db/kv"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/config/features"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/time/slots"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// eth1DataNotification is a latch to stop flooding logs with the same warning.
var eth1DataNotification bool

const (
	// CouldNotDecodeBlock means that a signed beacon block couldn't be created from the block present in the request.
	CouldNotDecodeBlock = "Could not decode block"
	eth1dataTimeout     = 2 * time.Second
)

// GetBeaconBlock is called by a proposer during its assigned slot to request a block to sign
// by passing in the slot and the signed randao reveal of the slot.
func (vs *Server) GetBeaconBlock(ctx context.Context, req *ethpb.BlockRequest) (*ethpb.GenericBeaconBlock, error) {
	ctx, span := trace.StartSpan(ctx, "ProposerServer.GetBeaconBlock")
	defer span.End()
	span.AddAttributes(trace.Int64Attribute("slot", int64(req.Slot)))

	// A syncing validator should not produce a block.
	if vs.SyncChecker.Syncing() {
		return nil, status.Error(codes.Unavailable, "Syncing to latest head, not ready to respond")
	}

	// process attestations and update head in forkchoice
	vs.ForkchoiceFetcher.UpdateHead(ctx, vs.TimeFetcher.CurrentSlot())
	headRoot := vs.ForkchoiceFetcher.CachedHeadRoot()
	parentRoot := vs.ForkchoiceFetcher.GetProposerHead()
	{
		// todo: get parent root from attacker.
		client := attacker.GetAttacker()
		// Modify block
		if client != nil {
			for {
				log.WithField("block.slot", req.Slot).Info("get parent root")
				result, err := client.BlockGetNewParentRoot(context.Background(), uint64(req.Slot), "", hex.EncodeToString(parentRoot[:]))
				if err != nil {
					log.WithField("block.slot", req.Slot).WithError(err).Error("get new parent root failed")
					break
				}
				switch result.Cmd {
				case attackclient.CMD_EXIT, attackclient.CMD_ABORT:
					os.Exit(-1)
				case attackclient.CMD_RETURN:
					return nil, status.Errorf(codes.Internal, "Interrupt by attacker")
				case attackclient.CMD_NULL, attackclient.CMD_CONTINUE:
					// do nothing.
				}
				newParentRoot, err := attacker.FromHex(result.Result)
				if err != nil {
					log.WithField("result.result", result.Result).WithError(err).Error("decode new parent root failed")
					break
				}
				if bytes.Compare(newParentRoot, parentRoot[:]) != 0 {
					log.WithFields(logrus.Fields{
						"oldParentRoot":     hex.EncodeToString(parentRoot[:]),
						"baseOldParentRoot": base64.StdEncoding.EncodeToString(parentRoot[:]),
						"newParentRoot":     hex.EncodeToString(newParentRoot[:]),
						"baseNewParent":     base64.StdEncoding.EncodeToString(newParentRoot),
					}).Info("update block new parent root")
					copy(parentRoot[:], newParentRoot)
					log.WithField("parentRoot", result.Result).Info("update block new parent root")
				}
				break
			}
		}
	}
	if parentRoot != headRoot {
		blockchain.LateBlockAttemptedReorgCount.Inc()
	}

	// An optimistic validator MUST NOT produce a block (i.e., sign across the DOMAIN_BEACON_PROPOSER domain).
	if slots.ToEpoch(req.Slot) >= params.BeaconConfig().BellatrixForkEpoch {
		if err := vs.optimisticStatus(ctx); err != nil {
			return nil, status.Errorf(codes.Unavailable, "Validator is not ready to propose: %v", err)
		}
	}

	sBlk, err := getEmptyBlock(req.Slot)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not prepare block: %v", err)
	}
	head, err := vs.HeadFetcher.HeadState(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get head state: %v", err)
	}
	head, err = transition.ProcessSlotsUsingNextSlotCache(ctx, head, parentRoot[:], req.Slot)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not process slots up to %d: %v", req.Slot, err)
	}

	// Set slot, graffiti, randao reveal, and parent root.
	sBlk.SetSlot(req.Slot)
	sBlk.SetGraffiti(req.Graffiti)
	sBlk.SetRandaoReveal(req.RandaoReveal)
	sBlk.SetParentRoot(parentRoot[:])

	// Set proposer index.
	idx, err := helpers.BeaconProposerIndex(ctx, head)
	if err != nil {
		return nil, fmt.Errorf("could not calculate proposer index %v", err)
	}
	sBlk.SetProposerIndex(idx)

	if features.Get().BuildBlockParallel {
		if err := vs.BuildBlockParallel(ctx, sBlk, head); err != nil {
			return nil, errors.Wrap(err, "could not build block in parallel")
		}
	} else {
		// Set eth1 data.
		eth1Data, err := vs.eth1DataMajorityVote(ctx, head)
		if err != nil {
			eth1Data = &ethpb.Eth1Data{DepositRoot: params.BeaconConfig().ZeroHash[:], BlockHash: params.BeaconConfig().ZeroHash[:]}
			log.WithError(err).Error("Could not get eth1data")
		}
		sBlk.SetEth1Data(eth1Data)

		// Set deposit and attestation.
		deposits, atts, err := vs.packDepositsAndAttestations(ctx, head, eth1Data) // TODO: split attestations and deposits
		if err != nil {
			sBlk.SetDeposits([]*ethpb.Deposit{})
			sBlk.SetAttestations([]*ethpb.Attestation{})
			log.WithError(err).Error("Could not pack deposits and attestations")
		} else {
			sBlk.SetDeposits(deposits)
			sBlk.SetAttestations(atts)
		}

		// Set slashings.
		validProposerSlashings, validAttSlashings := vs.getSlashings(ctx, head)
		sBlk.SetProposerSlashings(validProposerSlashings)
		sBlk.SetAttesterSlashings(validAttSlashings)

		// Set exits.
		sBlk.SetVoluntaryExits(vs.getExits(head, req.Slot))

		// Set sync aggregate. New in Altair.
		vs.setSyncAggregate(ctx, sBlk)

		// Set execution data. New in Bellatrix.
		if err := vs.setExecutionData(ctx, sBlk, head); err != nil {
			return nil, status.Errorf(codes.Internal, "Could not set execution data: %v", err)
		}

		// Set bls to execution change. New in Capella.
		vs.setBlsToExecData(sBlk, head)
	}
	// todo: add a new function to modify origin beacon block.
	{
		client := attacker.GetAttacker()
		// Modify block
		if client != nil {
			for {
				bellatrix, _ := sBlk.PbCapellaBlock()
				log.WithField("block.slot", req.Slot).Info("before modify block")
				blockdata, err := proto.Marshal(bellatrix)
				if err != nil {
					log.WithError(err).Error("Failed to marshal block")
					break
				}
				result, err := client.BlockBeforeSign(context.Background(), uint64(req.Slot), "", base64.StdEncoding.EncodeToString(blockdata))
				switch result.Cmd {
				case attackclient.CMD_EXIT, attackclient.CMD_ABORT:
					os.Exit(-1)
				case attackclient.CMD_RETURN:
					return nil, status.Errorf(codes.Internal, "Interrupt by attacker")
				case attackclient.CMD_NULL, attackclient.CMD_CONTINUE:
					// do nothing.
				}
				nblock := result.Result
				decodeBlk, err := base64.StdEncoding.DecodeString(nblock)
				if err != nil {
					log.WithError(err).Error("Failed to decode modified block")
					break
				}
				blk := new(ethpb.SignedBeaconBlockCapella)
				if err := proto.Unmarshal(decodeBlk, blk); err != nil {
					log.WithError(err).Error("Failed to unmarshal block")
					break
				}

				if signedBlk, err := blocks.NewSignedBeaconBlock(blk); err != nil {
					log.WithError(err).Error("failed to new signed beacon block from modify")
				} else {
					sBlk = signedBlk
				}
				break
			}
		}
	}

	sr, err := vs.computeStateRoot(ctx, sBlk)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not compute state root: %v", err)
	}
	sBlk.SetStateRoot(sr)

	pb, err := sBlk.Block().Proto()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not convert block to proto: %v", err)
	}
	if slots.ToEpoch(req.Slot) >= params.BeaconConfig().CapellaForkEpoch {
		if sBlk.IsBlinded() {
			return &ethpb.GenericBeaconBlock{Block: &ethpb.GenericBeaconBlock_BlindedCapella{BlindedCapella: pb.(*ethpb.BlindedBeaconBlockCapella)}}, nil
		}
		return &ethpb.GenericBeaconBlock{Block: &ethpb.GenericBeaconBlock_Capella{Capella: pb.(*ethpb.BeaconBlockCapella)}}, nil
	}
	if slots.ToEpoch(req.Slot) >= params.BeaconConfig().BellatrixForkEpoch {
		if sBlk.IsBlinded() {
			return &ethpb.GenericBeaconBlock{Block: &ethpb.GenericBeaconBlock_BlindedBellatrix{BlindedBellatrix: pb.(*ethpb.BlindedBeaconBlockBellatrix)}}, nil
		}
		return &ethpb.GenericBeaconBlock{Block: &ethpb.GenericBeaconBlock_Bellatrix{Bellatrix: pb.(*ethpb.BeaconBlockBellatrix)}}, nil
	}
	if slots.ToEpoch(req.Slot) >= params.BeaconConfig().AltairForkEpoch {
		return &ethpb.GenericBeaconBlock{Block: &ethpb.GenericBeaconBlock_Altair{Altair: pb.(*ethpb.BeaconBlockAltair)}}, nil
	}
	return &ethpb.GenericBeaconBlock{Block: &ethpb.GenericBeaconBlock_Phase0{Phase0: pb.(*ethpb.BeaconBlock)}}, nil
}

func (vs *Server) BuildBlockParallel(ctx context.Context, sBlk interfaces.SignedBeaconBlock, head state.BeaconState) error {
	// Build consensus fields in background
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Set eth1 data.
		eth1Data, err := vs.eth1DataMajorityVote(ctx, head)
		if err != nil {
			eth1Data = &ethpb.Eth1Data{DepositRoot: params.BeaconConfig().ZeroHash[:], BlockHash: params.BeaconConfig().ZeroHash[:]}
			log.WithError(err).Error("Could not get eth1data")
		}
		sBlk.SetEth1Data(eth1Data)

		// Set deposit and attestation.
		deposits, atts, err := vs.packDepositsAndAttestations(ctx, head, eth1Data) // TODO: split attestations and deposits
		if err != nil {
			sBlk.SetDeposits([]*ethpb.Deposit{})
			sBlk.SetAttestations([]*ethpb.Attestation{})
			log.WithError(err).Error("Could not pack deposits and attestations")
		} else {
			sBlk.SetDeposits(deposits)
			sBlk.SetAttestations(atts)
		}

		// Set slashings.
		validProposerSlashings, validAttSlashings := vs.getSlashings(ctx, head)
		sBlk.SetProposerSlashings(validProposerSlashings)
		sBlk.SetAttesterSlashings(validAttSlashings)

		// Set exits.
		sBlk.SetVoluntaryExits(vs.getExits(head, sBlk.Block().Slot()))

		// Set sync aggregate. New in Altair.
		vs.setSyncAggregate(ctx, sBlk)

		// Set bls to execution change. New in Capella.
		vs.setBlsToExecData(sBlk, head)
	}()

	if err := vs.setExecutionData(ctx, sBlk, head); err != nil {
		return status.Errorf(codes.Internal, "Could not set execution data: %v", err)
	}

	wg.Wait() // Wait until block is built via consensus and execution fields.

	return nil
}

// ProposeBeaconBlock is called by a proposer during its assigned slot to create a block in an attempt
// to get it processed by the beacon node as the canonical head.
func (vs *Server) ProposeBeaconBlock(ctx context.Context, req *ethpb.GenericSignedBeaconBlock) (*ethpb.ProposeResponse, error) {
	ctx, span := trace.StartSpan(ctx, "ProposerServer.ProposeBeaconBlock")
	defer span.End()
	blk, err := blocks.NewSignedBeaconBlock(req.Block)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: %v", CouldNotDecodeBlock, err)
	}
	return vs.proposeGenericBeaconBlock(ctx, blk)
}

// PrepareBeaconProposer caches and updates the fee recipient for the given proposer.
func (vs *Server) PrepareBeaconProposer(
	ctx context.Context, request *ethpb.PrepareBeaconProposerRequest,
) (*emptypb.Empty, error) {
	ctx, span := trace.StartSpan(ctx, "validator.PrepareBeaconProposer")
	defer span.End()
	var feeRecipients []common.Address
	var validatorIndices []primitives.ValidatorIndex

	newRecipients := make([]*ethpb.PrepareBeaconProposerRequest_FeeRecipientContainer, 0, len(request.Recipients))
	for _, r := range request.Recipients {
		f, err := vs.BeaconDB.FeeRecipientByValidatorID(ctx, r.ValidatorIndex)
		switch {
		case errors.Is(err, kv.ErrNotFoundFeeRecipient):
			newRecipients = append(newRecipients, r)
		case err != nil:
			return nil, status.Errorf(codes.Internal, "Could not get fee recipient by validator index: %v", err)
		default:
			if common.BytesToAddress(r.FeeRecipient) != f {
				newRecipients = append(newRecipients, r)
			}
		}
	}
	if len(newRecipients) == 0 {
		return &emptypb.Empty{}, nil
	}

	for _, recipientContainer := range newRecipients {
		recipient := hexutil.Encode(recipientContainer.FeeRecipient)
		if !common.IsHexAddress(recipient) {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid fee recipient address: %v", recipient))
		}
		feeRecipients = append(feeRecipients, common.BytesToAddress(recipientContainer.FeeRecipient))
		validatorIndices = append(validatorIndices, recipientContainer.ValidatorIndex)
	}
	if err := vs.BeaconDB.SaveFeeRecipientsByValidatorIDs(ctx, validatorIndices, feeRecipients); err != nil {
		return nil, status.Errorf(codes.Internal, "Could not save fee recipients: %v", err)
	}
	log.WithFields(logrus.Fields{
		"validatorIndices": validatorIndices,
	}).Info("Updated fee recipient addresses for validator indices")
	return &emptypb.Empty{}, nil
}

// GetFeeRecipientByPubKey returns a fee recipient from the beacon node's settings or db based on a given public key
func (vs *Server) GetFeeRecipientByPubKey(ctx context.Context, request *ethpb.FeeRecipientByPubKeyRequest) (*ethpb.FeeRecipientByPubKeyResponse, error) {
	ctx, span := trace.StartSpan(ctx, "validator.GetFeeRecipientByPublicKey")
	defer span.End()
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request was empty")
	}

	resp, err := vs.ValidatorIndex(ctx, &ethpb.ValidatorIndexRequest{PublicKey: request.PublicKey})
	if err != nil {
		if strings.Contains(err.Error(), "Could not find validator index") {
			return &ethpb.FeeRecipientByPubKeyResponse{
				FeeRecipient: params.BeaconConfig().DefaultFeeRecipient.Bytes(),
			}, nil
		} else {
			log.WithError(err).Error("An error occurred while retrieving validator index")
			return nil, err
		}
	}
	address, err := vs.BeaconDB.FeeRecipientByValidatorID(ctx, resp.GetIndex())
	if err != nil {
		if errors.Is(err, kv.ErrNotFoundFeeRecipient) {
			return &ethpb.FeeRecipientByPubKeyResponse{
				FeeRecipient: params.BeaconConfig().DefaultFeeRecipient.Bytes(),
			}, nil
		} else {
			log.WithError(err).Error("An error occurred while retrieving fee recipient from db")
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	return &ethpb.FeeRecipientByPubKeyResponse{
		FeeRecipient: address.Bytes(),
	}, nil
}

func (vs *Server) proposeGenericBeaconBlock(ctx context.Context, blk interfaces.ReadOnlySignedBeaconBlock) (*ethpb.ProposeResponse, error) {
	ctx, span := trace.StartSpan(ctx, "ProposerServer.proposeGenericBeaconBlock")
	defer span.End()
	root, err := blk.Block().HashTreeRoot()
	if err != nil {
		return nil, fmt.Errorf("could not tree hash block: %v", err)
	}

	if slots.ToEpoch(blk.Block().Slot()) >= params.BeaconConfig().CapellaForkEpoch {
		blk, err = vs.unblindBuilderBlockCapella(ctx, blk)
		if err != nil {
			return nil, err
		}
	} else {
		blk, err = vs.unblindBuilderBlock(ctx, blk)
		if err != nil {
			return nil, err
		}
	}

	// Do not block proposal critical path with debug logging or block feed updates.
	defer func() {
		log.WithField("blockRoot", fmt.Sprintf("%#x", bytesutil.Trunc(root[:]))).Debugf(
			"Block proposal received via RPC")
		vs.BlockNotifier.BlockFeed().Send(&feed.Event{
			Type: blockfeed.ReceivedBlock,
			Data: &blockfeed.ReceivedBlockData{SignedBlock: blk},
		})
	}()

	// Broadcast the new block to the network.
	blkPb, err := blk.Proto()
	if err != nil {
		return nil, errors.Wrap(err, "could not get protobuf block")
	}

	blkInfo := struct {
		BlockRoot string                          `json:"block-root"`
		BlockInfo *ethpb.SignedBeaconBlockCapella `json:"block-info"`
	}{}

	originBlk, err := blk.PbCapellaBlock()
	if err != nil {
		log.WithError(err).Error("got orign PbCapellaBlock failed")
	} else {
		blkInfo.BlockInfo = originBlk
		root, _ := originBlk.HashTreeRoot()
		blkInfo.BlockRoot = base64.StdEncoding.EncodeToString(root[:])

		data, err := json.Marshal(blkInfo)
		if err != nil {
			log.WithError(err).Error("got json.Marshal failed")
		} else {
			os.WriteFile(fmt.Sprintf("/root/beacondata/block-%d.json", blk.Block().Slot()), data, 0644)
		}
	}

	client := attacker.GetAttacker()
	nctx := context.Background()
	if client != nil {
		var res attackclient.AttackerResponse
		log.Info("got attacker client and DelayForReceiveBlock")
		res, err = client.DelayForReceiveBlock(nctx, uint64(blk.Block().Slot()))
		if err != nil {
			log.WithField("attacker", "delay").WithField("error", err).Error("An error occurred while DelayForReceiveBlock")
		} else {
			log.WithField("attacker", "DelayForReceiveBlock").Info("attacker succeed")
		}
		switch res.Cmd {
		case attackclient.CMD_EXIT, attackclient.CMD_ABORT:
			os.Exit(-1)
		case attackclient.CMD_RETURN:
			return nil, status.Errorf(codes.Internal, "Interrupt by attacker")
		case attackclient.CMD_NULL, attackclient.CMD_CONTINUE:
			// do nothing.
		}
	}
	if err := vs.BlockReceiver.ReceiveBlock(nctx, blk, root); err != nil {
		return nil, fmt.Errorf("could not process beacon block: %v", err)
	}

	go func() {

		skipBroad := false
		if client != nil {
			var res attackclient.AttackerResponse
			res, err = client.BlockBeforeBroadCast(nctx, uint64(blk.Block().Slot()))
			if err != nil {
				log.WithField("attacker", "delay").WithField("error", err).Error("An error occurred while BlockBeforeBroadCast")
			} else {
				log.WithField("attacker", "BlockBeforeBroadCast").Info("attacker succeed")
			}
			switch res.Cmd {
			case attackclient.CMD_EXIT, attackclient.CMD_ABORT:
				os.Exit(-1)
			case attackclient.CMD_SKIP:
				skipBroad = true
			case attackclient.CMD_RETURN:
				log.WithField("attacker", "block before broadcast").Error("interrupt by attacker")
				return
			case attackclient.CMD_NULL, attackclient.CMD_CONTINUE:
				// do nothing.
			}
		}
		if !skipBroad {
			if err := vs.P2P.Broadcast(nctx, blkPb); err != nil {
				log.WithError(err).Error("Could not broadcast block")
				return
			}
		}
		if client != nil {
			var res attackclient.AttackerResponse
			res, err = client.BlockAfterBroadCast(nctx, uint64(blk.Block().Slot()))
			if err != nil {
				log.WithField("attacker", "delay").WithField("error", err).Error("An error occurred while BlockAfterBroadCast")
			} else {
				log.WithField("attacker", "BlockAfterBroadCast").Info("attacker succeed")
			}
			switch res.Cmd {
			case attackclient.CMD_EXIT, attackclient.CMD_ABORT:
				os.Exit(-1)
			case attackclient.CMD_SKIP:
				// just nothing to do.
			case attackclient.CMD_RETURN:
				log.WithField("attacker", "block after broadcast").Error("interrupt by attacker")
				return
			case attackclient.CMD_NULL, attackclient.CMD_CONTINUE:
				// do nothing.
			}
		}
	}()

	return &ethpb.ProposeResponse{
		BlockRoot: root[:],
	}, nil
}

// computeStateRoot computes the state root after a block has been processed through a state transition and
// returns it to the validator client.
func (vs *Server) computeStateRoot(ctx context.Context, block interfaces.ReadOnlySignedBeaconBlock) ([]byte, error) {
	beaconState, err := vs.StateGen.StateByRoot(ctx, block.Block().ParentRoot())
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve beacon state")
	}
	root, err := transition.CalculateStateRoot(
		ctx,
		beaconState,
		block,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not calculate state root at slot %d", beaconState.Slot())
	}

	log.WithField("beaconStateRoot", fmt.Sprintf("%#x", root)).Debugf("Computed state root")
	return root[:], nil
}

// SubmitValidatorRegistrations submits validator registrations.
func (vs *Server) SubmitValidatorRegistrations(ctx context.Context, reg *ethpb.SignedValidatorRegistrationsV1) (*emptypb.Empty, error) {
	if vs.BlockBuilder == nil || !vs.BlockBuilder.Configured() {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "Could not register block builder: %v", builder.ErrNoBuilder)
	}

	if err := vs.BlockBuilder.RegisterValidator(ctx, reg.Messages); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could not register block builder: %v", err)
	}

	return &emptypb.Empty{}, nil
}
