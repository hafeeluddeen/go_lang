package main

import (
	"log"
	"time"

	pb "github.com/hafeeluddeen/go_lang/go_grpc_demo/proto"
)

func (s *helloServer) SayHelloServerStreaming(req *pb.NamesList, stream pb.GreetService_SayHelloServerStreamingServer) error {
	log.Printf("Got Request with names: %v", req.Names)

	for _, name := range req.Names {
		res := &pb.HelloResponse{
			Messsage: "Hello" + name,
		}

		if err := stream.Send(res); err != nil {
			return err
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}
