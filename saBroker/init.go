package saBroker

import (
	"encoding/json"
	"errors"
	"fmt"
)

/*
使用请参考brokerTest
通过Redis，分发事务给不同的实例使用
并发安全
只有多实例时才有用，如果只部署了一个实例，应当使用saGo/saChannel去处理
*/

var _manager *BrokerManager

//Init host: redis://pass@127.0.0.1:6379  queue: tioBroker  concurrency:并发数
func Init(host string, queue string, concurrency int) *BrokerManager {
	if host == "" || queue == "" {
		return nil
	}

	return initInstance(host, queue, concurrency)
}

//RegisterJobs 必须一次性注册所有任务
func RegisterJobs(jobs ...RemoteJobModel) error {
	if _manager == nil {
		return errors.New("未注册remote job")
	}

	if len(jobs) == 0 {
		return errors.New("缺少必要参数")
	}

	err := _manager.RegisterJob(jobs...)
	if err != nil {
		fmt.Println("broker jobs init error.")
		return err
	}

	go _manager.run()
	return nil
}

func Send(name string, params interface{}) error {
	bAry, err := json.Marshal(params)
	if err != nil {
		return err
	}
	str := string(bAry)
	if str == "null" {
		str = ""
	}

	err = _manager.do(name, str)
	return err
}
