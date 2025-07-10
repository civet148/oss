package aliyun

import (
	"bytes"
	"context"
	"fmt"
	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/civet148/oss/common"
	"github.com/civet148/oss/types"
	"io"
)

type OssAliyunClient struct {
	types.UnimplementedOssService
	opts   *types.OssOptions
	client *alioss.Client
}

func init() {
}

func NewOssClient(opfs ...types.Option) types.IOssClient {
	opts, err := common.ParseAndCheckOptions(types.OssType_Aliyun, opfs...)
	if err != nil {
		panic(err.Error())
	}
	client, err := alioss.New(opts.GetEndpoint(), opts.GetAccessKey(), opts.GetSecretKey())
	if err != nil {
		panic(err.Error())
	}
	return &OssAliyunClient{
		opts:   opts,
		client: client,
	}
}

func (m *OssAliyunClient) PutFile(ctx context.Context, objectName, filePath string) (path string, err error) {
	bucketName := m.opts.GetBucketName()
	var bucket *alioss.Bucket
	bucket, err = m.client.Bucket(bucketName)
	if err != nil {
		return path, err
	}
	err = bucket.PutObjectFromFile(objectName, filePath)
	if err != nil {
		return path, err
	}
	return m.makeObjectUrl(bucketName, objectName), nil
}

func (m *OssAliyunClient) PutBuf(ctx context.Context, objectName string, buf *bytes.Buffer) (path string, err error) {
	bucketName := m.opts.GetBucketName()
	var bucket *alioss.Bucket
	bucket, err = m.client.Bucket(bucketName)
	if err != nil {
		return path, err
	}
	err = bucket.PutObject(objectName, buf)
	if err != nil {
		return path, err
	}
	return m.makeObjectUrl(bucketName, objectName), nil
}

func (m *OssAliyunClient) GetFile(ctx context.Context, writer io.Writer, objectName string) (err error) {
	var bucket *alioss.Bucket
	bucketName := m.opts.GetBucketName()
	bucket, err = m.client.Bucket(bucketName)
	if err != nil {
		return err
	}
	var reader io.ReadCloser
	reader, err = bucket.GetObject(objectName)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}
	return nil
}

func (m *OssAliyunClient) makeObjectUrl(bucketName, objectName string) string {
	var fileUrl string
	if m.opts.GetDomain() == "" {
		fileUrl = fmt.Sprintf("%s/%s", bucketName, objectName)
	} else {
		fileUrl = fmt.Sprintf("%s/%s/%s", m.opts.GetDomain(), bucketName, objectName)
	}
	return fileUrl
}
