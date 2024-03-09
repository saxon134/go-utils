package saHttp

import (
	"fmt"
	"testing"
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
