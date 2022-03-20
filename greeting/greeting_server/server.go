package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"praktek-grpc-go/greeting/greeting_pb"

	"google.golang.org/grpc"
)

type server struct {
	greeting_pb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greeting_pb.GreetingRequest) (*greeting_pb.GreetingResponse, error) {

	fmt.Printf("Greet function was invoked with %v \n", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName
	res := &greeting_pb.GreetingResponse{
		Result: result,
	}
	return res, nil
}
func main() {
	fmt.Println("hello world ini adalah server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()
	greeting_pb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
