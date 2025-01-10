package fileutils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/h2non/bimg"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

// CompressImage compresses the image using bimg library
// if original file size is less than targetSizeMB, return the original file path
func CompressImage(filePath string, outputSuffix string, targetSizeMB int64, targetHeight int, alwaysCompress bool) (string, error) {

	originalFileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}

	// Read the file content
	buffer, err := bimg.Read(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Get image size
	size, err := bimg.NewImage(buffer).Size()
	if err != nil {
		return "", fmt.Errorf("failed to get image size: %v", err)
	}

	needResize := size.Height > targetHeight

	if !needResize && originalFileInfo.Size() <= int64(targetSizeMB*1024*1024) {
		return filePath, nil
	}

	// Generate new file path
	compressedPath := GetCompressedFilePath(filePath, outputSuffix)

	// if compressed file exists, return the original path
	if !alwaysCompress {
		if _, err := os.Stat(compressedPath); err == nil {
			return compressedPath, nil
		}
	}

	// Prepare compression options
	options := bimg.Options{
		Quality: 100,
		Type:    bimg.JPEG,
	}

	// Add resize if needed, maintaining aspect ratio
	if needResize {
		ratio := float64(size.Width) / float64(size.Height)
		options.Height = targetHeight
		options.Width = int(float64(targetHeight) * ratio)
	}

	startTime := time.Now()

	// 预估文件系统开销（稍微增加系数）
	const fileSystemOverhead = 1.2
	memoryLimit := int(float64(targetSizeMB*1024*1024) / fileSystemOverhead)

	for options.Quality > 20 {

		compressed, err := bimg.NewImage(buffer).Process(options)
		if err != nil {
			return "", fmt.Errorf("failed to compress image: %v", err)
		}

		// 先检查内存中的大小
		if len(compressed) <= memoryLimit {
			// 只在预计大小合适时写入磁盘
			if err := bimg.Write(compressedPath, compressed); err != nil {
				return "", fmt.Errorf("failed to write compressed file: %v", err)
			}

			// 最后验证一次实际文件大小
			fileInfo, err := os.Stat(compressedPath)
			if err != nil {
				os.Remove(compressedPath)
				options.Quality -= 10
				continue
			}

			if fileInfo.Size() <= targetSizeMB*1024*1024 {
				logrus.WithField("compressedPath", compressedPath).
					WithField("quality", options.Quality).
					WithField("resized", needResize).
					WithField("newSize", fileInfo.Size()).
					WithField("timeCost", time.Since(startTime)).
					Info("[FileUtils] compressed image")
				return compressedPath, nil
			}
			// 如果实际大小超出预期，删除文件并继续尝试
			os.Remove(compressedPath)
		}
		options.Quality -= 10
	}
	return "", fmt.Errorf("failed to compress image")
}

func GetCompressedFilePath(filePath string, outputSuffix string) string {
	return fmt.Sprintf("%s%s.%s", strings.TrimSuffix(filePath, filepath.Ext(filePath)), outputSuffix, "jpg")
}
