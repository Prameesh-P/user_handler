package grpcs

import (
	"fmt"
	"log"
	"net"
	"github.com/Prameesh-P/user-handler/internal/user"
	pb "github.com/Prameesh-P/user-handler/pkg/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct {
	userUsecase user.UserUsecase
	pb.UnimplementedUserServiceServer
}

func NewGRPCServer(userUsecase user.UserUsecase) pb.UserServiceServer {
	return &server{
		userUsecase: userUsecase,
	}
}


func StartGRPCServer(userUsecase user.UserUsecase) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, NewGRPCServer(userUsecase))
	fmt.Println("gRPC User Service is running on :8080")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
func (s *server)UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error){
	user := &user.User{
		ID: uint(req.ID),
		FirstName: req.GetFirstName(),
		LastName: req.GetLastName(),
		Email: req.GetEmail(),
		Age: int(req.GetAge()),
		Phone: req.GetPhone(),
	}

	updatedUser, err := s.userUsecase.UpdateUser(user)
	if err!=nil {
		return nil,err
	}
	return &pb.User{
		ID: uint64(updatedUser.ID),
		FirstName:updatedUser.FirstName ,
		LastName:updatedUser.LastName,
		Email: updatedUser.Email,
		Age: int64(updatedUser.Age),
		Phone: updatedUser.Phone,
	},nil
	// if err != nil {
}
func (s *server)DeleteUserByID(ctx context.Context, req *pb.UserIDRequest) (*pb.UserDeleteResponse, error){
	userid:=req.GetUserId()
	msg,err:=s.userUsecase.DeleteUser(userid)
	if err !=nil{
		return nil,err
	}
	return &pb.UserDeleteResponse{
		Msg: msg,
	},nil
}
func (s *server) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	user := &user.User{
		FirstName: req.GetFirstName(),
		LastName: req.GetLastName(),
		Email: req.GetEmail(),
		Age: int(req.GetAge()),
		Phone: req.GetPhone(),
	}
	
	createdUser, err := s.userUsecase.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &pb.User{
	//	Username: createdUser.Username,
		Email:    createdUser.Email,
	}, nil
}

func (s *server) GetUserByID(ctx context.Context, req *pb.UserIDRequest) (*pb.User, error) {
	
	userByID, err := s.userUsecase.GetUserByID(req.UserId)
	if err != nil {
		return nil, err
	}
	if userByID  !=nil{
		return &pb.User{
		
			FirstName:userByID.FirstName ,
			LastName: userByID.LastName,
			Phone: userByID.Phone,
			Age:int64(userByID.Age),
			Email:    userByID.Email,
		}, nil
		
	}

	return new(pb.User),nil
}