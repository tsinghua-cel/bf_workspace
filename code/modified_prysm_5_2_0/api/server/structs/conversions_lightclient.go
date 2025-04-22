package structs

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	"github.com/prysmaticlabs/prysm/v5/runtime/version"
)

func LightClientUpdateFromConsensus(update interfaces.LightClientUpdate) (*LightClientUpdate, error) {
	attestedHeader, err := lightClientHeaderToJSON(update.AttestedHeader())
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal attested light client header")
	}
	finalizedHeader, err := lightClientHeaderToJSON(update.FinalizedHeader())
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal finalized light client header")
	}

	var scBranch [][32]byte
	var finalityBranch [][32]byte
	if update.Version() >= version.Electra {
		scb, err := update.NextSyncCommitteeBranchElectra()
		if err != nil {
			return nil, err
		}
		scBranch = scb[:]
		fb, err := update.FinalityBranchElectra()
		if err != nil {
			return nil, err
		}
		finalityBranch = fb[:]
	} else {
		scb, err := update.NextSyncCommitteeBranch()
		if err != nil {
			return nil, err
		}
		scBranch = scb[:]
		fb, err := update.FinalityBranch()
		if err != nil {
			return nil, err
		}
		finalityBranch = fb[:]
	}

	return &LightClientUpdate{
		AttestedHeader:          attestedHeader,
		NextSyncCommittee:       SyncCommitteeFromConsensus(update.NextSyncCommittee()),
		NextSyncCommitteeBranch: branchToJSON(scBranch),
		FinalizedHeader:         finalizedHeader,
		FinalityBranch:          branchToJSON(finalityBranch),
		SyncAggregate:           SyncAggregateFromConsensus(update.SyncAggregate()),
		SignatureSlot:           fmt.Sprintf("%d", update.SignatureSlot()),
	}, nil
}

func LightClientFinalityUpdateFromConsensus(update interfaces.LightClientFinalityUpdate) (*LightClientFinalityUpdate, error) {
	attestedHeader, err := lightClientHeaderToJSON(update.AttestedHeader())
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal attested light client header")
	}
	finalizedHeader, err := lightClientHeaderToJSON(update.FinalizedHeader())
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal finalized light client header")
	}

	var finalityBranch [][32]byte
	if update.Version() >= version.Electra {
		b, err := update.FinalityBranchElectra()
		if err != nil {
			return nil, err
		}
		finalityBranch = b[:]
	} else {
		b, err := update.FinalityBranch()
		if err != nil {
			return nil, err
		}
		finalityBranch = b[:]
	}

	return &LightClientFinalityUpdate{
		AttestedHeader:  attestedHeader,
		FinalizedHeader: finalizedHeader,
		FinalityBranch:  branchToJSON(finalityBranch),
		SyncAggregate:   SyncAggregateFromConsensus(update.SyncAggregate()),
		SignatureSlot:   fmt.Sprintf("%d", update.SignatureSlot()),
	}, nil
}

func LightClientOptimisticUpdateFromConsensus(update interfaces.LightClientOptimisticUpdate) (*LightClientOptimisticUpdate, error) {
	attestedHeader, err := lightClientHeaderToJSON(update.AttestedHeader())
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal attested light client header")
	}

	return &LightClientOptimisticUpdate{
		AttestedHeader: attestedHeader,
		SyncAggregate:  SyncAggregateFromConsensus(update.SyncAggregate()),
		SignatureSlot:  fmt.Sprintf("%d", update.SignatureSlot()),
	}, nil
}

func branchToJSON[S [][32]byte](branchBytes S) []string {
	if branchBytes == nil {
		return nil
	}
	branch := make([]string, len(branchBytes))
	for i, root := range branchBytes {
		branch[i] = hexutil.Encode(root[:])
	}
	return branch
}

func lightClientHeaderToJSON(header interfaces.LightClientHeader) (json.RawMessage, error) {
	// In the case that a finalizedHeader is nil.
	if header == nil {
		return nil, nil
	}

	var result any

	switch v := header.Version(); v {
	case version.Altair:
		result = &LightClientHeader{Beacon: BeaconBlockHeaderFromConsensus(header.Beacon())}
	case version.Capella:
		exInterface, err := header.Execution()
		if err != nil {
			return nil, err
		}
		ex, ok := exInterface.Proto().(*enginev1.ExecutionPayloadHeaderCapella)
		if !ok {
			return nil, fmt.Errorf("execution data is not %T", &enginev1.ExecutionPayloadHeaderCapella{})
		}
		execution, err := ExecutionPayloadHeaderCapellaFromConsensus(ex)
		if err != nil {
			return nil, err
		}
		executionBranch, err := header.ExecutionBranch()
		if err != nil {
			return nil, err
		}
		result = &LightClientHeaderCapella{
			Beacon:          BeaconBlockHeaderFromConsensus(header.Beacon()),
			Execution:       execution,
			ExecutionBranch: branchToJSON(executionBranch[:]),
		}
	case version.Deneb:
		exInterface, err := header.Execution()
		if err != nil {
			return nil, err
		}
		ex, ok := exInterface.Proto().(*enginev1.ExecutionPayloadHeaderDeneb)
		if !ok {
			return nil, fmt.Errorf("execution data is not %T", &enginev1.ExecutionPayloadHeaderDeneb{})
		}
		execution, err := ExecutionPayloadHeaderDenebFromConsensus(ex)
		if err != nil {
			return nil, err
		}
		executionBranch, err := header.ExecutionBranch()
		if err != nil {
			return nil, err
		}
		result = &LightClientHeaderDeneb{
			Beacon:          BeaconBlockHeaderFromConsensus(header.Beacon()),
			Execution:       execution,
			ExecutionBranch: branchToJSON(executionBranch[:]),
		}
	case version.Electra:
		exInterface, err := header.Execution()
		if err != nil {
			return nil, err
		}
		ex, ok := exInterface.Proto().(*enginev1.ExecutionPayloadHeaderElectra)
		if !ok {
			return nil, fmt.Errorf("execution data is not %T", &enginev1.ExecutionPayloadHeaderElectra{})
		}
		execution, err := ExecutionPayloadHeaderElectraFromConsensus(ex)
		if err != nil {
			return nil, err
		}
		executionBranch, err := header.ExecutionBranch()
		if err != nil {
			return nil, err
		}
		result = &LightClientHeaderDeneb{
			Beacon:          BeaconBlockHeaderFromConsensus(header.Beacon()),
			Execution:       execution,
			ExecutionBranch: branchToJSON(executionBranch[:]),
		}
	default:
		return nil, fmt.Errorf("unsupported header version %s", version.String(v))
	}

	return json.Marshal(result)
}

func LightClientBootstrapFromConsensus(bootstrap interfaces.LightClientBootstrap) (*LightClientBootstrap, error) {
	header, err := lightClientHeaderToJSON(bootstrap.Header())
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal light client header")
	}

	var scBranch [][32]byte
	if bootstrap.Version() >= version.Electra {
		b, err := bootstrap.CurrentSyncCommitteeBranchElectra()
		if err != nil {
			return nil, err
		}
		scBranch = b[:]
	} else {
		b, err := bootstrap.CurrentSyncCommitteeBranch()
		if err != nil {
			return nil, err
		}
		scBranch = b[:]
	}

	return &LightClientBootstrap{
		Header:                     header,
		CurrentSyncCommittee:       SyncCommitteeFromConsensus(bootstrap.CurrentSyncCommittee()),
		CurrentSyncCommitteeBranch: branchToJSON(scBranch),
	}, nil
}
