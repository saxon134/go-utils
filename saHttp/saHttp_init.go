package saHttp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saError"
	"github.com/saxon134/go-utils/saLog"
	"strings"
)

var _groups map[string]map[string]Router
var _root string

func InitRouters(g *gin.Engine, root string, groups map[string]map[string]Router) {
	_groups = groups
	_root = root
	for k, conf := range _groups {
		toAdd := map[string]Router{}
		for path, r := range conf {
			full := ConnPath(root, k)
			full = ConnPath(full, path)
			r.Set = strings.ToLower(r.Set)

			if strings.Index(r.Set, "get") >= 0 {
				g.GET(full, _get)
			}

			if strings.Index(r.Set, "post") >= 0 {
				g.POST(full, _post)
			}

			if strings.Index(r.Set, "auto") >= 0 {
				if strings.HasSuffix(path, "/") == false {
					g.GET(full, _get)

					g.GET(full+".list", _get)
					toAdd[path+".list"] = Router{Set: r.Set, Handle: r.Handle}

					g.POST(full+".add", _post)
					toAdd[path+".add"] = Router{Set: r.Set + ",ms,user", Handle: r.Handle}

					g.POST(full+".update", _post)
					toAdd[path+".update"] = Router{Set: r.Set + ",ms,user", Handle: r.Handle}

					g.POST(full+".update.status", _post)
					toAdd[path+".update.status"] = Router{Set: r.Set + ",ms,user", Handle: r.Handle}
				} else {
					panic(errors.New("路由设置有误"))
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
	for _, router := range _groups {
		for path, router := range router {
			full := ConnPath(_root, path)
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
		Context: c,
		User:    UserJwt{},
		Admin:   AdminJwt{},
	}
	ctx.Headers.MediaId, _ = saData.Stoi64(c.GetHeader("media-id"))
	ctx.Headers.AppId, _ = saData.Stoi64(c.GetHeader("app-id"))
	ctx.Headers.Product, _ = saData.Stoi(c.GetHeader("product"))
	ctx.Headers.Scene, _ = saData.ToInt(ctx.GetHeader("scene"))

	if r.Handle == nil {
		ResErr(ctx, "接口有误")
		return
	}

	//权限校验
	scene := ctx.Headers.Scene
	if scene == 1 || scene == 2 {
		if strings.Index(r.Set, "ms") >= 0 {
			if PrivilegeCheck(ctx, "ms") == false {
				err = saError.Error{Code: 1104, Msg: ""}
				ResErr(ctx, err)
				return
			}
		}

		//获取“我”的信息
		if ctx.Admin.AccountId <= 0 {
			token := c.GetHeader("Authorization")
			_ = ParseAdminJwt(token, &ctx.Admin)
		}
	} else {
		if strings.Index(r.Set, "user") >= 0 {
			if PrivilegeCheck(ctx, "user") == false {
				err = saError.Error{Code: 1104, Msg: ""}
				ResErr(ctx, err)
				return
			}
		}

		//获取“我”的信息
		if ctx.User.UserId <= 0 {
			token := c.GetHeader("Authorization")
			_ = ParseUserJwt(token, &ctx.User)
		}
	}

	if strings.Index(r.Set, "sign") >= 0 {
		if PrivilegeCheck(ctx, "sign") == false {
			err = saError.Error{Code: 1104, Msg: ""}
			ResErr(ctx, err)
			return
		}
	}

	//分页
	ctx.Paging.Limit = 20
	ctx.Paging.Offset = 0
	ctx.Paging.Valid = false
	if s, ok := ctx.GetQuery("pageSize"); ok {
		ctx.Paging.Valid = true

		if ctx.Paging.Limit, _ = saData.ToInt(s); ctx.Paging.Limit <= 0 {
			ctx.Paging.Limit = 20
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
			ctx.Paging.Limit = 20
		}

		s, _ = ctx.GetQuery("offset")
		ctx.Paging.Offset, _ = saData.ToInt(s)
		if ctx.Paging.Offset < 0 {
			ctx.Paging.Offset = 0
		}
	}
	if ctx.Paging.Limit > 400 {
		ctx.Paging.Limit = 20
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
	r.Handle(ctx)

	//error时，打印请求参数
	if c.Writer.Status() != 200 {
		args := map[string]interface{}{}
		_ = Bind(ctx, args)

		saLog.Err("Url:", c.Request.URL.Path)
		saLog.Err("Args:", args)
	}
}

func _post(c *gin.Context) {
	var err error
	var r Router
	for k, router := range _groups {
		for path, router := range router {
			full := ConnPath(_root, k)
			full = ConnPath(full, path)
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
		Context: c,
		User:    UserJwt{},
	}
	ctx.MediaId, _ = saData.Stoi64(c.GetHeader("media-id"))
	ctx.AppId, _ = saData.Stoi64(c.GetHeader("app-id"))
	ctx.Product, _ = saData.Stoi(c.GetHeader("product"))
	ctx.Scene, _ = saData.ToInt(ctx.GetHeader("scene"))

	if r.Handle == nil {
		ResErr(ctx, "接口有误")
		return
	}

	//权限校验
	scene := ctx.Headers.Scene
	if scene == 1 || scene == 2 {
		if strings.Index(r.Set, "ms") >= 0 {
			if PrivilegeCheck(ctx, "ms") == false {
				err = saError.Error{Code: 1104, Msg: ""}
				ResErr(ctx, err)
				return
			}
		}

		//获取“我”的信息
		if ctx.Admin.AccountId <= 0 {
			token := c.GetHeader("Authorization")
			_ = ParseAdminJwt(token, &ctx.Admin)
		}
	} else {
		if strings.Index(r.Set, "user") >= 0 {
			if PrivilegeCheck(ctx, "user") == false {
				err = saError.Error{Code: 1104, Msg: ""}
				ResErr(ctx, err)
				return
			}
		}

		//获取“我”的信息
		if ctx.User.UserId <= 0 {
			token := c.GetHeader("Authorization")
			_ = ParseUserJwt(token, &ctx.User)
		}
	}

	if strings.Index(r.Set, "sign") >= 0 {
		if PrivilegeCheck(ctx, "sign") == false {
			err = saError.Error{Code: 1104, Msg: ""}
			ResErr(ctx, err)
			return
		}
	}
	r.Handle(ctx)

	//error时，打印请求参数
	if c.Writer.Status() != 200 {
		args := map[string]interface{}{}
		_ = Bind(ctx, args)

		saLog.Err("Url:", c.Request.URL.Path)
		saLog.Err("Args:", args)
	}
}
