package selfishhonest

import (
	"fmt"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/types"
)

func GenSlotStrategy(duties []interface{}) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	for i := 0; i < len(duties); i++ {
		duty := duties[i].([]types.ProposerDuty)
		if len(duty) != 3 {
			continue
		}
		//a := duty[0]
		b := duty[1]
		//c := duty[2]
		s1 := types.SlotStrategy{
			Slot:    b.Slot,
			Level:   1,
			Actions: make(map[string]string),
		}

		s1.Actions[string(types.BlockBeforeBroadCast)] = fmt.Sprintf("delayWithSecond:%d", common.GetChainBaseInfo().SecondsPerSlot+2)
		strategys = append(strategys, s1)
	}

	return strategys

}
