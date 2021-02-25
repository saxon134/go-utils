package saMicro

import (
	"context"
	"github.com/saxon134/go-utils/saError"
)

type HandleFunc func(ctx context.Context, in *Request, out *Response) error

type Server int

var _handlers = map[string]HandleFunc{}

func (t *Server) RegisterHandlers(handlers map[string]HandleFunc) {
	if len(handlers) > 0 {
		_handlers = handlers
	}
}

func (t *Server) Api(c context.Context, args *Request, res *Response) error {
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
