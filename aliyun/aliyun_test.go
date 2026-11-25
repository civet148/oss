package aliyun

import (
	"strings"
	"testing"
	"time"

	"github.com/civet148/oss/types"
)

func newTestClient(authType types.CDNAuthType) *OssAliyunClient {
	opts := types.ParseOptions(&types.OssOptions{},
		types.WithDomain("http://domain.example.com"),
		types.WithBucketName("mybucket"),
		types.WithCDNAuthKey("aliyuncdnexp1234"),
		types.WithCDNAuthType(authType),
	)
	return &OssAliyunClient{opts: opts}
}

func TestSignCDNUrlTypeA(t *testing.T) {
	c := newTestClient(types.CDNAuthType_A)
	now := time.Unix(1444435200, 0)
	got, err := c.signCDNUrl("http://domain.example.com/video/standard/test.mp4", now, "0", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://domain.example.com/video/standard/test.mp4?auth_key=1444435200-0-0-23bf85053008f5c0e791667a313e28ce"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeB(t *testing.T) {
	c := newTestClient(types.CDNAuthType_B)
	now := time.Date(2015, 8, 15, 8, 0, 0, 0, time.FixedZone("CST", 8*3600))
	got, err := c.signCDNUrl("http://domain.example.com/4/44/44c0909bcfc20a01afaf256ca99a8b8b.mp3", now, "0", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://domain.example.com/201508150800/9044548ef1527deadafa49a890a377f0/4/44/44c0909bcfc20a01afaf256ca99a8b8b.mp3"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeC(t *testing.T) {
	c := newTestClient(types.CDNAuthType_C)
	now := time.Unix(0x55CE8100, 0)
	got, err := c.signCDNUrl("http://domain.example.com/test.flv", now, "0", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://domain.example.com/c6880e19a04f71f9a585d0394cf0794e/55ce8100/test.flv"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlUnsupportedType(t *testing.T) {
	c := newTestClient(types.CDNAuthType_D)
	_, err := c.signCDNUrl("http://domain.example.com/test.flv", time.Now(), "0", "0")
	if err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("expect not supported error, got %v", err)
	}
}

func TestSignCDNUrlWithoutKey(t *testing.T) {
	opts := types.ParseOptions(&types.OssOptions{}, types.WithDomain("http://domain.example.com"))
	c := &OssAliyunClient{opts: opts}
	_, err := c.signCDNUrl("http://domain.example.com/test.flv", time.Now(), "0", "0")
	if err == nil || !strings.Contains(err.Error(), "auth key is empty") {
		t.Fatalf("expect empty key error, got %v", err)
	}
}
