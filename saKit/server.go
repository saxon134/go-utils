package saKit

import (
	"context"
	"fmt"
	"github.com/saxon134/go-utils/saKit/proto"
)

type Service interface {
	Api(ctx context.Context, in *proto.Request) (ack *proto.Response, err error)
}

type baseServer struct {
}

func NewService() Service {
	var server Service
	server = &baseServer{}
	return server
}

func (s baseServer) Api(ctx context.Context, in *proto.Request) (ack *proto.Response, err error) {
	fmt.Println("called:", in.Data)
	ack = new(proto.Response)
	ack.Data = "hello, this is kit server."
	return
}
