package common

import (
	"net/url"
	"strings"
)

// ParseUrl 解析URL，未指定协议时默认使用https
func ParseUrl(rawUrl string) (*url.URL, error) {
	if !strings.Contains(rawUrl, "://") {
		rawUrl = "https://" + strings.TrimLeft(rawUrl, "/")
	}
	return url.Parse(rawUrl)
}
