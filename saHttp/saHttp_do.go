package saHttp

import (
	"bytes"
	"errors"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saHit"
	"github.com/saxon134/go-utils/saLog"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Limiter struct {
	Key         string
	MilliSecond int
}

type Params struct {
	Method          string                 //默认GET，仅支持GET/POST方法
	Url             string                 //不能空
	Query           map[string]interface{} //interface部分json序列化后进行UrlEncode
	QueryValues     url.Values
	Header          map[string]interface{} //interface部分会json序列化
	Body            map[string]interface{} //会进行json序列化或者query序列化（form表单），取决于content-type；默认query序列化
	BodyString      string                 //string类型body，仅content-type为application-json，且Body为空时有效
	BodyValues      url.Values             //仅content-type为application/x-www-form-urlencoded，且body为空时有效
	Timeout         time.Duration          //默认60秒
	CallbackWhenErr bool                   //是否在失败时回调，默认关闭
	Retry           func(retry int, v interface{}, err error) bool
}

type FormParams struct {
	Url             string //不能空
	Key             string //默认file
	LocalPath       string //本地文件路径，优先级高于File
	FileName        string //文件名称，带后缀
	FileReader      io.Reader
	Query           map[string]interface{} //interface部分json序列化后进行UrlEncode
	Header          map[string]interface{} //interface部分会json序列化
	Body            map[string]interface{}
	Timeout         time.Duration //默认60秒
	CallbackWhenErr bool          //是否在失败时回调，默认关闭
	Limiter         Limiter
}

type CallbackFun func(request string)

var _errCallbackFunc CallbackFun

// SetErrCallback
// @Description: 设置error时回调
func SetErrCallback(handle CallbackFun) {
	if handle != nil {
		_errCallbackFunc = handle
	}
}

// @Description: 发送请求
// @param params 请求参数
// @param resPtr 返回结果接收对象的指针，必须是指针或者空
// @return err
func Do(in Params, resPtr interface{}) (err error) {
	//最多重试100次
	for i := 0; i <= 100; i++ {
		err = _do(in, resPtr)
		if in.Retry == nil || in.Retry(i+1, resPtr, err) == false {
			break
		}
		time.Sleep(time.Millisecond * 1500)
	}
	return err
}

func _do(in Params, resPtr interface{}) (err error) {
	//接口调用失败时，回调
	if _errCallbackFunc != nil && in.CallbackWhenErr == true {
		defer func() {
			if err != nil {
				_errCallbackFunc(saData.String(map[string]string{
					"err": err.Error(),
				}))
			}
		}()
	}

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
		in.Timeout = time.Second * 60
	}
	client := &http.Client{Timeout: in.Timeout}
	var request *http.Request

	//绑定query参数
	var urlAry = strings.Split(in.Url, "#")
	var queryValues = url.Values{}
	{
		if len(in.Query) > 0 {
			for k, v := range in.Query {
				queryValues.Add(k, saData.String(v))
			}
		} else if len(in.QueryValues) > 0 {
			queryValues = in.QueryValues
		}

		if len(urlAry) == 2 {
			if strings.Contains(urlAry[1], "?") {
				in.Url += "&" + queryValues.Encode()
			} else {
				in.Url += "?" + queryValues.Encode()
			}
		} else {
			if strings.Contains(in.Url, "?") {
				in.Url += "&" + queryValues.Encode()
			} else {
				in.Url += "?" + queryValues.Encode()
			}
		}
	}

	//绑定body参数
	var bodyStr = ""
	var contentType = "application/json"
	if in.BodyString != "" {
		bodyStr = in.BodyString
	} else {
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
			if len(in.Body) == 0 && in.BodyString == "" {
				bodyStr = in.BodyValues.Encode()
			} else {
				var bodyValues = url.Values{}
				for k, v := range in.Body {
					if k != "" {
						bodyValues.Add(k, saData.String(v))
					}
				}
				bodyStr = bodyValues.Encode()
			}
		} else if strings.Contains(contentType, "application/json") {
			if in.Body != nil && len(in.Body) > 0 {
				bodyStr = saData.String(in.Body)
			}
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

			err = saData.BytesToModel(bAry, resPtr)
			if err != nil {
				saLog.Err(err)
				saLog.Err(string(bAry))
			}
			return nil
		}
	} else {
		err = &url.Error{Op: in.Method, URL: in.Url, Err: errors.New(resp.Status)}
		if resPtr == nil {
			return err
		}

		if bAry, e := io.ReadAll(resp.Body); e == nil {
			if b, ok := resPtr.(*[]byte); ok {
				*b = bAry
			} else {
				e = saData.BytesToModel(bAry, resPtr)
				if e != nil {
					saLog.Err(e)
					saLog.Err(string(bAry))
				}
			}
		}
		return err
	}
}

