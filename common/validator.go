package common

import (
	"fmt"
	"github.com/civet148/oss/types"
)

func ParseAndCheckOptions(ossType types.OssType, opfs ...types.Option) (*types.OssOptions, error) {
	var err error
	var opts = &types.OssOptions{}
	opts = types.ParseOptions(opts, opfs...)
	switch ossType {
	case types.OssType_Minio:
		opts, err = checkMinioOptions(opts)
		if err != nil {
			return nil, err
		}
	case types.OssType_COS:
		opts, err = checkCosOptions(opts)
		if err != nil {
			return nil, err
		}
	case types.OssType_Aliyun:
		opts, err = checkAliyunOptions(opts)
		if err != nil {
			return nil, err
		}
	}
	return opts, nil
}

func checkMinioOptions(opts *types.OssOptions) (*types.OssOptions, error) {

	if opts.GetEndpoint() == "" {
		return opts, fmt.Errorf("minio endpoint is empty")
	}
	if opts.GetAccessKey() == "" {
		return opts, fmt.Errorf("minio access key is empty")
	}
	if opts.GetSecretKey() == "" {
		return opts, fmt.Errorf("minio secret key is empty")
	}
	if opts.GetBucketName() == "" {
		return opts, fmt.Errorf("minio bucket name is empty")
	}
	return opts, nil
}

func checkCosOptions(opts *types.OssOptions) (*types.OssOptions, error) {
	if opts.GetEndpoint() == "" {
		return opts, fmt.Errorf("cos endpoint is empty")
	}
	if opts.GetAccessKey() == "" {
		return opts, fmt.Errorf("cos secret id is empty")
	}
	if opts.GetSecretKey() == "" {
		return opts, fmt.Errorf("cos secret key is empty")
	}
	bucket, ok := CheckCosUrl(opts.GetEndpoint())
	if ok {
		//通过COS服务地址获取桶名和文件访问域名前缀
		if opts.GetBucketName() == "" {
			opts = types.ParseOptions(opts, types.WithBucketName(bucket))
		}
		if opts.GetDomain() == "" {
			opts = types.ParseOptions(opts, types.WithDomain(opts.GetEndpoint()))
		}
	} else {
		if opts.GetBucketName() == "" {
			return opts, fmt.Errorf("cos bucket name is empty")
		}
	}
	return opts, nil
}

func checkAliyunOptions(opts *types.OssOptions) (*types.OssOptions, error) {
	if opts.GetEndpoint() == "" {
		return opts, fmt.Errorf("aliyun oss endpoint is empty")
	}
	if opts.GetAccessKey() == "" {
		return opts, fmt.Errorf("aliyun oss access key is empty")
	}
	if opts.GetSecretKey() == "" {
		return opts, fmt.Errorf("aliyun oss secret key is empty")
	}
	if opts.GetBucketName() == "" {
		return opts, fmt.Errorf("aliyun oss bucket name is empty")
	}
	return opts, nil
}
