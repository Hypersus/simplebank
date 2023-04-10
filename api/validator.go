package api

import (
	"github.com/Hypersus/simplebank/util"
	"github.com/go-playground/validator/v10"
)

func currencyValidator(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return util.IsValidCurrency(currency)
}
