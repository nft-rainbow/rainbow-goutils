package fileutils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/pkg/errors"
)

func IsPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Md5File(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return Md5(file)
}

func Md5(src io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetFileType(fileName string, fileReader io.Reader) (*types.Type, error) {
	var fileContent bytes.Buffer
	if _, err := io.CopyN(&fileContent, fileReader, 100); err != nil {
		return nil, err
	}

	fileType, err := filetype.Match(fileContent.Bytes())
	if err != nil {
		return nil, errors.WithMessage(err, "failed get file type from content")
	}

	if fileType != filetype.Unknown {
		return &fileType, nil
	}

	tmp := strings.Split(fileName, ".")
	return &types.Type{
		Extension: tmp[len(tmp)-1],
	}, nil
}

// Note: will return if file exists
func SaveFile(src io.Reader, saveFilePath string) error {
	isExists, err := IsPathExist(saveFilePath)
	if err != nil {
		return err
	}

	if isExists {
		return nil
	}

	dirPath := filepath.Dir(saveFilePath)
	if err := os.MkdirAll(dirPath, 0777); err != nil {
		return err
	}

	f, err := os.Create(saveFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, src); err != nil {
		return err
	}
	return nil
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func ReadFileInBase64(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return toBase64(bytes), nil
}

// check if file extension is allowed
func IsAllowedExt(fileName string, allowedExts []string) bool {
	if len(allowedExts) == 0 {
		return true
	}
	ext := filepath.Ext(fileName)
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
