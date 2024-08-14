package main

import (
	"context"
	"log"
	"net"

	pb "github.com/SaYaku64/hillel/product/grpc/proto"

	"google.golang.org/grpc"
)

type server struct {
	cache map[string]*pb.User

	pb.UnimplementedUserServiceServer
}

func (s *server) RegisterUser(ctx context.Context, user *pb.User) (*pb.UserResponse, error) {
	// Реєстрація користувача
	s.cache[user.GetId()] = user

	return &pb.UserResponse{Status: "User registered successfully", UserID: user.Id}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
	// Отримання інформації про користувача
	return s.cache[req.GetId()], nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{cache: make(map[string]*pb.User)})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
