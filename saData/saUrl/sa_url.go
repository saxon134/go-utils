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

// AppendQuery uri后面拼接query参数
func AppendQuery(urlStr string, query map[string]string) string {
	if urlStr == "" {
		return QueryFromMap(query)
	}

	var ary = strings.Split(urlStr, "#")
	if strings.Contains(ary[len(ary)-1], "?") {
		urlStr += "&" + QueryFromMap(query)
	} else {
		urlStr += "?" + QueryFromMap(query)
	}
	return urlStr
}

// ConnPath 返回结果是： /r/path
func ConnPath(r string, paths ...string) (full string) {
	full = strings.TrimSuffix(r, "/")
	for _, path := range paths {
		path = strings.TrimPrefix(path, "/")
		path = strings.TrimSuffix(path, "/")
		if path != "" && path != "/" {
			full += "/" + path
		}
	}

	if strings.HasPrefix(full, "/") == false {
		full = "/" + full
	}
	return full
}

// ConnectUri host如果不包含http开头，则直接返回
// 返回格式：http://xxxx/xxx
func ConnectUri(host string, paths ...string) string {
	if strings.HasPrefix(host, "http") {
		return strings.TrimPrefix(ConnPath(host, paths...), "/")
	}
	return host
}
