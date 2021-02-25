package saHttp

type Request struct {
	Id int64 `json:"id"`
}

type Response struct {
}

type ListRequest struct {
	Status int    `json:"status" api:"default:2"`
	Offset int    `json:"offset" api:"default:0;>=0"`
	Limit  int    `json:"limit" api:"default:20;>0"`
	Word   string `json:"word"`
}

type ListResponse struct {
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
	Cnt    int64 `json:"cnt"`
}

type UpdateStatusRequest struct {
	Status int     `json:"status" api:"default:2"`
	IdAry  []int64 `json:"idAry" api:"required"`
}
