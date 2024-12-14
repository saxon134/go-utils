package saOss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
)

func InitOss(region string, endpoint string, accessKeyId string, accessKeySecret string, bucket string) (*SaOss, error) {
	if len(endpoint) == 0 || len(accessKeyId) == 0 || len(accessKeySecret) == 0 || len(bucket) == 0 {
		return nil, errors.New("oss配置有误")
	}

	var err error
	var client *oss.Client

	client, err = oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}

	var aliOss SaOss
	aliOss.Bucket, err = client.Bucket(bucket)
	if err != nil {
		return nil, err
	}

	aliOss.region = region
	aliOss.endpoint = endpoint
	aliOss.accessKeyId = accessKeyId
	aliOss.accessKeySecret = accessKeySecret
	aliOss.bucket = bucket
	return &aliOss, nil
}
