package saMicro

import (
	"context"
	"encoding/json"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/pkg/errors"
	"time"
)

type Client struct {
	micro.Service
	App string
}

type Request struct {
	Method string
	Data   []byte
}

type Response struct {
	Data []byte
}

func (m *Request) Bind(ptr interface{}) error {
	if ptr == nil {
		return errors.New("ptr is empty")
	}

	err := json.Unmarshal(m.Data, ptr)
	return err
}

func (m *Response) JSON(v interface{}) {
	var err error
	m.Data, err = json.Marshal(v)
	if err != nil {
		return
	}
}

func (c *Client) Call(ctx context.Context, method string, in interface{}, resPtr interface{}) error {
	var req Request
	req.Method = method

	var err error

	if in != nil {
		if v, ok := in.([]byte); ok {
			req.Data = v
		} else {
			req.Data, err = json.Marshal(in)
			if err != nil {
				return err
			}
		}
	}

	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	}

	request := c.Client().NewRequest(c.App, "Server.Api", req, client.WithContentType("application/json"))
	var rsp Response
	if err = c.Client().Call(ctx, request, &rsp); err != nil {
		return err
	}

	if rsp.Data != nil && len(rsp.Data) > 0 {
		if _, ok := resPtr.([]byte); ok {
			resPtr = rsp.Data
			return nil
		} else if resPtr != nil {
			err = json.Unmarshal(rsp.Data, resPtr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
