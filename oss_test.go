package oss

import (
	"bytes"
	"context"
	"fmt"
	"github.com/civet148/oss/common"
	"github.com/civet148/oss/types"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const (
	uploadFile = "/tmp/Goland-protobuf-format.jpg"
	objectName = "tmp/Goland-protobuf-format.jpg"
)

var minioOpfs = []types.Option{
	types.WithBucketName("mybucket"),               //桶名
	types.WithEndpoint("127.0.0.1:9000"),           //OSS对象存储服务器地址
	types.WithDomain("https://files.mydomain.com"), //文件访问URL域名前缀
	types.WithAccessKey("b7UMyD............U4p"),   //Minio access key
	types.WithSecretKey("K3AJSV............vhjD"),  //Minio secret key
}

var cosOpfs = []types.Option{
	types.WithTimeout(100 * time.Second),                                           // 超时配置
	types.WithEndpoint("https://mybucket-123456789.cos.ap-guangzhou.myqcloud.com"), //OSS对象存储服务器地址(格式：https://<bucket>.cos.<region>.myqcloud.com 其中<bucket>是{name}-{appid}组合)
	types.WithAccessKey("AKIDx......................5MRKF"),                        //COS secret id
	types.WithSecretKey("83oOl......................RuGRO"),                        //COS secret key
	types.WithDomain("https://files.mydomain.com"),                                 //COS 文件访问域名(如果用CDN则填CDN域名，否则默认采用COS访问域名)
}

func TestGetInstance(t *testing.T) {
	tests := []struct {
		name    string
		ossType types.OssType
	}{
		{"Minio", types.OssType_Minio},
		{"COS", types.OssType_COS},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creator := getInstance(tt.ossType)
			if creator == nil {
				t.Errorf("get oss type [%v] instance error", tt.ossType)
			} else {
				fmt.Printf("get oss type [%v] instance success\n", tt.ossType)
			}
		})
	}
}

func TestPutMinioFile(t *testing.T) {
	var minioCli = NewOssClient(types.OssType_Minio, minioOpfs...)
	fileUrl, err := minioCli.PutFile(context.Background(), objectName, uploadFile)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Printf("minio file url: %s \n", fileUrl)
}

func TestPutCosFile(t *testing.T) {
	var cli = NewOssClient(types.OssType_COS, cosOpfs...)

	data, err := ioutil.ReadFile(uploadFile)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	//fileUrl, err := cli.PutFile(context.Background(), "tmp", "order.log", uploadFile, "application/octet-stream")
	buf := bytes.NewBuffer(data)
	fmt.Printf("file size: %d hex: %X\n", len(data), data[:32])
	fileUrl, err := cli.PutBuf(context.Background(), objectName, buf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Printf("cos file url: %s \n", fileUrl)
}

func TestCosUrl(t *testing.T) {
	url := "https://mybucket-12345678.cos.ap-guangzhou.myqcloud.com"
	bucket, ok := common.CheckCosUrl(url)
	if !ok {
		fmt.Printf("COS url %s not correct\n", url)
		return
	}
	fmt.Printf("bucket %s\n", bucket)
}

func TestGetMinioFile(t *testing.T) {
	var cli = NewOssClient(types.OssType_Minio, minioOpfs...)
	file, err := os.OpenFile("/tmp/download.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Errorf("create file error: %v", err)
		return
	}
	err = cli.GetFile(context.Background(), file, objectName)
	if err != nil {
		t.Errorf("download file error: %v", err)
		return
	}
}

func TestGetCosFile(t *testing.T) {
	var cli = NewOssClient(types.OssType_COS, cosOpfs...)
	file, err := os.OpenFile("/tmp/golang.jpg", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Errorf("create file error: %v", err)
		return
	}
	err = cli.GetFile(context.Background(), file, objectName)
	if err != nil {
		t.Errorf("download file error: %v", err)
		return
	}
}
