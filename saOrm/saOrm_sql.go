package saOrm

import (
	"github.com/pkg/errors"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saImg"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

func FromDb(obj interface{}) error {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()

	//反射，判断输入类型是否有误
	reflectType := reflect.TypeOf(obj)
	reflectValue := reflect.ValueOf(obj)
	{
		if reflectType.Kind() == reflect.Ptr {
			reflectType = reflectType.Elem()
			reflectValue = reflectValue.Elem()
		}

		if reflectType.Kind() != reflect.Struct {
			return errors.New("Error:类型有误，只能是Struct，或Struct指针")
		}
	}

	//通过反射获取结构体名称及元素名称
	fieldNum := reflectType.NumField()
	for i := 0; i < fieldNum; i++ {
		tags := reflectType.Field(i).Tag.Get("type")
		if tags == "" {
			tags = reflectType.Field(i).Tag.Get("gorm")
		}
		tagAry := strings.Split(tags, ";")

		var isOss = false
		columnKind := reflectType.Field(i).Type.Kind()
		columnName := reflectType.Field(i).Name
		for _, tag := range tagAry {
			tag = strings.ToLower(tag)
			if tag == "oss" || tag == "img" {
				isOss = true
				break
			}
		}

		if isOss == false {
			if columnName == "Img" || columnName == "Cover" || columnName == "Avatar" || columnName == "ImgAry" {
				isOss = true
			}
		}

		if isOss {
			if columnKind == reflect.String {
				str := reflectValue.Field(i).String()
				s := saImg.AddDefaultUriRoot(str)
				if s != str {
					reflectValue.Field(i).SetString(s)
				}
			} else if columnKind == reflect.Slice || columnKind == reflect.Array {
				v := reflectValue.Field(i)
				vLen := v.Len()
				for j := 0; j < vLen; j++ {
					str := v.Index(j).String()
					s := saImg.AddDefaultUriRoot(str)
					if s != str {
						v.Index(j).SetString(s)
					}
				}
			}
		}
	}
	return nil
}

func ToDB(obj interface{}) error {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()

	//反射，判断输入类型是否有误
	reflectType := reflect.TypeOf(obj)
	reflectValue := reflect.ValueOf(obj)
	{
		if reflectType.Kind() == reflect.Ptr {
			reflectType = reflectType.Elem()
			reflectValue = reflectValue.Elem()
		}

		if reflectType.Kind() != reflect.Struct {
			return errors.New("Error:类型有误，只能是Struct，或Struct指针")
		}
	}

	//通过反射获取结构体名称及元素名称
	fieldNum := reflectType.NumField()
	for i := 0; i < fieldNum; i++ {
		tags := reflectType.Field(i).Tag.Get("type")
		if tags == "" {
			tags = reflectType.Field(i).Tag.Get("orm")
		}
		if tags == "" {
			tags = reflectType.Field(i).Tag.Get("gorm")
		}
		tagAry := strings.Split(tags, ";")

		var in []string
		var isUpdated = false
		var isOss = false
		var isPhone = false
		var gte int64
		var gt int64
		var lte int64
		var lt int64
		var gtlt int64
		columnKind := reflectType.Field(i).Type.Kind()
		columnName := reflectType.Field(i).Name
		for _, tag := range tagAry {
			tag = strings.ToLower(tag)

			if strings.HasPrefix(tag, "varchar") {
				tag = strings.Replace(tag, "varchar", "", -1)
				tag = strings.Replace(tag, "(", "", -1)
				tag = strings.Replace(tag, ")", "", -1)
				lte, _ = saData.ToInt64(tag)
				if lte <= 0 {
					lte = 32
				}
			} else if strings.HasPrefix(tag, "char") {
				tag = strings.Replace(tag, "char", "", -1)
				tag = strings.Replace(tag, "(", "", -1)
				tag = strings.Replace(tag, ")", "", -1)
				lte, _ = saData.ToInt64(tag)
				if lte <= 0 {
					lte = 16
				}
			} else if strings.HasPrefix(tag, "in(") {
				if strings.Contains(tag, ":") {
					tag = strings.TrimPrefix(tag, "in(")
					tag = strings.TrimSuffix(tag, ")")
					in = strings.Split(tag, ",")
				}
			} else if tag == "updated" {
				isUpdated = true
			} else if tag == "oss" || tag == "img" {
				isOss = true
			} else if tag == "phone" {
				isPhone = true
			} else if strings.HasPrefix(tag, ">") {
				s := strings.TrimPrefix(tag, ">")
				gt, _ = saData.ToInt64(s)
			} else if strings.HasPrefix(tag, ">=") {
				s := strings.TrimPrefix(tag, ">=")
				gte, _ = saData.ToInt64(s)
			} else if strings.HasPrefix(tag, "<") {
				s := strings.TrimPrefix(tag, "<")
				lt, _ = saData.ToInt64(s)
			} else if strings.HasPrefix(tag, "<=") {
				s := strings.TrimPrefix(tag, "<=")
				lte, _ = saData.ToInt64(s)
			} else if strings.HasPrefix(tag, "<>") {
				s := strings.TrimPrefix(tag, "<>")
				gtlt, _ = saData.ToInt64(s)
			} else if strings.HasPrefix(tag, "enum(") {
				if ary := strings.Split(tag, ","); len(ary) > 0 {
					in := make([]string, 0, len(ary))
					for _, s := range ary {
						ary2 := strings.Split(s, ":")
						if len(ary2) > 0 {
							in = append(in, ary2[0])
						}
					}
				}
			}
		}

		switch columnKind {
		case reflect.Int, reflect.Int64, reflect.Uint64, reflect.Uint,
			reflect.Int8, reflect.Uint8,
			reflect.Int32, reflect.Uint32,
			reflect.Int16, reflect.Uint16:

			value := reflectValue.Field(i).Int()
			if lte > 0 {
				if value > lte {
					return errors.New(columnName + "超出长度")
				}
			}
			if lt > 0 {
				if value >= lte {
					return errors.New(columnName + "超出长度")
				}
			}
			if gte > 0 {
				if value < gte {
					return errors.New(columnName + "长度过短")
				}
			}
			if gt > 0 {
				if value <= gt {
					return errors.New(columnName + "长度过短")
				}
			}
			if gtlt > 0 {
				if value != gtlt {
					return errors.New(columnName + "长度不匹配")
				}
			}
			if len(in) > 0 {
				ok := false
				for _, v := range in {
					iv, _ := saData.ToInt64(v)
					if iv == value {
						ok = true
						break
					}
				}
				if ok == false {
					return errors.New(columnName + "有误")
				}
			}
		case reflect.String:
			if lte == 0 {
				if columnName == "Cover" || columnName == "Img" || columnName == "Avatar" {
					lte = 120
					isOss = true
				} else if columnName == "name" {
					lte = 60
				} else if columnName == "title" {
					lte = 250
				}
			}

			columnStr := reflectValue.Field(i).String()
			length := int64(saData.StrLen(columnStr))
			if lte > 0 {
				if length > lte {
					return errors.New(columnName + "超出长度")
				}
			}
			if lt > 0 {
				if length >= lte {
					return errors.New(columnName + "超出长度")
				}
			}
			if gte > 0 {
				if length < gte {
					return errors.New(columnName + "长度过短")
				}
			}
			if gt > 0 {
				if length <= gt {
					return errors.New(columnName + "长度过短")
				}
			}
			if gtlt > 0 {
				if length != gtlt {
					return errors.New(columnName + "长度不匹配")
				}
			}
			if isOss {
				s := saImg.DeleteUriRoot(columnStr)
				if s != columnStr {
					reflectValue.Field(i).SetString(s)
				}
			}
			if isPhone && saData.IsPhone(columnStr) == false {
				return errors.New("手机号格式有误")
			}
			if len(in) > 0 {
				ok := false
				for _, v := range in {
					if columnStr == v {
						ok = true
						break
					}
				}
				if ok == false {
					return errors.New(columnName + "有误")
				}
			}
		case reflect.Struct:
			if isUpdated {
				now := time.Now()
				reflectValue.Field(i).Set(reflect.ValueOf(now))
			}
		case reflect.Ptr:
			if isUpdated {
				now := time.Now()
				reflectValue.Field(i).Set(reflect.ValueOf(&now))
			}
		case reflect.Slice, reflect.Array:
			if isOss == false {
				if columnName == "imgAry" {
					isOss = true
				}
			}
			if isOss {
				v := reflectValue.Field(i)
				vLen := v.Len()
				v = v.Slice(0, vLen)
				for j := 0; j < vLen; j++ {
					str := v.String()
					s := saImg.DeleteUriRoot(str)
					if s != str {
						v.Field(j).SetString(s)
					}
				}
			}
		}
	}
	return nil
}

func (m *DB) Insert(obj interface{}) (tx *gorm.DB) {
	return m.Create(obj)
}
