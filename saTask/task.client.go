package saTask

import (
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saTime"
	"github.com/saxon134/go-utils/saTask/task"
)

// ///////////////////////////////////////////////////////
//
//	0/30 * * * * *                        every 30s
//	0 43 21 * * *                         21:43
//	0 15 05 * * * 　　                    05:15
//	0 0 17 * * *                          17:00
//	0 0 17 * * 1                          17:00 in every Monday
//	0 0,10 17 * * 0,2,3                   17:00 and 17:10 in every Sunday, Tuesday and Wednesday
//	0 0-10 17 1 * *                       17:00 to 17:10 in 1 min duration each time on the first day of month
//	0 0 0 1,15 * 1                        0:00 on the 1st day and 15th day of month
//	0 42 4 1 * * 　 　                    4:42 on the 1st day of month
//	0 0 21 * * 1-6　　                    21:00 from Monday to Saturday
//	0 0,10,20,30,40,50 * * * *　          every 10 min duration
//	0 */10 * * * * 　　　　　　           every 10 min duration
//	0 * 1 * * *　　　　　　　　           1:00 to 1:59 in 1 min duration each time
//	0 0 1 * * *　　　　　　　　           1:00
//	0 0 */1 * * *　　　　　　　           0 min of hour in 1 hour duration
//	0 0 * * * *　　　　　　　　           0 min of hour in 1 hour duration
//	0 2 8-20/3 * * *　　　　　　          8:02, 11:02, 14:02, 17:02, 20:02
//	0 30 5 1,15 * *　　　　　　           5:30 on the 1st day and 15th day of month
//
// ///////////////////////////////////////////////////////

type Case struct {
	Key     string //任务唯一标识
	Spec    string //非空
	Handler Handler
	Params  string
}

type Handler task.Handler

// Init 初始化
func Init(cases ...Case) {
	for _, c := range cases {
		if c.Spec == "" {
			panic("spec can not be empty for task")
		}
		task.AddTask(c.Key, task.NewTask(c.Key, c.Spec, c.Params, task.Handler(c.Handler)))
	}

	//开启任务
	task.StartTask()
}

// Event 事件请求
func Event(in *EventRequest) (err error) {
	if in == nil || in.Key == "" || in.Event == "" {
		return saError.ErrParams
	}

	if in.Event == "start" {
		if task.AdminTaskList[in.Key] == nil {
			return saError.New("事件方法不存在")
		}

		if in.Spec == "" {
			return saError.New("spec err")
		}

		task.AdminTaskList[in.Key].SetEnable(true)
		task.AdminTaskList[in.Key].SetSpec(in.Spec)
	} else if in.Event == "stop" {
		if task.AdminTaskList[in.Key] != nil {
			task.AdminTaskList[in.Key].SetEnable(false)
		}
	} else if in.Event == "once" {
		if task.AdminTaskList[in.Key] == nil {
			return saError.ErrNotExisted
		}

		err = task.AdminTaskList[in.Key].Run(in.Params)
		return err
	} else {
		return saError.ErrNotSupport
	}
	return nil
}

func Status(key string) (out map[string]string, err error) {
	if key == "" {
		return nil, saError.ErrParams
	}

	var t = task.AdminTaskList[key]
	if t == nil {
		return nil, saError.ErrNotExisted
	}

	return map[string]string{
		"nextTime": saTime.TimeToStr(t.GetNext(), saTime.FormatDefault),
		"preTime":  saTime.TimeToStr(t.GetPrev(), saTime.FormatDefault),
		"errMsg":   t.GetStatus(),
	}, nil
}

func AddCase(c *Case) (err error) {
	if c == nil || c.Key == "" || c.Spec == "" || c.Handler == nil {
		return saError.ErrParams
	}

	if task.AdminTaskList[c.Key] == nil {
		task.AddTask(c.Key, task.NewTask(c.Key, c.Spec, c.Params, task.Handler(c.Handler)))
	} else {
		task.AdminTaskList[c.Key].SetEnable(true)
	}
	return nil
}

func DelCase(c *Case) (err error) {
	if c == nil || c.Key == "" || c.Spec == "" || c.Handler == nil {
		return saError.ErrParams
	}

	if task.AdminTaskList[c.Key] == nil {
		return nil
	} else {
		task.AdminTaskList[c.Key].SetEnable(false)
	}
	return nil
}
