package validation

import "github.com/go-playground/validator/v10"

func minIDsRule(fl validator.FieldLevel) bool {
	ids, ok := fl.Field().Interface().([]uint)
	if ok {
		return len(ids) > 0
	}
	return false
}
