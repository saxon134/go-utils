package saData

import "regexp"

func IsPhone(str string) bool {
	isorno, _ := regexp.MatchString(`^(13[0-9]|14[5-9]|15[012356789]|166|17[0-8]|18[0-9]|19[8-9])[0-9]{8}`, str)
	return isorno
}
