package decimalutils

import (
	"math"

	"github.com/shopspring/decimal"
)

func EthToWei(valInEth decimal.Decimal) decimal.Decimal {
	return valInEth.Mul(decimal.NewFromInt(1e18))
}

func WeiToOnes(valInWei decimal.Decimal, _decimal uint) decimal.Decimal {
	return valInWei.Div(decimal.NewFromInt(int64(math.Pow(10, float64(_decimal)))))
}

// 转换以个位为单位到以最小单位为单位
func OnesToWei(val decimal.Decimal, _decimal uint) decimal.Decimal {
	return val.Mul(decimal.NewFromFloat(math.Pow(10, float64(_decimal))))
}
