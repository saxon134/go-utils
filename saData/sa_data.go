package saData

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/saxon134/go-utils/saHit"
	"reflect"
	"strconv"
)

func ToBytes(data interface{}) ([]byte, error) {
	if data == nil {
		return []byte{}, nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func String(data interface{}) string {
	s, _ := ToStr(data)
	return s
}

func ToStr(data interface{}) (string, error) {
	if data == nil {
		return "", nil
	}
	defer func() {
		_ = recover()
	}()

	switch v := data.(type) {
	case string:
		return v, nil
	case int, int8, int16, int64, float32, float64:
		return fmt.Sprint(data), nil
	case *string:
		return *v, nil
	}

	bAry, err := json.Marshal(data)
	if err == nil && bAry != nil {
		s := BytesToStr(bAry)
		s = saHit.Str(s == "null", "", s)
		return s, nil
	}
	return "", err
}

func ToMap(data interface{}) (map[string]interface{}, error) {
	if data == nil {
		return map[string]interface{}{}, nil
	}
	if d, ok := data.(map[string]interface{}); ok {
		return d, nil
	}
	if d, ok := data.(*map[string]interface{}); ok {
		return *d, nil
	}

	defer func() {
		_ = recover()
	}()

	//json
	if s, ok := data.(string); ok {
		var dic map[string]interface{}
		var err error
		if s != "" {
			err = json.Unmarshal(StrToBytes(s), &dic)
		}
		return dic, err
	}

	//反射
	v := reflect.ValueOf(data)
	vKind := v.Kind()
	if vKind == reflect.Map {
		keyAry := v.MapKeys()
		cnt := len(keyAry)
		ret := make(map[string]interface{}, cnt)
		for _, key := range keyAry {
			key_s := fmt.Sprint(key.Convert(v.Type().Key()))
			ret[key_s] = v.MapIndex(key).Interface()
		}
		return ret, nil
	} else if vKind == reflect.Struct || vKind == reflect.Ptr || vKind == reflect.Interface {
		var dic map[string]interface{}
		bAry, err := json.Marshal(data)
		if err != nil {
			return map[string]interface{}{}, err
		}
		err = json.Unmarshal(bAry, &dic)
		return dic, err
	}
	return map[string]interface{}{}, errors.New("类型不匹配")
}

func ToStrMap(data interface{}) (map[string]string, error) {
	if data == nil {
		return map[string]string{}, nil
	}
	if d, ok := data.(map[string]string); ok {
		return d, nil
	}
	if d, ok := data.(*map[string]string); ok {
		return *d, nil
	}

	defer func() {
		_ = recover()
	}()

	//json
	if s, ok := data.(string); ok {
		var dic map[string]string
		var err error
		if s != "" {
			err = json.Unmarshal(StrToBytes(s), &dic)
		}
		return dic, err
	}

	//反射
	v := reflect.ValueOf(data)
	vKind := v.Kind()
	if vKind == reflect.Map {
		keyAry := v.MapKeys()
		cnt := len(keyAry)
		ret := make(map[string]string, cnt)
		for i := 0; i < cnt; i++ {
			key := keyAry[i].String()
			if key != "" {
				val := v.MapIndex(keyAry[i]).Interface()
				s, _ := ToStr(val)
				ret[key] = s
			}
		}
		return ret, nil
	} else if vKind == reflect.Struct || vKind == reflect.Ptr || vKind == reflect.Interface {
		var dic map[string]string
		bAry, err := json.Marshal(data)
		if err != nil {
			return map[string]string{}, err
		}
		err = json.Unmarshal(bAry, &dic)
		return dic, err
	}
	return map[string]string{}, errors.New("类型不匹配")
}

func ToAry(data interface{}) ([]interface{}, error) {
	if data == nil {
		return []interface{}{}, nil
	}
	if v, ok := (data).([]interface{}); ok {
		return v, nil
	}
	if v, ok := (data).(*[]interface{}); ok {
		return *v, nil
	}

	defer func() {
		_ = recover()
	}()

	//json
	if s, ok := data.(string); ok {
		var ary []interface{}
		var err error
		if s != "" {
			err = json.Unmarshal(StrToBytes(s), &ary)
		}
		return ary, err
	}

	//反射
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		cnt := v.Len()
		ret := make([]interface{}, cnt)
		for i := 0; i < cnt; i++ {
			ret[i] = v.Index(i).Interface()
		}
		return ret, nil
	}
	return []interface{}{}, errors.New("类型不匹配")
}

func ToMapAry(data interface{}) ([]map[string]interface{}, error) {
	if data == nil {
		return []map[string]interface{}{}, nil
	}
	if ary, ok := data.([]map[string]interface{}); ok {
		return ary, nil
	}
	if ary, ok := data.(*[]map[string]interface{}); ok {
		return *ary, nil
	}

	defer func() {
		_ = recover()
	}()

	//json转ary
	if s, ok := data.(string); ok {
		var ary []map[string]interface{}
		var err error
		if s != "" {
			err = json.Unmarshal(StrToBytes(s), &ary)
		}
		return ary, err
	}

	//反射map
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		cnt := v.Len()
		ary := make([]map[string]interface{}, cnt)
		for i := 0; i < cnt; i++ {
			d := v.Index(i).Interface()
			m, _ := ToMap(d)
			ary[i] = m
		}
		return ary, nil
	}
	return []map[string]interface{}{}, errors.New("类型不匹配")
}

