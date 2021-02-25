package saImg

import (
	"net/url"
	"strings"
)

var ImgUriRoot string
var StyleStrAry = []string{
	"?x-oss-process=style/default-style",
	"?x-oss-process=style/small-style",
	"?x-oss-process=style/cover-style",
	"?x-oss-process=style/imgtxt-style",
	"?x-oss-process=style/banner-style",
}

func AddDefaultUriRoot(s string) string {
	return AddUriRoot(s, NullImgStyle)
}

func AddUriRoot(s string, style ImgStyle) string {
	if s == "" {
		return s
	}

	if strings.HasPrefix(s, "http") == false {
		r := ImgUriRoot
		if r == "" {
			return s
		}

		s = r + s
	}

	for _, v := range StyleStrAry {
		s = strings.Replace(s, v, "", -1)
	}

	if style == NullImgStyle {
		return s
	}

	var index = int(style) - 1
	if index < len(StyleStrAry) {
		return s + StyleStrAry[index]
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
	if r := ImgUriRoot; r != "" {
		if root == r {
			s = strings.Replace(s, root, "", 1)
			for _, v := range StyleStrAry {
				s = strings.Replace(s, v, "", -1)
			}
		}
	}
	return s
}
