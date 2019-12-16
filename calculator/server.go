package main

import (
	"fmt"
	"context"
	"net"
	"log"
	"google.golang.org/grpc"
	"github.com/itzmanish/grpc-go/calculator/calculatorpb"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error){
	fmt.Printf("request object: %v\n",req)
	first_number := req.GetFirstNumber()
	second_number := req.GetSecondNumber()
	sum := first_number + second_number
	res := &calculatorpb.SumResponse{
		Result: sum,
	}
	return res, nil
}

func main(){
	fmt.Println("Server is starting at localhost:50051")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil{
		log.Fatalf("Failed to listen server %v",err)
	}
	s:=grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err!=nil {
		log.Fatalf("Failed to start serve: %v", err)
	}
}
