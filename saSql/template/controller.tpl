package controller

import (
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saError"
	"github.com/saxon134/go-utils/saHttp"
)

func {{Struct}}(c *saHttp.Context) {
	switch c.Automatic {
	case saHttp.ListRouter:
		_{{StructLower}}List(c)
	case saHttp.AddRouter:
		_add{{Struct}}(c)
	case saHttp.UpdateRouter:
		_update{{Struct}}(c)
	case saHttp.UpdateStatusRouter:
		_update{{Struct}}Status(c)
	default:
		saHttp.ResErr(c, saError.ErrorIo)
		return
	}
}

func _{{StructLower}}List(c *saHttp.Context) {
	db := common.DB.Model(&{{StructLower}}.{{TblModel}}{})
	if c.Paging.Valid {
		db.Offset(c.Paging.Offset).Limit(c.Paging.Limit)
	}
	if len(c.Order.Key) > 0 {
		db.Order(c.Order.Key + " " + saHit.Str(c.Order.Desc, "desc", "asc"))
	}

	var resp = tio.{{Struct}}ListResp{}
	err := db.Find(&resp.Ary).Error
	if saError.DbErr(err) {
		saHttp.ResErr(c, err)
		return
	}

	if c.Paging.Valid {
		db.Count(&resp.Cnt)
		if resp.Cnt <= 0 {
			resp.Cnt = int64(len(resp.Ary))
		}
	}

	saHttp.ResAry(c, resp.Ary, resp.ListResponse)
}

func _add{{Struct}}(c *saHttp.Context) {
	obj := new({{StructLower}}.{{TblModel}})
	err := saHttp.Bind(c, obj)
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	if obj.Id > 0 {
		saHttp.ResErr(c, saError.ErrorID)
		return
	}

	err = common.DB.Create(obj).Error
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	saHttp.Res(c, "ok", nil)
}

func _update{{Struct}}(c *saHttp.Context) {
	obj := new({{StructLower}}.{{TblModel}})
	err := saHttp.Bind(c, obj)
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	if obj.Id <= 0 {
		saHttp.ResErr(c, saError.ErrorID)
		return
	}

	err = common.DB.Save(obj).Error
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	saHttp.Res(c, "ok", nil)
}

func _update{{Struct}}Status(c *saHttp.Context) {
	in := new(saHttp.UpdateStatusRequest)
	err := saHttp.Bind(c, in)
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	err = common.DB.Model(&{{StructLower}}.{{TblModel}}{}).
		Where("id in ?", in.IdAry).
		Updates(map[string]int{"status": int(in.Status)}).Error
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	saHttp.Res(c, "ok", nil)
}