package saData

import (
	"strings"
)

// SQL转义
func SQLEspace(s string) string {
	//s = strings.ReplaceAll(s, `'`, ``)
	s = strings.ReplaceAll(s, `'`, `\'`)

	//尾部有多个转义字符，就得加多少个转义字符
	if s != "" {
		var count = 0
		var end = StrLen(s)
		for i := end; i > 0; i-- {
			var c = SubStr(s, i-1, 1)
			if c == `\` {
				count++
			} else {
				break
			}
		}

		for i := 0; i < count; i++ {
			s += `\`
		}
	}
	return s
}
