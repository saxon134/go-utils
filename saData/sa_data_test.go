package saData

import (
	"fmt"
	"testing"
)

func TestDataConvert(t *testing.T) {
	dic, err := ToMap([]map[string]string{})
	fmt.Println(dic)
	fmt.Println(err)
}
