package saHttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saError"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Get(url string, params map[string]string) (res string, err error) {
	return ToRequest("GET", url, params, nil)
}

func Post(url string, params map[string]string) (res string, err error) {
	return ToRequest("POST", url, params, nil)
}

func PostJson(uri string, obj interface{}) (res string, contentType string, err error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", "", err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)

	body := bytes.NewBuffer([]byte(jsonData))
	response, err := http.Post(uri, "application/json;charset=utf-8", body)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}

	if responseData, err := ioutil.ReadAll(response.Body); err == nil {
		contentType = response.Header.Get("Content-Type")
		return string(responseData), contentType, err
	} else {
		return "", "", err
	}
}

func PostRequest(uri string, obj interface{}, headers map[string]string) (res string, err error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)

	body := bytes.NewBuffer([]byte(jsonData))

	client := &http.Client{}
	var req *http.Request
	req, err = http.NewRequest("POST", uri, io.Reader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	for k, v := range headers {
		if k != "" && v != "" {
			req.Header.Set(k, v)
		}
	}

	var httpRes *http.Response
	httpRes, err = client.Do(req)
	if err != nil {
		return "", err
	}

	status := httpRes.StatusCode
	if status == 200 {
		var resData []byte
		resData, err = ioutil.ReadAll(httpRes.Body)
		if err != nil {
			return "", err
		} else {
			return string(resData), nil
		}
	} else {
		return "", saError.StackError("error:" + saData.Itos(status))
	}
}

func Download(url string) (localFilePath string, err error) {
	var (
		buf     = make([]byte, 32*1024)
		written int64
	)

	tmpFilePath := saData.RandomStr() + ".download"

	//创建一个http client
	client := new(http.Client)

	//get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	//创建文件
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if resp.Body == nil {
		return "", errors.New("body is null")
	}
	defer resp.Body.Close()

	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}

	_ = file.Close()

	if err == nil {
		return tmpFilePath, nil
	}
	return "", err
}

func ToRequest(method string, url string, params map[string]string, header map[string]string) (res string, err error) {
	client := &http.Client{}

	if method == "" {
		method = "GET"
	}
	method = strings.ToUpper(method)

	var req *http.Request
	if method == "GET" {
		if len(params) > 0 {
			paramsStr := QueryEncode(params)
			if strings.HasSuffix(url, "?") == false {
				url += "?"
			}
			url += paramsStr
		}
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return "", err
		}
	} else if method == "POST" {
		bodyStr := strings.TrimSpace(QueryEncode(params))
		req, err = http.NewRequest("POST", url, strings.NewReader(bodyStr))
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		return "", errors.New("暂只支持GET/POST方法")
	}

	for k, v := range header {
		if k != "" && v != "" {
			req.Header.Set(k, v)
		}
	}

	var httpRes *http.Response
	httpRes, err = client.Do(req)
	if err != nil {
		return "", err
	}

	status := httpRes.StatusCode
	if status == 200 {
		var resData []byte
		resData, err = ioutil.ReadAll(httpRes.Body)
		if err != nil {
			return "", err
		} else {
			return string(resData), nil
		}
	} else {
		return "", saError.StackError("error:" + saData.Itos(status))
	}
}

func StrEncode(s string) string {
	if s != "" {
		v := url.Values{}
		v.Add("k", s)
		s = v.Encode()
		return string([]rune(s)[2:])
	}
	return ""
}

func QueryEncode(m map[string]string) string {
	if m != nil {
		urlV := url.Values{}
		for k, v := range m {
			if k != "" {
				urlV.Add(k, v)
			}
		}
		return urlV.Encode()
	}
	return ""
}

func QueryDecode(urlStr string) map[string]string {
	values, _ := url.ParseQuery(urlStr)
	m := map[string]string{}
	for k, v := range values {
		if k != "" {
			m[k] = v[0]
		}
	}

	return m
}

//返回结果是： /a/b/c
func ConnPath(r string, path string) (full string) {
	r = strings.TrimSuffix(r, "/")
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	if len(r) > 0 {
		full = r + "/" + path
	} else {
		full = path
	}

	return full
}
