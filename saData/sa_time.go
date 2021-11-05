package saData

import (
	"time"
)

type SaTimeFormat int8

const (
	TimeFormat_Default SaTimeFormat = iota //"2006-01-02 15:04:05"

	TimeFormat_hhmm        //"15:04"
	TimeFormat_sys_default //"2006-01-02T15:04:05Z"

	TimeFormat_yymm_Dotted       //"2006.01"
	TimeFormat_yymmdd_Dotted     //"2006.01.02"
	TimeFormat_yymmddhhmm_Dotted //"2006.01.02 15:04"
	TimeFormat_Dotted            //"2006.01.02 15:04:05"

	TimeFormat_yymm_Line       //"2006-01"
	TimeFormat_yymmdd_Line     //"2006-01-02"
	TimeFormat_yymmddhhmm_Line //"2006-01-02 15:04"
	TimeFormat_yymmddhhmmss_Line //"2006-01-02 15:04:05"

	TimeFormat_yymm_Chinese       //"2006年01月"
	TimeFormat_yymmdd_Chinese     //"2006年01月02日"
	TimeFormat_yymmddhhmm_Chinese //"2006年01月02日 15点04分"
	TimeFormat_Chinese            //"2006年01月02日 15时04分05秒"

	TimeFormate_yearStr   //2006
	TimeFormate_monthStr  //200601
	TimeFormate_dayStr    //20060102
	TimeFormate_hourStr   //2006010215
	TimeFormate_minuteStr //200601021504
	TimeFormate_secondStr //20060102150405
)

func (m SaTimeFormat) Format() string {
	switch m {
	case TimeFormat_Default:
		return "2006-01-02 15:04:05"
	case TimeFormat_hhmm:
		return "15:04"
	case TimeFormat_sys_default:
		return "2006-01-02T15:04:05Z"
	case TimeFormat_yymm_Dotted:
		return "2006.01"
	case TimeFormat_yymmdd_Dotted:
		return "2006.01.02"
	case TimeFormat_yymmddhhmm_Dotted:
		return "2006.01.02 15:04"
	case TimeFormat_Dotted:
		return "2006.01.02 15:04:05"
	case TimeFormat_yymm_Line:
		return "2006-01"
	case TimeFormat_yymmdd_Line:
		return "2006-01-02"
	case TimeFormat_yymmddhhmm_Line:
		return "2006-01-02 15:04"
	case TimeFormat_yymmddhhmmss_Line:
		return "2006-01-02 15:04:05"
	case TimeFormat_yymm_Chinese:
	return "2006年01月"
	case TimeFormat_yymmdd_Chinese:
		return "2006年01月02日"
	case TimeFormat_yymmddhhmm_Chinese:
	return "2006年01月02日 15点04分"
	case TimeFormat_Chinese:
		return "2006年01月02日 15时04分05秒"
	case TimeFormate_yearStr:
		return "2006"
	case TimeFormate_monthStr:
		return "200601"
	case TimeFormate_dayStr:
		return "20060102"
	case TimeFormate_hourStr:
	return "2006010215"
	case TimeFormate_minuteStr:
		return "200601021504"
	case TimeFormate_secondStr:
		return "20060102150405"
	}
	return "2006-01-02 15:04:05"
}

func Now() *time.Time {
	t := time.Now()
	return &t
}

