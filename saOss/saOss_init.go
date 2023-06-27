package saOss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"net/http"
)

func InitOss(ossType OssType, endpoint string, accessKeyId string, accessKeySecret string, bucket string) (SaOss, error) {
	if len(endpoint) == 0 || len(accessKeyId) == 0 || len(accessKeySecret) == 0 || len(bucket) == 0 {
		return nil, errors.New("oss配置有误")
	}

	var err error
	if ossType == AliOssType {
		var client *oss.Client

		client, err = oss.New(endpoint, accessKeyId, accessKeySecret)
		if err != nil {
			return nil, err
		}

		var aliOss aliOss
		aliOss.Bucket, err = client.Bucket(bucket)
		if err != nil {
			return nil, err
		}

		return &aliOss, nil
	} else if ossType == CtOssType {
		region := "cn"
		config := &aws.Config{
			Credentials:      credentials.NewStaticCredentials(accessKeyId, accessKeySecret, ""),
			Endpoint:         aws.String(endpoint),
			S3ForcePathStyle: aws.Bool(true),
			DisableSSL:       aws.Bool(true),
			LogLevel:         aws.LogLevel(aws.LogDebug),
			Region:           aws.String(region),
			HTTPClient: &http.Client{
				Transport: &http.Transport{
					MaxConnsPerHost: 1000,
				},
			},
		}

		var ctOss ctOss
		ctOss.S3 = s3.New(session.Must(session.NewSession(config)))
		ctOss.Bucket = bucket
		return &ctOss, nil
	} else {
		return nil, errors.New("暂只支持阿里云、天翼云oss")
	}
}
