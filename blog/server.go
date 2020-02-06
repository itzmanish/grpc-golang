package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	ID      primitive.ObjectID `bson:"_id, omitempty"`
	Author  string             `bson:"author_id"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("Create blog request..")
	blog := req.GetBlog()

	data := blogItem{
		Author:  blog.GetAuthor(),
		Title:   blog.GetTitle(),
		Content: blog.GetContent(),
	}
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("Internal server error: %v", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	fmt.Println(oid, reflect.TypeOf(oid))
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error occured while converting objectid: %v", err))
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:      oid.Hex(),
			Author:  blog.GetAuthor(),
			Title:   blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
}

func main() {
	// For getting location of raised error or warnings with line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("Connecting to mongoDB...")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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
