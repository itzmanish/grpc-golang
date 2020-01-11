package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc/reflection"

	"github.com/itzmanish/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	// For getting location of raised error or warnings with line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(s, &server{})
	reflection.Register(s)

	go func() {
		fmt.Println("Starting Blog Server....")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve : %v", err)
		}
	}()

	// Wait for ctrl+c
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block untill signal is recieved
	<-ch
	fmt.Printf("\nStopping server..... \n ")
	s.Stop()
	lis.Close()

}
