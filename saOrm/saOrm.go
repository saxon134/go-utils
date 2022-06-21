package saOrm

import (
	"github.com/saxon134/go-utils/saHit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

var _db *DB

type Conf struct {
	MaxIdleConns int
	MaxOpenConns int
}

func Open(dsn string, conf Conf) *DB {
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	var db *gorm.DB
	var err error
	db, err = gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
	})
	if err != nil {
		panic("MySQL启动异常" + err.Error())
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(saHit.Int(conf.MaxIdleConns > 0, conf.MaxIdleConns, 10))
	sqlDB.SetMaxOpenConns(saHit.Int(conf.MaxOpenConns > 0, conf.MaxOpenConns, 5))

	_db = &DB{DB: db}
	return _db
}

// Session
// @Description: 每次调用都会生成一个新的会话，为了避免全局会话影响
func (m *DB) Session() *gorm.DB {
	tx := _db.DB.Session(&gorm.Session{})
	return tx
}

//
//func Scan(rows *sql.Rows, aryPtr QueryAry) error {
//	defer rows.Close()
//
//	var cnt int64
//	nInds := make([]reflect.Value, len(sInds))
//	sInd := sInds[0]
//	for rows.Next() {
//		columns, err := rows.Columns()
//		if err != nil {
//			return err
//		}
//
//		columnsMp := make(map[string]interface{}, len(columns))
//		refs := make([]interface{}, 0, len(columns))
//		for _, col := range columns {
//			var ref interface{}
//			columnsMp[col] = &ref
//			refs = append(refs, &ref)
//		}
//
//		if err := rows.Scan(refs...); err != nil {
//			return err
//		}
//
//		if cnt == 0 && !sInd.IsNil() {
//			sInd.Set(reflect.New(sInd.Type()).Elem())
//		}
//
//		var ind reflect.Value
//		if eTyps[0].Kind() == reflect.Ptr {
//			ind = reflect.New(eTyps[0].Elem())
//		} else {
//			ind = reflect.New(eTyps[0])
//		}
//
//		if ind.Kind() == reflect.Ptr {
//			ind = ind.Elem()
//		}
//
//		if sMi != nil {
//			for _, col := range columns {
//				if fi := sMi.fields.GetByColumn(col); fi != nil {
//					value := reflect.ValueOf(columnsMp[col]).Elem().Interface()
//					field := ind.FieldByIndex(fi.fieldIndex)
//					if fi.fieldType&IsRelField > 0 {
//						mf := reflect.New(fi.relModelInfo.addrField.Elem().Type())
//						field.Set(mf)
//						field = mf.Elem().FieldByIndex(fi.relModelInfo.fields.pk.fieldIndex)
//					}
//					if fi.isFielder {
//						fd := field.Addr().Interface().(Fielder)
//						err := fd.SetRaw(value)
//						if err != nil {
//							return 0, errors.Errorf("set raw error:%s", err)
//						}
//					} else {
//						o.setFieldValue(field, value)
//					}
//				}
//			}
//		} else {
//			var recursiveSetField func(rv reflect.Value)
//			recursiveSetField = func(rv reflect.Value) {
//				for i := 0; i < rv.NumField(); i++ {
//					f := rv.Field(i)
//					fe := rv.Type().Field(i)
//
//					// check if the field is a Struct
//					// recursive the Struct type
//					if fe.Type.Kind() == reflect.Struct {
//						recursiveSetField(f)
//					}
//
//					_, tags := parseStructTag(fe.Tag.Get(defaultStructTagName))
//					var col string
//					if col = tags["column"]; col == "" {
//						col = nameStrategyMap[nameStrategy](fe.Name)
//					}
//					if v, ok := columnsMp[col]; ok {
//						value := reflect.ValueOf(v).Elem().Interface()
//						o.setFieldValue(f, value)
//					}
//				}
//			}
//
//			// init call the recursive function
//			recursiveSetField(ind)
//		}
//
//		if eTyps[0].Kind() == reflect.Ptr {
//			ind = ind.Addr()
//		}
//	}
//}
