package ext_staircase

import (
	"fmt"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/types"
	"strconv"
)

func getSlotStrategy(slot string, cas int, isLatestHackSlot bool) types.SlotStrategy {
	strategy := types.SlotStrategy{
		Slot:    slot,
		Level:   0,
		Actions: make(map[string]string),
	}
	base := common.GetChainBaseInfo()
	secondsPerSlot := base.SecondsPerSlot
	slotsPerEpoch := base.SlotsPerEpoch
	switch cas {
	case 0:
		strategy.Actions["BlockBeforeSign"] = "return"
		strategy.Actions["AttestBeforeSign"] = fmt.Sprintf("return")

	case 1:
		if isLatestHackSlot {
			islot, _ := strconv.Atoi(slot)
			stageI := (slotsPerEpoch - islot%slotsPerEpoch) * secondsPerSlot
			stageII := 12 * secondsPerSlot

			strategy.Actions["AttestBeforeSign"] = fmt.Sprintf("return")

			strategy.Actions["BlockBeforeSign"] = "packPooledAttest"
			strategy.Actions["BlockDelayForReceiveBlock"] = fmt.Sprintf("%s:%d", "delayWithSecond", stageI)
			strategy.Actions["BlockBeforeBroadCast"] = fmt.Sprintf("%s:%d", "delayWithSecond", stageII)
		} else {
			strategy.Actions["BlockBeforeSign"] = "return"
			strategy.Actions["AttestAfterSign"] = fmt.Sprintf("addAttestToPool")
			strategy.Actions["AttestBeforePropose"] = fmt.Sprintf("return")
		}
	}
	return strategy

}

func GenSlotStrategy(duties []types.ProposerDuty, hackduties []types.ProposerDuty, cas int) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	for i := 0; i < len(duties); i++ {
		duty := duties[i]
		s := getSlotStrategy(duty.Slot, cas, duty.Slot == hackduties[len(hackduties)-1].Slot)
		strategys = append(strategys, s)
	}
	return strategys
}
