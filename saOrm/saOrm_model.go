package saOrm

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/saxon134/go-utils/saData"
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

type Ids string

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
				tmp := ""
				for i, v := range ary {
					tmp += saData.I64tos(saData.CharbaseToi64(v))
					if i+1 < len(ary) {
						tmp += ","
					}
				}
				*m = Ids(tmp)
			}
		}
	}

	return nil
}

func (m Ids) Value() (driver.Value, error) {
	if len(m) > 0 {
		tmp := ""
		ary := strings.Split(string(m), ",")
		if len(ary) > 0 {
			for _, v := range ary {
				i64, _ := saData.ToInt64(v)
				if i64 > 0 {
					tmp += saData.I64ToCharbase(i64) + ","
				}
			}

			tmp = strings.TrimSuffix(tmp, ",")
		}
		return tmp, nil
	}

	return "", nil
}
