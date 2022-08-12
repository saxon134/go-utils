package saImg

import (
	"github.com/saxon134/go-utils/saData"
	"net/url"
	"strings"
)

var _mainDomain string
var _domainAry []string //包含 _mainDomain
var _styleStrAry = []string{
	"?x-oss-process=style/default-style",
	"?x-oss-process=style/small-style",
	"?x-oss-process=style/cover-style",
	"?x-oss-process=style/imgtxt-style",
	"?x-oss-process=style/banner-style",
}

// Init
// @Description: 初始化
// @param domains 多个域名，分号隔开。第一个为默认输出的域名，其他域名在删除时会自动删除
// @param styleAry 样式组，删除时会自动移除
func Init(domains string, styleAry []string) {
	_domainAry = strings.Split(domains, ";")
	if len(_domainAry) > 0 {
		for i, v := range _domainAry {
			_domainAry[i] = strings.TrimSuffix(v, "/")
		}
		_mainDomain = _domainAry[0]
	}

	if len(styleAry) > 0 {
		_styleStrAry = styleAry
	}
}

func AddDefaultUriRoot(s string) string {
	return AddUriRoot(s, NullImgStyle)
}

func AddUriRoot(s string, style ImgStyle) string {
	if s == "" || _mainDomain == "" {
		return s
	}

	if strings.HasPrefix(s, "http") == false {
		s = saData.ConnPath(_mainDomain, s)
	}

	for _, v := range _styleStrAry {
		s = strings.Replace(s, v, "", -1)
	}

	if style == NullImgStyle {
		return s
	}

	var index = int(style) - 1
	if index < len(_styleStrAry) {
		return s + _styleStrAry[index]
	}
	return s
}

func DeleteUriRoot(s string) string {
	if s == "" || len(_domainAry) == 0 {
		return s
	}

	u, err := url.Parse(s)
	if err != nil {
		return s
	}

	root := u.Scheme + "://" + u.Host
	root = strings.TrimSuffix(root, "/")
	for _, r := range _domainAry {
		if r == root {
			s = strings.Replace(s, root, "", 1)
			for _, v := range _styleStrAry {
				s = strings.Replace(s, v, "", -1)
			}
		}
	}
	return s
}

func ConnectUri(host string, url string) string {
	if strings.HasPrefix(host, "http") &&
		url != "" &&
		strings.HasPrefix(url, "http") == false {
		url = saData.ConnPath(host, url)
	}
	return url
}
