package brokerTest

import (
	"fmt"
	"github.com/saxon134/go-utils/saBroker"
)

type PrintJob struct {
	saBroker.Job
}

func NewPrintJob(name string) *PrintJob {
	if name == "" {
		panic("job name不能空")
	}

	j := &PrintJob{}
	j.Name = name
	return j
}

func (t *PrintJob) Handle(bAry []byte) error {
	fmt.Println("print task handle:", bAry)
	return nil
}
