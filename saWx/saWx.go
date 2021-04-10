package saWx

import (
	"errors"
	"github.com/saxon134/go-utils/saOss"
	"github.com/saxon134/go-utils/saRedis"
	//"github.com/silenceper/wechat/v2"
	//"github.com/silenceper/wechat/v2/cache"
	//"github.com/silenceper/wechat/v2/openplatform"
	//"github.com/silenceper/wechat/v2/openplatform/config"
)

var Gzh GzhServer
var Xcx XcxServer
var Pay PayServer

type Conf struct {
	Redis *saRedis.Redis
	Oss   saOss.SaOss
}

func Init(conf *Conf) error {
	if conf.Redis == nil || conf.Oss == nil {
		return errors.New("saWx初始化失败")
	}

	return nil
}

//func NewOpenPlatform() *openplatform.OpenPlatform {
//	var (
//		wc         *wechat.Wechat
//		redisCache *cache.Redis
//		redisOpts  *cache.RedisOpts
//
//		cfg *config.Config
//	)
//	redisOpts = &cache.RedisOpts{
//		Host:     beego.AppConfig.String("redis::address"),
//		Password: beego.AppConfig.String("redis::password"),
//	}
//	redisCache = cache.NewRedis(redisOpts)
//	cfg = &config.Config{
//		AppID:          beego.AppConfig.String("wechat_openplatform::app_id"),
//		AppSecret:      beego.AppConfig.String("wechat_openplatform::app_secret"),
//		Token:          beego.AppConfig.String("wechat_openplatform::token"),
//		EncodingAESKey: beego.AppConfig.String("wechat_openplatform::encoding_aes_key"),
//		Cache:          redisCache,
//	}
//
//	return wc.GetOpenPlatform(cfg)
//}
//
//func RefreshGzhToken(open *openplatform.OpenPlatform, gzh *models.OpenWxApp) {
//	res, err := open.GetAuthrAccessToken(gzh.AuthorizerAppid)
//	if err != nil {
//		logs.Error(err)
//	}
//	logs.Info("RefreshGzhToken", res, err)
//	if res == "" {
//		res, _ := open.RefreshAuthrToken(gzh.AuthorizerAppid, gzh.AuthorizerRefreshToken)
//		logs.Info("RefreshGzhToken:" + utils.GetString(gzh.Id))
//		if res != nil && res.RefreshToken != "" {
//			gzh.AuthorizerRefreshToken = res.RefreshToken
//			_ = new(models.BaseModel).Update(gzh, "authorizer_refresh_token")
//			logs.Info("RefreshGzhTokenSave:" + res.RefreshToken)
//
//		}
//	}
//}
//
////单刷新token
//func SingleRefreshGzhToken(open *openplatform.OpenPlatform, gzh *models.OpenWxApp) {
//	res, _ := open.RefreshAuthrToken(gzh.AuthorizerAppid, gzh.AuthorizerRefreshToken)
//	logs.Info("RefreshGzhToken:" + utils.GetString(gzh.Id))
//	if res != nil && res.RefreshToken != "" {
//		gzh.AuthorizerRefreshToken = res.RefreshToken
//		_ = new(models.BaseModel).Update(gzh, "authorizer_refresh_token")
//		logs.Info("SingleRefreshGzhToken:", res.RefreshToken, gzh)
//	}
//
//}
//
//func RegexMiniProgramAppId(appid string) (match bool, err error) {
//	match, err = regexp.MatchString(`^wx[a-zA-Z0-9]+$`, appid)
//	return
//}
//
//type MpBody struct {
//	Base64img string `json:"base64img"`
//	Ret       int    `json:"ret"`
//}
//
//// 通过微信后台生成小程序码
//func GenerateMpQrCodeFromWeixin(appId string, path string) (err error, qrCode string) {
//	// 获取配置
//	tokenConfig := &models.Config{}
//	cookieConfig := &models.Config{}
//	tokenConfig.Column = "WX_TOKEN"
//	_ = new(models.BaseModel).First(tokenConfig, "column")
//	cookieConfig.Column = "WX_COOKIE"
//	_ = new(models.BaseModel).First(cookieConfig, "column")
//	if tokenConfig.Id == 0 || cookieConfig.Id == 0 {
//		err = errors.WechatErrGenerateQrCode1
//		return
//	}
//	if path != "" && string(path[0]) == "/" {
//		path = path[1:]
//	}
//	// 拼接参数
//	token := tokenConfig.ColumnValue
//	cookie := cookieConfig.ColumnValue
//	path = url.QueryEscape(url.QueryEscape(url.QueryEscape(url.QueryEscape(path))))
//	requestUrl := "https://mp.weixin.qq.com/wxamp/cgi/route?path=%2Fwxopen%2Fwxaqrcode%3Faction%3Dgetqrcode%26f%3Djson%26appid%3D" + appId + "%26path%3D" + path + "%26token%3D1817808811%26lang%3Dzh_CN&token=" + token + "&lang=zh_CN&random=0.786240458954204"
//	res, _ := utils.Request.Get(requestUrl, nil, map[string]string{
//		"Host":            "mp.weixin.qq.com",
//		"Accept":          "application/json",
//		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
//		"Sec-Fetch-Site":  "same-origin",
//		"Sec-Fetch-Mode":  "cors",
//		"Sec-Fetch-Dest":  "empty",
//		"Referer":         "https://mp.weixin.qq.com/wxamp/wxaqrcode/weappcode?simple=1&token=1817808811&lang=zh_CN",
//		"Accept-Language": "zh-CN,zh;q=0.9",
//		"Cookie":          cookie,
//	})
//
//	var mpbody MpBody
//	err = json.Unmarshal(res, &mpbody)
//	if err != nil {
//		err = errors.WechatErrGenerateQrCode2
//		return
//	}
//	if mpbody.Base64img == "" {
//		err = notice.NewDingService().Select(common.WarningRobot).Text("微信COOKIE过期")
//		err = errors.WechatErrGenerateQrCode3
//		return
//	}
//	dist, _ := base64.StdEncoding.DecodeString(mpbody.Base64img)
//	//写入新文件
//	fileName := fmt.Sprintf("union_qrcode/%s.png", utils.Md5(dist))
//	ossService := upload.NewOssService(upload.NewFileService(bytes.NewBuffer(dist), fileName))
//	err = ossService.Upload(file.NewOssHandle())
//	if err != nil {
//		err = errors.WechatErrGenerateQrCode4
//	}
//
//	qrCode = utils.FormatImageUrl(fileName, 0)
//
//	return
//}
