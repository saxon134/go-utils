package saOss

import (
	"gitee.com/go-utils/saError"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
)

func InitOss(ossType OssType, endpoint string, accessKeyId string, accessKeySecret string, bucket string) (SaOss, error) {
	if ossType != AliOssType {
		return nil, errors.New("暂只支持阿里云oss")
	}

	if len(endpoint) == 0 || len(accessKeyId) == 0 || len(accessKeySecret) == 0 || len(bucket) == 0 {
		return nil, saError.StackError("oss配置有误")
	}

	var err error
	var client *oss.Client

	client, err = oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, saError.StackError(err)
	}

	var aliOss aliOss
	aliOss.Bucket, err = client.Bucket(bucket)
	if err != nil {
		return nil, saError.StackError(err)
	}

	return &aliOss, nil
}
