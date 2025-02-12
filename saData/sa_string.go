package saData

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func Stoi64(s string) (int64, error) {
	if i64, err := strconv.ParseInt(s, 10, 64); err == nil {
		return int64(i64), nil
	} else {
		return 0, err
	}
}

func Stoui32(s string) (uint32, error) {
	if i32, err := strconv.ParseInt(s, 10, 64); err == nil {
		return uint32(i32), nil
	} else {
		return 0, err
	}
}

func Stoi32(s string) (int32, error) {
	if i32, err := strconv.ParseInt(s, 10, 32); err == nil {
		return int32(i32), nil
	} else {
		return 0, err
	}
}

func Stoi16(s string) (int16, error) {
	if i16, err := strconv.ParseInt(string(s), 10, 16); err == nil {
		return int16(i16), nil
	} else {
		return 0, err
	}
}

func Stoi(s string) (int, error) {
	if i, err := strconv.Atoi(string(s)); err == nil {
		return i, nil
	} else {
		return 0, err
	}
}

func Itos(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func I64tos(i int64) string {
	return strconv.FormatInt(i, 10)
}

func F32tos(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', 2, 32)
}

func Btos(b bool) string {
	if b == true {
		return "1"
	} else {
		return "0"
	}
}

// 去除 ' ' '\n' '\r' '\t'前缀，如果有多个也会去除
func TrimPrefixSpace(s string) string {
	if s != "" {
		var start int = 0
		for i := 0; i < StrLen(s); i++ {
			var c = SubStr(s, i, 1)
			if c != " " && c != "\n" && c != "\t" && c != "\r" {
				start = i
				break
			}
		}
		return SubStr(s, start, StrLen(s)-start)
	}
	return ""
}

// 去除 ' ' '\n' '\r' '\t'后缀，如果有多个也会去除
func TrimSuffixSpace(s string) string {
	if s != "" {
		var end int = StrLen(s)
		for i := end; i > 0; i-- {
			var c string = SubStr(s, i-1, 1)
			if c != " " && c != "\n" && c != "\t" && c != "\r" {
				end = i
				break
			}
		}
		return SubStr(s, 0, end)
	}
	return string(s)
}

// 去除 ' ' '\n' '\r' '\t' 前缀和后缀
func TrimPreSuffixSpace(s string) string {
	s = TrimPrefixSpace(s)
	s = TrimSuffixSpace(s)
	return string(s)
}

// 去除所有 ' ' '\n' '\r' '\t'
func TrimSpace(s string) string {
	if s != "" {
		var i = 0
		for {
			var c = string(s[i : i+1])
			if c == " " || c == "\n" || c == "\t" || c == "\r" {
				if i == 0 {
					s = s[i+1:]
				} else {
					s = s[0:i] + s[i+1:]
				}
			} else {
				i++
			}
			if i >= len(s) {
				break
			}
		}
	}
	return string(s)
}

// 不包含点，如 abc.xml 返回 xml
func GetSuffix(s string, defaultSuffix string) string {
	if s == "" {
		return ""
	}

	strLen := StrLen(s)
	for i := strLen; i > 0; i-- {
		if SubStr(s, i-1, 1) == "." {
			var suffix = SubStr(s, i, strLen-i)
			if len(suffix) < 10 {
				return suffix
			}
			return defaultSuffix
		}
	}
	return defaultSuffix
}

func SubIndex(s string, sub string) int {
	if s == "" || sub == "" {
		return -1
	}

	cnt := StrLen(s)
	subCnt := StrLen(sub)
	tmp := ""
	for i := 0; i < cnt; i++ {
		tmp = SubStr(s, i, subCnt)
		if tmp == sub {
			return i
		}
	}
	return -1
}

func SubStr(s string, start int, cnt int) string {
	var r = []rune(string(s))
	var strLen = len(r)
	if start < 0 || cnt <= 0 || start >= strLen {
		return ""
	}

	if start+cnt > strLen {
		cnt = strLen - start
	}

	return string(r[start : start+cnt])
}

// rune长度
func StrLen(s string) int {
	var r = []rune(string(s))
	return len(r)
}

// rune长度，支持中文
func LenCheck(m interface{}, max int) error {
	str, _ := ToStr(m)
	if StrLen(str) <= max {
		return nil
	}

	return errors.New("超出范围")
}

/*
	去除字符串中H5的style、script；

将标签转换为回车，去除连续回车，去除每段开始、结尾空格
*/
func TrimH5Tags(src string) (str string) {
	s := string(src)

	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	s = re.ReplaceAllStringFunc(s, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	s = re.ReplaceAllString(s, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	s = re.ReplaceAllString(s, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	s = re.ReplaceAllString(s, "\n")

	//去除连续的空白（包括换行）
	re, _ = regexp.Compile("\\s{2,}")
	s = re.ReplaceAllString(s, "\n")

	//去除开头、结尾的空白
	s = strings.TrimSpace(s)

	return s
}

/** 返回16位字符串 */
func RandomStr() string {
	return IdToCode(rand.Int63n(1000), 3) + I64tos(time.Now().UnixMilli())
}

// 通过内存操作，效率极高，但是有风险。只在数据量很大、效率要求高的场景使用
func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// JoinString 字符串拼接 通过strings.Builder，效率更高
func JoinStr(org string, elems ...string) string {
	if len(elems) > 0 {
		var b strings.Builder
		b.WriteString(org)
		for _, s := range elems {
			b.WriteString(s)
		}
		return b.String()
	}
	return org
}

// 根据后缀判断是否是图片
func IsImg(url string) bool {
	if url == "" {
		return false
	}

	var ary = strings.Split(url, "?")
	if len(ary) == 0 {
		return false
	}

	ary = strings.Split(ary[0], ".")
	if len(ary) == 0 {
		return false
	}

	suffix := strings.ToLower(ary[len(ary)-1])
	switch suffix {
	case "png", "jpg", "jpeg", "bmp", "gif", "tif":
		return true
	}
	return false
}

// 根据后缀判断是否是视频
func IsVideo(url string) bool {
	if url == "" {
		return false
	}

	var ary = strings.Split(url, "?")
	if len(ary) == 0 {
		return false
	}

	ary = strings.Split(ary[0], ".")
	if len(ary) == 0 {
		return false
	}

	suffix := strings.ToLower(ary[len(ary)-1])
	switch suffix {
	case "mov", "mp4", "3gp", "flv", "rm", "rmvb", "avi", "mpg", "mlv", "mpe", "mpeg", "dat":
		return true
	}
	return false
}
