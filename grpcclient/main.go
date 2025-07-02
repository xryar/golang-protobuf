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
	stream, err := chatClient.Chat(context.Background())
	if err != nil {
		log.Fatal("Failed to send message", err)
	}

	err = stream.Send(&chat.ChatMessage{
		Message: "ACumalaka",
		Content: "Hello this is client",
	})
	if err != nil {
		log.Fatalf("Failed to send message %v", err)
	}

	msg, err := stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive message %v", err)
	}
	log.Printf("Got reply from server %s content %s", msg.Message, msg.Content)

	msg, err = stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive message %v", err)
	}
	log.Printf("Got reply from server %s content %s", msg.Message, msg.Content)

	err = stream.Send(&chat.ChatMessage{
		Message: "ACumalaka",
		Content: "Hello this is client again rawr",
	})
	if err != nil {
		log.Fatalf("Failed to send message %v", err)
	}
	msg, err = stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive message %v", err)
	}
	log.Printf("Got reply from server %s content %s", msg.Message, msg.Content)
	msg, err = stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive message %v", err)
	}
	log.Printf("Got reply from server %s content %s", msg.Message, msg.Content)

}
