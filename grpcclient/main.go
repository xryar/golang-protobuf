package main

import (
	"context"
	"golang-protobuf/pb/chat"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	clientConnection, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to create client ", err)
	}

	chatClient := chat.NewChatServiceClient(clientConnection)
	stream, err := chatClient.SendMessage(context.Background())
	if err != nil {
		log.Fatal("Failed to send message", err)
	}

	err = stream.Send(&chat.ChatMessage{
		Message: "Acumalaka",
		Content: "Hello from client",
	})
	if err != nil {
		log.Fatal("Failed to send via stream ", err)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Failed to close", err)
	}
	log.Println("Connection is closed. Message: ", response)
}
