package saOrm

import "time"

type BaseModel struct {
	Id        int64 `json:"id"`
	CreatedAt *Time `json:"createdAt" orm:"datetime;created"`
	DeletedAt *Time `json:"deletedAt" orm:"datetime;default:null"`
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
