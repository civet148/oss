package aliyun

import (
	"context"
	"testing"
	"time"

	"github.com/civet148/oss/types"
)

func TestPresignGetUrlLocal(t *testing.T) {
	cli := NewOssClient(
		types.WithEndpoint("https://oss-cn-hangzhou.aliyuncs.com"),
		types.WithBucketName("mybucket"),
		types.WithAccessKey("ak"),
		types.WithSecretKey("sk"),
	)
	u, err := cli.PresignGetUrl(context.Background(), "tmp/a.jpg", time.Hour)
	if err != nil {
		t.Fatalf("presign error: %v", err)
	}
	t.Logf("url=%v", u)
}
