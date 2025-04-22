package utils

import (
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/types"
	"strconv"
)

type attackerInc struct {
	backend types.ServiceBackend
}

func (a attackerInc) UpdateStrategy(strategy types.Strategy) error {
	return a.backend.UpdateStrategy(strategy)
}

func (a attackerInc) GetStrategyFeedback(uid string) (types.FeedBackInfo, error) {
	return a.backend.GetFeedBack(uid)
}

func (a attackerInc) GetChainBaseInfo() types.ChainBaseInfo {
	return common.GetChainBaseInfo()
}

func (a attackerInc) GetCurSlot() int64 {
	h, err := a.backend.GetLatestBeaconHeader()
	if err != nil {
		return a.backend.GetCurSlot()
	} else {
		slot, _ := strconv.ParseInt(h.Header.Message.Slot, 10, 64)
		return slot
	}
}

func (a attackerInc) GetEpochDuties(epoch int64) ([]types.ProposerDuty, error) {
	return a.backend.GetProposeDuties(int(epoch))
}

func WrapToAttacker(backend types.ServiceBackend) types.AttackerInc {
	return &attackerInc{
		backend: backend,
	}
}
