package saData

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
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
	return strconv.FormatInt(int64(i), 10)
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

//去除 ' ' '\n' '\r' '\t'前缀，如果有多个也会去除
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

//去除 ' ' '\n' '\r' '\t'后缀，如果有多个也会去除
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

func GetSuffix(s string) string {
	if s == "" {
		return ""
	}

	strLen := StrLen(s)
	for i := strLen; i > 0; i-- {
		if SubStr(s, i-1, 1) == "." {
			return SubStr(s, i, strLen-i)
		}
	}
	return ""
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

//rune长度
func StrLen(s string) int {
	var r = []rune(string(s))
	return len(r)
}

//rune长度，支持中文
func LenCheck(m interface{}, max int) error {
	str, _ := ToStr(m)
	if StrLen(str) <= max {
		return nil
	}

	return errors.New("超出范围")
}

/* 去除字符串中H5的style、script；
将标签转换为回车，去除连续回车，去除每段开始、结尾空格 */
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

func RandomStr() string {
	t := time.Now().UnixNano() / 1000
	r := rand.Intn(1000)
	s := I64tos(t) + Itos(r)

	h := md5.New()
	_, _ = io.WriteString(h, s)
	s = fmt.Sprintf("%x", h.Sum(nil))
	return s
}

// int64和88进制字符转换
func I64ToCharBase(v int64) string {
	if v <= 0 {
		return ""
	}

	ary := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^*()_=+<>.?/[]{}|`~"
	var axis string

	if v%int64(len(ary)) == 0 {
		axis = string(ary[len(ary)-1]) + axis
		v -= int64(len(ary))
	} else {
		axis = string(ary[(v%int64(len(ary)))-1]) + axis
	}
	v /= int64(len(ary))

	return axis
}
func CharBaseToI64(str string) int64 {
	if str == "" {
		return 0
	}

	var v int64
	ary := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^*()_=+<>.?/[]{}|`~"
	for i := len(str) - 1; i >= 0; i-- {
		for j := 0; j < len(ary); j++ {
			if ary[j] == str[i] {
				var tv int64 = int64(j) + 1
				for k := 0; k < len(str)-1-i; k++ {
					tv *= int64(len(ary))
				}
				v += tv
				break
			}
		}
	}
	return v
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

func StrEncode(s string) string {
	if s != "" {
		v := url.Values{}
		v.Add("k", s)
		s = v.Encode()
		return string([]rune(s)[2:])
	}
	return ""
}

//query & map 互转
func MapToQuery(m map[string]string) string {
	if m != nil {
		urlV := url.Values{}
		for k, v := range m {
			if k != "" {
				urlV.Add(k, v)
			}
		}
		return urlV.Encode()
	}
	return ""
}
func QueryToMap(urlStr string) map[string]string {
	values, _ := url.ParseQuery(urlStr)
	m := map[string]string{}
	for k, v := range values {
		if k != "" {
			m[k] = v[0]
		}
	}

	return m
}

//返回结果是： /r/path
func ConnPath(r string, path string) (full string) {
	r = strings.TrimSuffix(r, "/")
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	if len(r) > 0 {
		full = r + "/" + path
	} else {
		full = path
	}

	return full
}
