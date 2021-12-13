package saImg

import (
	"github.com/saxon134/go-utils/saData"
	"net/url"
	"strings"
)

var _imgUriRoot string
var _styleStrAry = []string{
	"?x-oss-process=style/default-style",
	"?x-oss-process=style/small-style",
	"?x-oss-process=style/cover-style",
	"?x-oss-process=style/imgtxt-style",
	"?x-oss-process=style/banner-style",
}

func Init(urlRoot string, styleAry []string) {
	_imgUriRoot = urlRoot
	if len(styleAry) > 0 {
		_styleStrAry = styleAry
	}
}

func AddDefaultUriRoot(s string) string {
	return AddUriRoot(s, NullImgStyle)
}

func AddUriRoot(s string, style ImgStyle) string {
	if s == "" {
		return s
	}

	if strings.HasPrefix(s, "http") == false {
		r := _imgUriRoot
		if r == "" {
			return s
		}

		s = r + s
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
	if s == "" {
		return ""
	}

	u, err := url.Parse(s)
	if err != nil {
		return s
	}

	root := u.Scheme + "://" + u.Host + "/"
	if r := _imgUriRoot; r != "" {
		if root == r {
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
