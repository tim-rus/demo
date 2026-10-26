package main

import (
	"context"
	"demo/interal/proto/tggen"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TGService struct {
}

func NewTGServce() *TGService {
	return &TGService{}
}

//

func (svc *TGService) Ping(context.Context, *connect.Request[tggen.PingRequest]) (*connect.Response[tggen.PingResponse], error) {
	return &connect.Response[tggen.PingResponse]{
		Msg: &tggen.PingResponse{
			Response:  "pong",
			Timestamp: timestamppb.New(time.Now()),
		},
	}, nil
}
