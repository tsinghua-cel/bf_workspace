package unrealized

import (
	"context"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/types"
	"time"
)

type Instance struct {
}

func (o *Instance) Name() string {
	return "unrealized"
}

func (o *Instance) Description() string {
	desc_eng := `Unrealized justification reorg attack`
	return desc_eng
}

func (o *Instance) Run(ctx context.Context, params types.LibraryParams, feedbacker types.FeedBacker) {
	log.WithField("name", o.Name()).Info("start to run strategy")
	attacker := params.Attacker
	history := make(map[int]bool)
	ticker := time.NewTicker(time.Second * 3)
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
			// get next epoch duties
			duties, err := attacker.GetEpochDuties(nextEpoch)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"epoch": nextEpoch,
				}).Error("failed to get next duties")
				continue
			}
			if hackDuties, happen := CheckDuties(params, duties); happen {
				strategy := types.Strategy{}
				strategy.Uid = uuid.NewString()
				strategy.Slots = GenSlotStrategy(hackDuties)
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
				history[int(nextEpoch)] = true
			}
		}
	}
}
