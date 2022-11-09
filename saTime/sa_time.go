package saTime

import (
	"time"
)

const (
	FormatDefault = "2006-01-02 15:04:05"
	FormatSystem  = "2006-01-02T15:04:05Z"

	FormatYMDHMS = "2006-01-02 15:04:05"
	FormatYMDHM  = "2006-01-02 15:04"
	FormatYMD    = "2006-01-02"
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
	FormatYMDSimple    = "20060102"
)

func Now() *time.Time {
	t := time.Now()
	return &t
}

func TimeToStr(t *time.Time, format string) string {
	if t == nil || t.IsZero() {
		return ""
	}
	if format == "" {
		format = FormatDefault
	}
	return t.Format(format)
}

func StrToTime(s string, format string) *time.Time {
	if format == "" {
		format = FormatDefault
	}

	location, _ := time.LoadLocation("Asia/Shanghai")
	if location == nil {
		location = &time.Location{}
	}
	t, _ := time.ParseInLocation(format, s, location)
	if t.IsZero() == false {
		return &t
	}
	return nil
}

func AdaptTime(s string) *time.Time {
	if s == "" {
		return &time.Time{}
	}

	switch len(s) {
	case 0:
		return &time.Time{}
	case 20: //2006-01-02T15:04:05Z
		return StrToTime(FormatSystem, s)
	case 19: //2006-01-02 15:04:05
		return StrToTime(FormatDefault, s)
	case 16: //2006-01-02 15:04
		return StrToTime(FormatYMDHM, s)
	case 10: //2006-01-02
		return StrToTime(FormatYMD, s)
	case 8: //20060102
		return StrToTime(FormatYMDSimple, s)
	}
	return nil
}

//当年第几周
func WeekIndex(t *time.Time) int {
	if t == nil || t.IsZero() == true {
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
