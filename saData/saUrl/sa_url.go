package saUrl

import (
	"net/url"
	"strings"
)

// QueryEncode 将字符串进行query编码
func QueryEncode(s string) string {
	if s != "" {
		return url.QueryEscape(s)
	}
	return ""
}

// QueryDecode 对字符串进行query解码
func QueryDecode(s string) string {
	res, err := url.QueryUnescape(s)
	if err != nil {
		return ""
	}
	return res
}

// QueryFromMap query & map 互转
func QueryFromMap(m map[string]string) string {
	if m != nil {
		urlV := url.Values{}
		for k, v := range m {
			if k != "" {
				urlV.Add(k, v)
			}
		}
		return urlV.Encode()
	}
	return ""
}
func QueryToMap(urlStr string) map[string]string {
	values, _ := url.ParseQuery(urlStr)
	m := map[string]string{}
	for k, v := range values {
		if k != "" {
			m[k] = v[0]
		}
	}

	return m
}

// ConnPath 返回结果是： /r/path
func ConnPath(r string, path string) (full string) {
	r = strings.TrimSuffix(r, "/")
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	if len(r) > 0 {
		full = r + "/" + path
	} else {
		full = path
	}

	return full
}
