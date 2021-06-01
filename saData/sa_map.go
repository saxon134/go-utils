package saData

import (
	"github.com/pkg/errors"
	"strings"
)

//注意：使用时嵌套的字典必须是SaMap，否则会导致数据无法被保存
type SaMap map[string]interface{}

//a:b:c = v
func (m SaMap) Set(key string, val interface{}) (err error) {
	if m == nil || key == "" {
		return errors.New("map为空")
	}

	ary := strings.Split(key, ":")
	if len(ary) > 0 {
		//保证链路数据都存在
		//var tmp interface{}
		//tmp = m
		//for i, k := range ary {
		//	if i+1 < len(ary) {
		//		if k[0:1] == "[" {
		//			k = strings.TrimPrefix(k, "[")
		//			k = strings.TrimSuffix(k, "]")
		//			idx, _ := ToInt(k)
		//			if idx >= 0 {
		//
		//			}
		//		}
		//	} else {
		//
		//	}
		//if _, ok := tmp[k]; ok == false {
		//	if i+1 == len(ary) {
		//		m[k] = val
		//	} else {
		//		m[k] = SaMap{}
		//	}
		//}
		//}
		//
		//for _, key = range ary {
		//
		//}
	}
	return nil
}

func (m SaMap) Get(key string) interface{} {
	return nil
}