func ToStrAry(data interface{}) ([]string, error) {
	if data == nil {
		return []string{}, nil
	}
	if v, ok := (data).([]string); ok {
		return v, nil
	}
	if v, ok := (data).(*[]string); ok {
		return *v, nil
	}

	defer func() {
		_ = recover()
	}()

	//json
	if s, ok := data.(string); ok {
		var ary = []string{}
		var err error
		if s != "" {
			err = json.Unmarshal(StrToBytes(s), &ary)
		}
		return ary, err
	}

	//反射
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		cnt := v.Len()
		ary := make([]string, cnt)
		for i := 0; i < cnt; i++ {
			d := v.Index(i).Interface()
			s, _ := ToStr(d)
			ary[i] = s
		}
		return ary, nil
	}
	return []string{}, errors.New("类型不匹配")
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

	if s, e := ToStr(d); e == nil {
		if s == "" {
			return 0, nil
		}
		return strconv.Atoi(s)
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

	if s, ok := d.(string); ok {
		f, err := strconv.ParseFloat(s, 32)
		if err == nil {
			return float32(f), nil
		}
	}

	f, err := ToInt(d)
	if err == nil {
		return float32(f), nil
	}
	return 0, errors.New("类型不匹配")
}

func ToFloat64(d interface{}) (float64, error) {
	if f, ok := d.(float64); ok {
		return f, nil
	}
	if f, ok := d.(float32); ok {
		return float64(f), nil
	}
	if s, ok := d.(string); ok {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f, nil
		}
	}
	if f, err := ToInt(d); err == nil {
		return float64(f), nil
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
		return strconv.ParseInt(s, 10, 64)
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

	if s, e := ToStr(d); e == nil {
		if s == "" {
			return 0, nil
		}
		return strconv.ParseInt(s, 10, 64)
	}

	return 0, errors.New("类型不匹配")
}

func StrToData(str string) (interface{}, error) {
	defer func() {
		_ = recover()
	}()

	var err error
	bAry := StrToBytes(str)

	dic := map[string]interface{}{}
	err = json.Unmarshal(bAry, &dic)
	if err == nil {
		return dic, nil
	}

	ary := make([]interface{}, 0, 10)
	err = json.Unmarshal(bAry, &ary)
	if err == nil {
		return ary, nil
	}
	return str, nil
}

func StrToModel(str string, m interface{}) error {
	defer func() {
		_ = recover()
	}()

	err := json.Unmarshal(StrToBytes(str), m)
	return err
}

func BytesToModel(b []byte, m interface{}) error {
	defer func() {
		_ = recover()
	}()

	err := json.Unmarshal(b, m)
	return err
}
