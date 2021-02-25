/**
微信错误码（转换后的）
code：
0，		正常
1,      access token 有误
2,		其他错误
1000,   敏感信息校验不通过
1100,   链接被封
*/
package saWx

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saHttp"
	"github.com/saxon134/go-utils/saWx/model"
	"os"
	"strconv"
	"strings"
	"time"
)

func Code2Session(appId string, key string, code string) (sessionKey string, openid string, uniId string, err error) {
	if appId == "" || key == "" || code == "" {
		return "", "", "", errors.New("缺少参数")
	}

	params := map[string]string{
		"appid":      appId,
		"secret":     key,
		"grant_type": "authorization_code",
		"js_code":    code,
	}

	var resStr string
	resStr, err = saHttp.Get("https://api.weixin.qq.com/sns/jscode2session", params)
	if resStr == "" || err != nil {
		return "", "", "", err
	}
	dic, _ := saData.JsonToMap(resStr)

	if errCode, _ := saData.ToInt(dic["errcode"]); errCode == 0 {
		openid, _ = saData.ToStr(dic["openid"])
		sessionKey, _ = saData.ToStr(dic["session_key"])
		uniId, _ = saData.ToStr(dic["unionid"])
		if openid != "" && sessionKey != "" {
			return sessionKey, openid, uniId, nil
		} else {
			return "", "", "", errors.New("微信返回数据空")
		}
	} else {
		errMsg, ok := dic["errmsg"].(string)
		if ok == false || errMsg == "" {
			errMsg = "微信接口报错"
		}
		return "", "", "", errors.New(errMsg)
	}
}

func AccessToken(appId string, key string, ignoreCache bool) string {
	var token string

	if ignoreCache == false {
		tk, err := _redis.Get("accessToken" + appId)
		if err != nil {
			return ""
		}

		if v, _ := tk.(string); len(v) > 0 {
			token = v
		}
	}

	if len(token) == 0 {
		var resStr string
		var err error
		url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appId, key)
		resStr, err = saHttp.Get(url, nil)
		if resStr == "" || err != nil {
			return ""
		}

		dic, _ := saData.JsonToMap(resStr)
		if dic == nil || len(dic) == 0 {
			return ""
		}

		token, _ = saData.ToStr(dic["access_token"])
		if len(token) == 0 {
			return ""
		}

		expiresIn, _ := saData.ToInt(dic["expires_in"])
		if expiresIn <= 0 {
			expiresIn = 7200
		}

		_ = _redis.Set("accessToken"+appId, token, time.Second*time.Duration(int64(expiresIn)))
	}

	return token
}

func DecryptWxOpenData(appId, sessionKey, encryptData, iv string) (*map[string]interface{}, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return nil, err
	}

	sessionKeyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	dataBytes, err := AesDecrypt(decodeBytes, sessionKeyBytes, ivBytes)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(dataBytes, &m)
	if err != nil {
		return nil, err
	}

	temp := m["watermark"].(map[string]interface{})
	aid := temp["appid"].(string)
	if aid != appId {
		return nil, fmt.Errorf("appId不匹配")
	}

	return &m, nil
}

/**
code：
0，		正常
1,      access token 有误
2,		其他错误
1000，	内容审核不通过
*/
func WxTxtSecCheck(txt string, accessToken string) (code int) {
	if txt != "" {
		if accessToken == "" {
			return 1
		}

		params := map[string]interface{}{"content": txt}
		if res, _, err := saHttp.PostJson("https://api.weixin.qq.com/wxa/msg_sec_check?access_token="+accessToken, params); err == nil {
			if dic, err := saData.JsonToMap(res); err == nil {
				if errcode, err := saData.ToInt(dic["errcode"]); err == nil {
					if errcode == 0 {
						return 0
					} else if errcode == 87014 {
						return 1000
					} else if errcode == 42001 {
						return 1
					}
				}
			}
		}
	}
	return 2
}

func FetchQrCode(accessToken string, params model.QrCoder) (fileUri string, err error) {
	if accessToken == "" {
		return "", errors.New("缺少参数")
	}

	urlStr := "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=" + accessToken
	response, contentType, err := saHttp.PostJson(urlStr, params)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(contentType, "application/json") {
		// 返回错误信息
		var result model.WxError
		err = json.Unmarshal([]byte(response), &result)
		if err == nil && result.ErrCode != 0 {
			err = fmt.Errorf("fetchCode error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
			return "", err
		}
	} else if contentType == "image/jpeg" {
		// 返回文件
		f_name := strconv.FormatInt(time.Now().UnixNano(), 10)

		if fw, err := os.Create(f_name); err == nil {
			if _, err := fw.Write([]byte(response)); err == nil {
				ossPath := "images/tmp/wx/qrcode/" + saData.TimeStr(time.Now(), saData.TimeFormate_dayStr) + "/" + f_name + ".jpg"
				if err := _oss.UploadFromLocalFile(ossPath, f_name); err == nil {
					return ossPath, nil
				}
			}
		}
		_ = os.Remove(f_name)
		return "", errors.New("图片转存失败")
	}

	err = fmt.Errorf("fetchCode error : unknown response content type - %v", contentType)
	return "", err
}

/**
code：
0，		正常
1,      access token 有误
2,		其他错误
1100,   链接被封
*/
func WxShortUri(uri string, accessToken string) (code int, shortUri string) {
	if uri != "" && strings.HasPrefix(uri, "http") {
		if accessToken == "" {
			return 1, ""
		}

		params := map[string]interface{}{"long_url": uri, "action": "long2short"}
		if res, _, err := saHttp.PostJson("https://api.weixin.qq.com/cgi-bin/shorturl?access_token="+accessToken, params); err == nil {
			if dic, err := saData.JsonToMap(res); err == nil {
				if errcode, err := saData.ToInt(dic["errcode"]); err == nil {
					if errcode == 0 {
						shortUri, _ = saData.ToStr(dic["short_url"])
						return 0, shortUri
					} else if errcode == 42001 {
						return 1, ""
					} else if errcode == 44001 { //todo
						return 1100, ""
					} else {
						return 2, ""
					}
				}
			}
		}
	}
	return 2, ""
}
