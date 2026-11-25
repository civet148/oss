package common

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/civet148/oss/types"
)

// HashString 使用指定算法计算字符串的签名(小写十六进制)
func HashString(algo types.CDNSignAlgorithm, str string) (string, error) {
	switch algo {
	case types.CDNSignAlgorithm_MD5:
		sum := md5.Sum([]byte(str))
		return hex.EncodeToString(sum[:]), nil
	case types.CDNSignAlgorithm_SHA256:
		sum := sha256.Sum256([]byte(str))
		return hex.EncodeToString(sum[:]), nil
	default:
		return "", fmt.Errorf("unsupported sign algorithm %v", algo)
	}
}

// RandomString 生成指定长度的随机字符串(由大小写字母和数字组成)
func RandomString(length int) (string, error) {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if length <= 0 {
		return "", nil
	}
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i := range buf {
		buf[i] = charset[int(buf[i])%len(charset)]
	}
	return string(buf), nil
}
