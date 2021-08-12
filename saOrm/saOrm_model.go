package saOrm

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saImg"
	"strings"
)

/******** StringAry ********/

type StringAry []string

func (m *StringAry) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bAry, ok := value.([]byte)
	if ok {
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

/******** Ids ********/

type Ids []int64

func (m *Ids) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bAry, ok := value.([]byte)
	if ok {
		s := string(bAry)
		if len(s) > 0 {
			ary := strings.Split(s, ",")
			if len(ary) > 0 {
				for _, v := range ary {
					if i64, _ := saData.Stoi64(v); i64 > 0 {
						*m = append(*m, i64)
					}
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

/******** CompressIds ********/

type CompressIds []int64

func (m *CompressIds) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bAry, ok := value.([]byte)
	if ok {
		s := string(bAry)
		if len(s) > 0 {
			ary := strings.Split(s, ",")
			if len(ary) > 0 {
				for _, v := range ary {
					if i64 := saData.CharBaseToI64(v); i64 > 0 {
						*m = append(*m, i64)
					}
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

/******** RichTxt ********/

type RichTxt struct {
	Path string
	Md5  string
}

func (m *RichTxt) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bAry, ok := value.([]byte)
	if ok {
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
