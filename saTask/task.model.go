package saTask

type TaskInfo struct {
	Id     int64  `orm:"" json:"id" form:"id"`
	App    string `orm:"varchar(20);comment:系统名" json:"app" form:"app"`
	Key    string `orm:"varchar(30);comment:任务唯一标识" json:"key" form:"key"`
	Name   string `orm:"varchar(40);comment:任务名" json:"name" form:"name"`
	Status int    `orm:"tinyint;comment:-1-已删除 1-暂停 2-正常" json:"status" form:"status"`
	Spec   string `orm:"varchar(40)" json:"spec" form:"spec"`
	Params string `orm:"varchar(500);comment:执行任务时参数" json:"params" form:"params"`
}

type EventRequest struct {
	Key    string                 `json:"key"`
	Event  string                 `json:"event"` //start, stop, once
	Spec   string                 `json:"spec"`
	Params map[string]interface{} `json:"params"`
}
