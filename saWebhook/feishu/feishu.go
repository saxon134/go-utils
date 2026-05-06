package feishu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saHttp"
	"time"
)

type FeiShu struct {
	webhookUrl string
	secret     string
}

func New(webhookUrl string, secret string) *FeiShu {
	return &FeiShu{webhookUrl: webhookUrl, secret: secret}
}

func (m *FeiShu) SendTxt(txt string) error {
	var args = map[string]interface{}{
		"msg_type": "text",
		"content":  map[string]string{"text": txt},
	}
	if m.secret != "" {
		var timestamp = time.Now().Unix()
		var sign = sign(timestamp, m.secret)
		args["timestamp"] = timestamp
		args["sign"] = sign
	}

	var response = new(Response)
	var err = saHttp.Do(saHttp.Params{
		Method: saHttp.MethodPost,
		Url:    m.webhookUrl,
		Body:   args,
		Header: map[string]interface{}{"Content-Type": "application/json"},
	}, response)
	if err != nil || response.Code != 0 {
		return saError.OrErr(response.Msg, err)
	}
	return nil
}

func (m *FeiShu) SendTxtWithTitle(title string, txt string) error {
	var contents = make([]*PostItem, 0, 1)
	contents = append(contents, &PostItem{Tag: "text", Text: txt})
	var args = map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title":   title,
					"content": [][]*PostItem{contents},
				},
			},
		},
	}
	if m.secret != "" {
		var timestamp = time.Now().Unix()
		var sign = sign(timestamp, m.secret)
		args["timestamp"] = timestamp
		args["sign"] = sign
	}

	var response = new(Response)
	var err = saHttp.Do(saHttp.Params{Method: saHttp.MethodPost, Url: m.webhookUrl, Body: args, Header: map[string]interface{}{"Content-Type": "application/json"}}, response)
	if err != nil || response.Code != 0 {
		return saError.OrErr(response.Msg, err)
	}
	return nil
}

// SendPostTxt
// @Description: 发送富文本
func (m *FeiShu) SendPostTxt(title string, contents []*PostItem) error {
	var args = map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title":   title,
					"content": [][]*PostItem{contents},
				},
			},
		},
	}
	if m.secret != "" {
		var timestamp = time.Now().Unix()
		var sign = sign(timestamp, m.secret)
		args["timestamp"] = timestamp
		args["sign"] = sign
	}

	var response = new(Response)
	var err = saHttp.Do(saHttp.Params{Method: saHttp.MethodPost, Url: m.webhookUrl, Body: args, Header: map[string]interface{}{"Content-Type": "application/json"}}, response)
	if err != nil || response.Code != 0 {
		return saError.OrErr(response.Msg, err)
	}
	return nil
}

func SendCardMsg(webhookUrl string, content string) error {
	var response = &Response{}
	var err = saHttp.Do(saHttp.Params{
		Method: "POST",
		Url:    webhookUrl,
		Header: map[string]interface{}{"Content-Type": "application/json"},
		Body: map[string]interface{}{
			"msg_type": "interactive",
			"card": map[string]interface{}{
				"elements": []map[string]interface{}{
					{
						"tag": "div",
						"text": map[string]string{
							"content": content,
							"tag":     "lark_md",
						},
					},
				},
			},
		},
	}, response)
	if err != nil || response.Code != 0 {
		return saError.OrErr(response.Msg, err)
	}
	return nil
}

func sign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
