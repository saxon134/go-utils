package saHttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saUrl"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Get(url string, params map[string]string) (res string, err error) {
	return ToRequest("GET", url, params, nil)
}

func Post(url string, params map[string]string) (res string, err error) {
	return ToRequest("POST", url, params, nil)
}

func PostRequest(uri string, obj interface{}, headers map[string]string) (res string, err error) {
	var data = saData.String(obj)
	var body = bytes.NewBuffer([]byte(data))

	client := &http.Client{}
	var request *http.Request
	request, err = http.NewRequest("POST", uri, io.Reader(body))
	if err != nil {
		return "", err
	}

	for k, v := range headers {
		if k != "" && v != "" {
			request.Header.Set(k, v)
		}
	}

	var response *http.Response
	response, err = client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	status := response.StatusCode
	if status == 200 {
		var resData []byte
		resData, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		} else {
			return string(resData), nil
		}
	} else {
		return "", saError.New("error:" + saData.Itos(status))
	}
}

func PostJson(uri string, obj interface{}) (res string, err error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)

	body := bytes.NewBuffer(jsonData)
	response, err := http.Post(uri, "application/json; charset=UTF-8", body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseData []byte
	responseData, err = io.ReadAll(response.Body)
	if err == nil {
		res = string(responseData)
	}

	if response.StatusCode != http.StatusOK {
		return res, fmt.Errorf("http error : uri=%v , statusCode=%v , response=%v", uri, response.StatusCode, string(responseData))
	} else {
		return res, nil
	}
}

func Download(url string, suffix string) (localFilePath string, err error) {
	var (
		buf     = make([]byte, 32*1024)
		written int64
	)

	suffix = strings.TrimPrefix(suffix, ".")
	if suffix != "" && len(suffix) < 10 {
		suffix = "." + suffix
	}
	tmpFilePath := saData.RandomStr() + suffix

	//创建一个http client
	client := new(http.Client)

	//get方法获取资源
	response, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	//创建文件
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if response.Body == nil {
		return "", errors.New("body is null")
	}

	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := response.Body.Read(buf)
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

	var request *http.Request
	if method == "GET" {
		if len(params) > 0 {
			paramsStr := saUrl.QueryFromMap(params)
			if strings.HasSuffix(url, "?") == false {
				url += "?"
			}
			url += paramsStr
		}
		request, err = http.NewRequest(method, url, nil)
		if err != nil {
			return "", err
		}
	} else if method == "POST" {
		bodyStr := strings.TrimSpace(saUrl.QueryFromMap(params))
		request, err = http.NewRequest("POST", url, strings.NewReader(bodyStr))
		if err != nil {
			return "", err
		}

		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		return "", errors.New("暂只支持GET/POST方法")
	}

	for k, v := range header {
		if k != "" && v != "" {
			request.Header.Set(k, v)
		}
	}

	var httpRes *http.Response
	httpRes, err = client.Do(request)
	if err != nil {
		return "", err
	}
	defer httpRes.Body.Close()

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
		return "", saError.New("error:" + saData.Itos(status))
	}
}
