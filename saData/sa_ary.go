package saData

import (
	"reflect"
	"strings"
)

func GetIdByKeys(args map[string]interface{}, key1 string, key2 string) int64 {
	if args == nil {
		return 0
	}

	if key1 == "" {
		return 0
	}

	id, _ := ToInt64(args[key1])
	if id > 0 {
		return id
	}

	if key2 == "" {
		return 0
	}

	id, _ = ToInt64(args[key2])
	return id
}

// []{1,2,3} => 1,2,3
func AryToIds(ary []int64) string {
	if ary == nil || len(ary) == 0 {
		return ""
	}

	var ids = ""
	for _, v := range ary {
		if v > 0 {
			ids += I64tos(v) + ","
		}
	}
	if ids != "" {
		ids = strings.TrimSuffix(ids, ",")
	}

	return ids
}

// 1,2,3 => []{1,2,3}
func IdsToAry(str string) []int64 {
	if str == "" {
		return []int64{}
	}

	ary := strings.Split(str, ",")
	idAry := make([]int64, 0, len(ary))
	for _, v := range ary {
		id, _ := ToInt64(v)
		if id > 0 {
			idAry = append(idAry, id)
		}
	}

	return idAry
}

func FormatIds(str string) string {
	if str == "" {
		return ""
	}

	ary := strings.Split(str, ",")
	str = ""
	for _, v := range ary {
		id, _ := ToInt64(v)
		if id > 0 {
			str += I64tos(id) + ","
		}
	}
	str = strings.TrimSuffix(str, ",")
	return str
}

//注意：只支持基础类型数据
func InArray(item interface{}, ary interface{}) (exist bool) {
	v1 := String(item)
	switch vv := ary.(type) {
	case []int64:
		for _, v := range vv {
			v2 := I64tos(v)
			if v1 == v2 {
				return true
			}
		}
		return false
	case []int8:
		for _, v := range vv {
			v2 := Itos(int(v))
			if v1 == v2 {
				return true
			}
		}
		return false
	case []int:
		for _, v := range vv {
			v2 := Itos(v)
			if v1 == v2 {
				return true
			}
		}
		return false
	case []string:
		for _, v := range vv {
			if v1 == v {
				return true
			}
		}
		return false
	case []interface{}:
		for _, v := range vv {
			v2 := String(v)
			if v1 == v2 {
				return true
			}
		}
		return false
	}
	return false
}

func InArrayFun(ary interface{}, fun func(i int) bool) bool {
	val := reflect.ValueOf(ary)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if fun(i) == true {
				return true
			}
		}
	}
	return false
}

//注意：只支持基础类型数据
func AppendId(ary []int64, id int64) []int64 {
	if id > 0 {
		exist := false
		for _, v := range ary {
			if v == id {
				exist = true
				break
			}
		}

		if exist == false {
			ary = append(ary, id)
		}
	}
	return ary
}
