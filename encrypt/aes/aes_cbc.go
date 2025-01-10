package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"io"

	"github.com/sirupsen/logrus"
)

type AesCbcEncryptor struct {
	iv []byte
}

func NewAesCbcEncryptor(iv []byte) *AesCbcEncryptor {
	return &AesCbcEncryptor{iv: iv}
}

func (a *AesCbcEncryptor) Encrypt(input io.Reader, output io.Writer, key []byte) error {
	logrus.Info("encrypt by method aes-cbc")
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	output.Write(a.iv)

	for {
		buf := make([]byte, 4096)
		n, err := input.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// fmt.Println("read", n)

		paded := pad(buf[:n])
		cipher.NewCBCEncrypter(block, a.iv).CryptBlocks(paded, paded)
		_, err = output.Write(paded)
		if err != nil {
			return err
		}
		// fmt.Println("write", n)
	}

	return nil
}

func (a *AesCbcEncryptor) Decrypt(input io.Reader, output io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = input.Read(iv)
	if err != nil {
		return err
	}

	end := false
	for {
		if end {
			break
		}

		buf := make([]byte, 4096)
		n, err := input.Read(buf)
		if err == io.EOF || n < 4096 {
			end = true
		} else if err != nil {
			return err
		}

		logrus.WithField("size", n).WithField("end", end).Info("read data")

		chunk := buf[:n]
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(chunk, chunk)

		if end {
			chunk = trimTailZeros(chunk)
		}

		n, err = output.Write(chunk)
		if err != nil {
			return err
		}

		// fmt.Println("write", n)
		logrus.WithField("size", n).Info("write to output")
	}
	return nil
}

func pad(input []byte) []byte {
	padLen := aes.BlockSize - len(input)%aes.BlockSize
	input = append(input, make([]byte, padLen)...)
	return input
}

func trimTailZeros(input []byte) []byte {
	lastNonZeroIndex := len(input) - 1
	for lastNonZeroIndex >= 0 && input[lastNonZeroIndex] == 0 {
		lastNonZeroIndex--
	}

	result := input[:lastNonZeroIndex+1]
	return result
}
