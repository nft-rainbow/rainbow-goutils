package encrypt

import (
	"bytes"
	"path"
	"path/filepath"

	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Encryptor interface {
	Encrypt(input io.Reader, output io.Writer, key []byte) error
	Decrypt(input io.Reader, output io.Writer, key []byte) error
}

func EncryptBytes(e Encryptor, input, key []byte) ([]byte, error) {
	inputBuf := bytes.NewBuffer(input)
	outputBuf := bytes.NewBuffer(make([]byte, 0))

	if err := e.Encrypt(inputBuf, outputBuf, key); err != nil {
		return nil, err
	}
	return io.ReadAll(outputBuf)
}

func DecryptBytes(e Encryptor, input, key []byte) ([]byte, error) {
	inputBuf := bytes.NewBuffer(input)
	outputBuf := bytes.NewBuffer(make([]byte, 0))

	if err := e.Decrypt(inputBuf, outputBuf, key); err != nil {
		return nil, err
	}
	return io.ReadAll(outputBuf)
}

func EncryptFile(e Encryptor, source, outputDirPath string, key []byte) (string, error) {
	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		return "", errors.WithMessage(err, "Failed to create directory")
	}

	sf, err := os.OpenFile(source, os.O_RDONLY, 0666)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to open source file")
	}
	defer sf.Close()

	fileName := filepath.Base(source)

	outputhPath := path.Join(outputDirPath, fileName+".encrypt")

	of, err := os.OpenFile(outputhPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to create output file")
	}
	defer of.Close()

	logrus.WithField("source", source).WithField("output", outputhPath).Info("encrypt file")

	return outputhPath, e.Encrypt(sf, of, key)
}

func DecryptFile(e Encryptor, source, outputDirPath string, key []byte) (string, error) {
	// fmt.Printf("decrypt file source %s, out %s\n", source, outputDirPath)
	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		return "", errors.WithMessage(err, "Failed to create directory")
	}

	sf, err := os.OpenFile(source, os.O_RDONLY, 0666)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to open source file")
	}
	defer sf.Close()

	fileName := filepath.Base(source)

	outputhPath := path.Join(outputDirPath, fileName+".decrypt")

	of, err := os.OpenFile(outputhPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to create output file")
	}
	defer of.Close()

	return outputhPath, e.Decrypt(sf, of, key)
}
