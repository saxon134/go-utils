package saVideo

type Ordered interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | string
}

func InArray[T Ordered](needle T, haystack []T) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}
	return false
}
