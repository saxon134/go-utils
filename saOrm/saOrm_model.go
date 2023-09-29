package saOrm

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saUrl"
	"strings"
)

/******* StringAry *********/
/* 存储格式：json */

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

/******** Dic ********/
/* 存储格式：json */

type Dic map[string]interface{}

func (m *Dic) Scan(value interface{}) error {
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

func (m Dic) Value() (driver.Value, error) {
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

/*********** Ids **********/
/* 存储格式： id1,id2,id3 */

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
				if i64, e := saData.Stoi64(v); e == nil {
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

/*********** Ints **********/
/* 存储格式： int1,int2,int3 */

type Ints []int

func (m *Ints) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bAry, ok := value.([]byte)
	if ok && len(bAry) > 0 {
		s := string(bAry)
		ary := strings.Split(s, ",")
		if len(ary) > 0 {
			for _, v := range ary {
				if i, e := saData.Stoi(v); e == nil {
					*m = append(*m, i)
				}
			}
		}
	}

	return nil
}

func (m Ints) Value() (driver.Value, error) {
	if len(m) > 0 {
		tmp := ""
		for _, v := range m {
			tmp += saData.Itos(v) + ","
		}
		tmp = strings.TrimSuffix(tmp, ",")
		return tmp, nil
	}

	return "", nil
}

/**************** CompressIds ******************/
/*存储格式： 字符串，ID转换为88进制，逗号分隔 */

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

/**************** Price ******************/
/*存储格式： 整数，分为单位，故不存在四舍五入一说 */

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

/********** LiPrice **********/
/*存储格式： 整数，厘为单位 */

type LiPrice float32

func (m *LiPrice) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	i, err := saData.ToInt(value)
	if err == nil {
		*m = LiPrice(saData.Li2Yuan(i, saData.RoundTypeDefault))
	}
	return err
}

func (m LiPrice) Value() (driver.Value, error) {
	i := saData.Yuan2Li(float32(m), saData.RoundTypeDefault)
	return saData.Itos(i), nil
}

/********** PriceDigit4 **********/
/*存储格式： 4位小数点钱 */

type PriceDigit4 float32

func (m *PriceDigit4) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	i, err := saData.ToInt(value)
	if err == nil {
		*m = PriceDigit4(saData.IntToFloat(i, 4, saData.RoundTypeDefault))
	}
	return err
}

func (m PriceDigit4) Value() (driver.Value, error) {
	i := saData.FloatToInt64(float32(m), 4, saData.RoundTypeDefault)
	return saData.I64tos(i), nil
}

/********** Weight **********/
/*存储格式： 整数，克为单位 */

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

/****************** Rate ******************/
/*存储格式： 整数，百分比为单位，如：23% */

type Rate float32

func (m *Rate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	i, err := saData.ToInt(value)
	if err == nil {
		*m = Rate(i / 100)
	}
	return err
}

func (m Rate) Value() (driver.Value, error) {
	i := m * 100
	return saData.Itos(int(i)), nil
}
