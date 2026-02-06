package saHttp

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/saxon134/go-utils/saData/saError"
	"time"
)

func JwtGenerate(ptr interface{}, key string) (j string, err error) {
	var bAry []byte
	bAry, err = json.Marshal(ptr)
	if err != nil {
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": string(bAry),
		"iat": time.Now().Unix(),                        //nbf生效时间
		"edt": time.Now().Unix() + int64(time.Hour*240), //edt失效时间，尚未做失效控制
	})
	j, err = token.SignedString([]byte(key))
	return
}

func JwtParse(token string, key string, ptr interface{}) (err error) {
	var (
		t      *jwt.Token
		claims jwt.MapClaims
		ok     bool
	)

	t, err = jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return
		}
		return []byte(key), nil
	})
	if err != nil {
		return saError.Error{Code: saError.UnAuthedErrorCode}
	}

	claims, ok = t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return saError.Error{Code: saError.UnAuthedErrorCode}
	}

	value := claims["jti"]
	var bAry []byte
	if bAry, ok = value.([]byte); ok == false {
		var str string
		if str, ok = value.(string); ok == true {
			bAry = []byte(str)
		}
	}

	err = json.Unmarshal(bAry, ptr)
	if err != nil {
		return saError.Error{Code: saError.UnAuthedErrorCode}
	}
	return
}
