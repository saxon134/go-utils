package saHttp

type Request struct {
	Id int64 `json:"id"`
}

type Response struct {
}

type ListRequest struct {
	Status   int    `form:"status" type:"default:2"`
	Statuses string `form:"statuses"`
	Offset   int    `form:"offset" type:"default:0;>=0"`
	Limit    int    `form:"limit" type:"default:20;>0"`
	Word     string `form:"word"`
}

type ListResponse struct {
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
	Cnt    int64 `json:"cnt"`
}

type UpdateStatusRequest struct {
	Status int     `form:"status"`
	IdAry  []int64 `form:"idAry"`
}
