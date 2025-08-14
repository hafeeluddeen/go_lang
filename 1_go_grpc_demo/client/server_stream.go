package main

import (
	"context"
	"io"
	"log"

	pb "github.com/hafeeluddeen/go_lang/go_grpc_demo/proto"
)

func callSayHelloServerStream(client pb.GreetServiceClient, names *pb.NamesList) {

	log.Printf("Stream Started")

	stream, err := client.SayHelloServerStreaming(context.Background(), names)

	if err != nil {
		log.Fatalf("could not send names: %v", err)
	}

	for {

		// listen for stream data.
		message, err := stream.Recv()

		// end of stream.
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while streaming %v", err)
		}

		log.Println(message)

	}

	log.Printf("Streamiong finishes")

}
