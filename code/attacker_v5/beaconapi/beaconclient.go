package beaconapi

import (
	"context"
	"errors"
	"fmt"
	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	lru "github.com/hashicorp/golang-lru"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/types"
	"strconv"
	"time"
)

const (
	SLOTS_PER_EPOCH  = "SLOTS_PER_EPOCH"
	SECONDS_PER_SLOT = "SECONDS_PER_SLOT"
)

var (
	validatorListCacheKey = "validator_list"
)

type BeaconGwClient struct {
	endpoint string
	config   map[string]string
	service  eth2client.Service
	cache    *lru.Cache
}

func NewBeaconGwClient(endpoint string) *BeaconGwClient {
	cache, _ := lru.New(100)
	return &BeaconGwClient{
		endpoint: endpoint,
		config:   make(map[string]string),
		cache:    cache,
	}
}

func (b *BeaconGwClient) GetIntConfig(key string) (int, error) {
	config := b.GetBeaconConfig()
	if v, exist := config[key]; !exist {
		return 0, nil
	} else {
		return strconv.Atoi(v)
	}
}

func (b *BeaconGwClient) GetBeaconConfig() map[string]string {
	if len(b.config) == 0 {
		config, err := b.GetSpec()
		if err != nil {
			log.WithError(err).Error("get beacon spec failed")
			return nil
		}
		b.config = make(map[string]string)
		for key, v := range config {
			switch v.(type) {
			case time.Duration:
				b.config[key] = strconv.FormatFloat(v.(time.Duration).Seconds(), 'f', -1, 64)
			case time.Time:
				b.config[key] = strconv.FormatInt(v.(time.Time).Unix(), 10)
			case []uint8:
				b.config[key] = fmt.Sprintf("0x%#x", v.([]uint8))
			case int:
				b.config[key] = strconv.Itoa(v.(int))
			case uint64:
				b.config[key] = strconv.FormatUint(v.(uint64), 10)
			case int64:
				b.config[key] = strconv.FormatInt(v.(int64), 10)
			case float64:
				b.config[key] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
			case string:
				b.config[key] = v.(string)
			case phase0.Version:
				b.config[key] = fmt.Sprintf("%#x", v.(phase0.Version))
			case phase0.DomainType:
				b.config[key] = fmt.Sprintf("%#x", v.(phase0.DomainType))
			default:
				log.Warnf("unknown beacon config key %s type %T", key, v)
			}

		}
	}
	return b.config
}

func (b *BeaconGwClient) getLatestBeaconHeader() (*apiv1.BeaconBlockHeader, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.BeaconBlockHeadersProvider).BeaconBlockHeader(context.Background(), &api.BeaconBlockHeaderOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		Block: "head",
	})
	if err != nil {
		log.WithError(err).Error("get latest beacon header failed")
		return nil, err
	}
	return res.Data, nil
}

func (b *BeaconGwClient) GetValidatorsList() ([]*phase0.Validator, error) {
	if v, ok := b.cache.Get(validatorListCacheKey); ok {
		return v.([]*phase0.Validator), nil
	}
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.BeaconStateProvider).BeaconState(context.Background(), &api.BeaconStateOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		State: "head",
	})
	if err != nil {
		log.WithError(err).Error("get beacon state failed")
		return nil, err
	}
	vals, err := res.Data.Validators()
	if err != nil {
		log.WithError(err).Error("get validators failed")
		return nil, err
	}
	b.cache.Add(validatorListCacheKey, vals)

	return vals, nil
}

func (b *BeaconGwClient) GetLatestValidators() (*spec.VersionedBeaconState, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.BeaconStateProvider).BeaconState(context.Background(), &api.BeaconStateOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		State: "head",
	})
	if err != nil {
		log.WithError(err).Error("get beacon state failed")
		return nil, err
	}

	return res.Data, nil
}

