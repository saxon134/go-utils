package saOss

import (
	"errors"
)

func InitOss(provider Provider, region string, endpoint string, accessKeyId string, accessKeySecret string, bucket string) (SaOss, error) {
	if len(endpoint) == 0 || len(accessKeyId) == 0 || len(accessKeySecret) == 0 || len(bucket) == 0 {
		return nil, errors.New("oss配置有误")
	}

	switch provider {
	case ProviderAliyun:
		return newAliyunOss(region, endpoint, accessKeyId, accessKeySecret, bucket)
	case ProviderTos:
		if region == "" {
			return nil, errors.New("oss配置有误")
		}
		return newTosOss(region, endpoint, accessKeyId, accessKeySecret, bucket)
	default:
		return nil, errors.New("oss provider不支持")
	}
}