func TimeStr(t time.Time, format SaTimeFormat) string {
	s := ""
	y := t.Year()
	m := t.Month()
	d := t.Day()
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()

	yS := Itos(y)

	mS := Itos(int(m))
	if m < 10 {
		mS = "0" + mS
	}

	dS := Itos(d)
	if d < 10 {
		dS = "0" + dS
	}

	hourS := Itos(hour)
	if hour < 10 {
		hourS = "0" + hourS
	}

	minuteS := Itos(minute)
	if minute < 10 {
		minuteS = "0" + minuteS
	}

	secondS := Itos(second)
	if second < 10 {
		secondS = "0" + secondS
	}

	switch format {
	case TimeFormat_yymm_Dotted:
		s = yS + "." + mS
	case TimeFormat_yymmdd_Dotted:
		s = yS + "." + mS + "." + dS
	case TimeFormat_yymmddhhmm_Dotted:
		s = yS + "." + mS + "." + dS + " " + hourS + ":" + minuteS
	case TimeFormat_Dotted:
		s = yS + "." + mS + "." + dS + " " + hourS + ":" + minuteS + ":" + secondS

	case TimeFormat_yymm_Line:
		s = yS + "-" + mS
	case TimeFormat_yymmdd_Line:
		s = yS + "-" + mS + "-" + dS
	case TimeFormat_yymmddhhmm_Line:
		s = yS + "-" + mS + "-" + dS + " " + hourS + ":" + minuteS
	case TimeFormat_yymmddhhmmss_Line:
		s = yS + "-" + mS + "-" + dS + " " + hourS + ":" + minuteS + ":" + secondS
	case TimeFormat_yymm_Chinese:
		s = yS + "年" + mS + "月"
	case TimeFormat_yymmdd_Chinese:
		s = yS + "年" + mS + "月" + dS + "日"
	case TimeFormat_yymmddhhmm_Chinese:
		s = yS + "年" + mS + "月" + dS + "日 " + hourS + "点" + minuteS + "分"
	case TimeFormat_Chinese:
		s = yS + "年" + mS + "月" + dS + "日 " + hourS + "点" + minuteS + "分" + secondS + "秒"

	case TimeFormate_yearStr:
		s = yS
	case TimeFormate_monthStr:
		s = yS + mS
	case TimeFormate_dayStr:
		s = yS + mS + dS
	case TimeFormate_hourStr:
		s = yS + mS + dS + hourS
	case TimeFormate_minuteStr:
		s = yS + mS + dS + hourS + minuteS
	case TimeFormate_secondStr:
		s = yS + mS + dS + hourS + minuteS + secondS

	case TimeFormat_hhmm:
		s = minuteS + ":" + secondS

	case TimeFormat_Default:
		s = yS + "-" + mS + "-" + dS + " " + hourS + ":" + minuteS + ":" + secondS
	case TimeFormat_sys_default:
		s = yS + "-" + mS + "-" + dS + "T" + hourS + ":" + minuteS + ":" + secondS + "Z"
	}

	return s
}

func AdaptToTime(s string) time.Time {
	if s==""{
		return time.Time{}
	}

	switch len(s) {
	case 0:
		return time.Time{}
	case 20: //2006-01-02T15:04:05Z
		var t = StrToTime(TimeFormat_sys_default, s)
		if t.IsZero() ==false {
			return t
		}
	case 19: //2006-01-02 15:04:05
		var t = StrToTime(TimeFormat_Default, s)
		if t.IsZero() ==false {
			return t
		}
	case 16: //2006-01-02 15:04
		var t = StrToTime(TimeFormat_yymmddhhmm_Line, s)
		if t.IsZero() ==false {
			return t
		}
	case 10: //2006-01-02
		var t = StrToTime(TimeFormat_yymmdd_Line, s)
		if t.IsZero() ==false {
			return t
		}
	case 8: //20060102
		var t = StrToTime(TimeFormate_dayStr, s)
		if t.IsZero() ==false {
			return t
		}
	default:
		return time.Time{}
	}
	return time.Time{}
}

func StrToTime(format SaTimeFormat, s string) time.Time {
	formatStr := "2006-01-02 15:04:05"
	switch format {
	case TimeFormat_yymm_Line:
		formatStr = "2006-01"
	case TimeFormat_yymmdd_Line:
		formatStr = "2006-01-02"
	case TimeFormat_yymmddhhmm_Line:
		formatStr = "2006-01-02 15:04"

	case TimeFormat_yymm_Dotted:
		formatStr = "2006.01"
	case TimeFormat_yymmdd_Dotted:
		formatStr = "2006.01.02"
	case TimeFormat_yymmddhhmm_Dotted:
		formatStr = "2006.01.02 15:04"
	case TimeFormat_Dotted:
		formatStr = "2006.01.02 15:04:05"

	case TimeFormat_yymm_Chinese:
		formatStr = "2006年01月"
	case TimeFormat_yymmdd_Chinese:
		formatStr = "2006年01月02日"
	case TimeFormat_yymmddhhmm_Chinese:
		formatStr = "2006年01月02日 15点04分"
	case TimeFormat_Chinese:
		formatStr = "2006年01月02日 15点04分"

	case TimeFormate_yearStr:
		formatStr = "2006"
	case TimeFormate_monthStr:
		formatStr = "200601"
	case TimeFormate_dayStr:
		formatStr = "20060102"
	case TimeFormate_hourStr:
		formatStr = "2006010215"
	case TimeFormate_minuteStr:
		formatStr = "200601021504"
	case TimeFormate_secondStr:
		formatStr = "20060102150405"

	case TimeFormat_hhmm:
		formatStr = "15:04"

	case TimeFormat_sys_default:
		formatStr = "2006-01-02T15:04:05Z"
	}

	location, _ := time.LoadLocation("Asia/Shanghai")
	if location == nil {
		location = &time.Location{}
	}
	t, _ := time.ParseInLocation(formatStr, s, location)
	return t
}

//当年第几周
func WeekIndex(t time.Time) int {
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
