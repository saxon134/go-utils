package brokerTest

import (
	"fmt"
	"gitee.com/go-utils/saBroker"
)

type PrintJob struct {
	saBroker.RemoteJob
}

func NewPrintJob(name string) *PrintJob {
	if name == "" {
		panic("job name不能空")
	}

	j := &PrintJob{}
	j.Name = name
	return j
}

func (t *PrintJob) GetHandle() interface{} {
	return func(bAry []byte) error {
		fmt.Println("print task handle:", bAry)
		return nil
	}
}
