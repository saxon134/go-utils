package saData

import (
	"os"
	"strings"
)

func ModifySysFilePath(f string) string {
	if f == "" {
		return ""
	}

	separator := "/"
	if os.IsPathSeparator('\\') {
		separator = "\\"
	}

	f = strings.Replace(f, "/", separator, -1)
	return f
}

/* 下划线命名转为大驼峰命名
XxYy to xx_yy , XxYY to xx_yy */
func SnakeStr(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

/* 下划线命名转化为小驼峰命名
xx_yy to xxYy */
func CamelStr(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}

	//首字母小写
	d := data[0]
	if d >= 'A' && d <= 'Z' {
		d += 32
		data[0] = d
	}
	return string(data[:])
}
