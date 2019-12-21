package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/itzmanish/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error occured while connecting to grpc server: %v", err)
	}

	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)
	// unaryDo(c)
	// serverStream(c)
	// clientStream(c)
	BiDiStreaming(c)
}

func unaryDo(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting to do unary rpc..")
	req := &calculatorpb.SumRequest{
		FirstNumber:  33,
		SecondNumber: 65,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error on doing sum : %v", err)
	}
	log.Printf("Result is %v", res)
}

func serverStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting server streaming")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 1245,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error on getting stream: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Some error while getting response from stream: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func clientStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("client streaming started")
	request := []int32{2, 3, 5, 6, 8, 45, 7, 5, 56, 67}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	for _, req := range request {
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: req,
		})
		time.Sleep(1000 * time.Millisecond)
		fmt.Println(req)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res)
}

func BiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("BiDirectional streaming started")
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error on creating stream: %v", err)
	}
	request := []int32{2, 41, 67, 2, 6, 4, 765, 3, 5, 6, 2465, 567463, 3453, 564, 9}
	waitc := make(chan struct{})

	// Don't forget to use go routine ever
	// Goroutine for sending  req
	go func() {
		for _, req := range request {
			log.Printf("sending: %v\n", req)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: req,
			})
			time.Sleep(500 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// Goroutine for recieving
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("err on recieving in client: %v", err)
				break
			}
			log.Printf("Maximum number: %v\n", res.GetMaximum())
		}
		close(waitc)
	}()

	<-waitc

}
