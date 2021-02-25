package saData

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

func DataToJson(m interface{}) (string, error) {
	switch v := m.(type) {
	case *string:
		var s string = *v
		return s, nil
	case string:
		var s string = v
		return s, nil
	case int, int8, int16, int64, float32, float64:
		return fmt.Sprint(m), nil
	}

	if bAry, err := json.Marshal(m); err == nil && bAry != nil {
		s := string(bAry)
		if s == "null" {
			return "", nil
		}
		return s, nil
	} else {
		return "", err
	}
}

func LenCheck(m interface{}, max int) error {
	str, _ := DataToJson(m)
	if StrLen(str) <= max {
		return nil
	}

	return errors.New("超出范围")
}

func DataToMap(v interface{}) (map[string]interface{}, error) {
	if d, ok := v.(map[string]interface{}); ok {
		return d, nil
	}
	if d, ok := v.(*map[string]interface{}); ok {
		return *d, nil
	}

	bStr, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var dic map[string]interface{}
	err = json.Unmarshal(bStr, &dic)
	if err != nil {
		return nil, err
	}
	return dic, nil
}

func JsonToStruct(s string, m interface{}) error {
	err := json.Unmarshal([]byte(s), m)
	if err != nil {
		return err
	}
	return nil
}

func JsonToMap(s string) (ma map[string]interface{}, err error) {
	var m map[string]interface{}
	err = json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func JsonToSlice(s string) (*[]interface{}, error) {
	var ary []interface{}
	err := json.Unmarshal([]byte(s), &ary)
	if err != nil {
		return nil, err
	}
	return &ary, nil
}
