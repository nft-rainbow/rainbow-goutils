package ginutils

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

func RegisterValidation(tag string, fn validator.Func) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation(tag, fn)
	}
}

func RegisterValidations() {
	RegisterValidation("ethAddress", EthAddressValidator)
	RegisterValidation("uppercase", UpperCaseValidator)
}
