package main

import (
	"context"
	"fmt"
	"log"
	"praktek-grpc-go/greet/greetpb"

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
	c := greetpb.NewGreetServiceClient(cc)

	// fmt.Printf("created client %f", c)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do a unary rpc")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
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
