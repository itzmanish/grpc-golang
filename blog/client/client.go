package main

import (
	"context"
	"fmt"
	"log"

	"github.com/itzmanish/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("Client is running...")
	tls := false
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/CA.crt"
		cred, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error loading ceritficate: %v", sslErr)
		}
		opts = grpc.WithTransportCredentials(cred)
	}

	conn, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Failed to connect to server %v", err)
	}

	defer conn.Close()
	c := blogpb.NewBlogServiceClient(conn)

	blog := &blogpb.Blog{
		Author:  "Manish",
		Title:   "First Blog",
		Content: "This is the content of my first blog.",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Error occured on saving blog: %v", err)
	}
	fmt.Printf("Created Blog: %v", res)
}
