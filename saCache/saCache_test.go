package saCache

import "testing"

func TestCache(t *testing.T) {
	Get("appInfo", "1", func() (interface{}, error) {
		return "测试应用", nil
	})
}
