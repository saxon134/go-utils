package saData

import "strings"

/**
int64和字符串互转，转成的字符串非固定长度，可指定最小长度
数字范围较大，超出范围字符长度自动加1位
emw都表示0，避免出现连续e的情况
*/

var defaultSource = "e8trxizqkp9bs2ng4uv5cjh3d6y7af"
var zeroAry = []string{"e", "m", "w"}

//0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^*()_=+<>.?/[]{}|`~

func IdToChar(v int64) string {
	return IdToCharWithSource(v, 3, "")
}

func CharToId(str string) int64 {
	return CharToIdWithSource(str, "")
}

func IdToCharWithSource(v int64, minLen int, source string) string {
	if v <= 0 {
		return ""
	}

	if len(source) == 0 {
		source = defaultSource
	}

	var axis string
	var sLen = int64(len(source))
	var zeroIdx = 0 //控制零值时在zeroAry间轮询
	for {
		if v%sLen == 0 {
			axis = zeroAry[zeroIdx] + axis
			zeroIdx++
			if zeroIdx+1 >= len(zeroAry) {
				zeroIdx = 0
			}
		} else {
			axis = string(source[(v%sLen)]) + axis
		}

		v /= sLen
		if v <= 0 {
			break
		}
	}

	for i := 0; i < minLen; i++ {
		if len(axis) < minLen {
			axis = zeroAry[zeroIdx] + axis
			zeroIdx++
			if zeroIdx+1 >= len(zeroAry) {
				zeroIdx = 0
			}
		}
	}

	return axis
}

func CharToIdWithSource(str string, source string) int64 {
	if str == "" {
		return 0
	}

	if len(source) == 0 {
		source = defaultSource
	}

	//所有零值都替换成第一个
	for _, s := range zeroAry {
		str = strings.Replace(str, s, zeroAry[0], -1)
	}

	var v int64
	for i := 0; i < len(str); i++ {
		for j := 0; j < len(source); j++ {
			if source[j] == str[i] {
				r := j
				for k := 0; k < len(str)-1-i; k++ {
					r *= len(source)
				}
				v += int64(r)
				break
			}
		}
	}
	return v
}
