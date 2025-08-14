package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/hafeeluddeen/go_lang/go_grpc_demo/proto"
)

func callSayHelloBiDirectional(client pb.GreetServiceClient, names *pb.NamesList) {

	log.Printf("Client Streaming Started....")

	stream, err := client.SayHelloBiDirectional(context.Background())

	if err != nil {
		log.Fatalf("Something went wrong in client streaming -> %v", err)
	}

	// need to send a stream and receive a stream
	// need to use channels for that

	waitc := make(chan struct{})

	// go routine
	go func() {

		for {
			message, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while STreaming %v", err)
			}

			log.Println(message)
		}

		close(waitc)

	}()

	for _, name := range names.Names {
		// creating the request object
		req := &pb.HelloRequest{
			Name: name,
		}

		// send the request to the client
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending the  --> %v", err)
		}

		time.Sleep(2 * time.Second)

	}

	stream.CloseSend()
	<-waitc
	log.Printf("Bidirectional Streaming Done")

}
