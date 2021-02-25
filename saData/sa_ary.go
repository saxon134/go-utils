package saData

import (
	"reflect"
	"strings"
)

func RemoveDuplicate(a interface{}) (ret []interface{}) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}

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

// => []{1,2,3}
func ToIdAry(args map[string]interface{}, key string) []int64 {
	if key == "" {
		key = "idAry"
	}

	if idStrAry, _ := ToAry(args[key]); idStrAry != nil && len(idStrAry) > 0 {
		idAry := make([]int64, 0, len(idStrAry))
		for _, v := range idStrAry {
			id, _ := ToInt64(v)
			if id > 0 {
				idAry = append(idAry, id)
			}
		}
		return idAry
	}

	return []int64{}
}

// []{1,2,3} => 1,2,3
func ToIds(ary []int64) string {
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
