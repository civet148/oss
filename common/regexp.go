package common

import (
	"regexp"
)

// 编译正则表达式
const cosUrlPattern = `^https://([a-zA-Z0-9][a-zA-Z0-9_-]*-[0-9]+)\.cos\.([a-z0-9-]+)\.myqcloud\.com$`

// COS端点格式：https://<bucket>.cos.<region>.myqcloud.com
// bucket的命名规则为{name}-{appid} ，此处填写的存储桶名称必须为此格式(将<bucket>和<region>修改为真实的信息)
func CheckCosUrl(strUrl string) (bucket string, ok bool) {
	re := regexp.MustCompile(cosUrlPattern)

	// 匹配并提取bucket
	matches := re.FindStringSubmatch(strUrl)
	if len(matches) < 2 {
		return "", false
	}
	return matches[1], true
}
