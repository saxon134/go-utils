package client

import (
	"context"
	"fmt"
	"github.com/saxon134/go-utils/saKit/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestGrpcClient(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:8881", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	cl := NewGRPCClient(conn)
	ack, err := cl.Api(context.Background(), &proto.Request{Data: "this is kit client."})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("请求成功：", ack.Data)
}
