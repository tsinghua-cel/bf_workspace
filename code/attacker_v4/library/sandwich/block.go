package sandwich

import (
	"fmt"
	"github.com/tsinghua-cel/attacker-service/types"
)

func GenSlotStrategy(duties []interface{}) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	for i := 0; i < len(duties); i++ {
		duty := duties[i].([]types.ProposerDuty)
		if len(duty) != 3 {
			continue
		}
		a := duty[0]
		b := duty[1]
		c := duty[2]
		s1 := types.SlotStrategy{
			Slot:    a.Slot,
			Level:   1,
			Actions: make(map[string]string),
		}
		s1.Actions["AttestBeforeSign"] = fmt.Sprintf("modifyAttestHead:%s", a.Slot)
		strategys = append(strategys, s1)

		s2 := types.SlotStrategy{
			Slot:    b.Slot,
			Level:   1,
			Actions: make(map[string]string),
		}
		s2.Actions["AttestBeforeSign"] = fmt.Sprintf("modifyAttestHead:%s", a.Slot)
		strategys = append(strategys, s2)

		s3 := types.SlotStrategy{
			Slot:    fmt.Sprintf("%s", c.Slot),
			Level:   1,
			Actions: make(map[string]string),
		}
		s3.Actions["BlockGetNewParentRoot"] = fmt.Sprintf("modifyParentRoot:%s", a.Slot)
		strategys = append(strategys, s3)
	}

	return strategys

}
