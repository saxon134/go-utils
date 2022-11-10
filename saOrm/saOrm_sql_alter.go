package saOrm

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saHit"
	"reflect"
	"strings"
)

type Column struct {
	ColumnName    string `json:"column_name"`
	ColumnType    string `json:"column_type"`
	IsNullable    string `json:"is_nullable"`
	ColumnDefault string `json:"column_default_int"`
	ColumnComment string `json:"column_comment"`
}

// AlterTbl
// 修改数据库SQL
func AlterTbl(db *DB, tblName string, obj interface{}) {
	if db == nil || tblName == "" || obj == nil {
		return
	}

	//通过反射获取结构体名称及元素名称
	reflectType := reflect.TypeOf(obj).Elem()
	var columns []struct {
		name  string
		snake string
		tags  []string
	}
	fieldNum := reflectType.NumField()
	for i := 0; i < fieldNum; i++ {
		fieldName := reflectType.Field(i).Name
		tag := reflectType.Field(i).Tag.Get("orm")
		if tag == "" {
			tag = reflectType.Field(i).Tag.Get("type")
		}
		if tag == "" {
			tag = reflectType.Field(i).Tag.Get("gorm")
		}

		v := struct {
			name  string
			snake string
			tags  []string
		}{name: fieldName, snake: saData.SnakeStr(fieldName)}
		v.tags = strings.Split(tag, ";")
		columns = append(columns, v)
	}
	if len(columns) == 0 {
		fmt.Println("结构体空")
		return
	}

	//获取数据库表字段
	tblColumns := make([]*Column, 0, 20)
	//tblIndex := make([]string, 0, 5)
	rows, err := db.Raw(`
			SELECT
			  COLUMN_NAME column_name,
			  COLUMN_TYPE column_type,
			  IS_NULLABLE is_nullable,
			  COLUMN_DEFAULT column_default_int,
			  COLUMN_DEFAULT column_default_str,
			  COLUMN_COMMENT column_comment 
			FROM
			 INFORMATION_SCHEMA.COLUMNS
			where
			table_name  = '` + tblName + "'").Rows()
	if saError.DbErr(err) {
		return
	}
	db.ScanRows(rows, &tblColumns)
	for _, v := range tblColumns {
		if strings.Contains(v.ColumnType, "(") &&
			strings.Contains(v.ColumnType, "varchar") == false &&
			strings.Contains(v.ColumnType, "char") == false {
			v.ColumnType = strings.ReplaceAll(v.ColumnType, v.ColumnType[strings.Index(v.ColumnType, "("):strings.Index(v.ColumnType, ")")+1], "")
		}
		v.ColumnDefault, _ = saData.ToStr(v.ColumnDefault)
		v.IsNullable = strings.ToLower(v.IsNullable)
		v.IsNullable = saHit.Str(v.IsNullable == "yes", "true", "false")
	}

	//获取模型字段
	modelColumns := make([]*Column, 0, 20)
	modelIndex := make([]string, 0, 5)
	{
		fieldNum := reflectType.NumField()
		for i := 0; i < fieldNum; i++ {
			filed := Column{}
			filed.ColumnComment = ""
			filed.IsNullable = "true"
			if columns[i].snake == "base_model" {
				modelColumns = append(modelColumns, &Column{
					ColumnName:    "id",
					ColumnType:    "bigint unsigned not null auto_increment",
					IsNullable:    "false",
					ColumnDefault: "",
					ColumnComment: "",
				})
				modelColumns = append(modelColumns, &Column{
					ColumnName:    "created_at",
					ColumnType:    "datetime",
					IsNullable:    "true",
					ColumnDefault: "current_timestamp",
					ColumnComment: "",
				})
				continue
			} else if columns[i].snake == "base_model_with_delete" {
				modelColumns = append(modelColumns, &Column{
					ColumnName:    "id",
					ColumnType:    "bigint unsigned not null auto_increment",
					IsNullable:    "false",
					ColumnDefault: "",
					ColumnComment: "",
				})
				modelColumns = append(modelColumns, &Column{
					ColumnName:    "created_at",
					ColumnType:    "datetime",
					IsNullable:    "true",
					ColumnDefault: "current_timestamp",
					ColumnComment: "",
				})
				modelColumns = append(modelColumns, &Column{
					ColumnName:    "deleted_at",
					ColumnType:    "datetime",
					IsNullable:    "true",
					ColumnDefault: "default null",
					ColumnComment: "",
				})
				continue
			} else {
				filed.ColumnName = columns[i].snake
			}

			columnSigned := false
			columnKind := reflectType.Field(i).Type.Kind()
			for _, tag := range columns[i].tags {
				tag = strings.ToLower(tag)

				//索引
				if strings.HasPrefix(tag, "index") {
					tag = strings.Replace(tag, "index", "", 1)
					if tag == "" {
						tag = "(" + columns[i].snake + ")"
					}

					tagName := "IDX_" + strings.Replace(tag, "(", "", 1)
					tagName = strings.Replace(tag, ")", "", 1)
					tagName = strings.Replace(tag, ",", "_", 1)
					modelIndex = append(modelIndex, tagName+tag)
				} else if strings.HasPrefix(tag, "varchar") {
					tag = strings.Replace(tag, "varchar", "", -1)
					tag = strings.Replace(tag, "(", "", -1)
					tag = strings.Replace(tag, ")", "", -1)
					length, _ := saData.ToInt(tag)
					if length <= 0 {
						length = 32
					}
					filed.ColumnDefault = "''"
					filed.ColumnType = "varchar(" + saData.Itos(length) + ")"
				} else if strings.HasPrefix(tag, "char") {
					tag = strings.Replace(tag, "char", "", -1)
					tag = strings.Replace(tag, "(", "", -1)
					tag = strings.Replace(tag, ")", "", -1)
					length, _ := saData.ToInt(tag)
					if length <= 0 {
						length = 16
					}

					filed.ColumnDefault = "''"
					filed.ColumnType = "char(" + saData.Itos(length) + ")"
				} else if strings.HasPrefix(tag, "int") {
					filed.ColumnDefault = "0"
					tag = strings.TrimPrefix(tag, "int")
					if strings.HasPrefix(tag, "8") {
						filed.ColumnType = "tinyint"
						if columns[i].snake == "status" || columns[i].snake == "type" {
							filed.ColumnType = "tinyint"
							filed.ColumnDefault = saHit.Str(filed.ColumnDefault == "", "-1", filed.ColumnDefault)
						}
					} else if strings.HasPrefix(tag, "64") {
						filed.ColumnType = "bigint"
					} else {
						filed.ColumnType = "integer"
					}

					if strings.Contains(tag, "unsigned") {
						filed.ColumnType += " unsigned"
					}
				} else if tag == "float" || tag == "double" {
					filed.ColumnDefault = "0"
					filed.ColumnType = tag
				} else if strings.HasPrefix(tag, "tinyint") {
					filed.ColumnDefault = "0"
					filed.ColumnType = "tinyint"
					if strings.Contains(tag, "unsigned") {
						filed.ColumnType += " unsigned"
					}
				} else if strings.HasPrefix(tag, "in(") {
					if strings.Contains(tag, ":") {
						tag = strings.TrimPrefix(tag, "in(")
						tag = strings.TrimSuffix(tag, ")")
						filed.ColumnComment = tag
					}
				} else if strings.HasPrefix(tag, "decimal(") {
					filed.ColumnType = tag
				} else if tag == "signed" {
					columnSigned = true
				} else if tag == "required" || tag == "not null" {
					filed.IsNullable = "false"

				} else if strings.HasPrefix(tag, "comment") {
					filed.ColumnComment = strings.TrimPrefix(tag, "comment:")
				} else if strings.HasPrefix(tag, "default:") {
					tag = strings.TrimPrefix(tag, "default:")
					if columnKind == reflect.Bool {
						ok, _ := saData.ToBool(tag)
						filed.ColumnDefault = saData.Itos(saHit.Int(ok, 1, 0))
					} else if columnKind == reflect.String {
						filed.ColumnDefault = "'" + tag + "'"
					} else if columnKind >= reflect.Int && columnKind <= reflect.Uint64 {
						if i, err := saData.ToInt64(tag); err == nil {
							filed.ColumnDefault = saData.I64tos(i)
						}
					} else if columnKind >= reflect.Float32 && columnKind <= reflect.Complex128 {
						if f, err := saData.ToFloat32(tag); err == nil {
							filed.ColumnDefault = saData.F32tos(f)
						}
					}
				} else if tag == "created" {
					filed.ColumnType = "datetime"
					filed.ColumnDefault = "CURRENT_TIMESTAMP"
				} else if tag == "updated" {
					filed.ColumnType = "datetime"
					filed.ColumnDefault = "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"
				} else if tag == "datetime" {
					filed.ColumnType = "datetime"
				} else if tag == "oss" || tag == "img" {

				} else if tag == "phone" {

				} else if strings.HasPrefix(tag, "enum(") {
					filed.ColumnComment = strings.TrimPrefix(tag, "enum(")
					filed.ColumnComment = strings.TrimSuffix(filed.ColumnComment, ")")
				} else if tag == "time" {
					filed.ColumnType = "datetime"
				}
			}

			if filed.ColumnType == "" {
				if columnKind == reflect.Bool || columnKind == reflect.Int8 || columnKind == reflect.Uint8 {
					filed.ColumnType = saHit.Str(columnSigned, "tinyint", "tinyint unsigned")
				} else if columnKind == reflect.Int || columnKind == reflect.Int16 || columnKind == reflect.Uint16 || columnKind == reflect.Uint32 || columnKind == reflect.Int32 {
					filed.ColumnType = saHit.Str(columnSigned, "int", "int unsigned")
				} else if columnKind == reflect.Uint64 || columnKind == reflect.Int64 {
					filed.ColumnType = saHit.Str(columnSigned, "bigint", "bigint unsigned")
				} else if columnKind == reflect.Float32 || columnKind == reflect.Float64 || columnKind == reflect.Complex64 || columnKind == reflect.Complex128 {
					filed.ColumnType = "decimal(10,2)"
				} else {
					if columns[i].snake == "name" {
						filed.ColumnType = "varchar(60)"
					} else if columns[i].snake == "title" {
						filed.ColumnType = "varchar(250)"
					} else if columns[i].snake == "cover" {
						filed.ColumnType = "varchar(120)"
					} else if columns[i].snake == "img" {
						filed.ColumnType = "varchar(120)"
					} else {
						switch reflect.TypeOf(filed.ColumnType).Kind() {
						case reflect.Int, reflect.Int32:
							filed.ColumnType = "int"
							filed.ColumnType = saHit.Str(columnSigned, "int", "int unsigned")
						case reflect.Int64:
							filed.ColumnType = "bigint"
							filed.ColumnType = saHit.Str(columnSigned, "bigint", "bigint unsigned")
						case reflect.Int8:
							filed.ColumnType = "tinyint"
							filed.ColumnType = saHit.Str(columnSigned, "tinyint", "tinyint unsigned")
						}

						if filed.ColumnType == "" {
							filed.ColumnType = "varchar(64)"
						}
					}
				}
			} else if columnSigned {
				filed.ColumnType = strings.Replace(filed.ColumnType, " unsigned", "", -1)
			}

			modelColumns = append(modelColumns, &filed)
		}

		//比对数据库结构
		alterSql := ""
		for _, m := range modelColumns {
			var existed *Column
			for _, tbl := range tblColumns {
				if tbl.ColumnName == m.ColumnName {
					existed = tbl
					break
				}
			}

			if existed == nil || existed.ColumnName == "" {
				alterSql += fmt.Sprintf("alter table %s add column `%s` %s", tblName, m.ColumnName, m.ColumnType)
				if m.ColumnDefault != "" {
					alterSql += " default " + m.ColumnDefault
				}
				if m.ColumnComment != "" {
					alterSql += " comment '" + m.ColumnComment + "'"
				}
				alterSql += ", ALGORITHM=INPLACE, LOCK=NONE;\n"
			} else if m.ColumnName != "id" && m.ColumnName != "created_at" && m.ColumnName != "updated_at" && m.ColumnName != "deleted_at" {
				sql := ""
				m.ColumnDefault = saHit.Str(m.ColumnDefault == "''", "", m.ColumnDefault)
				if m.ColumnType != existed.ColumnType ||
					m.ColumnComment != existed.ColumnComment ||
					m.IsNullable != existed.IsNullable {

					sql = fmt.Sprintf("alter table %s modify column `%s` %s", tblName, m.ColumnName, m.ColumnType)
					if m.IsNullable == "false" {
						if strings.Contains(sql, "not null") == false {
							sql += " not null"
						}
					}
					if m.ColumnDefault != "" {
						sql += " default " + m.ColumnDefault
					}
					if m.ColumnComment != "" {
						sql += " comment '" + m.ColumnComment + "'"
					}
					sql += ";\n"
					alterSql += sql
				}
			}
		}
		if alterSql != "" {
			fmt.Println("-- 表结构变更：" + tblName)
			fmt.Println(alterSql)
		}

		//查找已删除字段
		var delColumns = ""
		for _, tbl := range tblColumns {
			existed := false
			for _, v := range modelColumns {
				if v.ColumnName == tbl.ColumnName {
					existed = true
					break
				}
			}
			if existed == false && tbl.ColumnName != "" {
				delColumns += tbl.ColumnName + ","
			}
		}
		if len(delColumns) > 0 {
			fmt.Println("-- 表结构已删除字段（注意：也有可能是修改了字段名称，请自行判断）")
			fmt.Println("-- ", tblName, ":", delColumns, "\n")
		}
	}
}
