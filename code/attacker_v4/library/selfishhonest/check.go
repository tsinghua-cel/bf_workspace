package selfishhonest

import (
	"github.com/tsinghua-cel/attacker-service/types"
	"strconv"
)

func CheckDuties(param types.LibraryParams, duties []types.ProposerDuty) ([]interface{}, bool) {
	result := make([]interface{}, 0)
	for i := 0; i < len(duties)-2; {
		a := duties[i]
		b := duties[i+1]
		c := duties[i+2]
		aValidatorIndex, _ := strconv.ParseInt(a.ValidatorIndex, 10, 64)
		bValidatorIndex, _ := strconv.ParseInt(b.ValidatorIndex, 10, 64)
		cValidatorIndex, _ := strconv.ParseInt(c.ValidatorIndex, 10, 64)
		if param.IsHackValidator(int(aValidatorIndex)) || !param.IsHackValidator(int(bValidatorIndex)) || param.IsHackValidator(int(cValidatorIndex)) {
			result = append(result, []types.ProposerDuty{a, b, c})
			i += 3
		} else {
			i++
		}
	}
	if len(result) > 0 {
		return result, true
	}
	return nil, false
}
