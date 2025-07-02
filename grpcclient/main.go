package main

import (
	"context"
	"errors"
	"golang-protobuf/pb/chat"
	"io"
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
	stream, err := chatClient.ReceiveMessage(context.Background(), &chat.ReceiveMessageRequest{
		UserId: 30,
	})
	if err != nil {
		log.Fatal("Failed to send message", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal("Failed to recieve message ", err)
		}

		log.Printf("Got message to %s content %s", msg.Message, msg.Content)
	}

}
