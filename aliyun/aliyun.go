package aliyun

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

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

// PresignGetUrl 生成阿里云OSS预签名下载URL
func (m *OssAliyunClient) PresignGetUrl(ctx context.Context, objectName string, expire time.Duration) (url string, err error) {
	var bucket *alioss.Bucket
	bucket, err = m.client.Bucket(m.opts.GetBucketName())
	if err != nil {
		return "", err
	}
	return bucket.SignURL(objectName, alioss.HTTPGet, int64(expire.Seconds()))
}

// GetCDNAuthUrl 生成阿里云CDN鉴权访问URL
// 支持的鉴权类型: TypeA、TypeB、TypeC；鉴权密钥通过 types.WithCDNAuthKey 配置
// 注意: 阿里云CDN鉴权的有效期由CDN控制台配置，expire参数不参与签名
func (m *OssAliyunClient) GetCDNAuthUrl(ctx context.Context, objectName string, expire time.Duration) (string, error) {
	_ = expire
	if m.opts.GetDomain() == "" {
		return "", fmt.Errorf("aliyun cdn domain is empty")
	}
	rand, err := common.RandomString(32)
	if err != nil {
		return "", fmt.Errorf("generate random string error: %v", err)
	}
	rawUrl := m.makeObjectUrl(m.opts.GetBucketName(), objectName)
	return m.signCDNUrl(rawUrl, time.Now(), rand, "0")
}

// signCDNUrl 按阿里云CDN鉴权规范对URL进行签名
// rand: 随机字符串(不能包含中划线) uid: 用户ID(固定为0)
func (m *OssAliyunClient) signCDNUrl(rawUrl string, now time.Time, rand, uid string) (string, error) {
	key := m.opts.GetCDNAuthKey()
	if key == "" {
		return "", fmt.Errorf("aliyun cdn auth key is empty")
	}
	u, err := common.ParseUrl(rawUrl)
	if err != nil {
		return "", fmt.Errorf("parse url %s error: %v", rawUrl, err)
	}
	path := u.EscapedPath()
	switch m.opts.GetCDNAuthType() {
	case types.CDNAuthType_A:
		// URL: http://DomainName/Filename?auth_key=<timestamp>-rand-uid-<md5hash>
		// md5hash = md5(URI-timestamp-rand-uid-PrivateKey)
		timestamp := strconv.FormatInt(now.Unix(), 10)
		signStr := fmt.Sprintf("%s-%s-%s-%s-%s", path, timestamp, rand, uid, key)
		hash, err := common.HashString(types.CDNSignAlgorithm_MD5, signStr)
		if err != nil {
			return "", err
		}
		query := u.Query()
		query.Set("auth_key", fmt.Sprintf("%s-%s-%s-%s", timestamp, rand, uid, hash))
		u.RawQuery = query.Encode()
	case types.CDNAuthType_B:
		// URL: http://DomainName/<timestamp>/<md5hash>/FileName
		// md5hash = md5(PrivateKey+timestamp+Path)，timestamp格式为UTC+8的YYYYMMDDHHMM
		timestamp := now.In(time.FixedZone("CST", 8*3600)).Format("200601021504")
		signStr := fmt.Sprintf("%s%s%s", key, timestamp, path)
		hash, err := common.HashString(types.CDNSignAlgorithm_MD5, signStr)
		if err != nil {
			return "", err
		}
		u.Path = fmt.Sprintf("/%s/%s%s", timestamp, hash, u.Path)
		u.RawPath = ""
	case types.CDNAuthType_C:
		// URL: http://DomainName/<md5hash>/<timestamp>/FileName
		// md5hash = md5(PrivateKey+URI+timestamp)，timestamp为十六进制Unix时间戳
		timestamp := strconv.FormatInt(now.Unix(), 16)
		signStr := fmt.Sprintf("%s%s%s", key, path, timestamp)
		hash, err := common.HashString(types.CDNSignAlgorithm_MD5, signStr)
		if err != nil {
			return "", err
		}
		u.Path = fmt.Sprintf("/%s/%s%s", hash, timestamp, u.Path)
		u.RawPath = ""
	default:
		return "", fmt.Errorf("aliyun cdn auth type %v not supported", m.opts.GetCDNAuthType())
	}
	return u.String(), nil
}
