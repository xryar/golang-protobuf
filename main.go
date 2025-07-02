package main

import (
	"context"
	"errors"
	"golang-protobuf/pb/chat"
	"golang-protobuf/pb/user"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type userService struct {
	user.UnimplementedUserServiceServer
}

type chatService struct {
	chat.UnimplementedChatServiceServer
}

func (us *userService) CreateUser(ctx context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	log.Println("User is Created")
	return &user.CreateResponse{
		Message: "User Created",
	}, nil
}

func (cs *chatService) SendMessage(stream grpc.ClientStreamingServer[chat.ChatMessage, chat.ChatResponse]) error {
	for {
		request, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return status.Errorf(codes.Unknown, "Error receiving message %v", err)
		}

		log.Printf("Receive message: %s, to %s", request.Content, request.Message)
	}

	return stream.SendAndClose(&chat.ChatResponse{
		Message: "Thanks for the messages!",
	})
}

func (cs *chatService) ReceiveMessage(request *chat.ReceiveMessageRequest, stream grpc.ServerStreamingServer[chat.ChatMessage]) error {
	log.Printf("Got connection request from %d\n", request.UserId)

	for i := 0; i < 10; i++ {
		err := stream.Send(&chat.ChatMessage{
			Message: "Acumalaka",
			Content: "Hi Rawr",
		})
		if err != nil {
			return status.Errorf(codes.Unknown, "error sending message to client %v", err)
		}
	}
	return nil
}

func (cs *chatService) Chat(stream grpc.BidiStreamingServer[chat.ChatMessage, chat.ChatMessage]) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return status.Errorf(codes.Unknown, "error receiving message")
		}

		log.Printf("Got message to %s content %s", msg.Message, msg.Content)

		err = stream.Send(&chat.ChatMessage{
			Message: "Acumalaka",
			Content: "Reply from server",
		})
		if err != nil {
			return status.Errorf(codes.Unknown, "error sending message")
		}

		err = stream.Send(&chat.ChatMessage{
			Message: "Acumalaka",
			Content: "Reply from server #2",
		})
		if err != nil {
			return status.Errorf(codes.Unknown, "error sending message")
		}
	}
	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("There is error in your net listen ", err)
	}

	server := grpc.NewServer()

	user.RegisterUserServiceServer(server, &userService{})
	chat.RegisterChatServiceServer(server, &chatService{})

	reflection.Register(server)

	if err := server.Serve(listen); err != nil {
		log.Fatal("Error running server", err)
	}
}
