package saOss

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/saxon134/go-utils/saData"
	"io"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
)

type ctOss struct {
	*s3.S3
	Bucket  string
	UrlRoot string
}

// Upload destination以"/"结尾，则认为是文件夹，会自动生成文件名；
func (m *ctOss) Upload(destination string, reader io.Reader) error {
	if strings.HasSuffix(destination, "/") {
		t := time.Now().Unix()
		r := rand.Intn(10000)
		destination += saData.I64tos(t) + saData.Itos(r)
	}

	var bAry, err = io.ReadAll(reader)
	if err != nil {
		return err
	}

	_, err = m.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(destination),
		Body:   bytes.NewReader(bAry),
	})
	return err
}

func (m *ctOss) UploadFromLocalFile(destination string, localPath string) error {
	if strings.HasSuffix(destination, "/") {
		t := time.Now().Unix()
		r := rand.Intn(10000)
		destination += saData.I64tos(t) + saData.Itos(r)
	}

	var bAry, err = os.ReadFile(localPath)
	if err != nil {
		return err
	}

	_, err = m.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(destination),
		Body:   bytes.NewReader(bAry),
	})
	return err
}

// Delete 支持文件、文件夹删除
func (m *ctOss) Delete(destination string) error {
	if destination == "" {
		return errors.New("path不能空")
	}

	_, err := m.S3.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(m.Bucket), Key: aws.String(destination)})
	return err
}

func (m *ctOss) SetUrlRoot(root string) {
	if len(root) == 0 {
		return
	}

	if strings.HasPrefix(root, "http") == false {
		return
	}

	if strings.HasSuffix(root, "/") == false {
		m.UrlRoot = root + "/"
	} else {
		m.UrlRoot = root
	}
}

func (m *ctOss) AddUrlRoot(url string) string {
	if m == nil || len(m.UrlRoot) == 0 || len(url) == 0 {
		return url
	}

	if strings.HasPrefix(url, "http") {
		return url
	}

	return strings.TrimSuffix(m.UrlRoot, "/") + "/" + strings.TrimPrefix(url, "/")
}

func (m *ctOss) DeleteUrlRoot(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return uri
	}

	root := u.Scheme + "://" + u.Host + "/"

	if root == m.UrlRoot {
		return strings.Replace(uri, root, "", 1)
	}

	return uri
}

// CopyWithBucket
// src为目录，则将src下的内容全部拷贝到destination目录下
// src为文件，如果dest后缀为/，则将src文件拷贝到destination目录下；如果destination后缀不是/，则将src拷贝成dest文件
func (m *ctOss) CopyWithBucket(src, destination string) error {
	if m == nil {
		return errors.New("bucket不存在")
	}

	if src == "" || destination == "" {
		return errors.New("路径不能空")
	}

	var err error
	_, err = m.S3.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(m.Bucket),
		CopySource: aws.String(src),
		Key:        aws.String(destination),
	})
	return err
}

func (m *ctOss) GetTxt(uri string) (res string, err error) {
	// 下载文件到流。
	var out *s3.GetObjectOutput
	out, err = m.S3.GetObject(&s3.GetObjectInput{Bucket: aws.String(m.Bucket), Key: aws.String(uri)})
	if err != nil {
		return "", err
	}

	var v []byte
	v, err = io.ReadAll(out.Body)
	if err != nil {
		return "", err
	}

	return string(v), nil
}

func (m *ctOss) UploadTxt(destination string, v string) (path string, err error) {
	if destination == "" || v == "" {
		return "", errors.New("缺参数")
	}

	if strings.HasSuffix(destination, "/") {
		t := time.Now().Unix()
		r := rand.Intn(10000)
		destination += saData.I64tos(t) + saData.Itos(r)
	}

	_, err = m.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(m.Bucket),
		Key:    aws.String(destination),
		Body:   strings.NewReader(v),
	})

	return destination, err
}
