package saHttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saError"
	"github.com/saxon134/go-utils/saHit"
	"github.com/saxon134/go-utils/saLog"
	"reflect"
	"runtime"
	"strconv"
)

type Router struct {
	Set    string //get,post,auto,sign,ms,user
	Handle func(c *Context)
}

type RouterType int8

const (
	NullRouter RouterType = iota
	ListRouter
	AddRouter
	UpdateRouter
	UpdateStatusRouter
)

type Headers struct {
	Scene   int //1-ms 2-md 3-wxApp 4-wxXcx 5-wxGzh 6-H5
	Product int
	MediaId int64 `form:"mediaId"`
	AppId   int64 `form:"appId"`
	Paging  struct {
		Limit  int //默认值为10，即便Valid为false，Limit也不会空
		Offset int
		Valid  bool //有些场景，不传分页参数，表示需要获取所有数据。具体业务代码控制
	}
	Order struct {
		Key  string
		Desc bool
	}
}

type Context struct {
	*gin.Context
	Headers
	User    UserJwt
	Admin   AdminJwt
	RawData []byte
}

/*
Content-Type为application/json时，备份rawData
目的是为了能够多次bind
*/
func Bind(c *Context, objPtr interface{}) (err error) {
	if c.Request.Method == "GET" {
		err = c.ShouldBindQuery(objPtr)
		if err != nil {
			return err
		}
	} else if c.Request.Method == "POST" {
		tp := c.Request.Header.Get("Content-Type")
		if tp == "" {
			tp = c.Request.Header.Get("content-type")
		}

		if tp == "application/json" {
			if len(c.RawData) == 0 {
				c.RawData, _ = c.GetRawData()
			}
			if len(c.RawData) > 0 {
				err = saData.JsonToStruct(c.RawData, &objPtr)
				if err != nil {
					return err
				}
			}
		} else {
			err = c.ShouldBind(objPtr)
			if err != nil {
				return err
			}
		}
	} else {
		err = errors.New("GET/POST之外不支持")
		return err
	}

	err = saData.TypeCheck(objPtr)
	return err
}

/**
权限校验
系统管理 和 媒体管理都会有isMs参数，媒体管理还会有isMd参数
*/
func PrivilegeCheck(c *Context, t string) bool {
	if t == "sign" {
		sign := c.GetHeader("Authorization")
		timestamp, _ := saData.ToInt64(c.GetHeader("timestamp"))
		if sign == "" || timestamp <= 0 {
			return false
		}

		sign2 := "ac79u%^yr!i" + strconv.FormatInt(timestamp, 10)
		sign2 = saData.Md5(sign2, true)
		if sign2 != sign {
			return false
		}
	}

	if t == "ms" {
		token := c.GetHeader("Authorization")
		if token == "" {
			return false
		}

		_ = ParseAdminJwt(token, &c.Admin)
		if c.Admin.AccountId <= 0 {
			return false
		}

		return true
	}

	if t == "user" {
		token := c.GetHeader("Authorization")
		if token == "" {
			return false
		}

		_ = ParseUserJwt(token, &c.User)
		if c.User.UserId <= 0 {
			return false
		}

		return true
	}

	return true
}

//返回正确时数据
func Res(c *Context, v interface{}, ext interface{}) {
	if c != nil {
		if v == nil || v == "" {
			v = "success"
		}

		var dic = map[string]interface{}{
			"code":   0,
			"result": v,
		}

		if ext != nil {
			dic["ext"] = ext
		}

		c.JSON(200, dic)
		c.Abort()
	}
}

//返回正确的数组数据
func ResAry(c *Context, ary interface{}, cnt int64) {
	if ary == nil || ary == "" {
		ary = []int{}
	}

	if reflect.ValueOf(ary).IsNil() {
		ary = []int{}
	}

	cnt = saHit.Int64(cnt > 0, cnt, int64(c.Paging.Offset+c.Paging.Limit))
	c.JSON(200, map[string]interface{}{
		"result": ary,
		"code":   0,
		"ext": map[string]interface{}{
			"cnt":     cnt,
			"hasNext": cnt > int64(c.Paging.Offset+c.Paging.Limit),
		},
	})
	c.Abort()
}

//返回error
func ResErr(c *Context, err interface{}) {
	if c == nil {
		return
	}

	var msg = "接口报错"
	var errMsg = ""
	var code = saError.NormalErrorCode
	var caller = ""

	if s, ok := err.(string); ok {
		if s != "" {
			msg = s
		}
	} else if e, ok := err.(saError.Error); ok {
		code = e.Code
		msg = e.Msg
		errMsg = e.Err
		caller = e.Caller
	} else if e, ok := err.(*saError.Error); ok {
		code = e.Code
		msg = e.Msg
		errMsg = e.Err
		caller = e.Caller
	} else if e, ok := err.(error); ok {
		err_s := e.Error()
		var dic map[string]interface{}
		_ = json.Unmarshal([]byte(err_s), &dic)

		//判断是否是micro error
		id, _ := saData.ToStr(dic["id"])
		if id == "go.micro.client" {
			err_s, _ = saData.ToStr(dic["detail"])
		}

		//如果是saError，则错误码、信息按saError输出；否则全部按照敏感信息处理
		sa_err := new(saError.Error)
		if saData.JsonToStruct([]byte(err_s), sa_err) == nil && sa_err.Code > 0 {
			code = sa_err.Code
			msg = sa_err.Msg
			errMsg = sa_err.Err
			caller = sa_err.Caller
		} else {
			code = saError.SensitiveErrorCode
			errMsg = err_s
		}
	}

	//重复信息不用打印
	if len(msg) == len(errMsg) && msg == errMsg {
		errMsg = ""
	}

	//获取调用文件及位置
	if len(caller) == 0 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			caller = file + ":" + saData.Itos(line)
			caller += ":" + saData.Itos(line)
		}
	}

	//打印错误信息
	if len(errMsg) > 0 {
		saLog.Err(fmt.Sprintf("%d %s\n%s\n%s", code, msg, errMsg, caller))
	} else {
		saLog.Err(fmt.Sprintf("%d %s\n%s", code, msg, caller))
	}

	if code != 0 {
		rsp_v := map[string]interface{}{"code": code}

		//过滤敏感信息
		if code == saError.UnauthorizedErrorCode {
			msg = "未授权"
		} else if code == saError.UnLoggedErrorCode {
			msg = "未登录"
		} else if code == saError.SensitiveErrorCode {
			//errMsg不返给前端，防止敏感信息
			msg = "error"
		}
		rsp_v["msg"] = msg

		c.JSON(400, rsp_v)
		c.Abort()
		return
	}

	//异常情况
	c.JSON(500, &map[string]interface{}{"code": saError.NormalErrorCode, "msg": "服务器开了个小差"})
	c.Abort()
}
