package saOrm

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saTime"
	"github.com/saxon134/go-utils/saData/saUrl"
	"strings"
	"time"
)

// StringAry
// 数据库存储格式：json
type StringAry []string

func (m StringAry) TrimSpace() StringAry {
	banners := make(StringAry, 0, len(m))
	for _, v := range m {
		if len(v) > 0 {
			banners = append(banners, v)
		}
	}
	return banners
}

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

func (m StringAry) IsSameImages(n StringAry) bool {
	if len(m) != len(n) {
		return false
	}
	for i, s := range m {
		s = saUrl.DeleteUriRoot(s)
		n[i] = saUrl.DeleteUriRoot(n[i])
		if s != n[i] {
			return false
		}
	}
	return true
}

// Ids
// 数据库存储格式： id1,id2,id3
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

// CompressIds
// 数据库存储格式：字符串，ID转换为88进制，逗号分隔
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
				if i64 := saData.CharToId(v); i64 > 0 {
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

// Price
// 数据库存储格式：整数，分为单位，故不存在四舍五入一说
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

// LiPrice
// 数据库存储格式：整数，厘为单位
type LiPrice float32

func (m *LiPrice) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	i, err := saData.ToInt(value)
	if err == nil {
		*m = LiPrice(saData.Fen2Yuan(i, saData.RoundTypeDefault))
	}
	return err
}

func (m LiPrice) Value() (driver.Value, error) {
	i := saData.Yuan2Fen(float32(m), saData.RoundTypeDefault)
	return saData.Itos(i), nil
}

// Weight
// 数据库存储格式：整数，克为单位
type Weight float32

func (m *Weight) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	i, err := saData.ToInt(value)
	if err == nil {
		*m = Weight(i / 1000)
	}
	return err
}

func (m Weight) Value() (driver.Value, error) {
	i := m * 1000
	return saData.Itos(int(i)), nil
}

// Rate
// 数据库存储格式：整数，万份之一为单位，如：0.23%
type Rate float32

func (m *Rate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	i, err := saData.ToInt(value)
	if err == nil {
		*m = Rate(i / 10000)
	}
	return err
}

func (m Rate) Value() (driver.Value, error) {
	i := m * 10000
	return saData.Itos(int(i)), nil
}

// Time
// 数据库存储格式：datetime
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
		t := saTime.TimeFromStr(str, saTime.FormatDefault)
		if t != nil {
			m.Time = *t
		}
	} else if str, ok = value.(string); ok && len(str) > 0 {
		t := saTime.TimeFromStr(str, saTime.FormatDefault)
		if t != nil {
			m.Time = *t
		}
	}
	return nil
}

func (m Time) Value() (driver.Value, error) {
	str := saTime.TimeToStr(&m.Time, saTime.FormatDefault)
	return str, nil
}

func (m *Time) SetNow() {
	now := time.Now()
	if m == nil {
		m = &Time{now}
	} else {
		m.Time = now
	}
}

func Now() *Time {
	return &Time{time.Now()}
}

func (m *Time) IsZero() bool {
	if m == nil || m.Time.IsZero() {
		return true
	}
	return false
}
