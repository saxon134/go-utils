package saTime

import (
	"time"
)

const (
	FormatDefault = "2006-01-02 15:04:05"
	FormatSystem  = "2006-01-02T15:04:05Z"

	FormatYMDHMS = "2006-01-02 15:04:05"
	FormatYMDHM  = "2006-01-02 15:04"
	FormatYMDH   = "2006-01-02 15"
	FormatYMD    = "2006-01-02"
	FormatMD     = "01-02"
	FormatYM     = "2006-01"
	FormatY      = "2006"
	FormatHMS    = "15:04:05"
	FormatHM     = "15:04"
	FormatH      = "15"
	FormatM      = "04"
	FormatS      = "05"

	FormatYMDHMSZh = "2006年01月02日 15点04分05秒"
	FormatYMDHMZh  = "2006年01月02日 15点04分"
	FormatYMDZh    = "2006年01月02日"

	FormatYMDHMSDot = "2006.01.02 15:04:05"
	FormatYMDHSDot  = "2006.01.02 15:04"
	FormatYMDDot    = "2006.01.02"

	FormatYMDHMSSimple = "20060102150405"
	FormatYMDHMSimple  = "200601021504"
	FormatYMDHimple    = "2006010215"
	FormatYMDSimple    = "20060102"
)

func Now() *time.Time {
	var t = time.Now().Local()
	return &t
}

func TimeToStr(t time.Time, format string) string {
	if t.IsZero() {
		return ""
	}
	if format == "" {
		format = FormatDefault
	}
	return t.Format(format)
}

func TimeFromStr(s string, format string) time.Time {
	if format == "" {
		format = FormatDefault
	}

	location, _ := time.LoadLocation("Asia/Shanghai")
	if location == nil {
		location = &time.Location{}
	}
	t, _ := time.ParseInLocation(format, s, location)
	if t.IsZero() == false {
		return t
	}

	return time.Time{}
}

func PtrFromStr(s string, format string) *time.Time {
	var t = TimeFromStr(s, format)
	if t.IsZero() {
		return nil
	}
	return &t
}

func PtrFromUnix(unix int64) *time.Time {
	var t = time.Unix(unix, 0)
	if t.IsZero() {
		return nil
	}
	return &t
}

func AdaptTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}

	switch len(s) {
	case 0:
		return time.Time{}
	case 20: //2006-01-02T15:04:05Z
		return TimeFromStr(s, FormatSystem)
	case 19: //2006-01-02 15:04:05
		return TimeFromStr(s, FormatDefault)
	case 16: //2006-01-02 15:04
		return TimeFromStr(s, FormatYMDHM)
	case 10: //2006-01-02
		return TimeFromStr(s, FormatYMD)
	case 8: //20060102
		return TimeFromStr(s, FormatYMDSimple)
	}
	return time.Time{}
}

// 当年第几周
func WeekIndex(t time.Time) int {
	if t.IsZero() == true {
		return -1
	}

	yearDay := t.YearDay()
	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
	firstDayInWeek := int(yearFirstDay.Weekday())

	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}
	var week int
	if yearDay <= firstWeekDays {
		week = 1
	} else {
		week = (yearDay-firstWeekDays)/7 + 2
	}
	return week
}

// LastDayOfMonth 当前月最后一天
func LastDayOfMonth(t time.Time) time.Time {
	var day = TimeFromStr(t.Format(time.DateOnly), time.DateOnly)
	var z = time.Date(day.Year(), day.Month()+1, 0, 23, 59, 59, 0, time.Local)
	return z
}
