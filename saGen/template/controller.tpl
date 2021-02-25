package controller

import (
	"git.efeng.co/server/yf-api/jsonApi"
	"git.efeng.co/server/yf-common/common"
	"git.efeng.co/server/yf-utils/yfData"
	"git.efeng.co/server/yf-utils/yfError"
	"git.efeng.co/server/yf-utils/yfHttp"
)

func {{Struct}}(c *yfHttp.Context) {
	switch c.Automatic {
	case yfHttp.ListRouter:
		_{{StructLower}}List(c)
	case yfHttp.AddRouter:
		_add{{Struct}}(c)
	case yfHttp.UpdateRouter:
		_update{{Struct}}(c)
	case yfHttp.UpdateStatusRouter:
		_update{{Struct}}Status(c)
	default:
		yfHttp.ResErr(c, yfError.ErrorIo)
		return
	}
}

func _{{StructLower}}List(c *yfHttp.Context) {
	db := common.DB.Model(&{{StructLower}}.{{TblModel}}{})
	if c.Paging.Valid {
		db.Offset(c.Paging.Offset).Limit(c.Paging.Limit)
	}
	if len(c.Order.Key) > 0 {
		db.Order(c.Order.Key + " " + yfHit.Str(c.Order.Desc, "desc", "asc"))
	}

	var resp = json{{Struct}}.{{Struct}}List{}
	err := db.Find(&resp.Ary).Error
	if yfError.DbErr(err) {
		yfHttp.ResErr(c, err)
		return
	}

	if c.Paging.Valid {
		db.Count(&resp.Cnt)
		if resp.Cnt <= 0 {
			resp.Cnt = int64(len(resp.Ary))
		}
	}

	yfHttp.ResAry(c, resp.Ary, resp.ListResponse)
}

func _add{{Struct}}(c *yfHttp.Context) {
	obj := new({{StructLower}}.{{TblModel}})
	err := c.ShouldBindJSON(obj)
	if err != nil {
		yfHttp.ResErr(c, err)
		return
	}

	if obj.Id > 0 {
		yfHttp.ResErr(c, yfError.ErrorID)
		return
	}

	err = common.DB.Create(obj).Error
	if err != nil {
		yfHttp.ResErr(c, err)
		return
	}

	yfHttp.Res(c, "ok", nil)
}

func _update{{Struct}}(c *yfHttp.Context) {
	obj := new({{StructLower}}.{{TblModel}})
	err := c.ShouldBindJSON(obj)
	if err != nil {
		yfHttp.ResErr(c, err)
		return
	}

	if obj.Id <= 0 {
		yfHttp.ResErr(c, yfError.ErrorID)
		return
	}

	err = common.DB.Save(obj).Error
	if err != nil {
		yfHttp.ResErr(c, err)
		return
	}

	yfHttp.Res(c, "ok", nil)
}

func _update{{Struct}}Status(c *yfHttp.Context) {
	in := new(jsonApi.UpdateStatusRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		yfHttp.ResErr(c, err)
		return
	}

	err = common.DB.Model(&{{StructLower}}.{{TblModel}}{}).
		Where("id in ?", in.IdAry).
		Updates(map[string]int{"status": int(in.Status)}).Error
	if err != nil {
		yfHttp.ResErr(c, err)
		return
	}

	yfHttp.Res(c, "ok", nil)
}