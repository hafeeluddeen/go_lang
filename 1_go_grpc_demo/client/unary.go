package main

import (
	"context"
	"log"
	"time"

	pb "github.com/hafeeluddeen/go_lang/go_grpc_demo/proto"
)

func callSayHello(client pb.GreetServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	res, err := client.SayHello(ctx, &pb.NoParam{})

	if err != nil {
		log.Fatalf("Could Not Greet: %v", err)
	}

	log.Printf("%s", res.Messsage)
}
