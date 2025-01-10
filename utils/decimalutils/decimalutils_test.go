package decimalutils

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOnesToWei(t *testing.T) {
	assert.Equal(t, decimal.NewFromInt(100000).String(), OnesToWei(decimal.NewFromInt(1), 5).String())
	assert.Equal(t, decimal.NewFromInt(2000).String(), OnesToWei(decimal.NewFromInt(2), 3).String())
	assert.Equal(t, decimal.NewFromInt(212340).String(), OnesToWei(decimal.NewFromFloat32(2.1234), 5).String())
}
