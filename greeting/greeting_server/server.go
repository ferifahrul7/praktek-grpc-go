package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"praktek-grpc-go/greeting/greeting_pb"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	greeting_pb.UnimplementedGreetServiceServer
}

func (*server) LongGreet(stream greeting_pb.GreetService_LongGreetServer) error {
	fmt.Printf("longreet function was invoked with streaming req")
	result := "xx"
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// finished reading stream
			return stream.SendAndClose(&greeting_pb.LongGreetingResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error while reading : %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		result += "Hello" + firstName + lastName + " !!"
	}
}

func (*server) GreetManyTimes(req *greeting_pb.GreetingManyTimesRequest, stream greeting_pb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetingManyTimes function was invoked with %v \n", stream)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	for i := 0; i < 99; i++ {
		result := "Hello " + firstName + " " + lastName + " " + strconv.Itoa(i)
		res := &greeting_pb.GreetingManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil

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