func (b *BeaconGwClient) GetLatestBeaconHeader() (types.BeaconHeaderInfo, error) {
	h, err := b.getLatestBeaconHeader()
	if err != nil {
		return types.BeaconHeaderInfo{}, err
	}
	header := types.BeaconHeaderInfo{
		Root:      h.Root.String(),
		Canonical: h.Canonical,
	}
	header.Header.Signature = h.Header.Signature.String()
	header.Header.Message.Slot = strconv.FormatInt(int64(h.Header.Message.Slot), 10)
	header.Header.Message.ProposerIndex = strconv.FormatInt(int64(h.Header.Message.ProposerIndex), 10)
	header.Header.Message.ParentRoot = h.Header.Message.ParentRoot.String()
	header.Header.Message.StateRoot = h.Header.Message.StateRoot.String()
	header.Header.Message.BodyRoot = h.Header.Message.BodyRoot.String()
	return header, nil
}

func (b *BeaconGwClient) getAllValReward(epoch int) (*apiv1.AttestationRewards, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.AttestationRewardsProvider).AttestationRewards(context.Background(), &api.AttestationRewardsOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		Epoch: phase0.Epoch(epoch),
	})
	if err != nil {
		log.WithField("epoch", epoch).WithError(err).Error("get val reward failed")
		return nil, err
	}

	return res.Data, nil
}

// default grpc-gateway port is 3500
func (b *BeaconGwClient) GetAllValReward(epoch int) (*apiv1.AttestationRewards, error) {
	info, err := b.getAllValReward(epoch)
	if err != nil {
		return nil, err
	}
	return info, err
}

func (b *BeaconGwClient) getProposerDuties(epoch int) ([]*apiv1.ProposerDuty, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.ProposerDutiesProvider).ProposerDuties(context.Background(), &api.ProposerDutiesOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		Epoch: phase0.Epoch(epoch),
	})
	if err != nil {
		log.WithError(err).Error("get proposer duties failed")
		return nil, err
	}

	return res.Data, nil
}

// /eth/v1/validator/duties/proposer/:epoch
func (b *BeaconGwClient) GetProposerDuties(epoch int) ([]types.ProposerDuty, error) {
	var duties = make([]types.ProposerDuty, 0)
	res, err := b.getProposerDuties(epoch)
	if err != nil {
		return duties, err
	}
	for _, duty := range res {
		duties = append(duties, types.ProposerDuty{
			Pubkey:         duty.PubKey.String(),
			Slot:           strconv.FormatInt(int64(duty.Slot), 10),
			ValidatorIndex: strconv.FormatInt(int64(duty.ValidatorIndex), 10),
		})
	}

	return duties, err
}

func (b *BeaconGwClient) getAttesterDuties(epoch int, vals []int) ([]*apiv1.AttesterDuty, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	indices := make([]phase0.ValidatorIndex, len(vals))
	for _, val := range vals {
		indices = append(indices, phase0.ValidatorIndex(val))
	}
	if len(indices) == 0 {
		// get validators list
		valList, err := b.GetValidatorsList()
		if err != nil {
			log.WithError(err).Error("get validators failed")
			return nil, err
		}
		indices = make([]phase0.ValidatorIndex, len(valList))
		for i, _ := range valList {
			indices[i] = phase0.ValidatorIndex(i)
		}
	}
	res, err := service.(eth2client.AttesterDutiesProvider).AttesterDuties(context.Background(), &api.AttesterDutiesOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		Epoch:   phase0.Epoch(epoch),
		Indices: indices,
	})
	if err != nil {
		log.WithError(err).Error("get attester duties failed")
		return nil, err
	}

	return res.Data, nil
}

