package saHit

func OrStr(a ...string) string {
	for _, v := range a {
		if v != "" {
			return v
		}
	}
	return ""
}

func OrInt(a ...int) int {
	for _, v := range a {
		if v != 0 {
			return v
		}
	}
	return 0
}

func MaxInt(a ...int) int {
	if len(a) == 0 {
		return 0
	}

	var max = a[0]
	for i := 1; i < len(a); i++ {
		if max < a[i] {
			max = a[i]
		}
	}
	return max
}

func OrInt64(a ...int64) int64 {
	for _, v := range a {
		if v != 0 {
			return v
		}
	}
	return 0
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
