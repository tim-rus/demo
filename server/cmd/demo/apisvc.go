package main

import (
	"context"
	"demo/interal/proto/apigen"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type APIService struct {
}

func NewAPIServce() *APIService {
	return &APIService{}
}

//

func (svc *APIService) Ping(context.Context, *connect.Request[apigen.PingRequest]) (*connect.Response[apigen.PingResponse], error) {
	return &connect.Response[apigen.PingResponse]{
		Msg: &apigen.PingResponse{
			Response:  "pong",
			Timestamp: timestamppb.New(time.Now()),
		},
	}, nil
}
