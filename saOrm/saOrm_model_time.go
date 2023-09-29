package saOrm

import (
	"database/sql/driver"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saTime"
	"time"
)

/****************** Time ******************/
/*存储格式： datetime */

type Time string

func (m *Time) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var str string

	if bAry, ok := value.([]byte); ok && len(bAry) > 0 {
		str = saData.BytesToStr(bAry)
		if m == nil {
			var t = Time(str)
			m = &t
		} else {
			*m = Time(str)
		}
	} else if str, ok = value.(string); ok && len(str) > 0 {
		var t = Time(str)
		if m == nil {
			m = &t
		} else {
			*m = t
		}
	} else if t, ok := value.(time.Time); ok && t.IsZero() == false {
		var t = Time(saTime.TimeToStr(t, saTime.FormatDefault))
		if m == nil {
			m = &t
		} else {
			*m = t
		}
	}
	return nil
}

func (m Time) Value() (driver.Value, error) {
	var t = m.T()
	if t.IsZero() == false {
		var str = saTime.TimeToStr(t, saTime.FormatDefault)
		return str, nil
	}
	return nil, nil
}

func (m *Time) SetNow() {
	now := time.Now()
	if m == nil {
		var t = Time(saTime.TimeToStr(now, saTime.FormatDefault))
		m = &t
	} else {
		*m = Time(saTime.TimeToStr(now, saTime.FormatDefault))
	}
}

func Now() *Time {
	var t = Time(saTime.TimeToStr(time.Now(), saTime.FormatDefault))
	return &t
}

func (m *Time) IsZero() bool {
	if m == nil {
		return true
	}

	var t = m.T()
	return t.IsZero()
}

func (m *Time) T() time.Time {
	var formatAry = []string{saTime.FormatDefault, saTime.FormatSystem, saTime.FormatYMD}
	for _, f := range formatAry {
		var t = saTime.TimeFromStr(string(*m), f)
		if t.IsZero() == false {
			return t
		}
	}
	return time.Time{}
}

func (m *Time) Str(format string) string {
	if format == "" {
		format = time.DateTime
	}

	var t = m.T()
	if t.IsZero() ==false {
		return t.Format(format)
	}

	return ""
}

func (m *Time) String() string {
	if m == nil {
		return ""
	}
	return string(*m)
}

func TimeFromStr(str string) *Time {
	if str == "" || str == "-" || str == "--" {
		return nil
	}
	var t = Time(str)
	return &t
}

func AfterLongTiem() *Time {
	var t = Time("9999-09-09 00:00:00")
	return &t
}

func BeforeLongTiem() *Time {
	var t = Time("1111-01-01 00:00:00")
	return &t
}