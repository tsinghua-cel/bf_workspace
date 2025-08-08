package unrealized

import (
	"fmt"
	"github.com/tsinghua-cel/attacker-service/types"
	"strconv"
)

func GenSlotStrategy(duties []interface{}) []types.SlotStrategy {
	strategies := make([]types.SlotStrategy, 0)
	duty := duties[0].([]types.ProposerDuty)

	frontCount := 8
	currentslot, _ := strconv.Atoi(duty[0].Slot)
	start := currentslot - frontCount
	for i := 1; i < frontCount; i++ {
		slot := start + i
		slotStrategy := types.SlotStrategy{
			Slot:    fmt.Sprintf("%d", slot),
			Level:   1,
			Actions: make(map[string]string),
		}
		// not generate new block.
		slotStrategy.Actions["BlockBeforeSign"] = "return"
		// don't broadcast attestations.
		slotStrategy.Actions["AttestBeforePropose"] = "return"
		// add attestations to pool.
		slotStrategy.Actions["AttestAfterSign"] = "addAttestToPool"
		strategies = append(strategies, slotStrategy)
	}
	{
		slotStrategy := types.SlotStrategy{
			Slot:    fmt.Sprintf("%d", currentslot),
			Level:   1,
			Actions: make(map[string]string),
		}
		// pack all attestation.
		slotStrategy.Actions["BlockBeforeSign"] = "packPooledAttest"
		// modify parent root to old slot.
		slotStrategy.Actions["BlockGetNewParentRoot"] = fmt.Sprintf("modifyParentRoot:%d", start)
		strategies = append(strategies, slotStrategy)
	}

	return strategies

}
