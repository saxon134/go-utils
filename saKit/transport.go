package saKit

import (
	"context"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/saxon134/go-utils/saKit/proto"
)

type grpcServer struct {
	server grpc.Handler
}

func NewGRPCServer(endpoint ApiEndPoint) proto.ApiServer {
	options := []grpc.ServerOption{
	//grpc.ServerBefore(func(ctx context.Context, md metadata.MD) context.Context {
	//	ctx = context.WithValue(ctx, v5_service.ContextReqUUid, md.Get(v5_service.ContextReqUUid))
	//	return ctx
	//}),
	//grpc.ServerErrorHandler(NewZapLogErrorHandler(log)),
	}
	return &grpcServer{server: grpc.NewServer(
		endpoint.EndPoint,
		RequestGrpcApi,
		ResponseGrpcApi,
		options...,
	)}
}

func (s grpcServer) RpcApi(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	_, ack, err := s.server.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return ack.(*proto.Response), nil
}

func RequestGrpcApi(_ context.Context, grpcReq interface{}) (interface{}, error) {
	//return &proto.Request{}, nil
	return grpcReq, nil
}

func ResponseGrpcApi(_ context.Context, response interface{}) (interface{}, error) {
	//return &proto.Response{}, nil
	return response, nil
}
