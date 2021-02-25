package saData

import (
	"encoding/base64"
)

func EncodeBase64(s string) string {
	if s != "" {
		return base64.StdEncoding.EncodeToString([]byte(s))
	}
	return ""
}

func DecodeBase64(s string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(s)
	if err == nil {
		return string(decodeBytes)
	}
	return ""
}
