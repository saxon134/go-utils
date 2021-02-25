package controller

import (
	"github.com/saxon134/go-utils/saData"
)

func {{FunDetail}}(c *gin.Context, args *map[string]interface{}) {
	id, _ := saData.ToInt64((*args)["id"])
	if id <= 0 {
		saHttp.ResErr(c, "缺少参数")
		return
	}

	obj, err := {{PkgName}}.GetByPk(nil, id)
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	resObj := new({{PkgName}}.{{TitlePkgName}})
	resObj.{{ModelName}} = obj

	if ary, er := imgTxt.GetListApi(id, enum.ContentParentType, nil); er == nil {
		resObj.ImgTxtAry = ary
	}

	saHttp.Res(c, resObj, nil)
}


func {{FunList}}(c *gin.Context, args *map[string]interface{}) {
	offset, limit := saHttp.DefaultListParams(args)
	params := map[string]string{
		"offset": saData.Itos(offset),
		"limit":  saData.Itos(limit),
	}

	isMs, _ := saData.ToBool((*args)["isMs"])
	if isMs {
		st := enum.ToStatus(args)
		if st != enum.NullStatus {
			params["status"] = saData.Itos(int(st))
		}
	} else {
		params["status"] = saData.Itos(int(enum.EnableStatus))
	}

	ary, err := {{PkgName}}.GetRows(nil, &params)
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	cnt, _ := {{PkgName}}.Count(nil, &params)
	if cnt <= 0 {
		cnt = len(ary)
	}

    saHttp.ResAry(c, ary, cnt, offset, limit)
}

func {{FunAdd}}(c *gin.Context, args *map[string]interface{}) {
	obj := new({{PkgName}}.{{ModelName}})
	err := obj.FromDic(args)
	if err != nil {
		saHttp.ResErr(c, err)
        return
	}

	if obj.Id > 0 {
		saHttp.ResErr(c, "id有误")
        return
	}

	//开启事务
	tx, _ := common.Sql.Begin()
	if tx == nil {
		saHttp.ResErr(c, "tx error")
        return
	}

	err = {{PkgName}}.Insert(tx, obj)
	if err != nil {
		saHttp.ResErr(c, err)
        return
	}

	obj.Id = {{PkgName}}.InsertLastId(tx)

	imgTxtAry, _ := saData.ToAryMap((*args)["imgTxtAry"])
	if imgTxtAry != nil && len(*imgTxtAry) > 0 {
		if err = imgTxt.AddListApi(tx, obj.Id, enum.ServeParentType, imgTxtAry); err != nil {
			_ = tx.Rollback()
			saHttp.ResErr(c, err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		saHttp.ResErr(c, err)
        return
	}

    saHttp.Res(c, "ok", nil)
}

func {{FunUpdate}}(c *gin.Context, args *map[string]interface{}) {
	id, _ := saData.ToInt64((*args)["id"])
	if id <= 0 {
		saHttp.ResErr(c, "缺少参数")
        return
	}

	obj, err := {{PkgName}}.GetByPk(nil, id)
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	if obj.Id <= 0 {
		saHttp.ResErr(c, "项目不存在")
		return
	}

	err = obj.FromDic(args)
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	//开启事务
	tx, _ := common.Sql.Begin()
	if tx == nil {
		saHttp.ResErr(c, "tx error")
		return
	}

	//图文信息
	imgTxtAry, _ := saData.ToAryMap((*args)["imgTxtAry"])
	if imgTxtAry != nil && len(*imgTxtAry) > 0 {
		if err = imgTxt.UpdateListApi(tx, obj.Id, enum.ContentParentType, imgTxtAry); err != nil {
			_ = tx.Rollback()
			saHttp.ResErr(c, err)
			return
		}
	}

	err = {{PkgName}}.Update(tx, obj)
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		saHttp.ResErr(c, err)
		return
	}

	saHttp.Res(c, "ok", nil)
}

func {{FunUpdateStatus}}(c *gin.Context, args *map[string]interface{}) {
	idAry := saData.ToIdAry(args, "idAry")
	status := enum.ToStatus(args)
	if status == enum.NullStatus || len(idAry) == 0 {
		saHttp.ResErr(c, "缺少参数")
		return
	}

	err := {{PkgName}}.UpdateRowsColumn(nil, "status", saData.Itos(int(status)), idAry)
	if err != nil {
		saHttp.ResErr(c, err)
		return
	}

	saHttp.Res(c, "ok", nil)
}

