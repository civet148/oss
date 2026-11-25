package cos

import (
	"context"
	"testing"
	"time"

	"github.com/civet148/oss/types"
)

func TestPresignGetUrlLocal(t *testing.T) {
	cli := NewOssClient(
		types.WithEndpoint("https://mybucket-123456789.cos.ap-guangzhou.myqcloud.com"),
		types.WithAccessKey("AKIDxxxxxxxx"),
		types.WithSecretKey("secret"),
	)
	u, err := cli.PresignGetUrl(context.Background(), "tmp/a.jpg", time.Hour)
	if err != nil {
		t.Fatalf("presign error: %v", err)
	}
	t.Logf("url=%v", u)
}
