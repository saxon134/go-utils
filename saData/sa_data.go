package saData

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

func ToBytes(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ToStr(d interface{}) (string, error) {
	s := ""
	switch v := d.(type) {
	case *string:
		var s string = *v
		return s, nil
	case string:
		var s string = v
		return s, nil
	case int:
	case int8:
	case int16:
	case int64:
	case float32:
	case float64:
		return fmt.Sprint(d), nil
	}

	s = fmt.Sprint(d)
	if s == "<nil>" {
		s = ""
		return s, errors.New("类型不匹配或者数据空")
	}

	return s, nil
}

func ToInt(d interface{}) (int, error) {
	if i, ok := d.(int); ok {
		return i, nil
	}

	if s, ok := d.(string); ok {
		if s == "" {
			return 0, nil
		}

		if i, err := strconv.Atoi(s); err == nil {
			return i, nil
		} else {
			return 0, err
		}
	}

	if f, ok := d.(float64); ok {
		return int(f), nil
	}

	if i, ok := d.(bool); ok {
		if i {
			return 1, nil
		} else {
			return 0, nil
		}
	}

	if i, ok := d.(int64); ok {
		return int(i), nil
	}

	if i, ok := d.(int32); ok {
		return int(i), nil
	}

	if i, ok := d.(int16); ok {
		return int(i), nil
	}

	if i, ok := d.(int8); ok {
		return int(i), nil
	}

	if i, ok := d.(uint8); ok {
		return int(i), nil
	}

	if i, ok := d.(uint16); ok {
		return int(i), nil
	}

	if i, ok := d.(uint32); ok {
		return int(i), nil
	}

	if i, ok := d.(uint64); ok {
		return int(i), nil
	}

	if i, ok := d.(uint64); ok {
		return int(i), nil
	}

	if f, ok := d.(float32); ok {
		return int(f), nil
	}

	return 0, errors.New("类型不匹配")
}

func ToInt8(d interface{}) (int8, error) {
	if v, err := ToInt(d); err == nil {
		return int8(v), nil
	} else {
		return 0, err
	}
}

func ToInt32(d interface{}) (int32, error) {
	if v, err := ToInt(d); err == nil {
		return int32(v), nil
	} else {
		return 0, err
	}
}

func ToFloat32(d interface{}) (float32, error) {
	if f, ok := d.(float32); ok {
		return f, nil
	}
	if f, ok := d.(float64); ok {
		return float32(f), nil
	}
	if f, err := ToInt(d); err == nil {
		return float32(f), nil
	}
	if s, ok := d.(string); ok {
		if f, err := strconv.ParseFloat(s, 32); err == nil {
			return float32(f), nil
		}
	}
	return 0, errors.New("类型不匹配")
}

func ToBool(d interface{}) (bool, error) {
	if i, ok := d.(bool); ok {
		if i {
			return true, nil
		} else {
			return false, nil
		}
	}

	if i, ok := d.(int); ok {
		return i != 0, nil
	}

	if s, ok := d.(string); ok {
		if i, err := strconv.Atoi(s); err == nil {
			return i != 0, nil
		} else {
			if s == "true" || s == "ok" {
				return true, nil
			}
			return false, nil
		}
	}

	if f, ok := d.(float64); ok {
		return int(f) != 0, nil
	}

	if i, ok := d.(int32); ok {
		return int(i) != 0, nil
	}

	if i, ok := d.(int16); ok {
		return int(i) != 0, nil
	}

	if i, ok := d.(int8); ok {
		return int(i) != 0, nil
	}

	if i, ok := d.(uint8); ok {
		return int(i) != 0, nil
	}

	if i, ok := d.(uint16); ok {
		return int(i) != 0, nil
	}

	if i, ok := d.(uint32); ok {
		return int(i) != 0, nil
	}

	if i, ok := d.(uint64); ok {
		return int(i) != 0, nil
	}

	if f, ok := d.(float32); ok {
		return int(f) != 0, nil
	}

	return false, errors.New("类型不匹配")
}

func ToInt64(d interface{}) (int64, error) {
	if i, ok := d.(int64); ok {
		return i, nil
	}

	if i, ok := d.(int); ok {
		return int64(i), nil
	}

	if s, ok := d.(string); ok {
		if s == "" {
			return 0, nil
		}

		if i64, err := strconv.ParseInt(s, 10, 64); err == nil {
			return i64, nil
		} else {
			return 0, err
		}
	}

	if f, ok := d.(float64); ok {
		return int64(f), nil
	}

	if i, ok := d.(bool); ok {
		if i {
			return 1, nil
		}
	}

	if i, ok := d.(int32); ok {
		return int64(i), nil
	}

	if i, ok := d.(int16); ok {
		return int64(i), nil
	}

	if i, ok := d.(int8); ok {
		return int64(i), nil
	}

	if i, ok := d.(uint8); ok {
		return int64(i), nil
	}

	if i, ok := d.(uint16); ok {
		return int64(i), nil
	}

	if i, ok := d.(uint32); ok {
		return int64(i), nil
	}

	if i, ok := d.(uint64); ok {
		return int64(i), nil
	}

	if f, ok := d.(float32); ok {
		return int64(f), nil
	}

	return 0, errors.New("类型不匹配")
}

func ToAry(arr interface{}) ([]interface{}, error) {
	if arr == nil {
		return nil, errors.New("empty")
	}

	if v, ok := (arr).([]interface{}); ok {
		return v, nil
	}

	if v, ok := (arr).(*[]interface{}); ok {
		return *v, nil
	}

	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("类型不匹配")
	}
	cnt := v.Len()
	ret := make([]interface{}, cnt)
	for i := 0; i < cnt; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret, nil
}

