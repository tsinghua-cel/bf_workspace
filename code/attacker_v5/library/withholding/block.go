package withholding

import (
	"fmt"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/types"
	"strconv"
)

func calcTargetTime(slot int, targetSlot int64, offset int) int64 {
	return common.TimeToSlot(targetSlot)*1000 - 3*1000 + 100*int64(offset)
}

func BlockStrategy(cur, start, end int, actions map[string]string) {
	// delay block to next epoch + 8 slots.
	targetSlot := common.EpochEnd(common.SlotToEpoch(int64(cur))) + 8
	targetTime := calcTargetTime(cur, targetSlot, cur-start)
	actions["BlockBeforeBroadCast"] = fmt.Sprintf("delayToMilliTime:%d", targetTime)
}

func GenSlotStrategy(allHacks []interface{}) []types.SlotStrategy {
	if len(allHacks) == 0 {
		return nil
	}
	strategys := make([]types.SlotStrategy, 0)

	// only use the last subduties
	subduties := allHacks[len(allHacks)-1]

	duties := subduties.([]types.ProposerDuty)
	start, _ := strconv.Atoi(duties[0].Slot)
	end, _ := strconv.Atoi(duties[len(duties)-1].Slot)

	for i := 0; i < len(duties); i++ {
		slot, _ := strconv.Atoi(duties[i].Slot)
		strategy := types.SlotStrategy{
			Slot:    duties[i].Slot,
			Level:   1,
			Actions: make(map[string]string),
		}
		BlockStrategy(slot, start, end, strategy.Actions)
		strategys = append(strategys, strategy)
	}

	return strategys

}
