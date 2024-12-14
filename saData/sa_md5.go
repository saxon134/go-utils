package saData

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// Md5 32位 lowercase = true 小写格式，否则大写
func Md5(s string, lowercase bool) string {
	h := md5.New()
	if _, err := io.WriteString(h, s); err == nil {
		s = strings.TrimSpace(fmt.Sprintf("%x", h.Sum(nil)))
		if lowercase {
			s = strings.ToLower(s)
		} else {
			s = strings.ToUpper(s)
		}
		return s
	}
	return ""
}

// Sha256 lowercase = true 小写格式，否则大写
func Sha256(s string, secret string, lowercase bool) string {
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := io.WriteString(h, s); err == nil {
		s = strings.TrimSpace(fmt.Sprintf("%x", h.Sum(nil)))
		if lowercase {
			return strings.ToLower(s)
		} else {
			return strings.ToUpper(s)
		}
	}
	return ""
}

// Sha1
func Sha1(str string, secret string) string {
	hash := sha1.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}
