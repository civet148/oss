package cos

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/civet148/oss/common"
	"github.com/civet148/oss/types"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
)

type OssCosClient struct {
	types.UnimplementedOssService
	opts   *types.OssOptions
	client *cos.Client
}

func init() {
}

func NewOssClient(opfs ...types.Option) types.IOssClient {

	opts, err := common.ParseAndCheckOptions(types.OssType_COS, opfs...)
	if err != nil {
		panic(err.Error())
	}
	var endpoint = opts.GetEndpoint()

	// COS端点格式：https://<bucket>.cos.<region>.myqcloud.com
	// bucket的命名规则为{name}-{appid} ，此处填写的存储桶名称必须为此格式(将<bucket>和<region>修改为真实的信息)
	u, _ := url.Parse(endpoint)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		//设置超时时间
		Timeout: opts.GetTimeout(),
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  opts.GetAccessKey(),
			SecretKey: opts.GetSecretKey(),
		},
	})

	return &OssCosClient{
		opts:   opts,
		client: client,
	}
}

func (m *OssCosClient) PutFile(ctx context.Context, objectName, filePath string) (path string, err error) {
	resp, err := m.client.Object.PutFromFile(ctx, objectName, filePath, nil)
	if err != nil {
		return path, err
	}
	err = m.checkResponse(resp)
	if err != nil {
		return path, err
	}
	return m.makeObjectUrl(objectName), nil
}

func (m *OssCosClient) PutBuf(ctx context.Context, objectName string, buf *bytes.Buffer) (path string, err error) {
	resp, err := m.client.Object.Put(ctx, objectName, buf, nil)
	if err != nil {
		return path, err
	}
	err = m.checkResponse(resp)
	if err != nil {
		return path, err
	}
	return m.makeObjectUrl(objectName), nil
}

func (m *OssCosClient) GetFile(ctx context.Context, writer io.Writer, objectName string) (err error) {
	var resp *cos.Response
	resp, err = m.client.Object.Get(ctx, objectName, nil)
	if err != nil {
		return errors.Wrap(err, "get object error")
	}
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return errors.Wrap(err, "copy object error")
	}
	return nil
}

func (m *OssCosClient) checkResponse(resp *cos.Response) (err error) {
	if resp.StatusCode != http.StatusOK {
		var msg []byte
		msg, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read http response body error: %s", err)
		}
		err = fmt.Errorf("upload buffer to COS return status code: %v msg: %s", resp.StatusCode, string(msg))
		return err
	}
	return nil
}

func (m *OssCosClient) makeObjectUrl(objectName string) string {
	var fileUrl string
	if m.opts.GetDomain() == "" {
		fileUrl = objectName
	} else {
		fileUrl = fmt.Sprintf("%s/%s", m.opts.GetDomain(), objectName)
	}
	return fileUrl
}

// PresignGetUrl 生成腾讯云COS预签名下载URL
func (m *OssCosClient) PresignGetUrl(ctx context.Context, objectName string, expire time.Duration) (url string, err error) {
	u, err := m.client.Object.GetPresignedURL(ctx, http.MethodGet, objectName,
		m.opts.GetAccessKey(), m.opts.GetSecretKey(), expire, nil)
	if err != nil {
		return "", errors.Wrap(err, "presign cos object error")
	}
	return u.String(), nil
}

// GetCDNAuthUrl 生成腾讯云CDN鉴权访问URL
// 支持的鉴权类型: TypeA、TypeB、TypeC、TypeD；签名算法: md5、sha256(通过 types.WithCDNSignAlgorithm 配置)
// 注意: 腾讯云CDN鉴权的有效期由CDN控制台配置，expire参数不参与签名
func (m *OssCosClient) GetCDNAuthUrl(ctx context.Context, objectName string, expire time.Duration) (string, error) {
	_ = expire
	if m.opts.GetDomain() == "" {
		return "", fmt.Errorf("cos cdn domain is empty")
	}
	rand, err := common.RandomString(16)
	if err != nil {
		return "", fmt.Errorf("generate random string error: %v", err)
	}
	return m.signCDNUrl(m.makeObjectUrl(objectName), time.Now(), rand, "0")
}

