package saBroker

import (
	"errors"
	"fmt"
)

/*
触发broker：server层 -> trigger
注册任务：jobs -> server层
使用参考brokerTest
注意：broker请求的参数模型，要单独放在一个包里，否则容易引发循环引用
*/

var Manager *BrokerManager

//host: redis://127.0.0.1:6379  queue: techioBroker
func Init(host string, queue string) *BrokerManager {
	if host == "" || queue == "" {
		return nil
	}

	return initInstance(host, queue)
}

//必须一次性注册所有任务
func RegisterRemoteJobs(jobs ...RemoteJobModel) error {
	if Manager == nil {
		return errors.New("未注册remote job")
	}

	if len(jobs) == 0 {
		return errors.New("缺少必要参数")
	}

	err := Manager.RegisterJob(jobs...)
	if err != nil {
		fmt.Println("broker jobs init error.")
		return err
	}

	go Manager.Run()
	return nil
}

//concurrent最大并行数
func RegisterLocalJobs(concurrent int, handle func(j *LocalJob)) error {
	_handle = handle
	initLocal(concurrent)
	return nil
}
