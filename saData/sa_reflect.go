package saData

import (
	"errors"
	"reflect"
)

func StructName(obj interface{}) (string, error) {
	var structName string

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", errors.New("Check elem type error not Struct")
	}

	structName = t.Name()
	if structName == "" {
		return "", errors.New("获取结构体信息失败")
	}

	return structName, nil
}

func AryStructName(ary interface{}) (string, error) {
	var structName string

	t := reflect.TypeOf(ary)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "", errors.New("Check elem type error not Struct")
	}

	structName = t.Name()
	if structName == "" {
		return "", errors.New("获取结构体信息失败")
	}
	return structName, nil
}

func StructElemValue(obj interface{}, k string) reflect.Value {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		e := reflect.ValueOf(obj).Elem()
		return e.FieldByName(k)
	}
	return reflect.Value{}
}

func ValueToMap(v reflect.Value) *map[string]interface{} {
	dic := map[string]interface{}{}
	refectValue := reflect.ValueOf(v)
	s := v.Type().Elem()
	sLen := s.NumField()
	for i := 0; i < sLen; i++ {
		field := s.Field(i)
		//d := refectValue.FieldByName(field.Name)
		d := refectValue.Field(i)
		switch field.Type.Kind() {
		case reflect.Int16, reflect.Int, reflect.Int64, reflect.Int8, reflect.Int32:
			dic[field.Name] = d.Int()
		case reflect.Float32, reflect.Float64:
			dic[field.Name] = d.Float()
		case reflect.Bool:
			dic[field.Name] = d.Bool()
		case reflect.String:
			dic[field.Name] = d.String()
		case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:

		case reflect.Array:

		case reflect.Map:

		case reflect.Struct:
		default:
		}
	}

	return &dic
}
