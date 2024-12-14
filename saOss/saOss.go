package saOss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type SaOss struct {
	*oss.Bucket
	UrlRoot string

	region          string
	endpoint        string
	accessKeyId     string
	accessKeySecret string
	bucket          string
}
