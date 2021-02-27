package saHttp

import (
	"github.com/gin-gonic/gin"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saError"
	"github.com/saxon134/go-utils/saLog"
	"strings"
)

var _groups map[string]map[string]Router

func InitRouters(g *gin.Engine, groups map[string]map[string]Router) {
	_groups = groups
	for k, conf := range _groups {
		toAdd := map[string]Router{}
		for path, r := range conf {
			full := ConnPath(k, path)

			if r.Method == GetMethod {
				g.GET(full, _get)
			} else if r.Method == PostMethod {
				g.POST(full, _post)
			} else if r.Method == AnyMethod {
				g.GET(full, _get)
				g.POST(full, _post)
			} else if r.Method == AutomaticMethod {
				if strings.HasSuffix(path, "/") == false {
					g.GET(full, _get)

					g.GET(full+".list", _get)
					toAdd[path+".list"] = Router{Method: r.Method, Check: r.Check, Handle: r.Handle}

					g.POST(full+".add", _post)
					toAdd[path+".add"] = Router{Method: r.Method, Check: MsOrUserCheck, Handle: r.Handle}

					g.POST(full+".update", _post)
					toAdd[path+".update"] = Router{Method: r.Method, Check: MsOrUserCheck, Handle: r.Handle}

					g.POST(full+".update.status", _post)
					toAdd[path+".update.status"] = Router{Method: r.Method, Check: MsOrUserCheck, Handle: r.Handle}
				} else {
					panic(saError.StackError("路由设置有误"))
				}
			}
		}

		for k, v := range toAdd {
			conf[k] = v
		}
	}
}

func _get(c *gin.Context) {
	var err error

	var r Router
	for k, routerDic := range _groups {
		for path, router := range routerDic {
			full := ConnPath(k, path)
			if full == c.Request.URL.Path {
				r = router
				break
			}
		}

		if r.Handle != nil {
			break
		}
	}

	ctx := &Context{
		Context:   c,
		Me:        JwtValue{},
		Automatic: NullRouter,
	}

	if r.Handle == nil {
		ResErr(ctx, "接口有误")
		return
	}

	//自动路由
	if r.Method == AutomaticMethod {
		ary := strings.Split(c.Request.URL.Path, "/")
		if len(ary) > 0 {
			source := ary[len(ary)-1]
			if strings.HasSuffix(source, ".list") {
				ctx.Automatic = ListRouter
			}
		}
	}

	//权限校验
	if r.Check != NullCheck {
		if PrivilegeCheck(ctx, r.Check) == false {
			err = saError.Error{Code: 1104, Msg: ""}
			ResErr(ctx, err)
			return
		}
	}

	//分页
	ctx.Paging.Limit = 10
	ctx.Paging.Offset = 0
	ctx.Paging.Valid = false
	if s, ok := ctx.GetQuery("pageSize"); ok {
		ctx.Paging.Valid = true

		if ctx.Paging.Limit, _ = saData.ToInt(s); ctx.Paging.Limit <= 0 {
			ctx.Paging.Limit = 10
		}

		s, _ = ctx.GetQuery("pageNumber")
		i, _ := saData.ToInt(s)
		if i < 1 {
			i = 1
		}

		ctx.Paging.Offset = ctx.Paging.Limit * (i - 1)
	} else if s, ok := ctx.GetQuery("limit"); ok {
		ctx.Paging.Valid = true

		if ctx.Paging.Limit, _ = saData.ToInt(s); ctx.Paging.Limit <= 0 {
			ctx.Paging.Limit = 10
		}

		s, _ = ctx.GetQuery("offset")
		ctx.Paging.Offset, _ = saData.ToInt(s)
		if ctx.Paging.Offset < 0 {
			ctx.Paging.Offset = 0
		}
	}
	if ctx.Paging.Limit > 400 {
		ctx.Paging.Limit = 10
	}

	//排序
	ctx.Order.Key = "id"
	ctx.Order.Desc = true
	order, _ := ctx.GetQuery("order")
	{
		ary := strings.Split(order, ".")
		if len(ary) > 0 {
			ctx.Order.Key = ary[0]
			if len(ary) > 1 {
				if ary[1] == "asc" || ary[1] == "ASC" {
					ctx.Order.Desc = false
				}
			}
		}
	}

	ctx.Scene = ctx.GetInt("scene")
	r.Handle(ctx)

	//error时，打印请求参数
	if c.Writer.Status() != 200 {
		args := map[string]interface{}{}
		_ = c.BindQuery(args)

		saLog.Err("Url:", c.Request.URL.Path)
		saLog.Err("Args:", args)
	}
}

func _post(c *gin.Context) {
	var err error
	var r Router
	for k, routerDic := range _groups {
		for path, router := range routerDic {
			full := ConnPath(k, path)
			if full == c.Request.URL.Path {
				r = router
				break
			}
		}

		if r.Handle != nil {
			break
		}
	}

	ctx := &Context{
		Context:   c,
		Me:        JwtValue{},
		Automatic: NullRouter,
	}

	if r.Handle == nil {
		ResErr(ctx, "接口有误")
		return
	}

	//自动路由
	if r.Method == AutomaticMethod {
		ary := strings.Split(c.Request.URL.Path, "/")
		if len(ary) > 0 {
			source := ary[len(ary)-1]
			if strings.HasSuffix(source, ".update") {
				ctx.Automatic = UpdateRouter
			} else if strings.HasSuffix(source, ".update.status") {
				ctx.Automatic = UpdateStatusRouter
			} else if strings.HasSuffix(source, ".add") {
				ctx.Automatic = AddRouter
			}

			//权限校验
			if ctx.Automatic != NullRouter {
				if PrivilegeCheck(ctx, MsOrUserCheck) == false {
					err = saError.Error{Code: saError.UnauthorizedErrorCode, Msg: ""}
					ResErr(ctx, err)
					return
				}
			}
		}
	} else {
		//权限校验
		if r.Check != NullCheck {
			if PrivilegeCheck(ctx, r.Check) == false {
				err = saError.Error{Code: saError.UnauthorizedErrorCode}
				ResErr(ctx, err)
				return
			}
		}
	}

	ctx.Scene = ctx.GetInt("scene")
	r.Handle(ctx)

	//error时，打印请求参数
	if c.Writer.Status() != 200 {
		args := map[string]interface{}{}
		_ = c.ShouldBind(args)

		saLog.Err("Url:", c.Request.URL.Path)
		saLog.Err("Args:", args)
	}
}
