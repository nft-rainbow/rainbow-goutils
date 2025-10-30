package ginutils

import (
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

type Req struct {
	Amount decimal.Decimal `binding:"decimalGt0"`
}

func TestDecimalGt0_WithGinBindingTag(t *testing.T) {
	RegisterValidations()

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		t.Fatal("binding.Validator.Engine() is not *validator.Validate")
	}

	if err := v.Struct(Req{Amount: decimal.NewFromInt(1)}); err != nil {
		t.Fatalf("1 should pass: %v", err)
	}
	if err := v.Struct(Req{Amount: decimal.Zero}); err == nil {
		t.Fatal("0 should not pass")
	}
	if err := v.Struct(Req{Amount: decimal.NewFromInt(-1)}); err == nil {
		t.Fatal("-1 should not pass")
	}
}
