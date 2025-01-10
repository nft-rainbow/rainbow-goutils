package aes

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("examplekey123456") // 16 bytes key for AES-128

	// 创建一个随机的明文数据
	plaintext := make([]byte, 1024) // 1KB of random data

	if _, err := io.ReadFull(rand.Reader, plaintext); err != nil {
		t.Fatalf("Failed to generate random plaintext: %v", err)
	}

	// 创建缓冲区用于加密输出
	ciphertextBuffer := new(bytes.Buffer)

	gcm := NewAesGcmEncryptor()

	// 加密
	if err := gcm.Encrypt(bytes.NewReader(plaintext), ciphertextBuffer, key); err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// 创建缓冲区用于解密输出
	decryptedBuffer := new(bytes.Buffer)

	// 解密
	if err := gcm.Decrypt(bytes.NewReader(ciphertextBuffer.Bytes()), decryptedBuffer, key); err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// 验证解密后的数据是否与原始明文相同
	if !bytes.Equal(plaintext, decryptedBuffer.Bytes()) {
		t.Fatal("Decrypted data does not match original plaintext")
	}
}

func TestDecryptCheckError(t *testing.T) {
	gcm := NewAesGcmEncryptor()

	plaintext := []byte("hello world")
	ciphertextBuffer := new(bytes.Buffer)

	err := gcm.Encrypt(bytes.NewReader(plaintext), ciphertextBuffer, []byte("1234567890123456"))
	assert.NoError(t, err)

	decryptedBuffer := new(bytes.Buffer)
	err = gcm.Decrypt(ciphertextBuffer, decryptedBuffer, []byte("1234567890123455"))
	assert.Error(t, err)

	fmt.Println(err)
}
