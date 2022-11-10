package saHit

func If(ok bool, a interface{}, b interface{}) interface{} {
	if ok {
		return a
	} else {
		return b
	}
}

func Str(ok bool, a string, b string) string {
	if ok {
		return a
	} else {
		return b
	}
}

func Int(ok bool, a int, b int) int {
	if ok {
		return a
	} else {
		return b
	}
}

func Int64(ok bool, a int64, b int64) int64 {
	if ok {
		return a
	} else {
		return b
	}
}
