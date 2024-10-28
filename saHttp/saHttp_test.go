package saHttp

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"testing"
	"time"
)

func TestUpload(t *testing.T) {
	res, err := Upload("https://v3.biyingniao.com/web/bang/upload", map[string]string{"name": "file", "path": "./saHttp_test.go"}, nil, nil)
	fmt.Println(err)
	fmt.Println(res)
}

func TestErrCallback(t *testing.T) {
	SetErrCallback(func(request string) {
		fmt.Println(request)
	})

	_ = Do(Params{
		Url:             "vv",
		Query:           nil,
		Header:          nil,
		Body:            nil,
		Timeout:         0,
		CallbackWhenErr: true,
	}, nil)
}

func TestTokenBucket(t *testing.T) {
	fmt.Println("开始：", time.Now().Format(time.DateTime))
	var buket = NewTokenBucket(10, 3)
	for i := 0; i < 26; i++ {
		DoWithTokenBucket(Params{
			Url: "http://127.0.0.1:1000/notfound?i=" + saData.String(i),
			Retry: func(retry int, v interface{}, err error) bool {
				return false
			},
			Timeout: time.Second,
		}, nil, buket)
	}
	buket.Wait()
	fmt.Println("结束：", time.Now().Format(time.DateTime))
}
