package minio

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"

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

// PresignGetUrl 生成Minio预签名下载URL(S3 V4签名)
func (m *OssMinioClient) PresignGetUrl(ctx context.Context, objectName string, expire time.Duration) (string, error) {
	u, err := m.client.PresignedGetObject(ctx, m.opts.GetBucketName(), objectName, expire, nil)
	if err != nil {
		return "", errors.Wrap(err, "presign minio object error")
	}
	return u.String(), nil
}

// GetCDNAuthUrl 生成Minio CDN访问URL
// Minio本身没有独立的CDN鉴权机制，采用S3 V4预签名实现CDN回源鉴权:
//   - 配置了CDN域名且开启预签名(types.WithPresign(true))时，返回替换为CDN域名后的预签名URL
//   - 配置了CDN域名但未开启预签名时，返回不带鉴权的CDN访问URL
//   - 未配置CDN域名时，直接返回预签名URL
func (m *OssMinioClient) GetCDNAuthUrl(ctx context.Context, objectName string, expire time.Duration) (string, error) {
	domain := m.opts.GetDomain()
	if domain == "" {
		return m.PresignGetUrl(ctx, objectName, expire)
	}
	if !m.opts.GetPresign() {
		return m.makeObjectUrl(m.opts.GetBucketName(), objectName), nil
	}
	u, err := m.client.PresignedGetObject(ctx, m.opts.GetBucketName(), objectName, expire, nil)
	if err != nil {
		return "", errors.Wrap(err, "presign minio object error")
	}
	return m.replaceUrlDomain(u, domain)
}

// replaceUrlDomain 将预签名URL的协议、域名替换为CDN域名，保留路径和签名参数
func (m *OssMinioClient) replaceUrlDomain(u *url.URL, domain string) (string, error) {
	base, err := common.ParseUrl(domain)
	if err != nil {
		return "", fmt.Errorf("parse cdn domain %s error: %v", domain, err)
	}
	if base.Host == "" {
		return "", fmt.Errorf("invalid cdn domain %s", domain)
	}
	cdnUrl := *u
	cdnUrl.Scheme = base.Scheme
	cdnUrl.Host = base.Host
	cdnUrl.Path = fmt.Sprintf("%s%s", strings.TrimRight(base.Path, "/"), u.Path)
	cdnUrl.RawPath = ""
	return cdnUrl.String(), nil
}
