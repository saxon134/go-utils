module github.com/saxon134/go-utils

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/RichardKnop/machinery v1.10.0
	github.com/aliyun/aliyun-oss-go-sdk v2.1.6+incompatible
	github.com/astaxie/beego v1.12.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/garyburd/redigo v1.6.0
	github.com/gin-gonic/gin v1.6.3
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/micro/go-micro/v2 v2.9.1
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	go.uber.org/zap v1.16.0
	gorm.io/gorm v1.20.12
)
