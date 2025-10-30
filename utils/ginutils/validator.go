package ginutils

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func EthAddressValidator(fl validator.FieldLevel) bool {
	addrStr, ok := fl.Field().Interface().(string)
	if ok {
		return common.IsHexAddress(addrStr)
	}
	return false
}

func UpperCaseValidator(fl validator.FieldLevel) bool {
	text, ok := fl.Field().Interface().(string)
	if ok {
		return strings.ToUpper(text) == text
	}
	return false
}

func DecimalGt0Validator(fl validator.FieldLevel) bool {
	num, ok := fl.Field().Interface().(decimal.Decimal)
	if ok {
		return num.GreaterThan(decimal.Zero)
	}
	return false
}

func RegisterValidation(tag string, fn validator.Func) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation(tag, fn)
	}
}

func RegisterValidations() {
	RegisterValidation("ethAddress", EthAddressValidator)
	RegisterValidation("uppercase", UpperCaseValidator)
	RegisterValidation("decimalGt0", DecimalGt0Validator)
}
