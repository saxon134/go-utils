package saOrm

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saImg"
	"strings"
	"time"
)

/* StringAry
数据库存储格式：json **/
type StringAry []string

func (m *StringAry) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bAry, ok := value.([]byte)
	if ok && len(bAry) > 0 {
		err := json.Unmarshal(bAry, m)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (m StringAry) Value() (driver.Value, error) {
	bAry, err := json.Marshal(m)
	if err == nil && len(bAry) > 0 {
		s := string(bAry)
		if s == "null" {
			return "", nil
		}
		return s, nil
	}

	return "", err
}

func (m StringAry) IsSame(n StringAry) bool {
	if len(m) != len(n) {
		return false
	}
	for i, s := range m {
		s = saImg.DeleteUriRoot(s)
		n[i] = saImg.DeleteUriRoot(n[i])
		if s != n[i] {
			return false
		}
	}
	return true
}

/* Ids
数据库存储格式： id1,id2,id3 **/
type Ids []int64

func (m *Ids) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bAry, ok := value.([]byte)
	if ok && len(bAry) > 0 {
		s := string(bAry)
		ary := strings.Split(s, ",")
		if len(ary) > 0 {
			for _, v := range ary {
				if i64, _ := saData.Stoi64(v); i64 > 0 {
					*m = append(*m, i64)
				}
			}
		}
	}

	return nil
}

func (m Ids) Value() (driver.Value, error) {
	if len(m) > 0 {
		tmp := ""
		for _, v := range m {
			tmp += saData.I64tos(v) + ","
		}
		tmp = strings.TrimSuffix(tmp, ",")
		return tmp, nil
	}

	return "", nil
}

/* CompressIds
数据库存储格式：字符串，ID转换为88进制，逗号分隔 **/
type CompressIds []int64

func (m *CompressIds) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bAry, ok := value.([]byte)
	if ok && len(bAry) > 0 {
		s := string(bAry)
		ary := strings.Split(s, ",")
		if len(ary) > 0 {
			for _, v := range ary {
				if i64 := saData.CharBaseToI64(v); i64 > 0 {
					*m = append(*m, i64)
				}
			}
		}
	}

	return nil
}

func (m CompressIds) Value() (driver.Value, error) {
	if len(m) > 0 {
		tmp := ""
		for _, v := range m {
			tmp += saData.I64tos(v) + ","
		}
		tmp = strings.TrimSuffix(tmp, ",")
		return tmp, nil
	}

	return "", nil
}

/* RichTxt
数据库存储格式：字符串，OSS路径和MD5空格隔开 **/
type RichTxt struct {
	Md5  string
	Path string
	//Type int 后续扩展类型的时候可以用
}

func (m *RichTxt) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bAry, ok := value.([]byte)
	if ok && len(bAry) > 0 {
		str := saData.BytesToStr(bAry)
		ary := strings.Split(str, " ")
		if len(ary) == 2 {
			m.Path = saImg.AddDefaultUriRoot(ary[0])
			m.Md5 = ary[1]
		}
		return nil
	}

	return nil
}

func (m RichTxt) Value() (driver.Value, error) {
	m.Path = saImg.DeleteUriRoot(m.Path)
	if m.Path != "" {
		m.Md5 = saData.Md5(m.Path, true)
		m.Md5 = strings.TrimSpace(m.Md5)
		return m.Path + " " + m.Md5, nil
	}
	return "", nil
}

/* Price
数据库存储格式：整数，分为单位 **/
type Price float32

func (m *Price) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	i, err := saData.ToInt(value)
	if err == nil {
		*m = Price(saData.Fen2Yuan(i, saData.RoundTypeDefault))
	}
	return err
}

func (m Price) Value() (driver.Value, error) {
	i := saData.Yuan2Fen(float32(m), saData.RoundTypeDefault)
	return saData.Itos(i), nil
}

/* Time
数据库存储格式：datetime **/
type Time struct {
	time.Time
}

func (m *Time) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var str string

	if bAry, ok := value.([]byte); ok && len(bAry) > 0 {
		str = saData.BytesToStr(bAry)
		if str == "" {
			return nil
		}
		m.Time = saData.StrToTime(saData.TimeFormat_Default, str)
	} else if str, ok = value.(string); ok && len(str) > 0 {
		m.Time = saData.StrToTime(saData.TimeFormat_Default, str)
	}
	return nil
}

func (m Time) Value() (driver.Value, error) {
	str := saData.TimeStr(m.Time, saData.TimeFormat_Default)
	return str, nil
}

func (m *Time) Now() {
	now := time.Now()
	if m == nil {
		m = &Time{now}
	} else {
		m.Time = now
	}
}

func (m *Time) IsZero() bool {
	if m == nil || m.IsZero() {
		return true
	}
	return false
}
