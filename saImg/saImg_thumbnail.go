package saImg

import (
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
)

// 计算图片缩放后的尺寸
func calculateRatioFit(srcWidth, srcHeight, maxWidth int) (int, int) {
	if maxWidth < 0 {
		maxWidth = 400
	}

	if srcWidth < maxWidth {
		return srcWidth, srcHeight
	} else {
		ratio := float64(maxWidth) / float64(srcWidth)
		return int(math.Ceil(float64(srcWidth) * ratio)), int(math.Ceil(float64(srcHeight) * ratio))
	}
}

// 生成缩略图
func MakeThumbnail(imagePath, savePath string, maxWidth int) error {
	file, _ := os.Open(imagePath)
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	b := img.Bounds()
	width := b.Max.X
	height := b.Max.Y

	w, h := calculateRatioFit(width, height, maxWidth)

	// 调用resize库进行图片缩放
	m := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)

	// 需要保存的文件
	imgfile, _ := os.Create(savePath)
	defer imgfile.Close()

	// 以PNG格式保存文件
	err = png.Encode(imgfile, m)
	if err != nil {
		return err
	}

	return nil
}
