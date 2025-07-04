package main

import (
	"context"
	"golang-protobuf/pb/user"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	clientConnection, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to create client ", err)
	}

	now := time.Now()

	userCLient := user.NewUserServiceClient(clientConnection)
	response, err := userCLient.CreateUser(context.Background(), &user.User{
		Age:      -1,
		Birtdate: timestamppb.New(now),
	})
	if err != nil {
		// st, ok := status.FromError(err)
		// if ok {
		// 	if st.Code() == codes.InvalidArgument {
		// 		log.Println("There is validation error", st.Message())
		// 	} else if st.Code() == codes.Unknown {
		// 		log.Println("There is unknow error", st.Message())
		// 	} else if st.Code() == codes.Internal {
		// 		log.Println("There is internal error", st.Message())
		// 	}

		// 	return
		// }

		log.Println("Failed to send message ", err)
		return
	}
	if !response.Base.IsSuccess {
		switch response.Base.StatusCode {
		case 400:
			log.Println("There is validation error", response.Base.Message)
		case 500:
			log.Println("There is internal error", response.Base.Message)
		}
		return
	}

	log.Println("Response from server ", response.Base.Message)
}
