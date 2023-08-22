package grpcs

import (
	"context"
	"fmt"

	userpb "github.com/Prameesh-P/user-handler/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateGRPCConnection() *grpc.ClientConn {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Failed to connect to gRPC server: %v\n", err)
		return nil
	}
	return conn
}

func CreateUserWithGRPC(conn *grpc.ClientConn,user *userpb.User) (*userpb.User, error) {
	client := userpb.NewUserServiceClient(conn)
	fmt.Println("client created")

	req := &userpb.User{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email:     user.Email,
		Age: user.Age,
		Phone:user.Phone,
	}

	resp, err := client.CreateUser(context.Background(), req)
	if err != nil {
		fmt.Printf("Failed to create user via gRPC: %v\n", err)
		return nil, err
	}
	fmt.Println("requested")
	return resp, nil
}

func GetUserByIDWithGRPC(conn *grpc.ClientConn, userID string) (*userpb.User, error) {
	client := userpb.NewUserServiceClient(conn)

	req := &userpb.UserIDRequest{
		UserId: userID,
	}

	resp, err := client.GetUserByID(context.Background(), req)
	if err != nil {
		fmt.Printf("Failed to get user via gRPC: %v\n", err)
		return nil, err
	}

	return resp, nil
}
func DeleteUserWithHTTP(conn *grpc.ClientConn, userid string)(*userpb.UserDeleteResponse, error){

	client := userpb.NewUserServiceClient(conn)
	req := &userpb.UserIDRequest{
		UserId: userid,
	}

	resp ,err := client.DeleteUserByID(context.Background(),req)
	if err != nil {
		fmt.Printf("Failed to get user via gRPC: %v\n", err)
		return nil, err
	}

	return resp, nil
}
func UpdateUserWithHTTP(conn *grpc.ClientConn, user *userpb.User) (*userpb.User, error) {
	client := userpb.NewUserServiceClient(conn)
	
	req := &userpb.User{
		ID: user.ID,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email:     user.Email,
		Age: user.Age,
		Phone:user.Phone,
	}

	resp, err := client.UpdateUser(context.Background(), req)
	if err != nil {
		fmt.Printf("Failed to get user via gRPC: %v\n", err)
		return nil, err
	}

	return resp, nil
}
