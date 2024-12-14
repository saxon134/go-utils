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

// AryToIds
// fullComma为true时： []{1,2,3} => ,1,2,3,
// fullComma为false时： []{1,2,3} => 1,2,3
func AryToIds(ary []int64, fullComma bool) string {
	if ary == nil || len(ary) == 0 {
		return ""
	}

	var ids = ""
	for i, v := range ary {
		if v > 0 {
			ids += I64tos(v)
		}
		if i+1 < len(ary) {
			ids += ","
		}
	}

	if fullComma {
		if ids != "" {
			ids = "," + ids + ","
		}
	}

	return ids
}

// ToIds
// fullComma为true时： []{1,2,3} => ,1,2,3,
// fullComma为false时： []{1,2,3} => 1,2,3
// ary支持类型：[]string []int64 []int
func ToIds(ary interface{}, fullComma bool) string {
	var ids = ""
	switch vv := ary.(type) {
	case []int64:
		for _, v := range vv {
			if v > 0 {
				ids += I64tos(v) + ","
			}
		}
	case []int:
		for _, v := range vv {
			if v > 0 {
				ids += Itos(v)+ ","
			}
		}
	case []string:
		for _, v := range vv {
			if v != "" {
				ids += v+ ","
			}
		}
	}

	if ids!= "" {
		if fullComma {
			ids = "," + ids
		} else {
			ids = strings.TrimSuffix(ids, ",")
		}
	}
	return ids
}

// []{'1','2','3'} => '1','2','3'
func ToSQLIds(ary []string) string {
	if len(ary) == 0 {
		return ""
	}

	var sql = ""
	for _, v := range ary {
		if v != "" {
			sql += "'" + v + "',"
		}
	}
	sql = strings.TrimSuffix(sql, ",")
	return sql
}

// []{'1','2','3'} => '1','2','3'
func AryToSQLIds(ary interface{}) string {
	var sql = ""
	switch vv := ary.(type) {
	case []int64:
		for _, v := range vv {
			if v > 0 {
				sql += "'" + String(v) + "',"
			}
		}
	case []int:
		for _, v := range vv {
			if v > 0 {
				sql += "'" + String(v) + "',"
			}
		}
	case []string:
		for _, v := range vv {
			if v !="" {
				sql += "'" + v + "',"
			}
		}
	}

	sql = strings.TrimSuffix(sql, ",")
	return sql
}

// Split  去空、去重 1,2,3 => []{"1","2","3"}
func Split(s, sep string) []string {
	var ary = strings.Split(s, sep)
	var resAry = make([]string, 0, len(ary))
	for _, v := range ary {
		if v != "" {
			var exist = false
			for _, e := range resAry {
				if e == v {
					exist = true
					break
				}
			}
			if exist == false {
				resAry = append(resAry, v)
			}
		}
	}
	return resAry
}

// IdsToAry 去零、去重 1,2,3 => []{1,2,3}
func IdsToAry(str string) []int64 {
	if str == "" {
		return []int64{}
	}

	var ary = strings.Split(str, ",")
	var resAry = make([]int64, 0, len(ary))
	for _, v := range ary {
		if id, _ := ToInt64(v); id != 0 {
			var exist = false
			for _, e := range resAry {
				if e == id {
					exist = true
					break
				}
			}
			if exist == false {
				resAry = append(resAry, id)
			}
		}
	}
	return resAry
}

func IdsToIntAry(str string) []int {
	if str == "" {
		return []int{}
	}

	var ary = strings.Split(str, ",")
	var resAry = make([]int, 0, len(ary))
	for _, v := range ary {
		if id, _ := ToInt(v); id != 0 {
			var exist = false
			for _, e := range resAry {
				if e == id {
					exist = true
					break
				}
			}
			if exist == false {
				resAry = append(resAry, id)
			}
		}
	}
	return resAry
}

// FormatIds 前后加逗号是为了方便SQL查询过滤
// 1,2,3 => ,1,2,3,
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

// InArray 注意：只支持基础类型数据
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

func InStrs(item string, ary []string) (exist bool) {
	for _, v := range ary {
		if item == v {
			return true
		}
	}
	return false
}

func ContainStrs(item string, ary []string) (exist bool) {
	for _, v := range ary {
		if strings.Contains(v, item) {
			return true
		}
	}
	return false
}

func InInt(item int, ary []int) (exist bool) {
	for _, v := range ary {
		if item == v {
			return true
		}
	}
	return false
}

func InInt64(item int64, ary []int64) (exist bool) {
	for _, v := range ary {
		if item == v {
			return true
		}
	}
	return false
}

// 注意：只支持基础类型数据，会排除已存在的
func AppendId(ary []int64, ids ...int64) []int64 {
	for _, id := range ids {
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
	}
	return ary
}

func AppendStr(ary []string, strs ...string) []string {
	for _, s := range strs {
		if s != "" {
			exist := false
			for _, v := range ary {
				if v == s {
					exist = true
					break
				}
			}

			if exist == false {
				ary = append(ary, s)
			}
		}
	}
	return ary
}
