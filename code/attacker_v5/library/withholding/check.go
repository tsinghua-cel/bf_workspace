package withholding

import (
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/types"
	"strconv"
)

func CheckDuties(param types.LibraryParams, duties []types.ProposerDuty) ([]interface{}, bool) {
	result := make([]interface{}, 0)

	tmpsub := make([]types.ProposerDuty, 0)
	for _, duty := range duties {
		valIdx, _ := strconv.Atoi(duty.ValidatorIndex)

		if param.IsHackValidator(valIdx) {
			tmpsub = append(tmpsub, duty)
		} else {
			if len(tmpsub) >= 5 {
				// if the latest slot is epochEnd, add it to response.
				latest, _ := strconv.Atoi(tmpsub[len(tmpsub)-1].Slot)
				if int64(latest) == (common.EpochEnd(common.SlotToEpoch(int64(latest)))) {
					result = append(result, tmpsub)
				}

			}
			tmpsub = make([]types.ProposerDuty, 0)
		}
	}
	if len(tmpsub) > 5 {
		result = append(result, tmpsub)
	}

	if len(result) > 0 {
		return result, true
	}

	return nil, false
}
