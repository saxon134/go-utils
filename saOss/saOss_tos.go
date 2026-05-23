package saOss

import (
	"context"
	"errors"
	"io"

	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

type tosProvider struct {
	client *tos.ClientV2
	bucket string
}

func newTosOss(region string, endpoint string, accessKeyId string, accessKeySecret string, bucket string) (SaOss, error) {
	client, err := tos.NewClientV2(
		endpoint,
		tos.WithRegion(region),
		tos.WithCredentials(tos.NewStaticCredentials(accessKeyId, accessKeySecret)),
	)
	if err != nil {
		return nil, err
	}

	return &ossClient{
		provider: &tosProvider{
			client: client,
			bucket: bucket,
		},
	}, nil
}

func (m *tosProvider) PutObject(destination string, reader io.Reader, options uploadOptions) error {
	_, err := m.client.PutObjectV2(context.Background(), &tos.PutObjectV2Input{
		PutObjectBasicInput: tosPutObjectBasicInput(m.bucket, destination, options),
		Content:             reader,
	})
	return err
}

func (m *tosProvider) PutObjectFromFile(destination string, localPath string, options uploadOptions) error {
	_, err := m.client.PutObjectFromFile(context.Background(), &tos.PutObjectFromFileInput{
		PutObjectBasicInput: tosPutObjectBasicInput(m.bucket, destination, options),
		FilePath:            localPath,
	})
	return err
}

func (m *tosProvider) IsObjectExist(key string) (bool, error) {
	return m.client.DoesObjectExist(context.Background(), &tos.DoesObjectExistInput{
		Bucket: m.bucket,
		Key:    key,
	})
}

func (m *tosProvider) DeleteObject(key string) error {
	_, err := m.client.DeleteObjectV2(context.Background(), &tos.DeleteObjectV2Input{
		Bucket: m.bucket,
		Key:    key,
	})
	return err
}

func (m *tosProvider) ListObjects(prefix string) ([]string, error) {
	lsRes, err := m.client.ListObjectsType2(context.Background(), &tos.ListObjectsType2Input{
		Bucket: m.bucket,
		Prefix: prefix,
	})
	if err != nil {
		return nil, err
	}

	objects := make([]string, 0, len(lsRes.Contents))
	for _, object := range lsRes.Contents {
		objects = append(objects, object.Key)
	}
	return objects, nil
}

func (m *tosProvider) CopyObject(src, destination string) error {
	_, err := m.client.CopyObject(context.Background(), &tos.CopyObjectInput{
		Bucket:    m.bucket,
		Key:       destination,
		SrcBucket: m.bucket,
		SrcKey:    src,
	})
	return err
}

func (m *tosProvider) GetObject(key string) (io.ReadCloser, error) {
	output, err := m.client.GetObjectV2(context.Background(), &tos.GetObjectV2Input{
		Bucket: m.bucket,
		Key:    key,
	})
	if err != nil {
		return nil, err
	}
	return output.Content, nil
}

func (m *tosProvider) StsToken(roleArn, roleSessionName string) (keyId, keySecret, token string, err error) {
	return "", "", "", errors.New("当前provider不支持StsToken")
}

func tosPutObjectBasicInput(bucket string, destination string, options uploadOptions) tos.PutObjectBasicInput {
	return tos.PutObjectBasicInput{
		Bucket:        bucket,
		Key:           destination,
		ObjectExpires: options.ObjectExpiresDays,
	}
}
