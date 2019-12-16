package main

import (
	"fmt"
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

	fmt.Printf("Successfully connected to server: %f", c)

}
