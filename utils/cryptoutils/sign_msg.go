package cryptoutils

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

func RecoverMessage(message string, signature string) (common.Address, error) {
	// 1. Hash the message with prefix
	hash := hashForSign([]byte(message))
	return RecoverHash(hash, signature)
}

func RecoverHash(hash []byte, sig string) (common.Address, error) {
	if sig[:2] == "0x" {
		sig = sig[2:]
	}
	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return common.Address{}, errors.WithMessage(err, "Invalid signature")
	}

	// 3. Adjust the v value in ECDSA signature
	if sigBytes[64] != 27 && sigBytes[64] != 28 {
		return common.Address{}, errors.New("Invalid Ethereum signature (v is not 27 or 28)")
	}
	sigBytes[64] -= 27 // Correct v value to 0 or 1

	// 4. Recover public key
	pubKey, err := crypto.SigToPub(hash, sigBytes)
	if err != nil {
		return common.Address{}, errors.WithMessage(err, "Failed to recover public key")
	}

	// 5. Convert public key to Ethereum address
	recoveredAddress := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddress, nil
}

// hashForSign 将消息加上以太坊前缀并进行 keccak256 哈希
func hashForSign(data []byte) []byte {
	prefix := []byte("\x19Ethereum Signed Message:\n" + fmt.Sprintf("%d", len(data)))
	return crypto.Keccak256(append(prefix, data...))
}
