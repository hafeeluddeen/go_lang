package main

import (
	"log"
	"net"

	pb "github.com/hafeeluddeen/go_lang/go_grpc_demo/proto"
	"google.golang.org/grpc"
)

const (
	port = ":8000"
)

type helloServer struct {
	pb.GreetServiceServer
}

func main() {
	listen, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen on port %v", port)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterGreetServiceServer(grpcServer, &helloServer{})

	log.Printf("Server started at %v", listen.Addr())

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
