package saMicro

import (
	"context"
	"errors"
	"github.com/saxon134/go-utils/saMicro/proto"
)

type Service interface {
	Api(ctx context.Context, in *proto.Request) (ack *proto.Response, err error) //RPC服务接口方法
}

type baseServer struct {
}

func NewService() Service {
	var server Service
	server = &baseServer{}
	return server
}

// Api
// @Description: client发起RPC接口请求，对外提供的服务
func (s baseServer) Api(ctx context.Context, in *proto.Request) (ack *proto.Response, err error) {
	if in == nil || in.Method == "" {
		err = errors.New("method is empty")
		return
	}

	for k, v := range _handlers {
		if k == in.Method {
			var out = new(Response)
			out, err = v(ctx, &Request{Method: in.Method, Data: in.Data})
			if err == nil {
				ack = &proto.Response{Data: out.Data}
			}
			return
		}
	}

	err = errors.New("no such handler")
	return
}

// Call client发起请求方法，由RPC内部调用，不对外服务
func (s baseServer) Call(ctx context.Context, in *Request) (ack *Response, err error) {
	var out *proto.Response
	out, err = s.Api(ctx, &proto.Request{Method: in.Method, Data: in.Data})
	if err != nil {
		return
	}

	ack = &Response{Data: out.Data}
	return
}
