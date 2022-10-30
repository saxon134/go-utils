package saBroker

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends/result"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
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

// run启动broker
func (m *BrokerManager) run() {
	if m.concurrency <= 0 {
		m.concurrency = 10
	}
	m.work = m.taskCenter.NewWorker("sa-broker", m.concurrency)
	err := m.work.Launch()
	if err != nil {
		fmt.Println(err)
	}
}

// do 执行任务
func (m *BrokerManager) do(name string, values ...interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	//任务名
	m.currentName = name

	//参数
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

	//发送任务
	if v, ok := m.taskMapSignature[m.currentName]; ok {
		res, err := m.taskCenter.SendTask(v)
		if err != nil {
			return err
		}

		m.taskResult[m.currentName] = res
	}

	return nil
}

// doDelay 延迟执行任务
func (m *BrokerManager) doDelay(name string, seconds int64, values ...interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	//任务名
	m.currentName = name

	//参数
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

	//发送任务
	if v, ok := m.taskMapSignature[m.currentName]; ok {
		delayTime := time.Now().UTC().Add(time.Second * time.Duration(seconds))
		v.ETA = &delayTime
		res, err := m.taskCenter.SendTask(v)
		if err != nil {
			return err
		}
		m.taskResult[m.currentName] = res
	}
	return nil
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
