/*
{{CreateSql}}*/

package {{Package}}


import (
    "gorm.io/gorm"
)

func (m *{{Model}}) TableName() string {
	return "{{TblName}}"
}

func (m *{{Model}}) AfterFind(tx *gorm.DB) (err error) {
    {{FromDBSql}}
	return nil
}

func (m *{{Model}}) BeforeCreate(tx *gorm.DB) (err error) {
	{{ToDBSql}}

	return nil
}

func (m *{{Model}}) BeforeSave(tx *gorm.DB) (err error) {
    {{ToDBSql}}

	return nil
}