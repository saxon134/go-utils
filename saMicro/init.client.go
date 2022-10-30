package saMicro

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/saxon134/go-utils/saMicro/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func InitClient(address string) (cl Service, err error) {
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()

	cl = NewGRPCClient(conn)
	return
}

func NewGRPCClient(conn *grpc.ClientConn) Service {
	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(func(ctx context.Context, md *metadata.MD) context.Context {
			ctx = metadata.NewOutgoingContext(context.Background(), *md)
			return ctx
		}),
	}
	var apiEndpoint endpoint.Endpoint
	{
		apiEndpoint = grpctransport.NewClient(
			conn,
			"proto.Api",
			"RpcApi",
			EncodeRequest,
			DecodeResponse,
			proto.Response{},
			options...).Endpoint()
	}
	return ApiEndPoint{
		EndPoint: apiEndpoint,
	}
}

func EncodeRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func DecodeResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response, nil
}
