package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type AesGcmEncryptor struct {
}

func NewAesGcmEncryptor() *AesGcmEncryptor {
	return &AesGcmEncryptor{}
}

func (a *AesGcmEncryptor) Encrypt(input io.Reader, output io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	nonce := make([]byte, 12) // 12 bytes nonce for GCM
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 将 nonce 写入输出
	if _, err := output.Write(nonce); err != nil {
		return err
	}

	buffer := make([]byte, 1024) // 1KB buffer
	for {
		n, err := input.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// 加密数据块
		ciphertext := aesgcm.Seal(nil, nonce, buffer[:n], nil)

		// 将密文写入输出
		if _, err := output.Write(ciphertext); err != nil {
			return err
		}
	}

	return nil
}

func (a *AesGcmEncryptor) Decrypt(input io.Reader, output io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := aesgcm.NonceSize()
	nonce := make([]byte, nonceSize)

	// 从输入中读取 nonce
	if _, err := io.ReadFull(input, nonce); err != nil {
		return err
	}

	buffer := make([]byte, 1024+aesgcm.Overhead()) // 1KB buffer + GCM overhead
	for {
		n, err := input.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// 解密数据块
		plaintext, err := aesgcm.Open(nil, nonce, buffer[:n], nil)
		if err != nil {
			return err // 如果密码错误或数据被篡改，解密会失败
		}

		// 将明文写入输出
		if _, err := output.Write(plaintext); err != nil {
			return err
		}
	}

	return nil
}
