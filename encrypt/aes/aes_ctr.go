package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type AesCtrEncryptor struct {
	iv []byte
}

func NewAesCtrEncryptor(iv string) *AesCtrEncryptor {
	return &AesCtrEncryptor{iv: []byte(iv)}
}

func (a *AesCtrEncryptor) Encrypt(input io.Reader, output io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	stream := cipher.NewCTR(block, a.iv)

	writer := &cipher.StreamWriter{S: stream, W: output}
	_, err = io.Copy(writer, input)
	return err
}

func (a *AesCtrEncryptor) Decrypt(input io.Reader, output io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	stream := cipher.NewCTR(block, a.iv)

	reader := &cipher.StreamReader{S: stream, R: input}
	_, err = io.Copy(output, reader)
	return err
}
