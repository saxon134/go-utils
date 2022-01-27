package saData

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)

//lowercase = true 小写格式
func Md5(s string, lowercase bool) string {
	h := md5.New()
	if _, err := io.WriteString(h, s); err == nil {
		if lowercase {
			s = fmt.Sprintf("%x", h.Sum(nil))
		} else {
			s = string(h.Sum(nil))
		}
		return strings.TrimSpace(s)
	}
	return ""
}

//lowercase = true 小写格式
func Sha256(s string, secret string, lowercase bool) string {
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := io.WriteString(h, s); err == nil {
		if lowercase {
			return fmt.Sprintf("%x", h.Sum(nil))
		} else {
			return string(h.Sum(nil))
		}
	}
	return ""
}
