package ding

import (
	"fmt"
	"github.com/saxon134/go-utils/saHttp"
)

func SendTxt(title, msg string, webhook string) error {
	//todo 优化 接受返回数据
	var err = saHttp.Do(saHttp.Params{
		Method: saHttp.MethodPost,
		Url:    webhook,
		Header: map[string]interface{}{"Content-Type": "application/json"},
		Body: map[string]interface{}{
			"msgtype": "markdown",
			"markdown": map[string]string{
				"title": title,
				"text":  msg,
			},
		},
	}, nil)
	return err
}

func SendImg(title, img string, webhook string) error {
	//todo 优化 接受返回数据

	var err = saHttp.Do(saHttp.Params{
		Method: saHttp.MethodPost,
		Url:    webhook,
		Header: map[string]interface{}{"Content-Type": "application/json"},
		Body: map[string]interface{}{
			"msgtype": "markdown",
			"markdown": map[string]string{
				"title": title,
				"text":  fmt.Sprintf(`![](%s)`, img),
			},
		},
	}, nil)
	return err
}
