package encryptutils

import (
	"testing"

	"github.com/nft-rainbow/rainbow-goutils/encrypt/aes"
	"gotest.tools/assert"
)

func TestBase64EncodeAndEncrypt(t *testing.T) {
	plaintext := "hello world"

	encrypted, err := Base64EncodeAndEncrypt(plaintext, aes.NewAesGcmEncryptor(), []byte("1234567812345678"))
	assert.NilError(t, err)

	decrypted, err := DecryptAndBase64Decode(encrypted, aes.NewAesGcmEncryptor(), []byte("1234567812345678"))
	assert.NilError(t, err)
	assert.Equal(t, plaintext, decrypted)
}
