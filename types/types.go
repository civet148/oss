package types

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"
)

var errNotSupported = errors.New("not supported")

type OssType int

const (
	OssType_Minio  OssType = 1 //minio
	OssType_COS    OssType = 2 //Tencent COS
	OssType_Aliyun OssType = 3 //Aliyun OSS
)

func (ot OssType) String() string {
	switch ot {
	case OssType_Minio:
		return "OSS_Minio"
	case OssType_COS:
		return "OSS_COS"
	case OssType_Aliyun:
		return "OSS_Aliyun"
	}
	return "<OSS_Unknown>"
}

// CDNAuthType CDN URL鉴权类型
// 不同云厂商支持的鉴权类型不同:
//   - 腾讯云COS: TypeA、TypeB、TypeC、TypeD
//   - 阿里云CDN: TypeA、TypeB、TypeC
type CDNAuthType int

const (
	CDNAuthType_A CDNAuthType = 1 //鉴权方式A
	CDNAuthType_B CDNAuthType = 2 //鉴权方式B
	CDNAuthType_C CDNAuthType = 3 //鉴权方式C
	CDNAuthType_D CDNAuthType = 4 //鉴权方式D
)

func (t CDNAuthType) String() string {
	switch t {
	case CDNAuthType_A:
		return "TypeA"
	case CDNAuthType_B:
		return "TypeB"
	case CDNAuthType_C:
		return "TypeC"
	case CDNAuthType_D:
		return "TypeD"
	}
	return "<CDNAuthType_Unknown>"
}

// CDNSignAlgorithm CDN URL鉴权签名算法
type CDNSignAlgorithm int

const (
	CDNSignAlgorithm_MD5    CDNSignAlgorithm = 1 //MD5
	CDNSignAlgorithm_SHA256 CDNSignAlgorithm = 2 //SHA256
)

func (a CDNSignAlgorithm) String() string {
	switch a {
	case CDNSignAlgorithm_MD5:
		return "md5"
	case CDNSignAlgorithm_SHA256:
		return "sha256"
	}
	return "<CDNSignAlgorithm_Unknown>"
}

// CDNTimestampFormat CDN鉴权时间戳格式
type CDNTimestampFormat int

const (
	CDNTimestampFormat_Decimal CDNTimestampFormat = 1 //十进制Unix时间戳
	CDNTimestampFormat_Hex     CDNTimestampFormat = 2 //十六进制Unix时间戳
)

func (f CDNTimestampFormat) String() string {
	switch f {
	case CDNTimestampFormat_Decimal:
		return "decimal"
	case CDNTimestampFormat_Hex:
		return "hex"
	}
	return "<CDNTimestampFormat_Unknown>"
}

type IOssClient interface {
	PutFile(ctx context.Context, objectName, filePath string) (path string, err error)
	PutBuf(ctx context.Context, objectName string, buf *bytes.Buffer) (path string, err error)
	GetFile(ctx context.Context, writer io.Writer, objectName string) (err error)
	// PresignGetUrl 生成预签名下载URL，expire为URL有效期
	PresignGetUrl(ctx context.Context, objectName string, expire time.Duration) (url string, err error)
	// GetCDNAuthUrl 生成带鉴权的CDN访问URL
	// 注意: 阿里云/腾讯云CDN鉴权的有效期由CDN控制台配置，expire仅用于Minio预签名
	GetCDNAuthUrl(ctx context.Context, objectName string, expire time.Duration) (url string, err error)
}

type UnimplementedOssService struct {
}

func (_ UnimplementedOssService) PutFile(ctx context.Context, objectName, filePath string) (path string, err error) {
	return path, errNotSupported
}

func (_ UnimplementedOssService) PutBuf(ctx context.Context, objectName string, buf *bytes.Buffer) (path string, err error) {
	return path, errNotSupported
}

func (_ UnimplementedOssService) GetFile(ctx context.Context, writer io.Writer, objectName string) (err error) {
	return errNotSupported
}

func (_ UnimplementedOssService) PresignGetUrl(ctx context.Context, objectName string, expire time.Duration) (url string, err error) {
	return url, errNotSupported
}

func (_ UnimplementedOssService) GetCDNAuthUrl(ctx context.Context, objectName string, expire time.Duration) (url string, err error) {
	return url, errNotSupported
}
