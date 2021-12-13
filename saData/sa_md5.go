package saData

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)

//x = true 十六进制小写格式
func Md5(s string, x bool) string {
	h := md5.New()
	if _, err := io.WriteString(h, s); err == nil {
		if x {
			s = fmt.Sprintf("%x", h.Sum(nil))
		} else {
			s = string(h.Sum(nil))
		}
		return strings.TrimSpace(s)
	}
	return ""
}

func Sha256(s string, secret string, x bool) string {
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := io.WriteString(h, s); err == nil {
		if x {
			return fmt.Sprintf("%x", h.Sum(nil))
		} else {
			return string(h.Sum(nil))
		}
	}
	return ""
}
