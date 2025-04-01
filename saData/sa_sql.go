package saData

import (
	"strings"
)

// SQL转义
func SQLEspace(s string) string {
	s = strings.ReplaceAll(s, `'`, `\'`)
	if strings.HasSuffix(s, `\`) {
		s += `\`
	}
	return s
}