// POST /eth/v1/validator/duties/attester/:epoch
func (b *BeaconGwClient) GetAttesterDuties(epoch int, vals []int) ([]types.AttestDuty, error) {
	res, err := b.getAttesterDuties(epoch, vals)
	if err != nil {
		return nil, err
	}
	duties := make([]types.AttestDuty, 0)
	for _, duty := range res {
		duties = append(duties, types.AttestDuty{
			Slot:                    strconv.FormatInt(int64(duty.Slot), 10),
			Pubkey:                  duty.PubKey.String(),
			ValidatorIndex:          strconv.FormatInt(int64(duty.ValidatorIndex), 10),
			CommitteeIndex:          strconv.FormatInt(int64(duty.CommitteeIndex), 10),
			CommitteeLength:         strconv.FormatInt(int64(duty.CommitteeLength), 10),
			CommitteesAtSlot:        strconv.FormatInt(int64(duty.CommitteesAtSlot), 10),
			ValidatorCommitteeIndex: strconv.FormatInt(int64(duty.ValidatorCommitteeIndex), 10),
		})
	}
	return duties, nil
}

func (b *BeaconGwClient) GetNextEpochProposerDuties() ([]types.ProposerDuty, error) {
	latestHeader, err := b.GetLatestBeaconHeader()
	if err != nil {
		return nil, err
	}
	slotPerEpoch, _ := b.GetIntConfig(SLOTS_PER_EPOCH)
	curSlot, _ := strconv.Atoi(latestHeader.Header.Message.Slot)
	epoch := curSlot / slotPerEpoch
	return b.GetProposerDuties(epoch + 1)
}

func (b *BeaconGwClient) GetCurrentEpochProposerDuties() ([]types.ProposerDuty, error) {
	latestHeader, err := b.GetLatestBeaconHeader()
	if err != nil {
		return nil, err
	}
	slotPerEpoch, _ := b.GetIntConfig(SLOTS_PER_EPOCH)
	curSlot, _ := strconv.Atoi(latestHeader.Header.Message.Slot)
	epoch := curSlot / slotPerEpoch
	return b.GetProposerDuties(epoch)
}

func (b *BeaconGwClient) GetCurrentEpochAttestDuties() ([]types.AttestDuty, error) {
	latestHeader, err := b.GetLatestBeaconHeader()
	if err != nil {
		return nil, err
	}
	slotPerEpoch, _ := b.GetIntConfig(SLOTS_PER_EPOCH)
	curSlot, _ := strconv.Atoi(latestHeader.Header.Message.Slot)
	epoch := curSlot / slotPerEpoch
	vals := make([]int, 64)
	for i := 0; i < len(vals); i++ {
		vals[i] = i
	}
	return b.GetAttesterDuties(epoch, vals)
}

func (b *BeaconGwClient) GetNextEpochAttestDuties() ([]types.AttestDuty, error) {
	latestHeader, err := b.GetLatestBeaconHeader()
	if err != nil {
		return nil, err
	}
	slotPerEpoch, _ := b.GetIntConfig(SLOTS_PER_EPOCH)
	curSlot, _ := strconv.Atoi(latestHeader.Header.Message.Slot)
	epoch := curSlot / slotPerEpoch
	vals := make([]int, 64)
	for i := 0; i < len(vals); i++ {
		vals[i] = i
	}
	return b.GetAttesterDuties(epoch+1, vals)
}

func (b *BeaconGwClient) getBlockReward(slot int) (*apiv1.BlockRewards, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.BlockRewardsProvider).BlockRewards(context.Background(), &api.BlockRewardsOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		Block: fmt.Sprintf("%d", slot),
	})
	if err != nil {
		log.WithField("slot", slot).WithError(err).Error("get block reward failed")
		return nil, err
	}
	return res.Data, nil
}

