package saImg

func GenPoster() {
	//var (
	//	err       error
	//	bgFile    *os.File
	//	bgImg     image.Image
	//	qrCodeImg image.Image
	//	offset    image.Point
	//)
	//
	//// 01: 打开背景图片
	//bgFile, err = os.Open("./bg.png")
	//if err != nil {
	//	fmt.Println("打开背景图片失败", err)
	//	return
	//}
	//defer bgFile.Close()
	//
	//// 02: 编码为图片格式
	//bgImg, err = png.Decode(bgFile)
	//if err != nil {
	//	fmt.Println("背景图片编码失败:", err)
	//	return
	//}
	//
	//// 03: 生成二维码
	//qrCodeImg, err = createAvatar()
	//if err != nil {
	//	fmt.Println("生成二维码失败:", err)
	//	return
	//}
	//
	//offset = image.Pt(426, 475)
	//
	//b := bgImg.Bounds()
	//
	//m := image.NewRGBA(b)
	//
	//draw.Draw(m, b, bgImg, image.Point{X: 0, Y: 0}, draw.Src)
	//
	//draw.Draw(m, qrCodeImg.Bounds().Add(offset), qrCodeImg, image.Point{X: 0, Y: 0}, draw.Over)
	//
	//// 上传至oss时这段要改
	//i, _ := os.Create(path.Base("a2.png"))
	//
	//_ = png.Encode(i, m)
	//
	//defer i.Close()

}
