package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
)

func CalculateMD5(file multipart.File) (string, error) {
	defer file.Seek(0, io.SeekStart) // 重置文件指针，确保后续可以再次读取

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