func (b *BeaconGwClient) GetBlockReward(slot int) (types.BlockRewardInfo, error) {
	res, err := b.getBlockReward(slot)
	if err != nil {
		return types.BlockRewardInfo{}, err
	}
	if res == nil {
		return types.BlockRewardInfo{}, errors.New("block reward not found")
	}
	reward := types.BlockRewardInfo{
		ProposerIndex:     uint64(res.ProposerIndex),
		Total:             uint64(res.Total),
		Attestations:      uint64(res.Attestations),
		SyncAggregate:     uint64(res.SyncAggregate),
		ProposerSlashings: uint64(res.ProposerSlashings),
		AttesterSlashings: uint64(res.AttesterSlashings),
	}

	return reward, nil
}

func (b *BeaconGwClient) getSlotRoot(slot int64) (*phase0.Root, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.BeaconBlockRootProvider).BeaconBlockRoot(context.Background(), &api.BeaconBlockRootOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		Block: fmt.Sprintf("%d", slot),
	})
	if err != nil {
		log.WithError(err).Error("getSlotRoot failed")
		return nil, err
	}
	return res.Data, nil
}

func (b *BeaconGwClient) GetSlotRoot(slot int64) (string, error) {
	root, err := b.getSlotRoot(slot)
	if err != nil {
		return "", err
	}
	if root == nil {
		return "0x", nil
	}
	return root.String(), nil
}

func (b *BeaconGwClient) getService() (eth2client.Service, error) {
	if b.service == nil {
		service, err := NewClient(context.Background(), b.endpoint)
		if err != nil {
			log.WithError(err).Error("create eth2client failed")
			return nil, err
		}
		b.service = service
	}
	return b.service, nil
}

func (b *BeaconGwClient) MonitorReorgEvent() chan *apiv1.ChainReorgEvent {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil
	}
	ch := make(chan *apiv1.ChainReorgEvent, 100)
	go func() {
		service.(eth2client.EventsProvider).Events(context.Background(), []string{"chain_reorg"}, func(event *apiv1.Event) {
			if ev, ok := event.Data.(*apiv1.ChainReorgEvent); !ok {
				log.Error("Failed to unmarshal reorg event")
				return
			} else {
				ch <- ev
			}
			return
		})
	}()
	return ch
}

func (b *BeaconGwClient) GetBlockHeaderById(id string) (*apiv1.BeaconBlockHeader, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	opts := &api.BeaconBlockHeaderOpts{
		Block: id,
	}
	res, err := service.(eth2client.BeaconBlockHeadersProvider).BeaconBlockHeader(context.Background(), opts)
	if err != nil {
		log.WithError(err).Error("get block header failed")
		return &apiv1.BeaconBlockHeader{}, err
	}
	return res.Data, nil
}

func (b *BeaconGwClient) GetDenebBlockBySlot(slot uint64) (*deneb.SignedBeaconBlock, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.SignedBeaconBlockProvider).SignedBeaconBlock(context.Background(), &api.SignedBeaconBlockOpts{
		Block: fmt.Sprintf("%d", slot),
	})
	if err != nil {
		log.WithError(err).Error("get block failed")
		return nil, err
	}
	return res.Data.Deneb, nil
}

func (b *BeaconGwClient) GetCapellaBlockBySlot(slot uint64) (*capella.SignedBeaconBlock, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}

	res, err := service.(eth2client.SignedBeaconBlockProvider).SignedBeaconBlock(context.Background(), &api.SignedBeaconBlockOpts{
		Block: fmt.Sprintf("%d", slot),
	})
	if err != nil {
		log.WithError(err).Error("get block failed")
		return nil, err
	}
	return res.Data.Capella, nil
}

func (b *BeaconGwClient) GetSpec() (map[string]any, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.SpecProvider).Spec(context.Background(), &api.SpecOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
	})
	if err != nil {
		log.WithError(err).Error("get genesis failed")
		return nil, err
	}
	return res.Data, nil
}

func (b *BeaconGwClient) GetGenesis() (*apiv1.Genesis, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.GenesisProvider).Genesis(context.Background(), &api.GenesisOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
	})
	if err != nil {
		log.WithError(err).Error("get genesis failed")
		return nil, err
	}
	return res.Data, nil
}
