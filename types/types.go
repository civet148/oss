package types

import (
	"bytes"
	"context"
	"errors"
	"io"
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

type IOssClient interface {
	PutFile(ctx context.Context, objectName, filePath string) (path string, err error)
	PutBuf(ctx context.Context, objectName string, buf *bytes.Buffer) (path string, err error)
	GetFile(ctx context.Context, writer io.Writer, objectName string) (err error)
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
