package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/reflection"

	"github.com/itzmanish/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function invoked with %v \n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet function invoked with %v \n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("long greet stream client invoked")
	result := ""
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		firstName := res.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("Bidirection server invoked")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Fatalf("end of file err: %v", err)
			return nil
		}
		if err != nil {
			log.Fatalf("Error on recieving: %v\n", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error on sending stream: %v\n", err)
			return err
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("Greet With Deadline function invoked with %v \n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// The client canceled the request
			fmt.Println("The client canceled the request")
			return nil, status.Errorf(codes.Canceled, "The client has canceled the request!")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Server is starting at 0.0.0.0:50051")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	opts := []grpc.ServerOption{}
	tls := true
	if tls {
		keyFile := "ssl/server.pem"
		certFile := "ssl/server.crt"

		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		cred, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Some error occured on loading ssl certificates: %v", sslErr)
		}
		opts = append(opts, grpc.Creds(cred))
	}

	s := grpc.NewServer(opts...)
	greetpb.RegisterGreetServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
