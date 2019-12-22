package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	// doStreaming(c)
	// doClientStreaming(c)
	// doBiDirectional(c)
	doUnaryWithDeadline(c, 5*time.Second) //Won't produce error
	doUnaryWithDeadline(c, 1*time.Second) //Deadline exceeded error
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

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Client streaming started")
	request := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greet{
				FirstName: "Manish",
			}},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greet{
				FirstName: "John",
			}},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greet{
				FirstName: "Stefen",
			}},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greet{
				FirstName: "Maarekk",
			}},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greet{
				FirstName: "josh",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	for _, req := range request {
		fmt.Println(req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(res)
}

func doBiDirectional(c greetpb.GreetServiceClient) {
	fmt.Println("Client streaming started")
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalln(err)
		return
	}
	request := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greet{
				FirstName: "Manish",
			}},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greet{
				FirstName: "John",
			}},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greet{
				FirstName: "Stefen",
			}},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greet{
				FirstName: "Maarekk",
			}},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greet{
				FirstName: "josh",
			},
		},
	}

	waitc := make(chan struct{})
	// send requests use goroutine
	go func() {
		for _, req := range request {
			fmt.Println(req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// recieve requests use go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
				break
			}
			log.Println(res)
		}
		close(waitc)
	}()
	// block until everything is done
	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a Unary Deadline RPC")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greet{
			FirstName: "Manish",
			LastName:  "Kumar",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		errMsg, ok := status.FromError(err)
		if ok {
			if errMsg.Code() == codes.DeadlineExceeded {
				fmt.Println("Request timeout.")
			}
			return
		} else {
			fmt.Printf("Something bad happened with grpc: %v", errMsg)
			return
		}
		log.Fatalf("error while calling GreetWithDeadline RPC... %v", err)
	}
	fmt.Printf("Response from GreetWithDeadline %v\n", res)
}