func MultiForm(in FormParams, resPtr interface{}) (err error) {
	//接口调用失败时，回调
	if _errCallbackFunc != nil && in.CallbackWhenErr == true {
		defer func() {
			if err != nil {
				_errCallbackFunc(saData.String(map[string]string{
					"err": err.Error(),
				}))
			}
		}()
	}

	if in.Url == "" {
		return errors.New("缺少URL")
	}

	if in.Timeout <= 0 {
		in.Timeout = time.Second * 60
	}
	in.Key = saHit.OrStr(in.Key, "file")

	//绑定query参数
	var urlAry = strings.Split(in.Url, "#")
	for k, v := range in.Query {
		var queryValues = url.Values{}
		queryValues.Add(k, saData.String(v))
		if len(urlAry) == 2 {
			if strings.Contains(urlAry[1], "?") {
				in.Url += "&" + queryValues.Encode()
			} else {
				in.Url += "?" + queryValues.Encode()
			}
		} else {
			if strings.Contains(in.Url, "?") {
				in.Url += "&" + queryValues.Encode()
			} else {
				in.Url += "?" + queryValues.Encode()
			}
		}
	}

	//新建请求body
	var requestBody = &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	{
		if in.FileName == "" {
			in.FileName = filepath.Base(in.LocalPath)
			if in.FileName == "" {
				in.FileName = saData.String(time.Now().UnixMilli())
			}
		}

		if in.LocalPath != "" {
			in.FileReader, err = os.Open(in.LocalPath)
			if err != nil {
				return err
			}

			var part io.Writer
			part, err = writer.CreateFormFile(in.Key, in.FileName)
			if err != nil {
				return err
			}
			_, err = io.Copy(part, in.FileReader)
		} else if in.FileReader != nil {
			var part io.Writer
			part, err = writer.CreateFormFile(in.Key, in.FileName)
			if err != nil {
				return err
			}
			_, err = io.Copy(part, in.FileReader)
		}

		// 其他参数列表写入 body
		for k, v := range in.Body {
			if err = writer.WriteField(k, saData.String(v)); err != nil {
				return err
			}
		}

		err = writer.Close()
		if err != nil {
			return err
		}
	}

	// 创建请求
	var req *http.Request
	req, err = http.NewRequest("POST", in.Url, requestBody)
	if err != nil {
		panic(err)
	}

	// 设置 header
	for k, v := range in.Header {
		if saData.InStrs(k, []string{"Content-Type", "content-type"}) {
			continue
		}
		req.Header.Set(k, saData.String(v))
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if resPtr == nil {
			return nil
		}

		var bAry []byte
		bAry, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		} else {
			if b, ok := resPtr.(*[]byte); ok {
				*b = bAry
				return nil
			}

			err = saData.BytesToModel(bAry, resPtr)
			if err != nil {
				saLog.Err(err)
				saLog.Err(string(bAry))
			}
			return nil
		}
	} else {
		err = &url.Error{URL: in.Url, Err: errors.New(resp.Status)}
		if resPtr == nil {
			return err
		}

		if bAry, e := io.ReadAll(resp.Body); e == nil {
			if b, ok := resPtr.(*[]byte); ok {
				*b = bAry
			} else {
				e = saData.BytesToModel(bAry, resPtr)
				if e != nil {
					saLog.Err(e)
					saLog.Err(string(bAry))
				}
			}
		}
		return err
	}
}
