package main

import (
	"context"
	"log"
	"time"

	pb "github.com/hafeeluddeen/go_lang/go_grpc_demo/proto"
)

/*

	1. "GreetServiceClient" check this in greet_grpc.pb.go to see how to implement Client Side Functioms

*/

// every client will receive this client OBJ
func callSayHelloClientStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Client Streaming Started....")

	stream, err := client.SayHelloClientStreaming(context.Background())

	if err != nil {
		log.Fatalf("Could not send names : %v", err)
	}

	// iterate over the names and send each name slowly as a stream one by one.
	for _, name := range names.Names {

		// creating the request object
		req := &pb.HelloRequest{
			Name: name,
		}

		// send the request to the client
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending the %s --> %v", name, err)
		}

		//
		log.Printf("Sent the request with name: %s", name)

		time.Sleep(2 * time.Second)
	}

	// all done close the stream

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving %v", err)
	}

	log.Printf("%v", res.Messages)

}
