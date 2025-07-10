package staircaseii

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
	return "staircaseii"
}

func (o *Instance) Description() string {
	desc_eng := `	Assume that the current epoch = 0, then in epoch = 1, the votes of all 
	malicious validators are not broadcast;
	In epoch = 2, the votes of all malicious validators are not broadcast;
	When the last malicious node in epoch = 2 produces a block, package the votes of
	all malicious validators before and broadcast the block at the last slot of epoch = 4.`
	return desc_eng
}

func (o *Instance) Run(ctx context.Context, params types.LibraryParams, feedbacker types.FeedBacker) {
	log.WithField("name", o.Name()).Info("start to run strategy")
	ticker := time.NewTicker(time.Second * 3)
	attacker := params.Attacker
	history := make(map[int]bool)
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

			{
				duties, err := attacker.GetEpochDuties(nextEpoch)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
						"epoch": nextEpoch,
					}).Error("failed to get duties")
					continue
				}
				if nextEpoch < 3 {
					log.WithField("epoch", nextEpoch).Info("skip to generate strategy")
					history[int(nextEpoch)] = true
					continue
				}
				strategy := types.Strategy{}
				strategy.Uid = uuid.NewString()
				strategy.Slots = GenSlotStrategy(params.FillterHackerDuties(duties), nextEpoch)
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
