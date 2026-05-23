package saOss

import (
	"io"
)

type Provider string

const (
	ProviderAliyun Provider = "aliyun"
	ProviderTos    Provider = "tos"
)

type SaOss interface {
	Upload(destination string, reader io.Reader, options ...interface{}) error
	UploadFromLocalFile(destination string, localPath string, options ...interface{}) error
	Delete(destination string) error
	SetUrlRoot(root string)
	AddUrlRoot(url string) string
	DeleteUrlRoot(uri string) string
	CopyWithBucket(src, destination string) error
	GetTxt(uri string) (res string, err error)
	UploadTxt(destination string, v string, options ...interface{}) (path string, err error)
	StsToken(roleArn, roleSessionName string) (keyId, keySecret, token string, err error)
}

type TageKey map[string]string

var TageDeleteDay = TageKey{"delete-day": "1"}
var TageDeleteWeek = TageKey{"delete-week": "7"}
var TageDeleteMonth = TageKey{"delete-month": "30"}
var TageDeleteQuarter = TageKey{"delete-quarter": "90"}
var TageDeleteYear = TageKey{"delete-year": "365"}

type uploadOptions struct {
	Tags              TageKey
	ObjectExpiresDays int64
}

type objectProvider interface {
	PutObject(destination string, reader io.Reader, options uploadOptions) error
	PutObjectFromFile(destination string, localPath string, options uploadOptions) error
	IsObjectExist(key string) (bool, error)
	DeleteObject(key string) error
	ListObjects(prefix string) ([]string, error)
	CopyObject(src, destination string) error
	GetObject(key string) (io.ReadCloser, error)
	StsToken(roleArn, roleSessionName string) (keyId, keySecret, token string, err error)
}

type ossClient struct {
	provider objectProvider
	urlRoot  string
}
