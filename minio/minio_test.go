package minio

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/civet148/oss/types"
)

func TestReplaceUrlDomain(t *testing.T) {
	c := &OssMinioClient{}
	signed, err := url.Parse("http://127.0.0.1:9000/mybucket/tmp/a.jpg?X-Amz-Signature=abc")
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.replaceUrlDomain(signed, "https://files.example.com/prefix")
	if err != nil {
		t.Fatalf("replace domain error: %v", err)
	}
	want := "https://files.example.com/prefix/mybucket/tmp/a.jpg?X-Amz-Signature=abc"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestGetCDNAuthUrlWithoutPresign(t *testing.T) {
	opts := types.ParseOptions(&types.OssOptions{},
		types.WithBucketName("mybucket"),
		types.WithDomain("https://files.example.com"),
	)
	c := &OssMinioClient{opts: opts}
	got, err := c.GetCDNAuthUrl(context.Background(), "tmp/a.jpg", time.Hour)
	if err != nil {
		t.Fatalf("get cdn url error: %v", err)
	}
	want := "https://files.example.com/mybucket/tmp/a.jpg"
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}
