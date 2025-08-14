package main

import (
	"context"

	pb "github.com/hafeeluddeen/go_lang/go_grpc_demo/proto"
)

// This means SayHello is a method on your helloServer struct.
func (s *helloServer) SayHello(ctx context.Context, req *pb.NoParam) (*pb.HelloResponse, error) {

	// return success_msg, nil_message
	return &pb.HelloResponse{
		Messsage: "Hello",
	}, nil
}
