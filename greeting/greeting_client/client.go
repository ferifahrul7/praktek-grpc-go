package main

import (
	"context"
	"fmt"
	"io"
	"log"
	greeting_pb "praktek-grpc-go/greeting/greeting_pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("hello i'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("couldn't connect : %v", err)
	}
	defer cc.Close()
	c := greeting_pb.NewGreetServiceClient(cc)
	// fmt.Printf("created client %f", c)
	// doUnary(c)
	// doServerStream(c)
	// doClientStream(c)
	doBidiStream(c)
}

func doUnary(c greeting_pb.GreetServiceClient) {
	fmt.Println("starting to do a unary rpc")
	req := &greeting_pb.GreetingRequest{
		Greeting: &greeting_pb.Greeting{
			FirstName: "Feri",
			LastName:  "Fahrul",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error when calling greet: %v", err)
	}
	log.Printf("Response from greet : %v", res.Result)
}
func doServerStream(c greeting_pb.GreetServiceClient) {
	fmt.Println("starting to do a server function rpc")

	req := &greeting_pb.GreetingManyTimesRequest{
		Greeting: &greeting_pb.Greeting{
			FirstName: "Feri",
			LastName:  "Fahrul",
		},
	}
	res, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error when calling greetingfrommany: %v", err)
	}
	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error when reading stream:%v", err)
		}
		fmt.Printf("response from greetingfrommanytimes %v \n", msg.GetResult())
	}
}
func doClientStream(c greeting_pb.GreetServiceClient) {
	fmt.Println("starting to do a server function rpc")
	requests := []*greeting_pb.LongGreetingRequest{
		{
			Greeting: &greeting_pb.Greeting{FirstName: "Feri", LastName: "Fahrul"},
		},
		{
			Greeting: &greeting_pb.Greeting{FirstName: "Feri", LastName: "Keren"},
		},
		{
			Greeting: &greeting_pb.Greeting{FirstName: "Feri", LastName: "Tampan"},
		},
		{
			Greeting: &greeting_pb.Greeting{FirstName: "Feri", LastName: "Ganteng"},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling long greet: %v", err)
	}
	for _, req := range requests {
		fmt.Printf("sending Request: %v", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving long greet: %v \n", err)
	}
	fmt.Printf("Long Greet Response : %v \n", res)
}

func doBidiStream(c greeting_pb.GreetServiceClient) {
	fmt.Println("starting to do a bidi streaming rpc")
	requests := []*greeting_pb.GreetingEveryoneRequest{
		{
			Greeting: &greeting_pb.Greeting{FirstName: "Feri", LastName: "Fahrul"},
		},
		{
			Greeting: &greeting_pb.Greeting{FirstName: "Feri", LastName: "Keren"},
		},
		{
			Greeting: &greeting_pb.Greeting{FirstName: "Feri", LastName: "Tampan"},
		},
		{
			Greeting: &greeting_pb.Greeting{FirstName: "Feri", LastName: "Ganteng"},
		},
	}
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error when sending from client: %v", err)
		return
	}
	waitc := make(chan struct{})

	go func() {
		//mengirim request
		for _, req := range requests {
			fmt.Printf("sennding message %v \n ", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error when receiving from server: %v", err)
				break
			}
			fmt.Printf("received: %v", res.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}
