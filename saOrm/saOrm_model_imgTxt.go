package saOrm

import (
	"database/sql/driver"
	"errors"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saError"
	"github.com/saxon134/go-utils/saOss"
	"time"
)

const _ValueMaxLength = 500

type ImgTxtItem struct {
	Title string `json:"title,omitempty" form:"title"`
	Desc  string `json:"desc,omitempty" form:"desc"`
	Img   string `json:"img,omitempty" form:"img"`
}

type ImgTxt struct {
	Ary     []ImgTxtItem `json:"ary,omitempty" form:"ary"`
	RichTxt string       `json:"richTxt,omitempty" form:"richTxt"`
	Txt     string       `json:"txt,omitempty" form:"txt"`
	Path    string       `json:"path,omitempty" form:"path"`
}

func (m *ImgTxt) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bAry, ok := value.([]byte)
	if ok && len(bAry) > 0 {
		_ = saData.BytesToModel(bAry, m)
	}

	return nil
}

func (m ImgTxt) Value() (driver.Value, error) {
	if m.Path != "" {
		m.Ary = []ImgTxtItem{}
		m.RichTxt = ""
		m.Txt = ""
	}

	str, _ := saData.ToStr(m)
	if saData.StrLen(str) > _ValueMaxLength {
		return nil, errors.New("ImgTxt长度超过限制，请先调用Save接口将数据存储到OSS")
	}
	return str, nil
}

func (m *ImgTxt) Save(oss saOss.SaOss, path string) (err error) {
	if m == nil || oss == nil {
		return errors.New("ImgTxt数据有误")
	}

	m.Path = "" //删除旧的路径

	var valueStr string
	valueStr, err = saData.ToStr(m)
	if err != nil {
		return
	}

	//数据较小，直接存数据库
	if len(valueStr) < _ValueMaxLength {
		return nil
	}

	//数据较大，存到oss
	if path == "" {
		path = "rt/imgTxt/" + saData.TimeStr(time.Now(), saData.TimeFormat_yymmdd_Line) + "/"
	}
	m.Path, err = oss.UploadTxt(path, valueStr)
	if err != nil {
		return saError.StackError(err)
	}
	return nil
}

func (m *ImgTxt) Get(oss saOss.SaOss) (err error) {
	if m == nil {
		return nil
	}

	if m.Path == "" {
		return nil
	}

	var valueStr string
	valueStr, err = oss.GetTxt(m.Path)
	if err != nil {
		return err
	}

	err = saData.StrToModel(valueStr, m)
	if err != nil {
		return
	}

	m.Path = ""
	return nil
}