// signCDNUrl 按腾讯云CDN鉴权规范对URL进行签名
// rand: 随机字符串 uid: 用户ID(固定为0)
func (m *OssCosClient) signCDNUrl(rawUrl string, now time.Time, rand, uid string) (string, error) {
	key := m.opts.GetCDNAuthKey()
	if key == "" {
		return "", fmt.Errorf("cos cdn auth key is empty")
	}
	u, err := common.ParseUrl(rawUrl)
	if err != nil {
		return "", fmt.Errorf("parse url %s error: %v", rawUrl, err)
	}
	path := u.EscapedPath()
	algo := m.opts.GetCDNSignAlgorithm()
	query := u.Query()
	switch m.opts.GetCDNAuthType() {
	case types.CDNAuthType_A:
		// URL: http://DomainName/Filename?sign=timestamp-rand-uid-<hash>
		// hash = HASH(文件路径-timestamp-rand-uid-自定义密钥)
		timestamp := m.formatTimestamp(now, types.CDNTimestampFormat_Decimal)
		signStr := fmt.Sprintf("%s-%s-%s-%s-%s", path, timestamp, rand, uid, key)
		hash, err := common.HashString(algo, signStr)
		if err != nil {
			return "", err
		}
		query.Set("sign", fmt.Sprintf("%s-%s-%s-%s", timestamp, rand, uid, hash))
	case types.CDNAuthType_B:
		// URL: http://DomainName/timestamp/<hash>/FileName
		// hash = HASH(自定义密钥+timestamp+文件路径)，timestamp格式为UTC+8的YYYYMMDDHHMM
		timestamp := now.In(time.FixedZone("CST", 8*3600)).Format("200601021504")
		signStr := fmt.Sprintf("%s%s%s", key, timestamp, path)
		hash, err := common.HashString(algo, signStr)
		if err != nil {
			return "", err
		}
		u.Path = fmt.Sprintf("/%s/%s%s", timestamp, hash, u.Path)
		u.RawPath = ""
	case types.CDNAuthType_C:
		// URL: http://DomainName/<hash>/timestamp/FileName
		// hash = HASH(自定义密钥+文件路径+timestamp)，timestamp默认为十六进制Unix时间戳
		timestamp := m.formatTimestamp(now, types.CDNTimestampFormat_Hex)
		signStr := fmt.Sprintf("%s%s%s", key, path, timestamp)
		hash, err := common.HashString(algo, signStr)
		if err != nil {
			return "", err
		}
		u.Path = fmt.Sprintf("/%s/%s%s", hash, timestamp, u.Path)
		u.RawPath = ""
	case types.CDNAuthType_D:
		// URL: http://DomainName/Filename?sign=<hash>&t=timestamp
		// hash = HASH(自定义密钥+文件路径+timestamp)
		timestamp := m.formatTimestamp(now, types.CDNTimestampFormat_Decimal)
		signStr := fmt.Sprintf("%s%s%s", key, path, timestamp)
		hash, err := common.HashString(algo, signStr)
		if err != nil {
			return "", err
		}
		query.Set("sign", hash)
		query.Set("t", timestamp)
	default:
		return "", fmt.Errorf("cos cdn auth type %v not supported", m.opts.GetCDNAuthType())
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}

// formatTimestamp 按配置格式化Unix时间戳；未显式配置时使用各鉴权模式的默认格式
func (m *OssCosClient) formatTimestamp(now time.Time, defaultFormat types.CDNTimestampFormat) string {
	format := defaultFormat
	if m.opts.IsCDNTimestampFormatSet() {
		format = m.opts.GetCDNTimestampFormat()
	}
	if format == types.CDNTimestampFormat_Hex {
		return strconv.FormatInt(now.Unix(), 16)
	}
	return strconv.FormatInt(now.Unix(), 10)
}
