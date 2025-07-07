package main

import (
	"context"
	"errors"
	"golang-protobuf/pb/chat"
	"golang-protobuf/pb/common"
	"golang-protobuf/pb/user"
	"io"
	"log"
	"net"
	"strings"

	protovalidate "buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func loggingMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("Masuk logging middleware")
	log.Println(info.FullMethod)
	res, err := handler(ctx, req)

	log.Println("Setelah request")
	return res, err
}

func authMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println("Masuk auth middleware")

	if info.FullMethod == "/user.UserService/Login" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unknown, "failed parsing metadata")
	}

	authToken, ok := md["authorization"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "token doesn't exist")
	}

	splitToken := strings.Split(authToken[0], " ")
	token := splitToken[1]

	if token != "secret" {
		return nil, status.Error(codes.Unauthenticated, "token is not valid")
	}

	return handler(ctx, req)
}

type userService struct {
	user.UnimplementedUserServiceServer
}

func (us *userService) Login(ctx context.Context, loginRequest *user.LoginRequest) (*user.LoginResponse, error) {
	return &user.LoginResponse{
		Base: &common.BaseResponse{
			StatusCode: 200,
			IsSuccess:  true,
			Message:    "Success",
		},
		AccessToken:  "secret",
		RefreshToken: "refresh secret",
	}, nil
}

type chatService struct {
	chat.UnimplementedChatServiceServer
}

func (us *userService) CreateUser(ctx context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	if err := protovalidate.Validate(userRequest); err != nil {
		if ve, ok := err.(*protovalidate.ValidationError); ok {
			var validations []*common.ValidationError = make([]*common.ValidationError, 0)
			for _, fieldErr := range ve.Violations {
				log.Printf("Field %s message %s", *fieldErr.Proto.Field.Elements[0].FieldName, *fieldErr.Proto.Message)

				validations = append(validations, &common.ValidationError{
					Field:   *fieldErr.Proto.Field.Elements[0].FieldName,
					Message: *fieldErr.Proto.Message,
				})
			}

			return &user.CreateResponse{
				Base: &common.BaseResponse{
					ValidationErrors: validations,
					StatusCode:       400,
					IsSuccess:        false,
					Message:          "validation error",
				},
			}, nil
		}
		return nil, status.Errorf(codes.InvalidArgument, "validation error %v", err)
	}

	log.Println("User is Created")
	return &user.CreateResponse{
		Base: &common.BaseResponse{
			StatusCode: 200,
			IsSuccess:  true,
			Message:    "User Created",
		},
		CreatedAt: timestamppb.Now(),
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

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingMiddleware,
			authMiddleware,
		),
	)

	user.RegisterUserServiceServer(server, &userService{})
	chat.RegisterChatServiceServer(server, &chatService{})

	reflection.Register(server)

	if err := server.Serve(listen); err != nil {
		log.Fatal("Error running server", err)
	}
}
