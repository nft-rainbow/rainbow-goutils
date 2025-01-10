package aes

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	"gotest.tools/assert"
)

func TestEncryptBytes(t *testing.T) {
	input := []byte("1")
	output := bytes.NewBuffer(make([]byte, 0))

	encryptor := NewAesCbcEncryptor([]byte("abcdef1234567890"))
	err := encryptor.Encrypt(bytes.NewBuffer(input), output, []byte("1234567812345678"))
	assert.NilError(t, err)

	r, _ := io.ReadAll(output)
	fmt.Printf("result %x\n", r)
	assert.Equal(t, hex.EncodeToString(r), "61626364656631323334353637383930cfab6d1815ef7a19aeba3b700c9d8c99")
}

func TestDecryptBytes(t *testing.T) {
	input, err := hex.DecodeString("61626364656631323334353637383930cfab6d1815ef7a19aeba3b700c9d8c99")
	assert.NilError(t, err)
	fmt.Printf("input %x\n", input)

	inputBuf := bytes.NewBuffer(input)
	outputBuf := bytes.NewBuffer(make([]byte, 0))

	encryptor := NewAesCbcEncryptor([]byte("abcdef1234567890"))
	err = encryptor.Decrypt(inputBuf, outputBuf, []byte("1234567812345678"))
	assert.NilError(t, err)

	r, _ := io.ReadAll(outputBuf)
	fmt.Printf("decrypt result %s\n", r)
	assert.Equal(t, "1", string(r))
}
