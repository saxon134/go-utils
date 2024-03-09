package saHit

func OrStr(a, b string) string {
	return Str(a != "", a, b)
}

func OrInt(a, b int) int {
	return Int(a != 0, a, b)
}

func OrInt64(a, b int64) int64 {
	return Int64(a != 0, a, b)
}

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

func Float(ok bool, a float32, b float32) float32 {
	if ok {
		return a
	} else {
		return b
	}
}

func Float64(ok bool, a float64, b float64) float64 {
	if ok {
		return a
	} else {
		return b
	}
}

