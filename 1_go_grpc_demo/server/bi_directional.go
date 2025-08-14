package main

import (
	"io"
	"log"

	pb "github.com/hafeeluddeen/go_lang/go_grpc_demo/proto"
)

func (s *helloServer) SayHelloBiDirectional(stream pb.GreetService_SayHelloBiDirectionalServer) error {

	// should respond immidieatly when received a message..should not send a cummulative response like CLient Streaming.

	log.Printf("Started to receive streaming i/p from client --> ")

	for {
		req, err := stream.Recv()

		// if the stream ends close the case
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		log.Printf("Got Request with Name --> %v", req.Name)

		// send the resp to client
		msg := "Hello " + req.Name

		res := &pb.HelloResponse{
			Messsage: msg,
		}

		if err := stream.Send(res); err != nil {
			return err
		}

	}
}
