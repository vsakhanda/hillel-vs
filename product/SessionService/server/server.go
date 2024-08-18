package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/SaYaku64/hillel/product/SessionService/proto"

	"google.golang.org/grpc"
)

type server struct {
	cache map[string]*pb.Session

	pb.UnimplementedSessionServiceServer
}

func (s *server) AddSession(ctx context.Context, session *pb.Session) (*pb.AddSessionResponse, error) {
	// Реєстрація користувача
	s.cache[session.GetId()] = session

	return &pb.AddSessionResponse{Status: "Session registered successfully", Session: session.Id}, nil
}

func (s *server) GetSession(ctx context.Context, req *pb.GetSessionRequest) (*pb.Session, error) {
	// Отримання інформації про користувача
	return s.cache[req.GetId()], nil
}

func (s *server) DelSession(ctx context.Context, req *pb.DelSessionRequest) (*pb.DelSessionResponse, error) {
	// Отримання інформації про користувача
	if _, exists := s.cache[req.GetId()]; exists {
		delete(s.cache, req.GetId())
		log.Printf("Session with ID %s deleted", req.GetId())
		return &pb.DelSessionResponse{}, nil
	} else {
		return nil, fmt.Errorf("session with ID %s not found", req.GetId())
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSessionServiceServer(s, &server{cache: make(map[string]*pb.Session)})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
