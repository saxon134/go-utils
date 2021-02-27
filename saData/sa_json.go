package saData

import (
	"fmt"
	json "github.com/json-iterator/go"
)

var j json.API

func init() {
	j = json.Config{
		EscapeHTML:             true,
		SortMapKeys:            false,
		ValidateJsonRawMessage: true,
		UseNumber:              true,
		CaseSensitive:          false,
	}.Froze()
}

func DataToJson(m interface{}) (string, error) {
	switch v := m.(type) {
	case string:
		return v, nil
	case int, int8, int16, int64, float32, float64:
		return fmt.Sprint(m), nil
	case *string:
		var s string = *v
		return s, nil
	}

	if bAry, err := j.Marshal(m); err == nil && bAry != nil {
		s := string(bAry)
		if s == "null" {
			return "", nil
		}
		return s, nil
	} else {
		return "", err
	}
}

func DataToMap(v interface{}) (map[string]interface{}, error) {
	if d, ok := v.(map[string]interface{}); ok {
		return d, nil
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

func JsonToStruct(data []byte, m interface{}) error {
	err := json.Unmarshal(data, m)
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
