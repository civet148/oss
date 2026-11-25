package types

import "time"

type OssOptions struct {
	endpoint    string              //服务器地址(示例：127.0.0.1:9000)
	accessKey   string              //访问Key
	secretKey   string              //密钥
	bucketName  string              //桶名称
	domain      string              //文件访问域名
	timeout     *time.Duration      //超时时间
	presign     *bool               //是否预签名(Minio CDN回源)
	cdnAuthKey  string              //CDN鉴权密钥
	cdnAuthType CDNAuthType         //CDN鉴权类型
	cdnSignAlgo CDNSignAlgorithm    //CDN鉴权签名算法
	cdnTsFormat *CDNTimestampFormat //CDN鉴权时间戳格式(十进制/十六进制)
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

func (o OssOptions) GetCDNAuthKey() string {
	return o.cdnAuthKey
}

func (o OssOptions) GetCDNAuthType() CDNAuthType {
	if o.cdnAuthType == 0 {
		return CDNAuthType_A
	}
	return o.cdnAuthType
}

func (o OssOptions) GetCDNSignAlgorithm() CDNSignAlgorithm {
	if o.cdnSignAlgo == 0 {
		return CDNSignAlgorithm_MD5
	}
	return o.cdnSignAlgo
}

// IsCDNTimestampFormatSet 是否显式设置了CDN鉴权时间戳格式
func (o OssOptions) IsCDNTimestampFormatSet() bool {
	return o.cdnTsFormat != nil
}

func (o OssOptions) GetCDNTimestampFormat() CDNTimestampFormat {
	if o.cdnTsFormat == nil {
		return CDNTimestampFormat_Decimal
	}
	return *o.cdnTsFormat
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

// WithPresign 设置Minio CDN回源是否使用S3预签名URL
func WithPresign(presign bool) Option {
	return func(v *OssOptions) {
		v.presign = &presign
	}
}

// WithCDNAuthKey 设置CDN URL鉴权密钥(与云厂商访问密钥无关)
func WithCDNAuthKey(key string) Option {
	return func(v *OssOptions) {
		v.cdnAuthKey = key
	}
}

// WithCDNAuthType 设置CDN URL鉴权类型(默认TypeA)
func WithCDNAuthType(authType CDNAuthType) Option {
	return func(v *OssOptions) {
		v.cdnAuthType = authType
	}
}

// WithCDNSignAlgorithm 设置CDN URL鉴权签名算法(默认MD5，目前仅腾讯云COS支持SHA256)
func WithCDNSignAlgorithm(algo CDNSignAlgorithm) Option {
	return func(v *OssOptions) {
		v.cdnSignAlgo = algo
	}
}

// WithCDNTimestampFormat 设置CDN URL鉴权Unix时间戳格式(十进制/十六进制)
// 默认十进制；TypeC按腾讯云规范默认使用十六进制，显式设置后以设置为准；TypeB固定为YYYYMMDDHHMM
func WithCDNTimestampFormat(format CDNTimestampFormat) Option {
	return func(v *OssOptions) {
		v.cdnTsFormat = &format
	}
}
