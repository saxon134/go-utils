package saMicro

import (
	"context"
	"github.com/saxon134/go-utils/saError"
)

type HandleFunc func(ctx context.Context, in *Request, out *Response) error

type ServiceHandle struct {}

var _handlers = map[string]HandleFunc{}

func (t *ServiceHandle) RegisterHandlers(handlers map[string]HandleFunc) {
	if len(handlers) > 0 {
		_handlers = handlers
	}
}

func (t *ServiceHandle) Api(c context.Context, args *Request, res *Response) error {
	if args == nil {
		return saError.StackError(saError.ErrorParams)
	}

	handle, _ := _handlers[args.Method]
	if handle == nil {
		return saError.StackError("RPC方法有误")
	}

	err := handle(c, args, res)
	if err != nil {
		return err
	}

	return nil
}
