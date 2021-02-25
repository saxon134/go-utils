package saHttp

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/saxon134/go-utils/saError"
	"time"
)

type JwtValue struct {
	UserId      int64
	AccountId   int64
	AccountType int8
}

const secret = "yfjwt*()tok*#en#^&"

func GenerateJwt(value *JwtValue) (j string, err error) {
	var bAry []byte
	bAry, err = json.Marshal(value)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": bAry,
		"nbf": time.Now().Unix(), //iat签发时间
		"iat": time.Now().Unix(), //nbf生效时间
	})
	j, err = token.SignedString([]byte(secret))
	return
}

func ParseJwt(token string, j *JwtValue) (err error) {
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
		return saError.Error{Code: saError.UnauthorizedErrorCode}
	}

	claims, ok = t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return saError.Error{Code: saError.UnauthorizedErrorCode}
	}

	value := claims["jti"]
	if bAry, ok := value.([]byte); ok {
		err = json.Unmarshal(bAry, j)
		if err != nil {
			return saError.Error{Code: saError.UnauthorizedErrorCode}
		}
	}

	return nil
}
