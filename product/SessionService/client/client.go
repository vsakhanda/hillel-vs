package main

import (
	"context"
	"log"
	"time"

	pb "github.com/SaYaku64/hillel/product/SessionService/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSessionServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.AddSession(ctx, &pb.Session{Id: "1", Session: "Session one for the first user"})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("Session registered: %s", r.Status)

	session, err := c.GetSession(ctx, &pb.GetSessionRequest{Id: "1"})
	if err != nil {
		log.Fatalf("could not get session: %v", err)
	}
	log.Printf("Session: %s", session.Session)

	_, err = c.DelSession(ctx, &pb.DelSessionRequest{Id: "1"})
	if err != nil {
		log.Fatalf("could not delete session: %v", err)
	}
	log.Println("Session deleted successfully")

}
