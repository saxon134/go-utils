package saHttp

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saError"
	"github.com/saxon134/go-utils/saLog"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Router struct {
	Method HttpMethod
	Check  CheckType
	Handle func(c *Context)
	Log    string //in-打印输入信息 out-打印输出信息 in;out-打印输入输出信息 注：error时一定会打印参数
}

type RouterType int8

const (
	NullRouter RouterType = iota
	ListRouter
	AddRouter
	UpdateRouter
	UpdateStatusRouter
)

type HttpMethod int

const (
	NullApiMethod HttpMethod = iota
	GetMethod
	PostMethod
	AnyMethod
	AutomaticMethod
)

type CheckType int

const (
	NullCheck     CheckType = 0
	MsCheck       CheckType = 1 //包含系统管理 和 媒体管理
	UserCheck     CheckType = 2
	MsOrUserCheck CheckType = 3
	ApiSignCheck  CheckType = 4
)

type Context struct {
	*gin.Context
	Scene  int
	Paging struct {
		Limit  int //默认值为10，基本Valid为false，Limit也不会空
		Offset int
		Valid  bool //有些场景，不传分页参数，表示需要获取所有数据。具体业务代码控制
	}
	Order struct {
		Key  string
		Desc bool
	}
	Me         JwtValue
	Automatic  RouterType //add、update时，一定做MsOrUserCheck，其他方法校验依据配置
	RawData    []byte
	CustomData map[string]interface{} //自定义参数
}

/*Content-Type为application/json时，备份rawData
目的是为了能够多次bind
*/
func Bind(c *Context, objPtr interface{}) (err error) {
	if c.Request.Method == "GET" {
		err = c.ShouldBindQuery(objPtr)
		return err
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
				return err
			}
			return nil
		} else {
			err = c.ShouldBind(objPtr)
			return err
		}
	}
	return errors.New("GET/POST之外不支持")
}

/**
权限校验
系统管理 和 媒体管理都会有isMs参数，媒体管理还会有isMd参数
*/
func PrivilegeCheck(c *Context, t CheckType) bool {
	if t == ApiSignCheck {
		sign := c.GetHeader("Authorization")
		accountId, _ := saData.ToInt64(c.GetHeader("account.id"))
		timestamp, _ := saData.ToInt64(c.GetHeader("timestamp"))
		if sign == "" || accountId <= 0 || timestamp <= 0 {
			return false
		}

		sign2 := "ac79u%^yr!i" + strconv.FormatInt(accountId, 10) + strconv.FormatInt(timestamp, 10)
		sign2 = saData.Md5(sign2, true)
		if sign2 != sign {
			return false
		}
	}

	if t == MsCheck {
		token := c.GetHeader("Authorization")
		accountId, _ := saData.ToInt64(c.GetHeader("account.id"))
		if accountId <= 0 || token == "" {
			return false
		}

		_ = ParseJwt(token, &c.Me)
		if c.Me.AccountId != accountId {
			return false
		}

		return true
	}

	if t == UserCheck {
		token := c.GetHeader("Authorization")
		userId, _ := saData.ToInt64(c.GetHeader("user.id"))
		if userId <= 0 || token == "" {
			return false
		}

		_ = ParseJwt(token, &c.Me)
		if c.Me.UserId != userId {
			return false
		}

		return true
	}

	if t == MsOrUserCheck {
		token := c.GetHeader("Authorization")
		if token == "" {
			return false
		}

		userId, _ := saData.ToInt64(c.GetHeader("user.id"))
		accountId, _ := saData.ToInt64(c.GetHeader("account.id"))

		_ = ParseJwt(token, &c.Me)
		if userId > 0 {
			return c.Me.UserId == userId
		}

		if accountId > 0 {
			return c.Me.AccountId == accountId
		}

		return false
	}

	return true
}

//返回正确时数据
func Res(c *Context, v interface{}, ext interface{}) {
	if c != nil {
		if v == nil || v == "" {
			v = map[string]int{}
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
func ResAry(c *Context, ary interface{}, paging ListResponse) {
	if ary == nil || ary == "" {
		ary = []int{}
	}

	if reflect.ValueOf(ary).IsNil() {
		ary = []int{}
	}

	c.JSON(200, map[string]interface{}{
		"result": ary,
		"code":   0,
		"ext": map[string]interface{}{
			"totalCount": paging.Cnt,
			"pageSize":   paging.Limit,
			"pageNumber": paging.Offset,
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

			_, file, line, ok := runtime.Caller(1)
			if ok {
				if ary := strings.Split(file, "/"); len(ary) > 0 {
					if len(ary) >= 2 {
						caller = ary[len(ary)-2] + "/" + ary[len(ary)-1]
					} else if len(ary) >= 1 {
						caller = ary[len(ary)-1]
					}
				}

				caller += ":" + saData.Itos(line)
			}
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

		//如果是saError，则错误码、信息按emerror输出；否则全部按照敏感信息处理
		yferr := new(saError.Error)
		if saData.JsonToStruct([]byte(err_s), yferr) == nil && yferr.Code > 0 {
			code = yferr.Code
			msg = yferr.Msg
			errMsg = yferr.Err
			caller = yferr.Caller
		} else {
			code = saError.SensitiveErrorCode
			errMsg = err_s

			_, file, line, ok := runtime.Caller(1)
			if ok {
				if ary := strings.Split(file, "/"); len(ary) > 0 {
					if len(ary) >= 2 {
						caller = ary[len(ary)-2] + "/" + ary[len(ary)-1]
					} else if len(ary) >= 1 {
						caller = ary[len(ary)-1]
					}
				}

				caller += ":" + saData.Itos(line)
			}
		}
	}
	saLog.Err(code, caller, msg, errMsg)

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
