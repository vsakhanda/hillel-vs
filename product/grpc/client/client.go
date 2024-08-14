package main

import (
	"context"
	"log"
	"time"

	pb "github.com/SaYaku64/hillel/product/grpc/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.RegisterUser(ctx, &pb.User{Id: "1", Username: "user1", Password: "pass1"})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("User registered: %s", r.Status)

	user, err := c.GetUser(ctx, &pb.UserRequest{Id: "1"})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("User: %s", user.Username)
}
