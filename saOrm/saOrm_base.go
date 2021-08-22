package saOrm

import "time"

type BaseModel struct {
	Id        int64 `json:"id"`
	CreatedAt *Time `orm:"datetime;created" json:"createdAt"`
}

type BaseModelWithDelete struct {
	Id        int64 `json:"id"`
	CreatedAt *Time `orm:"datetime;created" json:"createdAt"`
	DeletedAt *Time `orm:"datetime;default:null" json:"deletedAt"`
}

type Time time.Time

func (t *Time) IsZero() bool {
	if t == nil || t.IsZero() {
		return true
	}
	return false
}

func (t *Time) Now() {
	now := Time(time.Now())
	t = &now
}
