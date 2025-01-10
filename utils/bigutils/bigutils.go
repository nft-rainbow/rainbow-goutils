package bigutils

import (
	"math/big"
)

func MaxU256MinusOne() *big.Int {
	base := big.NewInt(2)
	exponent := big.NewInt(256)
	result := new(big.Int).Exp(base, exponent, nil)
	result.Sub(result, big.NewInt(1))
	return result
}
