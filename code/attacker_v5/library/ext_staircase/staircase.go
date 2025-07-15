package ext_staircase

import (
	"context"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/types"
	"strconv"
	"time"
)

type Instance struct {
}

func (o *Instance) Name() string {
	return "ext_staircase"
}

func (o *Instance) Description() string {
	desc_eng := `Extended staircase attack`
	return desc_eng
}

func (o *Instance) Run(ctx context.Context, params types.LibraryParams, feedbacker types.FeedBacker) {
	log.WithField("name", o.Name()).Info("start to run strategy")
	ticker := time.NewTicker(time.Second * 3)
	attacker := params.Attacker
	history := make(map[int]bool)
	started := false
	for {
		select {
		case <-ctx.Done():
			log.WithField("name", o.Name()).Info("stop to run strategy")
			return
		case <-ticker.C:
			slot := attacker.GetCurSlot()
			epoch := common.SlotToEpoch(int64(slot))
			nextEpoch := epoch + 1
			log.WithFields(log.Fields{
				"slot":      slot,
				"nextEpoch": nextEpoch,
			}).Info("get slot")

			if _, ok := history[int(nextEpoch)]; ok {
				continue
			}

			if nextEpoch < 3 {
				log.WithField("epoch", nextEpoch).Info("skip to generate strategy")
				history[int(nextEpoch)] = true
				continue
			}
			nextDuties, err := attacker.GetEpochDuties(nextEpoch)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"epoch": nextEpoch,
				}).Error("failed to get duties")
				continue
			}
			if nextEpoch%10 < 7 {
				cas := 0
				if started {
					cas = 1
					started = false
				} else {
					started = true
				}
				strategy := types.Strategy{}
				strategy.Uid = uuid.NewString()
				strategy.Slots = GenSlotStrategy(nextDuties, params.FillterHackerDuties(nextDuties), cas)
				strategy.Category = o.Name()
				if err = attacker.UpdateStrategy(strategy); err != nil {
					log.WithField("error", err).Error("failed to update strategy")
				} else {
					log.WithFields(log.Fields{
						"epoch":    nextEpoch,
						"strategy": strategy,
					}).Info("update strategy successfully")
					history[int(nextEpoch)] = true
				}
			} else {
				if started {
					started = false
				}
				strategy := types.Strategy{}
				strategy.Uid = uuid.NewString()
				strategy.Slots = GenSlotStrategyOnB(nextDuties, params.FillterHackerDuties(nextDuties), nextEpoch)
				strategy.Category = o.Name()
				if err = attacker.UpdateStrategy(strategy); err != nil {
					log.WithField("error", err).Error("failed to update strategy")
				} else {
					log.WithFields(log.Fields{
						"epoch":    nextEpoch,
						"strategy": strategy,
					}).Info("update strategy successfully")
					history[int(nextEpoch)] = true
				}
			}
		}
	}
}

func getLatestHackerSlot(duties []types.ProposerDuty, param types.LibraryParams) int {
	latest, _ := strconv.Atoi(duties[0].Slot)
	for _, duty := range duties {
		idx, _ := strconv.Atoi(duty.ValidatorIndex)
		slot, _ := strconv.Atoi(duty.Slot)
		if !param.IsHackValidator(idx) {
			continue
		}
		if slot > latest {
			latest = slot
		}
	}
	return latest

}

func checkFirstByzSlot(duties []types.ProposerDuty, param types.LibraryParams) bool {
	firstproposerindex, _ := strconv.Atoi(duties[0].ValidatorIndex)
	if !param.IsHackValidator(firstproposerindex) {
		return false
	}
	return true
}
