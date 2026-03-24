package saData

import (
	"errors"
	"io"
	"os"
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

// 文件排序
type SortFileByTime []os.DirEntry

func (f SortFileByTime) Less(i, j int) bool {
	var if1, _ = f[i].Info()
	var if2, _ = f[j].Info()
	if if1 != nil && if2 != nil {
		return if1.ModTime().After(if2.ModTime())
	}
	return false
}
func (f SortFileByTime) Len() int      { return len(f) }
func (f SortFileByTime) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

// 实现跨卷/分区移动文件
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return errors.New("Couldn't open source file: " + err.Error())
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return errors.New("Couldn't open dest file: " + err.Error())
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return errors.New("Writing to output file failed: " + err.Error())
	}

	err = os.Remove(sourcePath)
	if err != nil {
		return errors.New("Failed removing original file: " + err.Error())
	}
	return nil
}
