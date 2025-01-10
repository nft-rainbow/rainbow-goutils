package encryptutils

import (
	"encoding/base64"

	"github.com/nft-rainbow/rainbow-goutils/encrypt"
	"github.com/sirupsen/logrus"
)

func DecryptAndBase64Decode(content string, encryptor encrypt.Encryptor, password []byte) (string, error) {
	base64decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}

	decrypted, err := encrypt.DecryptBytes(encryptor, base64decoded, password)
	if err != nil {
		logrus.WithError(err).Error("decrypt failed")
		return "", err
	}
	return string(decrypted), nil
}

func Base64EncodeAndEncrypt(content string, encryptor encrypt.Encryptor, password []byte) (string, error) {
	encrypted, err := encrypt.EncryptBytes(encryptor, []byte(content), password)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}
