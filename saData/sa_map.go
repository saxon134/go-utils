package saData

type Map map[string]interface{}

func Merge(target map[string]string, obj map[string]string) {
	if obj == nil {
		return
	}

	if target == nil {
		target = map[string]string{}
	}

	for k, v := range obj {
		target[k] = v
	}
}
