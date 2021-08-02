package saBroker

import "github.com/RichardKnop/machinery/v1/tasks"

type RemoteJobModel interface {
	GetSignature() *tasks.Signature
	Handle(string) error
}

type Job struct {
	Name string
}

func (t *Job) GetSignature() (s *tasks.Signature) {
	s = &tasks.Signature{
		Name:         t.Name,
		RetryCount:   1,
		RetryTimeout: 2,
		Args: []tasks.Arg{
			{
				Name:  "params",
				Type:  "string",
				Value: "",
			},
		},
	}

	return
}
