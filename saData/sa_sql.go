package saData

import (
	"fmt"
	"strings"
)

// SQL转义
func SQLEspace(s string) string {
	s = strings.ReplaceAll(s, `'`, `''`)
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "`", "\\`")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// 去出输入空字符串
func TrimEmpty(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\t", "")
	str = strings.ReplaceAll(str, "#N/A", "")
	str = strings.ReplaceAll(str, "\u200B", "")
	str = strings.ReplaceAll(str, "\ufeff", "")
	str = strings.ReplaceAll(str, "(null)", "")
	if str == "-" || str == "--" {
		return ""
	}
	return str
}

// 格式化逗号
func FormatComma(str string) string {
	str = strings.ReplaceAll(str, " ", ",")
	str = strings.ReplaceAll(str, "\t", ",")
	str = strings.ReplaceAll(str, "，", ",")
	str = strings.ReplaceAll(str, "\n", ",")
	str = strings.Join(Split(str, ","), ",")
	return str
}

func FormatPrice(price string) string {
	//检查这个字符串是否带有正/负号。如果带有符号，就把符号先单独提取出来
	symbolString := ""
	if price[0] == '-' || price[0] == '+' {
		symbolString = string(price[0])
		price = price[1:]
	}

	//小数点前没有写0，就补一个0进去补齐，让数字字符串看起来更好看
	if price[0] == '.' {
		return "0" + price
	}

	//判断这个数字是不是浮点数值
	dotIndex, decimalString := strings.Index(price, "."), ""
	if dotIndex != -1 {
		decimalString = price[dotIndex:]
		price = price[:dotIndex]
	} else if dotIndex == -1 {
		dotIndex = len(price)
	}

	return fmt.Sprintf("%s%s%s", symbolString, _comma(price[:dotIndex]), decimalString)
}

func _comma(s string) string {
	if len(s) <= 3 {
		return s
	}

	return _comma(s[:len(s)-3]) + "," + _comma(s[len(s)-3:])
}
