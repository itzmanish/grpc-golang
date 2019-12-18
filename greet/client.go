package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/itzmanish/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("This is client ")

	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect to server %v", err)
	}

	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)

	// fmt.Printf("Successfully connected to server: %f", c)
	// doUnary(c)
	doStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greet{
			FirstName: "Manish",
			LastName:  "Kumar",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC... %v", err)
	}

	log.Printf("Response from greet %v", res)
}

func doStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting server streaming")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greet{
			FirstName: "Manish",
			LastName:  "Kumar",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error during streaming rpc: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Fatalln("end of streaming reached")
			break
		}
		if err != nil {
			log.Fatalf("something bad happened: %v", err)
		}
		fmt.Printf("recieved from server: %v\n", msg)
	}
}
