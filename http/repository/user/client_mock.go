package userrepository

import (
	"context"
	"errors"
	"log"
	"net"
	"os"

	gokitLog "github.com/go-kit/log"

	"github.com/jumaroar-globant/go-bootcamp/user/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var (
	// ErrForcedFailure is an error for a forced failure
	ErrForcedFailure = errors.New("forced failure")
	// ForceMockFail is a variable to fail a mock
	ForceMockFail = false
	// ForceBadAge is a variable to return an invalid age
	ForceBadAge = false
)

type grpcMock struct {
	pb.UnimplementedUserServiceServer
}

func initDialer(m *grpcMock) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, m)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func InitGRPCMock() (UserRepository, error) {
	var logger gokitLog.Logger
	{
		logger = gokitLog.NewLogfmtLogger(os.Stderr)
		logger = gokitLog.NewSyncLogger(logger)
		logger = gokitLog.With(logger,
			"service", "user_test",
			"time:", gokitLog.DefaultTimestampUTC,
			"caller", gokitLog.DefaultCaller,
		)
	}

	grpcServer := &grpcMock{}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(initDialer(grpcServer)))
	if err != nil {
		return nil, err
	}

	repo := NewUserRepository(conn, logger)

	return repo, nil
}

func (m *grpcMock) Authenticate(ctx context.Context, req *pb.UserAuthRequest) (*pb.UserAuthResponse, error) {
	if ForceMockFail {
		return nil, ErrForcedFailure
	}

	return &pb.UserAuthResponse{
		Message: "User Authenticated!",
	}, nil
}

func (m *grpcMock) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if ForceMockFail {
		return nil, ErrForcedFailure
	}

	response := &pb.CreateUserResponse{
		Id:   "USR123",
		Name: req.Name,
		Age:  "99",
	}

	if ForceBadAge {
		response.Age = "a"
	}

	return response, nil
}

func (m *grpcMock) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if ForceMockFail {
		return nil, ErrForcedFailure
	}

	response := &pb.GetUserResponse{
		Id:   "USR123",
		Name: "test",
		Age:  "99",
	}

	if ForceBadAge {
		response.Age = "a"
	}

	return response, nil
}

func (m *grpcMock) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if ForceMockFail {
		return nil, ErrForcedFailure
	}

	response := &pb.UpdateUserResponse{
		Id:   "USR123",
		Name: req.Name,
		Age:  "99",
	}

	if ForceBadAge {
		response.Age = "a"
	}

	return response, nil
}

func (m *grpcMock) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if ForceMockFail {
		return nil, ErrForcedFailure
	}

	return &pb.DeleteUserResponse{
		Message: "user deleted successfully",
	}, nil
}
