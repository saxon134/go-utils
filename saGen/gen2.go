package saGen

import (
	"errors"
	"fmt"
	"gitee.com/go-utils/saData"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func generate2(obj interface{}, dbPtr, pkName string, addImgRootFun string, deleteImgRootFun string) error {
	reflectType := reflect.TypeOf(obj)
	refectValue := reflect.ValueOf(obj)

	hasFromDB := false
	hasToDB := false

	if _, ok := reflectType.MethodByName("FromDB"); ok {
		hasFromDB = true
	}
	if _, ok := reflectType.MethodByName("ToDB"); ok {
		hasToDB = true
	}

	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	if reflectType.Kind() != reflect.Struct {
		return errors.New("Check type error not Struct")
	}

	if _, ok := reflectType.MethodByName("FromDB"); ok {
		hasFromDB = true
	}
	if _, ok := reflectType.MethodByName("ToDB"); ok {
		hasToDB = true
	}

	if pkName == "" {
		pkName = "Id"
	}
	pkNameSnake := saData.SnakeStr(pkName)

	//通过反射获取结构体名称及元素名称
	var pkgName string
	var tblName string
	var structName string
	var columns string
	var modelDestElements = ""
	var modelDestStrDef = ""
	var modelDestStrFormate = ""
	{
		structName = reflectType.Name()
		fieldNum := reflectType.NumField()
		for i := 0; i < fieldNum; i++ {
			fieldName := reflectType.Field(i).Name

			tags := reflectType.Field(i).Tag.Get("tbl")
			isTime := false
			if strings.Contains(tags, "time") ||
				strings.Contains(tags, "created") ||
				strings.Contains(tags, "updated") {
				isTime = true
			}

			if isTime {
				modelDestElements += "&" + fieldName + ", "
				modelDestStrDef += "var " + fieldName + " string\n"
				modelDestStrFormate += `m.` + fieldName + ` = saData.StrToTime(saData.TimeFormat_sys_default, ` + fieldName + ")\n"
			} else {
				switch reflectType.Field(i).Type.Kind() {
				case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
					modelDestElements += "&" + fieldName + ", "
					modelDestStrDef += fieldName + " := \"\"\n"
					modelDestStrFormate += `_ = saData.JsonToStruct(` + fieldName + ", &m." + fieldName + ")\n"
				default:
					modelDestElements += "&(m." + fieldName + "),"
				}
			}

			columns += saData.SnakeStr(fieldName)
			if i+1 < fieldNum {
				columns += ","
			}
		}

		if structName == "" || len(columns) <= 0 {
			return errors.New("获取结构体信息失败")
		}

		tblName = saData.SnakeStr(structName)
		pkgName = strings.Replace(reflectType.String(), "*", "", -1)
		pkgName = strings.Replace(reflectType.String(), "."+structName, "", -1)
	}

	//生成基础函数
	tplContent := ""
	{
		read_f, err := os.OpenFile("tbl_sql.go.tpl", os.O_RDONLY, 0600)
		if err != nil {
			return err
		}

		b, _ := ioutil.ReadAll(read_f)
		_ = read_f.Close()

		tplContent = string(b)
		tplContent = strings.Replace(tplContent, "{{Package}}", pkgName, -1)
		tplContent = strings.Replace(tplContent, "{{Model}}", structName, -1)
		tplContent = strings.Replace(tplContent, "{{TblName}}", saData.SnakeStr(structName), -1)
		tplContent = strings.Replace(tplContent, "{{Columns}}", columns, -1)
		tplContent = strings.Replace(tplContent, "{{SQL}}", dbPtr, -1)
		tplContent = strings.Replace(tplContent, "{{Id}}", pkName, -1)
		tplContent = strings.Replace(tplContent, "{{id}}", pkNameSnake, -1)
		tplContent = strings.Replace(tplContent, "{{ModelDestElements}}", modelDestElements, -1)
		tplContent = strings.Replace(tplContent, "{{ModelDestStrDef}}", modelDestStrDef, -1)
		tplContent = strings.Replace(tplContent, "{{ModelDestStrFormate}}", modelDestStrFormate, -1)
	}

	//from DB函数生成
	fromDBImgStr := ""
	fromDBTimeStr := ""
	{
		for i := 0; i < reflectType.NumField(); i++ {
			field := reflectType.Field(i)
			tblStr := field.Tag.Get("tbl")

			isImg := false
			if strings.Contains(tblStr, "img") {
				isImg = true
			}

			isTime := false
			if strings.Contains(tblStr, "time") {
				isImg = true
			}

			field_k := field.Type.Kind()
			field_n := field.Name
			if isImg && addImgRootFun != "" {
				if field_k == reflect.String {
					s := "m.FieldName = " + addImgRootFun + "(m.FieldName)\n"
					s = strings.Replace(s, "FieldName", field_n, -1)
					fromDBImgStr += s
				} else if field_k == reflect.Slice {
					v := refectValue.Elem().FieldByName(field_n)
					if _, ok := v.Interface().([]string); ok {
						s := "for i, v := range m.FieldName {\nm.FieldName[i] = " + addImgRootFun + "(v)\n}\n"
						s = strings.Replace(s, "FieldName", field_n, -1)
						fromDBImgStr += s
					}
					if _, ok := v.Interface().(*[]string); ok {
						s := "for i, v := range m.FieldName {\nm.FieldName[i] = " + addImgRootFun + "(v)\n}\n"
						s = strings.Replace(s, "FieldName", field_n, -1)
						fromDBImgStr += s
					}
				} else if field_k == reflect.Array {
					v := refectValue.Elem().FieldByName(field_n)
					if _, ok := v.Interface().([]string); ok {
						s := "for i, v := range m.FieldName {\nm.FieldName[i] = " + addImgRootFun + "(v)\n}\n"
						s = strings.Replace(s, "FieldName", field_n, -1)
						fromDBImgStr += s
					}
					if _, ok := v.Interface().(*[]string); ok {
						s := "for i, v := range m.FieldName {\nm.FieldName[i] = " + addImgRootFun + "(v)\n}\n"
						s = strings.Replace(s, "FieldName", field_n, -1)
						fromDBImgStr += s
					}
				}
			}

			if isTime {
				fromDBTimeStr = ""
			}
		}
	}
	tplContent = strings.Replace(tplContent, "{{FromDB_Img}}", fromDBImgStr, -1)
	tplContent = strings.Replace(tplContent, "{{FromDB_Time}}", fromDBTimeStr, -1)
	if hasFromDB {
		tplContent = strings.Replace(tplContent, "{{FromDB_Model}}", "	m.FromDB()", -1)
	} else {
		tplContent = strings.Replace(tplContent, "{{FromDB_Model}}", "", -1)
	}

	//to DB函数生成
	toDB_LenCheckStr := ""
	toDB_BtCheckStr := ""
	toDB_ImgStr := ""
	{
		for i := 0; i < reflectType.NumField(); i++ {
			field := reflectType.Field(i)
			field_name := field.Name

			tagStr := field.Tag.Get("tbl")
			tagAry := strings.Split(tagStr, " ")

			//数据最大长度校验，非结构性字段，检查序列化后的字符长度
			for _, v := range tagAry {
				if strings.HasPrefix(v, "len:") {
					maxLen := 0
					s := saData.SubStr(v, 4, saData.StrLen(v)-4)
					maxLen, _ = saData.Stoi(s)

					if maxLen > 0 {
						s := `
						
							if err = saData.LenCheck(m.FieldName, LEN); err != nil {
								err = saError.NewError(err.Error())
								return err
							}
						`
						s = strings.Replace(s, "FieldName", field_name, -1)
						s = strings.Replace(s, "LEN", saData.Itos(maxLen), -1)
						toDB_LenCheckStr += s
					}
					break
				}
			}

			//数据有效范围校验
			for _, v := range tagAry {
				if strings.HasPrefix(v, "bt:") {
					v = strings.TrimPrefix(v, "bt:")
					k := field.Type.Kind()
					if k >= reflect.Int && k <= reflect.Float64 {
						max := -1
						maxCompare := ""
						min := -1
						minCompare := ""
						ary2 := strings.Split(v, ":")
						if len(ary2) == 2 {
							if ary2[0] == "(" {
								minCompare = "<="
							} else if ary2[0] == "[" {
								minCompare = "<"
							}

							s := strings.Replace(ary2[0], "(", "", -1)
							s = strings.Replace(ary2[0], "[", "", -1)
							if i, err := saData.Stoi(s); err == nil {
								min = i
							}

							if ary2[1] == ")" {
								maxCompare = ">="
							} else if ary2[1] == "]" {
								maxCompare = ">"
							}
							s = strings.Replace(ary2[1], ")", "", -1)
							s = strings.Replace(ary2[1], "]", "", -1)
							if i, err := saData.Stoi(s); err == nil {
								max = i
							}
						}

						if max > 0 && min > 0 {
							if max == min {
								if maxCompare == ">=" && minCompare == "<=" {
									s := `
									
									if m.FieldName != Max {
										err = errors.New("FieldName越界" + err.Error())
									}
									`
									s = strings.Replace(s, "FieldName", field_name, -1)
									s = strings.Replace(s, "Max", saData.Itos(max), -1)
									toDB_BtCheckStr += s
								}
							} else if max > min {
								s := `
									
									if m.FieldName MinCompare Min || m.FieldName MaxCompare Max {
										err = errors.New("FieldName越界" + err.Error())
									}
									`
								s = strings.Replace(s, "FieldName", field_name, -1)
								s = strings.Replace(s, "MaxCompare", maxCompare, -1)
								s = strings.Replace(s, "MinCompare", minCompare, -1)
								s = strings.Replace(s, "Max", saData.Itos(max), -1)
								s = strings.Replace(s, "Min", saData.Itos(min), -1)
								toDB_BtCheckStr += s
							}
						} else if max > 0 {
							s := `
									
									if m.FieldName MaxCompare Max {
										err = errors.New("FieldName越界" + err.Error())
									}
									`
							s = strings.Replace(s, "FieldName", field_name, -1)
							s = strings.Replace(s, "MaxCompare", maxCompare, -1)
							s = strings.Replace(s, "Max", saData.Itos(max), -1)
							toDB_BtCheckStr += s
						} else if min >= 0 {
							s := `
									
									if m.FieldName MinCompare Min {
										err = errors.New("FieldName越界" + err.Error())
									}
									`
							s = strings.Replace(s, "FieldName", field_name, -1)
							s = strings.Replace(s, "MinCompare", minCompare, -1)
							s = strings.Replace(s, "Min", saData.Itos(min), -1)
							toDB_BtCheckStr += s
						}
					}
				}
			}

			//图片处理
			if deleteImgRootFun != "" {
				for _, v := range tagAry {
					if v == "img" {
						k := field.Type.Kind()
						if k == reflect.String {
							s := `m.FieldName = ` + deleteImgRootFun + "(m.FieldName)"
							s = strings.Replace(s, "FieldName", field_name, -1)
							toDB_ImgStr += s
						} else if k == reflect.Array || k == reflect.Slice {
							s := `
						
							for i, v := range m.FieldName {
								m.FieldName[i] = DeleteUriRootFun(v)
							}
						`
							s = strings.Replace(s, "FieldName", field_name, -1)
							s = strings.Replace(s, "DeleteUriRootFun", deleteImgRootFun, -1)
							toDB_ImgStr += s
						}
					}
				}
			}
		}
	}
	tplContent = strings.Replace(tplContent, "{{ToDB_Img}}", toDB_ImgStr, -1)
	tplContent = strings.Replace(tplContent, "{{ToDB_LenCheck}}", toDB_LenCheckStr, -1)
	tplContent = strings.Replace(tplContent, "{{ToDB_BtCheck}}", toDB_BtCheckStr, -1)
	if hasToDB {
		tplContent = strings.Replace(tplContent, "{{ToDB_Model}}", "	if err = m.ToDB(); err != nil {\n		return err\n	}", -1)
	} else {
		tplContent = strings.Replace(tplContent, "{{ToDB_Model}}", "", -1)
	}
	if len(toDB_ImgStr) > 0 || len(toDB_LenCheckStr) > 0 || len(toDB_BtCheckStr) > 0 || hasToDB {
		tplContent = strings.Replace(tplContent, "{{ToDB_Err}}", "	var err error", -1)
	} else {
		tplContent = strings.Replace(tplContent, "{{ToDB_Err}}", "", -1)
	}

	//insert row sql函数生成
	InsertRowColumnStr := "`" + pkNameSnake + "`, "
	InsertRowValuesStr := `sqlTxt += saData.I64tos(m.` + pkName + `) + ", "` + "\n"
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		field_name := field.Name

		if field_name != pkName {
			InsertRowColumnStr += "`" + saData.SnakeStr(field_name) + "`, "

			tagStr := field.Tag.Get("tbl")
			tagAry := strings.Split(tagStr, " ")

			k := field.Type.Kind()
			if k == reflect.Int64 {
				InsertRowValuesStr += `sqlTxt += saData.I64tos(m.FieldName) + ", "`
				InsertRowValuesStr += "\n"
			} else if k == reflect.Int {
				InsertRowValuesStr += `sqlTxt += saData.Itos(int(m.FieldName)) + ", "`
				InsertRowValuesStr += "\n"
			} else if k == reflect.Int8 {
				InsertRowValuesStr += `sqlTxt += saData.Itos(int(m.FieldName)) + ", "`
				InsertRowValuesStr += "\n"
			} else if k == reflect.String {
				InsertRowValuesStr += `sqlTxt += "'" + m.FieldName + "', "`
				InsertRowValuesStr += "\n"
			} else if k == reflect.Uint64 {
				InsertRowValuesStr += `sqlTxt += saData.I64tos(int64(m.FieldName)) + ", "`
				InsertRowValuesStr += "\n"
			} else if k >= reflect.Bool && k <= reflect.Uint32 {
				InsertRowValuesStr += `sqlTxt += saData.Itos(int(m.FieldName)) + ", "`
				InsertRowValuesStr += "\n"
			} else if k == reflect.Float32 {
				InsertRowValuesStr += `sqlTxt += saData.F32tos(m.FieldName) + ", "`
				InsertRowValuesStr += "\n"
			} else if k == reflect.Float64 {
				InsertRowValuesStr += `sqlTxt += saData.F32tos(float32(m.FieldName)) + ", "`
				InsertRowValuesStr += "\n"
			} else {
				isTime := false
				isCreated := false
				isUpdated := false
				for _, v := range tagAry {
					if v == "time" {
						isTime = true
					} else if v == "created" {
						isCreated = true
					} else if v == "updated" {
						isUpdated = true
					}
				}

				if isCreated || isUpdated {
					InsertRowValuesStr += "\nm.FieldName = time.Now()\n"
				}

				if isTime || isCreated || isUpdated {
					InsertRowValuesStr += `sqlTxt += "'" + saData.TimeStr(m.FieldName, saData.TimeFormat_Default) + "', "`
					InsertRowValuesStr += "\n"
				} else {
					s := `
						if s, err := saData.DataToJson(m.FieldName); err == nil {
							sqlTxt += "'" + s + "'" + ", "
						}
					`
					InsertRowValuesStr += s
				}
			}

			InsertRowValuesStr = strings.Replace(InsertRowValuesStr, "FieldName", field_name, -1)
			InsertRowColumnStr = strings.Replace(InsertRowColumnStr, "FieldName", field_name, -1)
		}
	}
	InsertRowColumnStr = strings.TrimSuffix(InsertRowColumnStr, ", ")
	InsertRowValuesStr = strings.TrimSuffix(InsertRowValuesStr, ", ")
	tplContent = strings.Replace(tplContent, "{{InsertRowValues}}", InsertRowValuesStr, -1)
	tplContent = strings.Replace(tplContent, "{{InsertRowColumns}}", InsertRowColumnStr, -1)

	//update row sql函数生成
	UpdateRowSqlStr := ""
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		field_name := field.Name

		if field_name != pkName {

			UpdateRowSqlStr += "sqlTxt += \"`" + saData.SnakeStr(field_name) + "` = "

			tagStr := field.Tag.Get("tbl")
			tagAry := strings.Split(tagStr, " ")

			k := field.Type.Kind()
			if k == reflect.Int64 {
				UpdateRowSqlStr += `" + saData.I64tos(m.FieldName) + ","`
				UpdateRowSqlStr += "\n"
			} else if k == reflect.Int {
				UpdateRowSqlStr += `" + saData.Itos(int(m.FieldName)) + ","`
				UpdateRowSqlStr += "\n"
			} else if k == reflect.Int8 {
				UpdateRowSqlStr += `" + saData.Itos(int(m.FieldName)) + ","`
				UpdateRowSqlStr += "\n"
			} else if k == reflect.String {
				UpdateRowSqlStr += `" + "'" + m.FieldName + "',"`
				UpdateRowSqlStr += "\n"
			} else if k == reflect.Uint64 {
				UpdateRowSqlStr += `" + saData.Itos(int64(m.FieldName)) + ","`
				UpdateRowSqlStr += "\n"
			} else if k >= reflect.Bool && k <= reflect.Uint32 {
				UpdateRowSqlStr += `" + saData.Itos(int(m.FieldName)) + ","`
				UpdateRowSqlStr += "\n"
			} else if k == reflect.Float32 {
				UpdateRowSqlStr += `" + saData.F32tos(m.FieldName) + ","`
				UpdateRowSqlStr += "\n"
			} else if k == reflect.Float64 {
				UpdateRowSqlStr += `" + saData.F32tos(float32(m.FieldName)) + ","`
				UpdateRowSqlStr += "\n"
			} else {
				UpdateRowSqlStr += "'\""

				isTime := false
				isUpdated := false
				isCreated := false
				for _, v := range tagAry {
					if v == "time" {
						isTime = true
					} else if v == "updated" {
						isUpdated = true
					} else if v == "created" {
						isCreated = true
					}
				}

				if isUpdated {
					UpdateRowSqlStr += "\nm.FieldName = time.Now()\n"
				}

				if isTime || isUpdated || isCreated {
					UpdateRowSqlStr += "\n" + `sqlTxt += saData.TimeStr(m.FieldName, saData.TimeFormat_Default) + "',"` + "\n"
				} else {
					s := `
						if s, err := saData.DataToJson(m.FieldName); err == nil {
							sqlTxt += s + "',"
						}					
					`
					UpdateRowSqlStr += s
				}
			}

			UpdateRowSqlStr = strings.Replace(UpdateRowSqlStr, "FieldName", field_name, -1)
		}
	}

	UpdateRowSqlStr = strings.TrimSuffix(UpdateRowSqlStr, ",")
	tplContent = strings.Replace(tplContent, "{{UpdateRowSql}}", UpdateRowSqlStr, -1)

	//生成数据库建表基本代码，需要自己去完善
	createSqlTxt := ""
	{
		structName = reflectType.Name()
		createSqlTxt = "CREATE TABLE IF NOT EXISTS `" + saData.SnakeStr(structName) + "` (\n"
		idxSqlTxt := ""
		fieldNum := reflectType.NumField()
		for i := 0; i < fieldNum; i++ {
			fieldName := reflectType.Field(i).Name
			tags := reflectType.Field(i).Tag.Get("tbl")
			tagAry := strings.Split(tags, " ")

			if fieldName == pkName {
				createSqlTxt += "  `" + pkNameSnake + "` bigint UNSIGNED AUTO_INCREMENT COMMENT '主键',\n"
			} else {
				if reflectType.Field(i).Type.String() == "time.Time" {
					createSqlTxt += "  `" + saData.SnakeStr(fieldName) + "` datetime DEFAULT NULL COMMENT '',\n"
				} else {
					switch reflectType.Field(i).Type.Kind() {
					case reflect.Bool, reflect.Int8, reflect.Uint8:
						createSqlTxt += "  `" + saData.SnakeStr(fieldName) + "` tinyint(1) DEFAULT 0,\n"
					case reflect.Int, reflect.Int64, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						createSqlTxt += "  `" + saData.SnakeStr(fieldName) + "` integer UNSIGNED DEFAULT 0 COMMENT '',\n"
					case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
						createSqlTxt += "  `" + saData.SnakeStr(fieldName) + "` decimal(20, 2) DEFAULT 0.00 COMMENT '',\n"
					case reflect.String, reflect.Array, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct:
						maxLen := 0
						for _, v := range tagAry {
							if strings.HasPrefix(v, "len:") {
								s := saData.SubStr(v, 4, saData.StrLen(v)-4)
								maxLen, _ = saData.Stoi(s)
								break
							}
						}

						if maxLen <= 0 {
							maxLen = 32
						}

						createSqlTxt += "  `" + saData.SnakeStr(fieldName) + "` varchar(" + saData.Itos(maxLen) + ") DEFAULT '' COMMENT '',\n"
					}
				}
			}

			for _, v := range tagAry {
				if v == "idx" {
					if idxSqlTxt != "" {
						idxSqlTxt += ",\n"
					}
					idxSqlTxt += "  KEY `IDX_" + structName + "_" + fieldName + "` (`" + saData.SnakeStr(fieldName) + "`)"
				}
			}
		}

		createSqlTxt += "  PRIMARY KEY (`" + pkNameSnake + "`) USING BTREE"
		if idxSqlTxt != "" {
			createSqlTxt += ",\n" + idxSqlTxt + "\n"
		} else {
			createSqlTxt += "\n"
		}
		createSqlTxt += ") ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '';"
		fmt.Println("建表参考语句：")
		fmt.Println(createSqlTxt)
	}
	tplContent = strings.Replace(tplContent, "{{CreateSql}}", createSqlTxt, -1)

	data := []byte(tplContent)
	f_n := tblName
	f_n = "output/" + f_n + "_sql"
	f_n = strings.Replace(f_n, "tbl_", "", -1)
	if ioutil.WriteFile(f_n, data, 0644) != nil {
		return errors.New("出错")
	}

	fmt.Println("代码生成成功")
	return nil
}
