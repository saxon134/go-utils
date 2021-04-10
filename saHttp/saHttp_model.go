package saHttp

type Request struct {
	Id int64 `json:"id"`
}

type Response struct {
}

type UpdateStatusRequest struct {
	Status int     `form:"status"`
	IdAry  []int64 `form:"idAry"`
}
