package saData

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

//验证参数
func ValidCheck(ptr interface{}) (err error) {
	objT := reflect.TypeOf(ptr)
	objV := reflect.ValueOf(ptr)
	if objT.Kind() == reflect.Struct || (objT.Kind() == reflect.Ptr && objT.Elem().Kind() == reflect.Struct) {
		objT = objT.Elem()
		objV = objV.Elem()
	} else {
		err = fmt.Errorf("%v must be a struct or a struct pointer", ptr)
		return
	}

	fieldNum := objT.NumField()
	for i := 0; i < fieldNum; i++ {
		tags := []string{}
		tag := objT.Field(i).Tag.Get("valid")
		if tag == "" {
			return nil
		}

		tags = strings.Split(tag, ";")
		for _, t := range tags {
			fieldV := objV.Field(i)
			if t == "required" {
				if fieldV.IsValid() == false {
					err = errors.New("数据有误")
					return
				}
			} else if strings.HasPrefix(t, "in(") {
				tmpStr := strings.TrimPrefix(t, "in(")
				tmpStr = strings.TrimSuffix(tmpStr, ")")
				ary := strings.Split(tmpStr, ",")
				s := fieldV.String()
				existed := false
				for _, v := range ary {
					if v == s {
						existed = true
						break
					}
				}
				if existed == false {
					err = errors.New("数据有误")
					return
				}
			} else if strings.HasPrefix(t, "enum(") {
				tmpStr := strings.TrimPrefix(t, "enum(")
				tmpStr = strings.TrimSuffix(tmpStr, ")")
				ary := strings.Split(tmpStr, ",")
				s := fieldV.String()
				existed := false
				for _, v := range ary {
					ary2 := strings.Split(v, ":")
					if len(ary2) >= 1 {
						if s == ary2[0] {
							existed = true
							break
						}
					}
				}
				if existed == false {
					err = errors.New("数据有误")
					return
				}
			} else if t == ">" || t == ">=" || t == "<" || t == "<=" || t == "<>" {
				fieldT := fieldV.Type().Kind()
				if fieldT == reflect.String {
					s := strings.TrimPrefix(t, ">=")
					s = strings.TrimPrefix(s, "<=")
					s = strings.TrimPrefix(s, "<>")
					s = strings.TrimPrefix(s, ">")
					s = strings.TrimPrefix(s, "<")
					l2, _ := Stoi(s)
					l1 := len(fieldV.String())
					if l2 > 0 {
						switch t {
						case ">":
							if l1 <= l2 {
								err = errors.New("数据格式有误")
								return
							}
						case ">=":
							if l1 < l2 {
								err = errors.New("数据格式有误")
								return
							}
						case "<":
							if l1 >= l2 {
								err = errors.New("数据格式有误")
								return
							}
						case "<=":
							if l1 > l2 {
								err = errors.New("数据格式有误")
								return
							}
						case "<>":
							if l1 == l2 {
								err = errors.New("数据格式有误")
								return
							}
						}
					}
				} else if fieldT == reflect.Int || fieldT == reflect.Int64 || fieldT == reflect.Int8 || fieldT == reflect.Int32 || fieldT == reflect.Int16 {
					s := strings.TrimPrefix(t, ">=")
					s = strings.TrimPrefix(s, "<=")
					s = strings.TrimPrefix(s, "<>")
					s = strings.TrimPrefix(s, ">")
					s = strings.TrimPrefix(s, "<")
					i2, _ := Stoi64(s)
					i1 := fieldV.Int()
					switch t {
					case ">":
						if i1 <= i2 {
							err = errors.New("数据格式有误")
							return
						}
					case ">=":
						if i1 < i2 {
							err = errors.New("数据格式有误")
							return
						}
					case "<":
						if i1 >= i2 {
							err = errors.New("数据格式有误")
							return
						}
					case "<=":
						if i1 > i2 {
							err = errors.New("数据格式有误")
							return
						}
					case "<>":
						if i1 == i2 {
							err = errors.New("数据格式有误")
							return
						}
					}
				}
			} else if t == "phone" {
				if fieldV.Kind() == reflect.String {
					if IsPhone(fieldV.String()) == false {
						err = errors.New("数据格式有误")
						return
					}
				}
			} else if strings.HasPrefix(t, "varchar(") {
				s := strings.TrimPrefix(t, "varchar(")
				s = strings.TrimSuffix(s, ")")
				i, _ := Stoi(s)
				if i > 0 {
					s := fieldV.String()
					if len(s) > i {
						err = errors.New("数据格式有误")
						return
					}
				}
			}
		}
	}
	return
}

func IsPhone(str string) bool {
	isorno, _ := regexp.MatchString(`^(13[0-9]|14[5-9]|15[012356789]|166|17[0-8]|18[0-9]|19[0-9])[0-9]{8}`, str)
	return isorno
}
