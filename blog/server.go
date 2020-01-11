package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/reflection"

	"github.com/itzmanish/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct{}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id, omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

func main() {
	// For getting location of raised error or warnings with line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("Connecting to mongoDB...")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://rudra:rudra108@ds361488.mlab.com:61488"))
	if err != nil {
		log.Fatalln(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	collection = client.Database("crudgrpc").Collection("blog")
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
	client.Disconnect(context.TODO())
}
