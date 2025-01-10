package addrutils

import "github.com/ethereum/go-ethereum/common"

func IsValidChecksumAddress(address string) bool {
	return common.IsHexAddress(address) && address == common.HexToAddress(address).Hex()
}
