package saGen

import (
	"fmt"
	"gitee.com/go-utils/saData"
	"gitee.com/go-utils/saHit"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

type Set struct {
	Obj              interface{}
	Options          string //设置：controller;ms;tbl;doc; 默认全部
	DB               string //数据库全局变量
	PK               string //主键名
	AddImgRootFun    string
	DeleteImgRootFun string
}

func GenerateTbl(set Set) {
	//默认值处理
	{
		set.Options = saHit.Str(set.Options == "", "controller;serve,;ms;doc;tbl", set.Options)
		set.DB = saHit.Str(set.DB == "", "common.DB", set.DB)
		set.PK = saHit.Str(set.PK == "", "Id", set.PK)
		set.AddImgRootFun = saHit.Str(set.AddImgRootFun == "", "yfImg.AddDefaultUriRoot", set.AddImgRootFun)
		set.DeleteImgRootFun = saHit.Str(set.DeleteImgRootFun == "", "yfImg.DeleteUriRoot", set.DeleteImgRootFun)
	}

	//反射，判断输入类型是否有误
	reflectType := reflect.TypeOf(set.Obj)
	reflectValue := reflect.ValueOf(set.Obj)
	{
		if reflectType.Kind() == reflect.Ptr {
			reflectType = reflectType.Elem()
		}

		if reflectType.Kind() != reflect.Struct {
			fmt.Println("Error:类型有误，只能是Struct，或Struct指针")
			return
		}
	}

	//通过反射获取结构体名称及元素名称
	pk := set.PK
	pkSnake := saData.SnakeStr(pk)
	var pkgName string
	var tblName string
	var structName string
	var columns []struct {
		name  string
		snake string
		tags  []string
	}
	hasFromDB := false
	hasToDB := false
	{
		//钩子函数
		{
			if _, ok := reflectType.MethodByName("FromDB"); ok {
				hasFromDB = true
			}
			if _, ok := reflectType.MethodByName("ToDB"); ok {
				hasToDB = true
			}
		}

		//获取结构体基本属性数据
		{
			structName = reflectType.Name()
			fieldNum := reflectType.NumField()
			for i := 0; i < fieldNum; i++ {
				fieldName := reflectType.Field(i).Name
				tag := reflectType.Field(i).Tag.Get("tbl")
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

			tblName = "t" + strings.TrimPrefix(saData.SnakeStr(structName), "tbl")
			if _, ok := reflectType.MethodByName("TableName"); ok {
				ary := reflectValue.Call([]reflect.Value{})
				if len(ary) > 0 {
					s := ary[0].String()
					if len(s) > 0 {
						tblName = s
					}
				}
			}

			pkgName = strings.Replace(reflectType.String(), "*", "", -1)
			pkgName = strings.Replace(reflectType.String(), "."+structName, "", -1)
		}
	}

	if len(columns) == 0 {
		fmt.Println("结构体空")
		return
	}

	// 生成数据库代码
	if strings.Contains(set.Options, "tbl") {
		createSqlTxt := ""
		indexSqlTxt := ""
		checkSqlTxt := ""
		fromDbSqlTxt := ""
		toDbSqlTxt := ""
		{
			structName = reflectType.Name()
			createSqlTxt = "CREATE TABLE IF NOT EXISTS `" + tblName + "` (\n"
			fieldNum := reflectType.NumField()
			for i := 0; i < fieldNum; i++ {
				columnType := ""
				columnDefault := ""
				columnComment := "''"
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
						indexSqlTxt += "\nALTER TABLE " + tblName + "ADD INDEX " + tagName + tag + ";"
					} else if strings.HasPrefix(tag, "varchar") {
						tag = strings.Replace(tag, "varchar", "", -1)
						tag = strings.Replace(tag, "(", "", -1)
						tag = strings.Replace(tag, ")", "", -1)
						length, _ := saData.ToInt(tag)
						if length <= 0 {
							length = 32
						}

						columnDefault = "''"
						columnType = "varchar(" + saData.Itos(length) + ")"
						checkSqlTxt += fmt.Sprintf(`
							if err:= saData.LenCheck(m.%s, %d); err != nil {
								return yfError.StackError(err, yfError.SensitiveErrorCode)
							}
						`, columns[i].name, length)
					} else if strings.HasPrefix(tag, "char") {
						tag = strings.Replace(tag, "char", "", -1)
						tag = strings.Replace(tag, "(", "", -1)
						tag = strings.Replace(tag, ")", "", -1)
						length, _ := saData.ToInt(tag)
						if length <= 0 {
							length = 16
						}

						columnDefault = "''"
						columnType = "char(" + saData.Itos(length) + ")"
						checkSqlTxt += fmt.Sprintf(`
							if err:= saData.LenCheck(m.%s, %d); err != nil {
								return yfError.StackError(err, yfError.SensitiveErrorCode)
							}
						`, columns[i].name, length)
					} else if strings.HasPrefix(tag, "int") {
						columnDefault = "0"
						tag = strings.TrimPrefix(tag, "int")
						if tag == "8" {
							columnType = "tinyint"
						} else if tag == "64" {
							columnType = "bigint unsigned"
						} else {
							columnType = "integer"
						}
					} else if strings.HasPrefix(tag, "in(") {
						if strings.Contains(tag, ":") {
							tag = strings.TrimPrefix(tag, "in(")
							tag = strings.TrimSuffix(tag, ")")
							columnComment = "'" + tag + "'"
						}
					} else if strings.HasPrefix(tag, "decimal(") {
						columnType = tag
					} else if tag == "required" {
						createSqlTxt += " NOT NULL"
					} else if strings.HasPrefix(tag, "comment") {
						columnComment = "'" + strings.TrimPrefix(tag, "comment:") + "'"
					} else if tag == "default" {
						tag = strings.TrimPrefix(tag, "default:")
						if columnKind == reflect.Bool {
							ok, _ := saData.ToBool(tag)
							columnDefault = saData.Itos(saHit.Int(ok, 1, 0))
						} else if columnKind == reflect.String {
							columnDefault = "'" + tag + "'"
						} else if columnKind == reflect.Float32 || columnKind == reflect.Float64 || columnKind == reflect.Complex64 || columnKind == reflect.Complex128 {
							if f, err := saData.ToFloat32(tag); err == nil {
								columnDefault = saData.F32tos(f)
							}
						}
					} else if tag == "created" || tag == "updated" {
						columnType = "datetime"
					} else if tag == "oss" {
						fromDbSqlTxt += fmt.Sprintf("\nm.%s = %s(m.%s)\n", columns[i].name, set.AddImgRootFun, columns[i].name)
						toDbSqlTxt += fmt.Sprintf("\nm.%s = %s(m.%s)\n", columns[i].name, set.DeleteImgRootFun, columns[i].name)
					} else if tag == "phone" {
						toDbSqlTxt += fmt.Sprintf(`
							if saData.IsPhone(m.%s) == false {
								return yfError.StackError("手机号格式有误", yfError.SensitiveErrorCode)
							}
							`, columns[i].name)
					} else if strings.HasPrefix(tag, ">") {
						var left = ""
						if columnKind <= reflect.Complex128 {
							left = "m." + columns[i].name
						} else if columnKind == reflect.String || columnKind == reflect.Slice || columnKind == reflect.Array || columnKind == reflect.Map {
							left = "len(m." + columns[i].name + ")"
						}
						if left != "" {
							tag = strings.TrimPrefix(tag, ">")
							toDbSqlTxt += fmt.Sprintf(`
								if %s <= %s {
									return  yfError.StackError(yfError.ErrorDate, yfError.SensitiveErrorCode)
								}
								`, left, tag)
						}
					} else if strings.HasPrefix(tag, ">=") {
						var left = ""
						if columnKind <= reflect.Complex128 {
							left = "m." + columns[i].name
						} else if columnKind == reflect.String || columnKind == reflect.Slice || columnKind == reflect.Array || columnKind == reflect.Map {
							left = "len(m." + columns[i].name + ")"
						}
						if left != "" {
							tag = strings.TrimPrefix(tag, ">=")
							toDbSqlTxt += fmt.Sprintf(`
								if m.%s < %s {
									return  yfError.StackError(yfError.ErrorDate, yfError.SensitiveErrorCode)
								}
								`, columns[i].name, tag)
						}
					} else if strings.HasPrefix(tag, "<") {
						var left = ""
						if columnKind <= reflect.Complex128 {
							left = "m." + columns[i].name
						} else if columnKind == reflect.String || columnKind == reflect.Slice || columnKind == reflect.Array || columnKind == reflect.Map {
							left = "len(m." + columns[i].name + ")"
						}
						if left != "" {
							tag = strings.TrimPrefix(tag, "<")
							toDbSqlTxt += fmt.Sprintf(`
								if m.%s >= %s {
									return  yfError.StackError(yfError.ErrorDate, yfError.SensitiveErrorCode)
								}
								`, columns[i].name, tag)
						}
					} else if strings.HasPrefix(tag, "<=") {
						var left = ""
						if columnKind <= reflect.Complex128 {
							left = "m." + columns[i].name
						} else if columnKind == reflect.String || columnKind == reflect.Slice || columnKind == reflect.Array || columnKind == reflect.Map {
							left = "len(m." + columns[i].name + ")"
						}
						if left != "" {
							tag = strings.TrimPrefix(tag, "<=")
							toDbSqlTxt += fmt.Sprintf(`
								if m.%s > %s {
									return  yfError.StackError(yfError.ErrorDate, yfError.SensitiveErrorCode)
								}
								`, columns[i].name, tag)
						}
					} else if strings.HasPrefix(tag, "<>") {
						var left = ""
						if columnKind <= reflect.Complex128 {
							left = "m." + columns[i].name
						} else if columnKind == reflect.String || columnKind == reflect.Slice || columnKind == reflect.Array || columnKind == reflect.Map {
							left = "len(m." + columns[i].name + ")"
						}
						if left != "" {
							tag = strings.TrimPrefix(tag, "<>")
							toDbSqlTxt += fmt.Sprintf(`
								if m.%s != %s {
									return  yfError.StackError(yfError.ErrorDate, yfError.SensitiveErrorCode)
								}
								`, columns[i].name, tag)
						}
					} else if tag == "time" {
						columnType = "datetime"
					}
				}

				if columnType == "" {
					if columnKind == reflect.Bool || columnKind == reflect.Int8 || columnKind == reflect.Uint8 {
						columnType = "tinyint"
					} else if columnKind == reflect.Int16 || columnKind == reflect.Uint16 || columnKind == reflect.Uint32 || columnKind == reflect.Int32 {
						columnType = "int"
					} else if columnKind == reflect.Uint64 || columnKind == reflect.Int64 {
						columnType = "bigint unsigned"
					} else if columnKind == reflect.Float32 || columnKind == reflect.Float64 || columnKind == reflect.Complex64 || columnKind == reflect.Complex128 {
						columnType = "decimal(10,2)"
					} else {
						if columns[i].snake == "name" {
							columnType = "varchar(64)"
						} else if columns[i].snake == "title" {
							columnType = "varchar(255)"
						} else if columns[i].snake == "cover" {
							columnType = "varchar(128)"
						} else {
							columnType = "varchar(64)"
						}
					}
				}

				createSqlTxt += "  `" + columns[i].snake + "` " + columnType
				if columnDefault != "" {
					createSqlTxt += " default " + columnDefault
				}
				if columnComment != "" {
					createSqlTxt += " comment " + columnComment
				}
				createSqlTxt += ",\n"
			}

			createSqlTxt += "  PRIMARY KEY (`" + pkSnake + "`) USING BTREE"
			createSqlTxt += "\n) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci;\n"
			fmt.Println("建表参考语句：")
			fmt.Println(createSqlTxt)
		}

		//生成tbl sql文件
		{
			//todo 换成网络地址
			tpl_f, err := os.OpenFile("/Users/jiang/yf.com/techio/server/yf-utils/yfGen/template/tbl_sql.tpl", os.O_RDONLY, 0600)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			b, _ := ioutil.ReadAll(tpl_f)
			_ = tpl_f.Close()
			tplStr := string(b)
			tplStr = strings.Replace(tplStr, "{{CreateSql}}", createSqlTxt, 1)
			tplStr = strings.Replace(tplStr, "{{Package}}", pkgName, 1)
			tplStr = strings.Replace(tplStr, "{{Model}}", structName, -1)
			tplStr = strings.Replace(tplStr, "{{TblName}}", tblName, -1)

			if hasFromDB {
				fromDbSqlTxt += `
				m.FromDB()
				`
			}
			if hasToDB {
				toDbSqlTxt += `
				if err = m.ToDB(); err != nil {
					return
				}
				`
			}

			tplStr = strings.Replace(tplStr, "{{FromDBSql}}", fromDbSqlTxt, -1)
			tplStr = strings.Replace(tplStr, "{{ToDBSql}}", toDbSqlTxt, -1)

			f_n := "./models/" + pkgName + "/" + strings.TrimPrefix(saData.SnakeStr(tblName), "t_") + "_sql.go"
			if err = createPath(f_n); err != nil {
				fmt.Println(err.Error())
				return
			}

			if err = ioutil.WriteFile(f_n, []byte(tplStr), 0644); err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}

	// 生成controller代码
	if strings.Contains(set.Options, "controller") {
		//todo 换成网络地址
		tpl_f, err := os.OpenFile("/Users/jiang/yf.com/techio/server/yf-utils/yfGen/template/controller.tpl", os.O_RDONLY, 0600)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		b, _ := ioutil.ReadAll(tpl_f)
		_ = tpl_f.Close()
		tplStr := string(b)

		tplStr = strings.Replace(tplStr, "{{Struct}}", strings.TrimPrefix(structName, "Tbl"), -1)
		tplStr = strings.Replace(tplStr, "{{StructLower}}", pkgName, -1)
		tplStr = strings.Replace(tplStr, "{{TblModel}}", structName, -1)

		f_n := "./http/controller/controller." + pkgName + "_gen.go"
		if err = createPath(f_n); err != nil {
			fmt.Println(err.Error())
			return
		}

		if err = ioutil.WriteFile(f_n, []byte(tplStr), 0644); err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	// 生成ms管理系统页面代码
	if strings.Contains(set.Options, "ms") {

	}

	// 生成doc接口文档
	if strings.Contains(set.Options, "doc") {

	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func createPath(file string) error {
	ary := strings.Split(file, "/")
	if strings.Contains(ary[len(ary)-1], ".") {
		ary = ary[:len(ary)-1]
	}
	if len(ary) > 0 {
		dir := strings.Join(ary, "/")
		if fileExists(dir) == false {
			err := os.MkdirAll(dir, 0700)
			return err
		}
	}

	return nil
}
