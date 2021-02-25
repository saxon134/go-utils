package saWx

import (
	"errors"
	"gitee.com/go-utils/saOss"
	"gitee.com/go-utils/saRedis"
)

var MApp MAppServer
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
