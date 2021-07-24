package saBroker

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends/result"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/astaxie/beego/logs"
	"github.com/saxon134/go-utils/saHit"
	"sync"
	"time"
)

var (
	Lock sync.Mutex
)

type BrokerManager struct {
	taskCenter       *machinery.Server
	tasks            []RemoteJobModel
	concurrency      int
	work             *machinery.Worker
	taskMapFunc      map[string]interface{}
	taskMapSignature map[string]*tasks.Signature
	taskResult       map[string]*result.AsyncResult
	currentName      string
	lock             sync.Mutex
}

// 单例
func initInstance(host string, queue string, concurrency int) *BrokerManager {
	if nil == _manager {
		Lock.Lock()
		if nil == _manager {
			_manager = &BrokerManager{
				taskCenter:       initMachineryServer(host, queue),
				taskResult:       make(map[string]*result.AsyncResult),
				taskMapFunc:      make(map[string]interface{}),
				taskMapSignature: make(map[string]*tasks.Signature),
				concurrency:      concurrency,
			}
			Lock.Unlock()
		} else {
			Lock.Unlock()
		}
	}
	return _manager
}

func (m *BrokerManager) RegisterJob(jobs ...RemoteJobModel) error {
	_manager.tasks = append(_manager.tasks, jobs...)
	_manager.parseTaskMapFuncAndSignature()
	err := _manager.taskCenter.RegisterTasks(_manager.taskMapFunc)
	return err
}

// 根据任务名称获取任务执行结果
func (m *BrokerManager) GetResultByTaskName(name string) *result.AsyncResult {
	return m.taskResult[name]
}

// 监听任务
func (m *BrokerManager) Run() {
	m.concurrency = saHit.Int(m.concurrency > 0, m.concurrency, 10)
	m.work = m.taskCenter.NewWorker("test", m.concurrency)
	err := m.work.Launch()
	if err != nil {
		fmt.Println(err)
	}
}

// 执行任务
func (m *BrokerManager) Do() error {
	if v, ok := m.taskMapSignature[m.currentName]; ok {
		res, err := m.taskCenter.SendTask(v)
		if err != nil {
			return err
		}

		m.taskResult[m.currentName] = res
	}

	return nil
}

func (m *BrokerManager) Switch(name string) (res *BrokerManager) {
	m.currentName = name
	res = m
	return
}
func (m *BrokerManager) SetParams(values ...interface{}) (res *BrokerManager) {
	if v, ok := m.taskMapSignature[m.currentName]; ok {
		for k, arg := range v.Args {
			if k < len(values) {
				arg.Value = values[k]
			}

			if k < len(v.Args) {
				v.Args[k] = arg
			}
		}
	}
	res = m
	return
}

// 延迟执行任务
func (m *BrokerManager) DoDelay(seconds int64) {
	if v, ok := m.taskMapSignature[m.currentName]; ok {
		delayTime := time.Now().UTC().Add(time.Second * time.Duration(seconds))
		v.ETA = &delayTime
		res, err := m.taskCenter.SendTask(v)
		if err != nil {
			logs.Error("DoDelay job任务执行错误", err, v)
		}
		m.taskResult[m.currentName] = res
	}
}

func (m *BrokerManager) parseTaskMapFuncAndSignature() {
	for _, t := range m.tasks {
		m.taskMapFunc[t.GetSignature().Name] = t.Handle
		m.taskMapSignature[t.GetSignature().Name] = t.GetSignature()
	}
}

func initMachineryServer(host string, queue string) *machinery.Server {
	server, err := machinery.NewServer(&config.Config{
		Broker:        host,
		DefaultQueue:  queue,
		ResultBackend: host,
	})
	if err != nil {
		panic(err)
	}

	return server
}
