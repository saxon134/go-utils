package saCache

import (
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	v, err := MGetWithFunc("appInfo", "1", func(id string) (interface{}, error) {
		return "测试应用", nil
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)
}
