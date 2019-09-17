package validator

import (
	"github.com/bytedance/go-tagexpr/validator"
)

func Validate(input interface{}) (bool, error) {
	var vd = validator.New("vd")

	result := vd.Validate(input)
	if result != nil {
		return false, result
	}

	return true, nil
}

//ValidateConds is validator of input conds
func ValidateConds(endpoint string, params *map[string]interface{}) error {
	return nil
}
