package saOss

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/saxon134/go-utils/saData"
	"io"
	"io/ioutil"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

type aliOss struct {
	*oss.Bucket
	UrlRoot string
}

//destination以"/"结尾，则认为是文件夹，会自动生成文件名；
func (m *aliOss) Upload(destination string, reader io.Reader) error {
	if strings.HasSuffix(destination, "/") {
		t := time.Now().UnixNano() / 1000
		r := rand.Intn(10000)
		destination += saData.I64tos(t) + saData.Itos(r)
	}

	err := m.PutObject(destination, reader)
	return err
}

func (m *aliOss) UploadFromLocalFile(destination string, localPath string) error {
	if strings.HasSuffix(destination, "/") {
		t := time.Now().UnixNano() / 1000
		r := rand.Intn(10000)
		destination += saData.I64tos(t) + saData.Itos(r)
	}

	err := m.PutObjectFromFile(destination, localPath)
	return err
}

//支持文件、文件夹删除
func (m *aliOss) Delete(destination string) error {
	if destination == "" {
		return errors.New("path不能空")
	}

	if isObj, _ := m.IsObjectExist(destination); isObj {
		err := m.DeleteObject(destination)
		if err != nil {
			return err
		}
	} else {
		lsRes, err := m.ListObjects(oss.Prefix(destination))
		if err != nil {
			return err
		}

		for _, object := range lsRes.Objects {
			err = m.DeleteObject(object.Key)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *aliOss) SetUrlRoot(root string) {
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

func (m *aliOss) AddUrlRoot(url string) string {
	if m == nil || len(m.UrlRoot) == 0 || len(url) == 0 {
		return url
	}

	if strings.HasPrefix(url, "http") {
		return url
	}

	return strings.TrimSuffix(m.UrlRoot, "/") + "/" + strings.TrimPrefix(url, "/")
}

func (m *aliOss) DeleteUrlRoot(uri string) string {
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

//src为目录，则将src下的内容全部拷贝到destination目录下
//src为文件，如果dest后缀为/，则将src文件拷贝到destination目录下；如果destination后缀不是/，则将src拷贝成dest文件
func (m *aliOss) CopyWithBucket(src, destination string) error {
	if m == nil {
		return errors.New("bucket不存在")
	}

	if src == "" || destination == "" {
		return errors.New("路径不能空")
	}

	var err error
	if isObj, _ := m.IsObjectExist(src); isObj {
		if strings.HasSuffix(destination, "/") {
			_ary := strings.Split(src, "/")
			destination += _ary[len(_ary)-1]
		} else {
			_, err = m.CopyObject(src, destination)
		}

		if err != nil {
			return err
		}
	} else {
		lsRes, err := m.ListObjects(oss.Prefix(src))
		if err != nil {
			return err
		}

		if strings.HasPrefix(destination, "/") == true {
			destination = saData.SubStr(destination, 1, saData.StrLen(destination)-1)
		}

		if strings.HasSuffix(destination, "/") == true {
			destination = saData.SubStr(destination, 0, saData.StrLen(destination)-1)
		}

		if strings.HasPrefix(src, "/") == true {
			src = saData.SubStr(src, 1, saData.StrLen(src)-1)
		}

		if strings.HasSuffix(src, "/") == true {
			src = saData.SubStr(src, 0, saData.StrLen(src)-1)
		}

		for _, object := range lsRes.Objects {
			if isObj, _ = m.IsObjectExist(object.Key); isObj {
				f_name := saData.SubStr(object.Key, saData.StrLen(src), saData.StrLen(object.Key)-saData.StrLen(src))
				_, err = m.CopyObject(object.Key, destination+f_name)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m *aliOss) GetTxt(uri string) (res string, err error) {
	// 下载文件到流。
	body, err := m.GetObject(uri)
	if err != nil {
		return "", err
	}
	defer body.Close()

	v, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(v), nil
}

func (m *aliOss) UploadTxt(destination string, v string) (url string, err error) {
	if destination == "" || v == "" {
		return "", errors.New("缺参数")
	}

	if strings.HasSuffix(destination, "/") {
		t := time.Now().UnixNano() / 1000
		r := rand.Intn(10000)
		destination += saData.I64tos(t) + saData.Itos(r)
	}

	err = m.PutObject(destination, strings.NewReader(v))
	return destination, err
}
