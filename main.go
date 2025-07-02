package main

import (
	"context"
	"golang-protobuf/pb/user"
	"log"
	"net"

	"google.golang.org/grpc"
)

type userService struct {
	user.UnimplementedUserServiceServer
}

func (us *userService) CreateUser(ctx context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	log.Println("User is Created")
	return &user.CreateResponse{
		Message: "User Created",
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("There is error in your net listen ", err)
	}

	server := grpc.NewServer()

	user.RegisterUserServiceServer(server, &userService{})

	if err := server.Serve(listen); err != nil {
		log.Fatal("Error running server", err)
	}
}
