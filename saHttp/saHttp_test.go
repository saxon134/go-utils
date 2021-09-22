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
