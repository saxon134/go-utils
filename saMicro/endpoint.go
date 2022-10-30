package saMicro

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/saxon134/go-utils/saMicro/proto"
	"golang.org/x/time/rate"
)

type ApiEndPoint struct {
	EndPoint endpoint.Endpoint
}

func NewEndPointServer(svc Service, limit *rate.Limiter) ApiEndPoint {
	var apiEndPoint endpoint.Endpoint
	{
		apiEndPoint = MakeApiEndPoint(svc)
	}
	return ApiEndPoint{EndPoint: apiEndPoint}
}

func (s ApiEndPoint) Api(ctx context.Context, in *proto.Request) (*proto.Response, error) {
	resp, err := s.EndPoint(ctx, in)
	if err != nil {
		fmt.Println("0:")
		return nil, err
	}
	return resp.(*proto.Response), nil
}

func MakeApiEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*proto.Request)
		return s.Api(ctx, req)
	}
}