func ToAryStr(arr interface{}) ([]string, error) {
	if arr == nil {
		return nil, errors.New("empty")
	}

	if v, ok := (arr).([]string); ok {
		return v, nil
	}

	if v, ok := (arr).(*[]string); ok {
		return *v, nil
	}

	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("类型不匹配")
	}
	cnt := v.Len()
	ret := make([]string, cnt)
	for i := 0; i < cnt; i++ {
		d := v.Index(i).Interface()
		s, _ := ToStr(d)
		ret[i] = s
	}
	return ret, nil
}

func ToAryMap(arr interface{}) ([]map[string]interface{}, error) {
	if arr == nil {
		return nil, errors.New("empty")
	}

	if v, ok := (arr).([]map[string]interface{}); ok {
		return v, nil
	}

	if v, ok := (arr).(*[]map[string]interface{}); ok {
		return *v, nil
	}

	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("类型不匹配")
	}
	cnt := v.Len()
	ret := make([]map[string]interface{}, cnt)
	for i := 0; i < cnt; i++ {
		d := v.Index(i).Interface()
		m, _ := ToMap(d)
		ret[i] = m
	}
	return ret, nil
}

func ToMap(ma interface{}) (map[string]interface{}, error) {
	if ma == nil {
		return nil, errors.New("empty")
	}

	if v, ok := (ma).(map[string]interface{}); ok {
		return v, nil
	}

	if v, ok := (ma).(*map[string]interface{}); ok {
		return *v, nil
	}

	v := reflect.ValueOf(ma)
	if v.Kind() != reflect.Map {
		return nil, errors.New("类型不匹配")
	}

	keyAry := v.MapKeys()
	cnt := len(keyAry)
	ret := make(map[string]interface{}, cnt)
	for i := 0; i < cnt; i++ {
		key := keyAry[i].String()
		if key != "" {
			ret[key] = v.MapIndex(keyAry[i]).Interface()
		}
	}
	return ret, nil
}

func ToMapStr(ma interface{}) (map[string]string, error) {
	if v, ok := (ma).(map[string]string); ok {
		return v, nil
	}

	v := reflect.ValueOf(ma)
	if v.Kind() != reflect.Map {
		return nil, errors.New("类型不匹配")
	}

	keyAry := v.MapKeys()
	cnt := len(keyAry)
	ret := make(map[string]string, cnt)
	for i := 0; i < cnt; i++ {
		key := keyAry[i].String()
		if key != "" {
			d := v.MapIndex(keyAry[i]).Interface()
			s, _ := ToStr(d)
			ret[key] = s
		}
	}
	return ret, nil
}

func MapToMapStr(ma map[string]interface{}) map[string]string {
	if ma == nil {
		return nil
	}

	dic := map[string]string{}
	for k, v := range ma {
		str, _ := ToStr(v)
		if v != "" {
			dic[k] = str
		}
	}
	return dic
}

func I64FromMap(ma map[string]interface{}, key string) NullInt64 {
	var res = NullInt64{Ok: false, V: 0}

	if ma == nil {
		return res
	}

	var err error
	res.V, err = ToInt64(ma[key])
	if err == nil {
		res.Ok = true
	}
	return res
}

func IFromMap(ma map[string]interface{}, key string) NullInt {
	var res = NullInt{Ok: false, V: 0}

	if ma == nil {
		return res
	}

	var err error
	res.V, err = ToInt(ma[key])

	if err == nil {
		res.Ok = true
	}
	return res
}

func StrFromMap(ma map[string]interface{}, key string) NullString {
	if ma == nil {
		return NullString{V: "", Ok: false}
	}

	var err error
	var res = NullString{Ok: false, V: ""}
	res.V, err = ToStr(ma[key])
	if err == nil {
		res.Ok = true
	}
	return res
}

func ToPrice(f float32) float32 {
	// 保留两位小数，四舍五入
	return float32(int(f*100)) / float32(100)
}
