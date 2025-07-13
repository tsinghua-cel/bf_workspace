package unrealizedhonest

import (
	"fmt"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/types"
	"strconv"
)

func calcTargetTime(slot int, targetSlot int64, offset int) int64 {
	return common.TimeToSlot(targetSlot)*1000 - 3*1000 + 100*int64(offset)
}

func GenSlotStrategy(duties []interface{}) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	duty := duties[0].([]types.ProposerDuty)
	currentslot, _ := strconv.Atoi(duty[0].Slot)
	begin := currentslot - 10
	// all honest proposer block broadcast delay to currentslot.
	for i := 0; i < 10; i++ {

		slotStrategy := types.SlotStrategy{
			Slot:    fmt.Sprintf("%d", begin+i),
			Level:   1,
			Actions: make(map[string]string),
		}
		targetTime := calcTargetTime(begin+i, int64(currentslot), i)
		slotStrategy.Actions["BlockBeforeBroadCast"] = fmt.Sprintf("delayToMilliTime:%d", targetTime)
		strategys = append(strategys, slotStrategy)
	}

	return strategys

}
