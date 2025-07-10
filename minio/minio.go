package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/civet148/oss/common"
	"github.com/civet148/oss/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"io"
	"strings"
)

const (
	ossRegion   = "us-east-1"
	octetStream = "application/octet-stream"
)

type OssMinioClient struct {
	types.UnimplementedOssService
	opts   *types.OssOptions
	client *minio.Client
}

func NewOssClient(opfs ...types.Option) types.IOssClient {
	opts, err := common.ParseAndCheckOptions(types.OssType_Minio, opfs...)
	if err != nil {
		panic(err.Error())
	}
	var endpoint = opts.GetEndpoint()
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.GetAccessKey(), opts.GetSecretKey(), ""),
		Secure: strings.Contains(endpoint, "https"), // 如果是http，则改为 false
	})
	if err != nil {
		panic(err.Error())
	}
	return &OssMinioClient{
		opts:   opts,
		client: client,
	}
}

func (m *OssMinioClient) PutFile(ctx context.Context, objectName, filePath string) (path string, err error) {
	bucketName := m.opts.GetBucketName()

	if err = m.makeBucket(bucketName); err != nil {
		return "", errors.Wrap(err, "make bucket error")
	}

	_, err = m.client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: octetStream})
	if err != nil {
		return "", errors.Wrap(err, "put object error")
	}
	return m.makeObjectUrl(bucketName, objectName), nil
}

func (m *OssMinioClient) PutBuf(ctx context.Context, objectName string, buf *bytes.Buffer) (path string, err error) {
	bucketName := m.opts.GetBucketName()

	if err = m.makeBucket(bucketName); err != nil {
		return "", errors.Wrap(err, "make bucket error")
	}

	_, err = m.client.PutObject(ctx, bucketName, objectName, buf, int64(buf.Len()), minio.PutObjectOptions{ContentType: octetStream})
	if err != nil {
		return "", errors.Wrap(err, "put object error")
	}
	return m.makeObjectUrl(bucketName, objectName), nil
}

func (m *OssMinioClient) GetFile(ctx context.Context, writer io.Writer, objectName string) (err error) {
	bucketName := m.opts.GetBucketName()
	var obj *minio.Object
	obj, err = m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("get file %s error: %v", objectName, err)
	}
	_, err = io.Copy(writer, obj)
	if err != nil {
		return fmt.Errorf("copy file %s error: %v", objectName, err)
	}
	return nil
}

func (m *OssMinioClient) makeBucket(bucketName string) (err error) {
	err = m.client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: ossRegion})
	if err != nil {
		exists, errBucketExists := m.client.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			err = nil
		} else {
			return err
		}
	}
	return nil
}

func (m *OssMinioClient) makeObjectUrl(bucketName, objectName string) string {
	var fileUrl string
	if m.opts.GetDomain() == "" {
		fileUrl = fmt.Sprintf("%s/%s", bucketName, objectName)
	} else {
		fileUrl = fmt.Sprintf("%s/%s/%s", m.opts.GetDomain(), bucketName, objectName)
	}
	return fileUrl
}
