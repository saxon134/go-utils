package saMicro

import (
	"context"
)

type Request struct {
	Method string
	Data   string
}

type Response struct {
	Data string
}

type Handler func(ctx context.Context, in *Request) (ack *Response, err error)

var _handlers map[string]Handler

func RegisterHandlers(handlers map[string]Handler) {
	_handlers = handlers
}
