package ginutils

import (
	"fmt"
	"math"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/nft-rainbow/rainbow-goutils/utils/fileutils"
)

// SaveFile returns file name which is md5+filetype, such as fe33dd18d2800788c8d39a844f18860c.png
func SaveFile(c *gin.Context, formKey string, saveDirPath string, maxUploadSizeMB int64, allowedExts []string) (string, error) {
	_maxSizeInMB := int64(math.MaxInt64 / 1024 / 1024)
	if maxUploadSizeMB > _maxSizeInMB {
		maxUploadSizeMB = _maxSizeInMB
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Check file size
	if fileHeader.Size > maxUploadSizeMB*1024*1024 {
		return "", fmt.Errorf("file size %d exceeds the maximum allowed size of %d MB", fileHeader.Size, maxUploadSizeMB)
	}

	// check file extension
	if !fileutils.IsAllowedExt(fileHeader.Filename, allowedExts) {
		return "", fmt.Errorf("file extension is not allowed: %s", fileHeader.Filename)
	}

	md5Sum, err := fileutils.Md5(file)
	if err != nil {
		return "", err
	}
	file.Seek(0, 0)

	fileType, err := fileutils.GetFileType(fileHeader.Filename, file)
	if err != nil {
		return "", err
	}
	file.Seek(0, 0)

	// create save path
	if err := os.MkdirAll(saveDirPath, 0777); err != nil {
		return "", err
	}

	saveFilePath := path.Join(saveDirPath, md5Sum+"."+fileType.Extension)
	if err = fileutils.SaveFile(file, saveFilePath); err != nil {
		return "", err
	}

	return saveFilePath, nil
}
