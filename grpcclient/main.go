package main

import (
	"context"
	"golang-protobuf/pb/user"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	clientConnection, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to create client ", err)
	}

	userClient := user.NewUserServiceClient(clientConnection)

	response, err := userClient.CreateUser(context.Background(), &user.User{
		Id:      1,
		Age:     13,
		Balance: 10000,
		Address: &user.Address{
			Id:          123,
			FullAddress: "Jln. Jaktim",
			Province:    "Jaktim",
			City:        "Jaktim",
		},
	})
	if err != nil {
		log.Fatal("Error calling user client", err)
	}

	log.Println("Got Message from server", response.Message)
}
