package cos

import (
	"bytes"
	"context"
	"fmt"
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
