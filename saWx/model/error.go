package model

import (
	"encoding/json"
	"fmt"
)

/******** 微信返回的通用错误 ********/
type WxError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// DecodeWxError 将返回值按照WxError解析
func DecodeWxError(response []byte) (err error) {
	var commError WxError
	err = json.Unmarshal(response, &commError)
	if err != nil {
		return
	}
	if commError.ErrCode != 0 {
		return fmt.Errorf("wx api Error , errcode=%d , errmsg=%s", commError.ErrCode, commError.ErrMsg)
	}
	return nil
}
