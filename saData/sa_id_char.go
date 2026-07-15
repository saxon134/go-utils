package saData

////////////////////////////////////////////////////////////////
// int64和字符串互转，转成的字符串
// 可指定最小长度，字符长度不固定
// 数字范围较大，超出范围字符长度自动加1位
// emw都表示0，避免出现连续e的情况
////////////////////////////////////////////////////////////////

var CodeSource = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"                              //方便给人识别的编码
var MaxSource = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" //最大范围的编码

var defaultSource = "e8trxizqkp9bs2ng4uv5cjh3d6y7af"
var zeroAry = []string{"e", "m", "w"}

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
	var zeroChars = zeroCharsForSource(source)
	var zeroIdx = 0 //控制零值时在zeroAry间轮询
	for {
		if v%sLen == 0 {
			axis = string(zeroChars[zeroIdx]) + axis
			zeroIdx++
			if zeroIdx+1 >= len(zeroChars) {
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
			axis = string(zeroChars[zeroIdx]) + axis
			zeroIdx++
			if zeroIdx+1 >= len(zeroChars) {
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

	var v int64
	var zeroChars = zeroCharsForSource(source)
	for i := 0; i < len(str); i++ {
		r := 0
		if !isZeroChar(str[i], zeroChars) {
			for j := 0; j < len(source); j++ {
				if source[j] == str[i] {
					r = j
					break
				}
			}
		}
		for k := 0; k < len(str)-1-i; k++ {
			r *= len(source)
		}
		v += int64(r)
	}
	return v
}

func zeroCharsForSource(source string) []byte {
	var zeroChars []byte
	for _, zero := range zeroAry {
		if len(zero) != 1 {
			continue
		}
		if !containsByte(source[1:], zero[0]) {
			zeroChars = append(zeroChars, zero[0])
		}
	}
	if len(zeroChars) == 0 {
		zeroChars = append(zeroChars, source[0])
	}
	return zeroChars
}

func isZeroChar(v byte, zeroChars []byte) bool {
	for _, zero := range zeroChars {
		if v == zero {
			return true
		}
	}
	return false
}

func containsByte(str string, v byte) bool {
	for i := 0; i < len(str); i++ {
		if str[i] == v {
			return true
		}
	}
	return false
}
