package saHttp

import (
	"errors"
	"github.com/saxon134/go-utils/saData"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Params struct {
	Method string                 //默认GET，仅支持GET/POST方法
	Url    string                 //不能空
	Query  map[string]interface{} //interface部分json序列化后进行UrlEncode
	Header map[string]interface{} //interface部分会json序列化
	Body   map[string]interface{} //会进行json序列化或者query序列化（form表单），取决于content-type；默认query序列化
	Timeout time.Duration //默认10秒
}

// Do
// @Description: 发送请求
// @param params 请求参数
// @param resPtr 返回结果接收对象的指针，必须是指针或者空
// @return err
func Do(in Params, resPtr interface{}) (err error) {
	if in.Url == "" {
		return errors.New("缺少URL")
	}

	if in.Method == "" {
		in.Method = "GET"
	}
	in.Method = strings.ToUpper(in.Method)
	if in.Method != "GET" && in.Method != "POST" {
		return errors.New("只支持GET/POST方法")
	}

	if in.Timeout == 0 {
		in.Timeout = time.Second*10
	}
	client := &http.Client{Timeout:in.Timeout}
	var request *http.Request

	//绑定query参数
	var urlAry = strings.Split(in.Url, "#")
	for k, v := range in.Query {
		values := url.Values{}
		values.Add(k, saData.String(v))
		if len(urlAry) == 2 {
			if strings.Contains(urlAry[1], "?") {
				in.Url += "&" + values.Encode()
			} else {
				in.Url += "?" + values.Encode()
			}
		} else {
			if strings.Contains(in.Url, "?") {
				in.Url += "&" + values.Encode()
			} else {
				in.Url += "?" + values.Encode()
			}
		}
	}

	//绑定body参数
	var bodyStr = ""
	var contentType = "application/x-www-form-urlencoded"
	{
		if in.Header != nil && len(in.Header) > 0 {
			var ct = saData.String(in.Header["content-type"])
			if ct == "" {
				ct = saData.String(in.Header["Content-Type"])
			}
			if ct != "" {
				contentType = ct
			}
		}

		if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			urlV := url.Values{}
			for k, v := range in.Body {
				if k != "" {
					urlV.Add(k, saData.String(v))
				}
			}
			bodyStr = urlV.Encode()
		} else if strings.Contains(contentType, "application/json") {
			bodyStr = saData.String(in.Body)
		}
	}

	//初始化请求，并传入body参数
	if bodyStr != "" {
		request, err = http.NewRequest(in.Method, in.Url, strings.NewReader(bodyStr))
	} else {
		request, err = http.NewRequest(in.Method, in.Url, nil)
	}
	if err != nil {
		return err
	}

	//绑定header参数
	for k, v := range in.Header {
		if k != "" && v != nil {
			request.Header.Set(k, saData.String(v))
		}
	}

	//发送请求
	var resp *http.Response
	resp, err = client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	if status == 200 {
		if resPtr == nil {
			return nil
		}

		var bAry []byte
		bAry, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		} else {
			if bytes, ok := resPtr.(*[]byte); ok {
				*bytes = bAry
				return nil
			}

			return saData.BytesToModel(bAry, resPtr)
		}
	} else {
		return errors.New("saHttp response error:" + saData.Itos(status))
	}
}
