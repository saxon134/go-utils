package saOrm

import (
	"fmt"
	"testing"
	"time"
)

type DBT struct {
	Id     int        `type:">10;<=30;"`
	T1     time.Time  `type:"updated"`
	T2     *time.Time `type:"updated"`
	S1     string     `type:"varchar(4)"`
	Img    string     `type:"oss"`
	ImgAry StringAry
}

func TestDB(t *testing.T) {
	var m = &DBT{
		Id:     12,
		T1:     time.Now(),
		T2:     nil,
		S1:     "222",
		Img:    "abc",
		ImgAry: StringAry{"1111", "2222"},
	}
	err := ToDB(m)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(m.T2)
	fmt.Println(m.Img)

	_ = FromDb(m)
	fmt.Println(m.Img)
	fmt.Println(m.ImgAry[0])
}
