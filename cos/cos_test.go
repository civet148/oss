package cos

import (
	"testing"
	"time"

	"github.com/civet148/oss/types"
)

func newTestClient(authType types.CDNAuthType, algo types.CDNSignAlgorithm, opfs ...types.Option) *OssCosClient {
	base := []types.Option{
		types.WithCDNAuthKey("dimtm5evg50ijsx2hvuwyfoiu65"),
		types.WithCDNAuthType(authType),
		types.WithCDNSignAlgorithm(algo),
	}
	opts := types.ParseOptions(&types.OssOptions{}, append(base, opfs...)...)
	return &OssCosClient{opts: opts}
}

func TestSignCDNUrlTypeAMD5(t *testing.T) {
	c := newTestClient(types.CDNAuthType_A, types.CDNSignAlgorithm_MD5)
	now := time.Unix(1656300600, 0)
	got, err := c.signCDNUrl("http://www.qcdntest.cn/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://www.qcdntest.cn/test.jpg?sign=1656300600-FnyigRo7yny-0-3276ffbdf30dc974ed955d7db761653a"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeAHexTimestamp(t *testing.T) {
	c := newTestClient(types.CDNAuthType_A, types.CDNSignAlgorithm_MD5, types.WithCDNTimestampFormat(types.CDNTimestampFormat_Hex))
	now := time.Unix(1656300600, 0)
	got, err := c.signCDNUrl("http://www.qcdntest.cn/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://www.qcdntest.cn/test.jpg?sign=62b92438-FnyigRo7yny-0-348e16216c54f5469f88d531665bfb75"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeDMD5(t *testing.T) {
	c := newTestClient(types.CDNAuthType_D, types.CDNSignAlgorithm_MD5)
	now := time.Unix(1656300600, 0)
	got, err := c.signCDNUrl("http://www.qcdntest.cn/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://www.qcdntest.cn/test.jpg?sign=708b9bf131071b2dd6d68729abe6dc3f&t=1656300600"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeASHA256(t *testing.T) {
	c := newTestClient(types.CDNAuthType_A, types.CDNSignAlgorithm_SHA256)
	now := time.Unix(1656300600, 0)
	got, err := c.signCDNUrl("http://www.qcdntest.cn/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://www.qcdntest.cn/test.jpg?sign=1656300600-FnyigRo7yny-0-0ead4878339e462062dfcb107af9e5b32909908b6ad2e9710404b7e116bd591f"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeDSHA256(t *testing.T) {
	c := newTestClient(types.CDNAuthType_D, types.CDNSignAlgorithm_SHA256)
	now := time.Unix(1656300600, 0)
	got, err := c.signCDNUrl("http://www.qcdntest.cn/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://www.qcdntest.cn/test.jpg?sign=77e7f28fe761ee7461b08968b729b72d59c7a85c425c8f42bbc96ec74d445a9d&t=1656300600"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeBMD5(t *testing.T) {
	c := newTestClient(types.CDNAuthType_B, types.CDNSignAlgorithm_MD5)
	now := time.Date(2020, 2, 27, 16, 10, 32, 0, time.FixedZone("CST", 8*3600))
	got, err := c.signCDNUrl("http://cloud.tencent.com/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://cloud.tencent.com/202002271610/2e03a07cfa55a47768226d3e5ea82a8d/test.jpg"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeBSHA256(t *testing.T) {
	c := newTestClient(types.CDNAuthType_B, types.CDNSignAlgorithm_SHA256)
	now := time.Date(2020, 2, 27, 16, 10, 32, 0, time.FixedZone("CST", 8*3600))
	got, err := c.signCDNUrl("http://cloud.tencent.com/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://cloud.tencent.com/202002271610/4444660cb3b623eb050fce4affe925448fd4c044831e7049e60b7b03cc8ee193/test.jpg"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeCMD5(t *testing.T) {
	c := newTestClient(types.CDNAuthType_C, types.CDNSignAlgorithm_MD5)
	now := time.Unix(0x5e577978, 0)
	got, err := c.signCDNUrl("http://cloud.tencent.com/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://cloud.tencent.com/7913fc0c5c9e92dd3633b7895152bbb2/5e577978/test.jpg"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeCSHA256(t *testing.T) {
	c := newTestClient(types.CDNAuthType_C, types.CDNSignAlgorithm_SHA256)
	now := time.Unix(0x5e577978, 0)
	got, err := c.signCDNUrl("http://cloud.tencent.com/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://cloud.tencent.com/2154423d5b271ee49130af4a4542132147c647d01a60422dd16ecaa26163cc5a/5e577978/test.jpg"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeCDecimalTimestamp(t *testing.T) {
	c := newTestClient(types.CDNAuthType_C, types.CDNSignAlgorithm_MD5, types.WithCDNTimestampFormat(types.CDNTimestampFormat_Decimal))
	now := time.Unix(0x5e577978, 0)
	got, err := c.signCDNUrl("http://cloud.tencent.com/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://cloud.tencent.com/900a5049aa8ac1ab144527d9c2be4cea/1582791032/test.jpg"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeDHexTimestamp(t *testing.T) {
	c := newTestClient(types.CDNAuthType_D, types.CDNSignAlgorithm_MD5, types.WithCDNTimestampFormat(types.CDNTimestampFormat_Hex))
	now := time.Unix(1656300600, 0)
	got, err := c.signCDNUrl("http://www.qcdntest.cn/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://www.qcdntest.cn/test.jpg?sign=588084276faa20472955092a72493f7e&t=62b92438"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlTypeBIgnoresTimestampFormat(t *testing.T) {
	c := newTestClient(types.CDNAuthType_B, types.CDNSignAlgorithm_MD5, types.WithCDNTimestampFormat(types.CDNTimestampFormat_Hex))
	now := time.Date(2020, 2, 27, 16, 10, 32, 0, time.FixedZone("CST", 8*3600))
	got, err := c.signCDNUrl("http://cloud.tencent.com/test.jpg", now, "FnyigRo7yny", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "http://cloud.tencent.com/202002271610/2e03a07cfa55a47768226d3e5ea82a8d/test.jpg"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestSignCDNUrlWithPathPrefix(t *testing.T) {
	opts := types.ParseOptions(&types.OssOptions{},
		types.WithDomain("https://cdn.example.com/prefix"),
		types.WithCDNAuthKey("secret"),
		types.WithCDNAuthType(types.CDNAuthType_D),
	)
	c := &OssCosClient{opts: opts}
	got, err := c.signCDNUrl(c.makeObjectUrl("tmp/a.jpg"), time.Unix(1656300600, 0), "rand", "0")
	if err != nil {
		t.Fatalf("sign url error: %v", err)
	}
	want := "https://cdn.example.com/prefix/tmp/a.jpg?sign=1dd23c5d0d0999d1ddeb31a11854e268&t=1656300600"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}
