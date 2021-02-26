package saWx

import (
	"errors"
	"github.com/saxon134/go-utils/saOss"
	"github.com/saxon134/go-utils/saRedis"
)

var Gzh GzhServer
var Xcx XcxServer
var Pay PayServer

type Conf struct {
	Redis *saRedis.Redis
	Oss   saOss.SaOss
}

func Init(conf *Conf) error {
	if conf.Redis == nil || conf.Oss == nil {
		return errors.New("saWx初始化失败")
	}

	return nil
}
