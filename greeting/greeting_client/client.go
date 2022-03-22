package main

import (
	"context"
	"fmt"
	"io"
	"log"
	greeting_pb "praktek-grpc-go/greeting/greeting_pb"

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
	// d := greeting_pb.new()
	doServerStream(c)
}

// func doUnary(c greeting_pb.GreetServiceClient) {
// 	fmt.Println("starting to do a unary rpc")
// 	req := &greeting_pb.GreetingRequest{
// 		Greeting: &greeting_pb.Greeting{
// 			FirstName: "Feri",
// 			LastName:  "Fahrul",
// 		},
// 	}
// 	res, err := c.Greet(context.Background(), req)
// 	if err != nil {
// 		log.Fatalf("error when calling greet: %v", err)
// 	}
// 	log.Printf("Response from greet : %v", res.Result)
// }
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
