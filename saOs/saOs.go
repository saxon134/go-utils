package saOs

import (
	"archive/zip"
	"errors"
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"log"
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

func Unzip(fileFullPath string) error {
	zipFile, err := zip.OpenReader(fileFullPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 第二步，遍历 zip 中的文件
	for _, f := range zipFile.File {
		filePath := f.Name
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		// 创建对应文件夹
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}
		// 解压到的目标文件
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		file, err := f.Open()
		if err != nil {
			return err
		}
		// 写入到解压到的目标文件
		if _, err := io.Copy(dstFile, file); err != nil {
			return err
		}
		dstFile.Close()
		file.Close()
	}
	return nil
}

func Zip(dir string, zipName string) error {
	zipFile, err := os.OpenFile(zipName, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	archive := zip.NewWriter(zipFile)
	defer archive.Close()
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if path == dir {
			return nil
		}

		info, _ := d.Info()
		h, _ := zip.FileInfoHeader(info)
		h.Name = strings.TrimPrefix(path, dir)
		h.Name = strings.TrimPrefix(h.Name, "/")
		h.Name = strings.TrimPrefix(h.Name, "\\")
		if info.IsDir() {
			h.Name += "/"
		} else {
			h.Method = zip.Deflate
		}
		writer, _ := archive.CreateHeader(h)
		if !info.IsDir() {
			srcFile, _ := os.Open(path)
			defer srcFile.Close()
			io.Copy(writer, srcFile)
		}
		return nil
	})
}

func ResizeImg(src string, width uint, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}

	img, err := png.Decode(file)
	if err != nil {
		return err
	}
	file.Close()

	m := resize.Resize(width, 0, img, resize.Lanczos3)
	out, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	jpeg.Encode(out, m, nil)
	return nil
}
