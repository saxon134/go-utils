package saOrm

import "time"

type BaseModel struct {
	Id        int64      `json:"id"`
	CreatedAt *time.Time `orm:"datetime;created" json:"createdAt"`
}

type BaseModelWithDelete struct {
	Id        int64      `json:"id"`
	CreatedAt *time.Time `orm:"datetime;created" json:"createdAt"`
	DeletedAt *time.Time `orm:"datetime;default:null" json:"deletedAt"`
}
