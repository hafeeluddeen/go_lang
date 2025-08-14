package main

import (
	"log"

	pb "github.com/hafeeluddeen/go_lang/go_grpc_demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8000"
)

func main() {

	// Client Uses GRPC.Dial to call the server
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	// close the connection
	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)

	names := &pb.NamesList{
		Names: []string{"Hafeel", "Jojo", "Gojo"},
	}

	// callSayHello(client)

	// callSayHelloServerStream(client, names)

	// callSayHelloClientStream(client, names)

	callSayHelloBiDirectional(client, names)

}
