package saOss

import (
	"errors"
	"io"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/saxon134/go-utils/saData"
)

func (m *ossClient) Upload(destination string, reader io.Reader, options ...interface{}) error {
	if err := m.checkProvider(); err != nil {
		return err
	}

	destination = buildDestination(destination)
	return m.provider.PutObject(destination, reader, parseUploadOptions(options))
}

func (m *ossClient) UploadFromLocalFile(destination string, localPath string, options ...interface{}) error {
	if err := m.checkProvider(); err != nil {
		return err
	}

	destination = buildDestination(destination)
	return m.provider.PutObjectFromFile(destination, localPath, parseUploadOptions(options))
}

func (m *ossClient) Delete(destination string) error {
	if err := m.checkProvider(); err != nil {
		return err
	}
	if destination == "" {
		return errors.New("path不能空")
	}

	if isObj, _ := m.provider.IsObjectExist(destination); isObj {
		return m.provider.DeleteObject(destination)
	}

	objects, err := m.provider.ListObjects(destination)
	if err != nil {
		return err
	}
	for _, object := range objects {
		if err = m.provider.DeleteObject(object); err != nil {
			return err
		}
	}
	return nil
}

func (m *ossClient) SetUrlRoot(root string) {
	if len(root) == 0 || !strings.HasPrefix(root, "http") {
		return
	}

	if strings.HasSuffix(root, "/") {
		m.urlRoot = root
		return
	}
	m.urlRoot = root + "/"
}

func (m *ossClient) AddUrlRoot(uri string) string {
	if m == nil || len(m.urlRoot) == 0 || len(uri) == 0 {
		return uri
	}
	if strings.HasPrefix(uri, "http") {
		return uri
	}
	return strings.TrimSuffix(m.urlRoot, "/") + "/" + strings.TrimPrefix(uri, "/")
}

func (m *ossClient) DeleteUrlRoot(uri string) string {
	if m == nil {
		return uri
	}
	if m.urlRoot != "" && strings.HasPrefix(uri, m.urlRoot) {
		return strings.Replace(uri, m.urlRoot, "", 1)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return uri
	}
	root := u.Scheme + "://" + u.Host + "/"
	if root == m.urlRoot {
		return strings.Replace(uri, root, "", 1)
	}
	return uri
}

func (m *ossClient) CopyWithBucket(src, destination string) error {
	if err := m.checkProvider(); err != nil {
		return err
	}
	if src == "" || destination == "" {
		return errors.New("路径不能空")
	}

	if isObj, _ := m.provider.IsObjectExist(src); isObj {
		if strings.HasSuffix(destination, "/") {
			parts := strings.Split(src, "/")
			destination += parts[len(parts)-1]
		}
		return m.provider.CopyObject(src, destination)
	}

	objects, err := m.provider.ListObjects(src)
	if err != nil {
		return err
	}

	src = trimPathSlash(src)
	destination = trimPathSlash(destination)
	for _, object := range objects {
		if isObj, _ := m.provider.IsObjectExist(object); !isObj {
			continue
		}

		fileName := strings.TrimPrefix(object, src)
		if err = m.provider.CopyObject(object, destination+fileName); err != nil {
			return err
		}
	}
	return nil
}

func (m *ossClient) GetTxt(uri string) (res string, err error) {
	if err = m.checkProvider(); err != nil {
		return "", err
	}

	body, err := m.provider.GetObject(uri)
	if err != nil {
		return "", err
	}
	defer body.Close()

	v, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (m *ossClient) UploadTxt(destination string, v string, options ...interface{}) (path string, err error) {
	if err = m.checkProvider(); err != nil {
		return "", err
	}
	if destination == "" || v == "" {
		return "", errors.New("缺参数")
	}

	destination = buildDestination(destination)
	err = m.provider.PutObject(destination, strings.NewReader(v), parseUploadOptions(options))
	return destination, err
}

func (m *ossClient) StsToken(roleArn, roleSessionName string) (keyId, keySecret, token string, err error) {
	if err = m.checkProvider(); err != nil {
		return "", "", "", err
	}
	return m.provider.StsToken(roleArn, roleSessionName)
}

func (m *ossClient) checkProvider() error {
	if m == nil || m.provider == nil {
		return errors.New("bucket不存在")
	}
	return nil
}

func buildDestination(destination string) string {
	if strings.HasSuffix(destination, "/") {
		t := time.Now().Unix()
		r := rand.Intn(10000)
		destination += saData.I64tos(t) + saData.Itos(r)
	}
	return destination
}

func trimPathSlash(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return path
}

func parseUploadOptions(options []interface{}) uploadOptions {
	res := uploadOptions{Tags: TageKey{}}
	for _, v := range options {
		tags, ok := v.(TageKey)
		if !ok || len(tags) == 0 {
			continue
		}

		if res.ObjectExpiresDays == 0 {
			res.ObjectExpiresDays = parseFirstExpiresDays(tags)
		}
		for k, value := range tags {
			res.Tags[k] = value
		}
	}

	if len(res.Tags) == 0 {
		res.Tags = nil
	}
	return res
}

func parseFirstExpiresDays(tags TageKey) int64 {
	keys := make([]string, 0, len(tags))
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		days, err := strconv.ParseInt(tags[key], 10, 64)
		if err == nil && days > 0 {
			return days
		}
	}
	return 0
}
