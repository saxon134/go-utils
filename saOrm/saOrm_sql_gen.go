package saOrm

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saError"
	"reflect"
)

// GenTblSql
// 生成修改数据库表结构语句
// 数据库不存在则生成建表语句，不存在则查看差异，生成差异代码；
// 建议如果返回值为true，则panic，因为数据库表结构跟model定义不一致；
// values项必须为指针
func GenTblSql(db *DB, values ...interface{}) {
	if db == nil || values == nil || len(values) == 0 {
		fmt.Println("Error:参数类型有误")
		return
	}

	for _, obj := range values {
		//反射，判断输入类型是否有误
		reflectType := reflect.TypeOf(obj)
		reflectValue := reflect.ValueOf(obj)
		{
			if reflectType.Kind() != reflect.Ptr {
				fmt.Println("Error:类型有误，只能是Struct指针")
			}

			if reflectType.Elem().Kind() != reflect.Struct {
				fmt.Println("Error:类型有误，只能是Struct指针")
			}
		}

		//通过反射获取数据库表的名称
		var tblName string
		{
			structName := reflectType.Name()
			tblName = saData.SnakeStr(structName)
			m := reflectValue.MethodByName("TableName")
			if m.IsValid() == false {
				m = reflectValue.Elem().MethodByName("TableName")
			}
			if m.IsValid() {
				v := m.Call(nil)
				s, _ := v[0].Interface().(string)
				if len(s) > 0 {
					tblName = s
				}
			}
		}

		//查询数据表是否已经存在
		colums := []string{}
		rows, err := db.Raw("SHOW TABLES LIKE '"+tblName+"'", &colums).Rows()
		if saError.DbErr(err) {
			fmt.Println("查询表是否存在失败：", tblName)
			fmt.Println(err)
			continue
		}
		db.ScanRows(rows, &colums)

		//创建表
		if len(colums) <= 2 {
			CreateTbl(obj)
		} else
		//修改表
		{
			AlterTbl(db, tblName, obj)
		}
	}
}
