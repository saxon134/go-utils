package saData

import (
	"path/filepath"
	"runtime"
	"strings"
)

// 获取绝对路径
func AbsPath(path string) string {
	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "/", "\\")
		path = strings.ReplaceAll(path, `\\`, `\`)
		path = strings.ReplaceAll(path, `\\`, `\`)
		path = strings.ReplaceAll(path, `\\`, `\`)
	} else {
		path = strings.ReplaceAll(path, "\\", "/")
		path = strings.ReplaceAll(path, `//`, `/`)
		path = strings.ReplaceAll(path, `//`, `/`)
		path = strings.ReplaceAll(path, `//`, `/`)
	}
	var s, _ = filepath.Abs(path)
	return s
}
