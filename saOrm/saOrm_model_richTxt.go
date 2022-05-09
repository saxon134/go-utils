package saOrm

import (
	"database/sql/driver"
	"errors"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saError"
	"github.com/saxon134/go-utils/saOss"
	"time"
)

type RichTxtType int

const (
	RichTxtTypeNull RichTxtType = iota
	RichTxtTypeTxt  RichTxtType = 1 //纯文本
	RichTxtTypeHtml RichTxtType = 2 //富文本
	RichTxtTypeJson RichTxtType = 3 //json
)

//json时对象格式应该是这样的，由前端控制，后端不做解析
type RtItem struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Img string `json:"img"`
}

/* 数据库存储格式：json或者内容字符，当内容小于250时，直接存入数据库；否则存入OSS，content存储路径
不管存储在哪里，MD5都是原始content的MD5 **/
type RichTxt struct {
	Type    RichTxtType `json:"type,omitempty"`
	Md5     string      `json:"md5,omitempty"`
	InOss   bool        `json:"inOss,omitempty"`
	Content string      `json:"content,omitempty"`
}

func (m *RichTxt) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bAry, ok := value.([]byte)
	if ok && len(bAry) > 0 {
		_ = saData.BytesToModel(bAry, m)
	}

	return nil
}

func (m RichTxt) Value() (driver.Value, error) {
	str, _ := saData.ToStr(m)
	if saData.StrLen(str) > 250 {
		return nil, errors.New("RichTxt长度超过限制")
	}
	return str, nil
}

func (m *RichTxt) Save(oss saOss.SaOss, t RichTxtType, txt string, path string) (err error) {
	if m == nil || oss == nil {
		return errors.New("RichTxt数据有误")
	}

	//默认为旧数据类型
	if t == RichTxtTypeNull {
		t = m.Type
	}

	//数据无变化
	md5 := saData.Md5(txt, true)
	if md5 == m.Md5 {
		return nil
	}

	//数据较小，直接存数据库
	if len(txt) < 200 {
		m.Md5 = md5
		m.Content = txt
		m.InOss = false
		m.Type = t
		return nil
	}

	//数据较大，存到oss
	if path == "" {
		path = "rt/default/" + saData.TimeStr(time.Now(), saData.TimeFormat_yymmdd_Line) + "/"
	}
	m.Content, err = oss.UploadTxt(path, txt)
	if err != nil {
		return saError.StackError(err)
	}
	m.Md5 = md5
	return nil
}

func (m *RichTxt) Get(oss saOss.SaOss) (txt string, err error) {
	if m == nil {
		return "", nil
	}

	if m.InOss && m.Content != "" {
		m.Content, err = oss.GetTxt(m.Content)
		if len(txt) > 0 && len(m.Md5) == 0 {
			m.Md5 = saData.Md5(txt, true)
		}
		return m.Content, err
	}
	return m.Content, nil
}
