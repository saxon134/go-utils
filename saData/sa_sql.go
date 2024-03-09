package saData

import "strings"

// SQLColumn 删除字符串里有影响SQL的字符
func SQLColumn(str string) string {
	str = strings.ReplaceAll(str, "'", "")
	str = strings.ReplaceAll(str, "\"", "")
	str = strings.ReplaceAll(str, "`", "")
	return str
}
