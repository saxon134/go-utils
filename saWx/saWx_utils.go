package saWx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"gitee.com/go-utils/saData"
	"gitee.com/go-utils/saLog"
	"github.com/nanjishidu/gomini/gocrypto"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

func AesDecrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)

	retData := make([]byte, 0, len(crypted))
	//获取的数据尾端有'/x0e'占位符,去除它
	for i, ch := range origData {
		if ch > '\x1f' && ch != '\x7f' {
			retData = append(retData, origData[i])
		}
	}
	return retData, nil
}

func WxSign(dic *map[string]string, appKey string) string {
	var keys []string
	for k := range *dic {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedStr := ""
	for _, k := range keys {
		v := (*dic)[k]
		if v != "" {
			sortedStr += k + "=" + (*dic)[k] + "&"
		}
	}
	sortedStr += "key=" + appKey
	sign := saData.Md5(sortedStr, true)
	sign = strings.ToUpper(sign)
	return sign
}

func Aes256Decode(mchKey string, info string) string {
	b, err := base64.StdEncoding.DecodeString(info)
	if err != nil {
		saLog.Err(b, err)
		return ""
	}

	_ = gocrypto.SetAesKey(strings.ToLower(gocrypto.Md5(mchKey)))
	plaintext, err := gocrypto.AesECBDecrypt(b)
	if err != nil {
		saLog.Err(err)
		return ""
	}
	return string(plaintext)
}

func KeyHttpsPost(url string, contentType string, body io.Reader) (*http.Response, error) {
	var wechatPayCert = ("./cert/cert.pem")
	var wechatPayKey = ("./cert/key.pem")
	var rootCa = ("./cert/rootCa.pem")
	var tr *http.Transport
	// 微信提供的API证书,证书和证书密钥 .pem格式
	if certs, err := tls.LoadX509KeyPair(wechatPayCert, wechatPayKey); err == nil {
		// 微信支付HTTPS服务器证书的根证书  .pem格式
		if rootCa, err := ioutil.ReadFile(rootCa); err == nil {
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(rootCa)

			tr = &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      pool,
					Certificates: []tls.Certificate{certs},
				},
			}

			client := &http.Client{Transport: tr}
			return client.Post(url, contentType, body)
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
