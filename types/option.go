package types

import "time"

type OssOptions struct {
	endpoint   string         //服务器地址(示例：127.0.0.1:9000)
	accessKey  string         //访问Key
	secretKey  string         //密钥
	bucketName string         //桶名称
	domain     string         //文件访问域名
	timeout    *time.Duration //超时时间
	presign    *bool          //是否预签名(CDN回源)
}

func (o OssOptions) GetEndpoint() string {
	return o.endpoint
}

func (o OssOptions) GetAccessKey() string {
	return o.accessKey
}

func (o OssOptions) GetSecretKey() string {
	return o.secretKey
}

func (o OssOptions) GetBucketName() string {
	return o.bucketName
}

func (o OssOptions) GetDomain() string {
	return o.domain
}

func (o OssOptions) GetTimeout() time.Duration {
	if o.timeout == nil {
		return 300 * time.Second
	}
	return *o.timeout
}

func (o OssOptions) GetPresign() bool {
	if o.presign == nil {
		return false
	}
	return *o.presign
}

type Option func(v *OssOptions)

func ParseOptions(opts *OssOptions, opfs ...Option) *OssOptions {
	for _, op := range opfs {
		op(opts)
	}
	return opts
}

func WithEndpoint(endpoint string) Option {
	return func(v *OssOptions) {
		v.endpoint = endpoint
	}
}

func WithAccessKey(ak string) Option {
	return func(v *OssOptions) {
		v.accessKey = ak
	}
}

func WithSecretKey(sk string) Option {
	return func(v *OssOptions) {
		v.secretKey = sk
	}
}

func WithBucketName(bucket string) Option {
	return func(v *OssOptions) {
		v.bucketName = bucket
	}
}

func WithDomain(domain string) Option {
	return func(v *OssOptions) {
		v.domain = domain
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(v *OssOptions) {
		v.timeout = &timeout
	}
}

func WithPresign(presign bool) Option {
	return func(v *OssOptions) {
		v.presign = &presign
	}
}
