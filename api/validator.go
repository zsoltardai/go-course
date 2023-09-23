package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/zsoltardai/simple_bank/util"
)

var validCurrency validator.Func = func(findLevel validator.FieldLevel) bool {
	if currency, ok := findLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}

	return false
}
