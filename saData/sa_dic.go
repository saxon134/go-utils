package saData

func DicDeepCopy(in map[string]interface{}) map[string]interface{} {
	if in == nil {
		return map[string]interface{}{}
	}

	var res = map[string]interface{}{}
	_ = StrToModel(String(in), &res)
	return res
}
