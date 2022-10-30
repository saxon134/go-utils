package saMicro

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/saxon134/go-utils/saMicro/proto"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"net"
)

func InitServer(address string) error {
	golangLimit := rate.NewLimiter(10, 1)
	server := NewService()
	endpoints := NewEndPointServer(server, golangLimit)
	grpcServer := NewGRPCServer(endpoints)
	grpcListener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(grpctransport.Interceptor))
	proto.RegisterApiServer(baseServer, grpcServer)
	if err = baseServer.Serve(grpcListener); err != nil {
		return err
	}
	return nil
}
