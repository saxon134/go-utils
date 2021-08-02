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

func (t *PrintJob) Handle(params string) error {
	fmt.Println("print task handle:", params)
	return nil
}
