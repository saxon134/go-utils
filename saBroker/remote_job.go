package saBroker

import "github.com/RichardKnop/machinery/v1/tasks"

type RemoteJobModel interface {
	GetSignature() *tasks.Signature
	GetHandle() interface{}
}

type RemoteJob struct {
	Name string
}

func (t *RemoteJob) GetSignature() (s *tasks.Signature) {
	s = &tasks.Signature{
		Name:         t.Name,
		RetryCount:   1,
		RetryTimeout: 2,
		Args: []tasks.Arg{
			{
				Name:  "params",
				Type:  "[]byte",
				Value: []byte{},
			},
		},
	}

	return
}
