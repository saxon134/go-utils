package saOss

import (
	"io"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type aliyunProvider struct {
	bucket          *oss.Bucket
	region          string
	accessKeyId     string
	accessKeySecret string
}

func newAliyunOss(region string, endpoint string, accessKeyId string, accessKeySecret string, bucket string) (SaOss, error) {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}

	ossBucket, err := client.Bucket(bucket)
	if err != nil {
		return nil, err
	}

	return &ossClient{
		provider: &aliyunProvider{
			bucket:          ossBucket,
			region:          region,
			accessKeyId:     accessKeyId,
			accessKeySecret: accessKeySecret,
		},
	}, nil
}

func (m *aliyunProvider) PutObject(destination string, reader io.Reader, options uploadOptions) error {
	return m.bucket.PutObject(destination, reader, aliyunOptions(options)...)
}

func (m *aliyunProvider) PutObjectFromFile(destination string, localPath string, options uploadOptions) error {
	return m.bucket.PutObjectFromFile(destination, localPath, aliyunOptions(options)...)
}

func (m *aliyunProvider) IsObjectExist(key string) (bool, error) {
	return m.bucket.IsObjectExist(key)
}

func (m *aliyunProvider) DeleteObject(key string) error {
	return m.bucket.DeleteObject(key)
}

func (m *aliyunProvider) ListObjects(prefix string) ([]string, error) {
	lsRes, err := m.bucket.ListObjects(oss.Prefix(prefix))
	if err != nil {
		return nil, err
	}

	objects := make([]string, 0, len(lsRes.Objects))
	for _, object := range lsRes.Objects {
		objects = append(objects, object.Key)
	}
	return objects, nil
}

func (m *aliyunProvider) CopyObject(src, destination string) error {
	_, err := m.bucket.CopyObject(src, destination)
	return err
}

func (m *aliyunProvider) GetObject(key string) (io.ReadCloser, error) {
	return m.bucket.GetObject(key)
}

func (m *aliyunProvider) StsToken(roleArn, roleSessionName string) (keyId, keySecret, token string, err error) {
	client, err := sts.NewClientWithAccessKey(m.region, m.accessKeyId, m.accessKeySecret)
	if err != nil {
		return "", "", "", err
	}

	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = roleArn
	request.RoleSessionName = roleSessionName

	response, err := client.AssumeRole(request)
	if err != nil {
		return "", "", "", err
	}

	return response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, response.Credentials.SecurityToken, nil
}

func aliyunOptions(options uploadOptions) []oss.Option {
	if len(options.Tags) == 0 {
		return nil
	}

	tags := make([]oss.Tag, 0, len(options.Tags))
	for tagK, tagV := range options.Tags {
		tags = append(tags, oss.Tag{Key: tagK, Value: tagV})
	}
	return []oss.Option{oss.SetTagging(oss.Tagging{Tags: tags})}
}
