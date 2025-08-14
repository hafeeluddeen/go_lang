package main

import (
	"io"
	"log"

	pb "github.com/hafeeluddeen/go_lang/go_grpc_demo/proto"
)

/*
	1. Refer "type GreetServiceServer" in "greet_grpc.pb.go" to find out how to define the function

*/

// this is the function name that we have defined in the proto buff file "SayHelloClientStreaming"
func (s *helloServer) SayHelloClientStreaming(stream pb.GreetService_SayHelloClientStreamingServer) error {

	// here the client will send the data in streams
	var messages []string

	// streram started
	log.Printf("Started to receive stream from client")

	// return success_msg, nil_message

	for {
		req, err := stream.Recv()

		// end of stream sent by client
		if err == io.EOF {
			return stream.SendAndClose(&pb.MessagesList{
				Messages: messages,
			})
		}

		// if anyother error other than EOF throw error
		if err != nil {
			return err

		}

		// Read the stream
		log.Printf("Got Request with the name -> %v", req.Name)
		messages = append(messages, "Hello ", req.Name)

	}

}
