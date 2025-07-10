package oss

import (
	"fmt"
	"github.com/civet148/oss/aliyun"
	"github.com/civet148/oss/cos"
	"github.com/civet148/oss/minio"
	"github.com/civet148/oss/types"
	"sync"
)

var (
	mtx sync.Mutex
)

func init() {
	register(types.OssType_COS, cos.NewOssClient)
	register(types.OssType_Minio, minio.NewOssClient)
	register(types.OssType_Aliyun, aliyun.NewOssClient)
}

func NewOssClient(ossType types.OssType, opfs ...types.Option) types.IOssClient {
	creator := getInstance(ossType)
	return creator(opfs...)
}

type OssCreator func(opfs ...types.Option) types.IOssClient

var ossInstances = make(map[types.OssType]OssCreator)

func register(ossType types.OssType, creator OssCreator) {
	mtx.Lock()
	defer mtx.Unlock()
	if _, ok := ossInstances[ossType]; ok {
		msg := fmt.Errorf("oss type %v is already registered", ossType)
		panic(msg)
	}
	ossInstances[ossType] = creator
}

func getInstance(ossType types.OssType) OssCreator {
	mtx.Lock()
	defer mtx.Unlock()
	creator, ok := ossInstances[ossType]
	if ok {
		return creator
	}
	msg := fmt.Sprintf("oss type %v is not registered", ossType)
	panic(msg)
	return nil
}
