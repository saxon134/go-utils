package saOss

import "io"

type OssType int8

const (
	NullOssType OssType = iota
	AliOssType
)

type SaOss interface {
	Upload(destination string, reader io.Reader) error
	UploadFromLocalFile(destination string, localPath string) error
	Delete(destination string) error
	SetUrlRoot(root string)
	AddUrlRoot(url string) string
	DeleteUrlRoot(url string) string
	CopyWithBucket(src, destination string) error
	GetTxt(uri string) (res string, err error)
	//destination已"/"结尾，则会自动加上随机名称，否则直接按照全路径保存
	UploadTxt(destination string, v string) (url string, err error)
}
