package saHttp

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/saxon134/go-utils/saError"
	"time"
)

type AdminJwt struct {
	MediaId   int64
	Product   int
	AccountId int64
}

type UserJwt struct {
	MediaId int64
	AppId   int64
	Product int
	UserId  int64
}

const secret = "yfjwt*()tok*#en#^&"

func GenerateAdminJwt(j *AdminJwt) (token string, err error) {
	var bAry []byte
	bAry, err = json.Marshal(j)
	if err != nil {
		return
	}

	token, err = generate(bAry)
	return
}

func ParseAdminJwt(token string, j *AdminJwt) (err error) {
	var bAry []byte
	bAry, err = parse(token)
	if err != nil {
		return
	}

	err = json.Unmarshal(bAry, j)
	if err != nil {
		return saError.Error{Code: saError.UnauthorizedErrorCode}
	}
	return
}

func GenerateUserJwt(j *UserJwt) (token string, err error) {
	var bAry []byte
	bAry, err = json.Marshal(j)
	if err != nil {
		return
	}

	token, err = generate(bAry)
	return
}

func ParseUserJwt(token string, j *UserJwt) (err error) {
	var bAry []byte
	bAry, err = parse(token)
	if err != nil {
		return
	}

	err = json.Unmarshal(bAry, j)
	if err != nil {
		return saError.Error{Code: saError.UnauthorizedErrorCode}
	}
	return
}

func generate(bAry []byte) (j string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": string(bAry),
		"nbf": time.Now().Unix(), //iat签发时间
		"iat": time.Now().Unix(), //nbf生效时间
	})
	j, err = token.SignedString([]byte(secret))
	return
}

func parse(token string) (bAry []byte, err error) {
	var (
		t      *jwt.Token
		claims jwt.MapClaims
		ok     bool
	)

	t, err = jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, saError.Error{Code: saError.UnauthorizedErrorCode}
	}

	claims, ok = t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return nil, saError.Error{Code: saError.UnauthorizedErrorCode}
	}

	value := claims["jti"]
	ok = false
	if bAry, ok = value.([]byte); ok == false {
		var str string
		if str, ok = value.(string); ok == true {
			bAry = []byte(str)
		}
	}
	return
}
