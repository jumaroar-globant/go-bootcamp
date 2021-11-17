package transports

import (
	"context"
	"strconv"

	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"

	"github.com/jumaroar-globant/go-bootcamp/user/endpoints"
	"github.com/jumaroar-globant/go-bootcamp/user/pb"

	sharedLib "github.com/jumaroar-globant/go-bootcamp/shared"
)

type gRPCServer struct {
	pb.UnimplementedUserServiceServer
	authenticate gt.Handler
	createUser   gt.Handler
	getUser      gt.Handler
	updateUser   gt.Handler
	deleteUser   gt.Handler
}

// NewGRPCServer initializes a new gRPC server
func NewGRPCServer(endpoints endpoints.UserEndpoints, logger log.Logger) pb.UserServiceServer {
	return &gRPCServer{
		authenticate: gt.NewServer(
			endpoints.Authenticate,
			decodeAuthenticateRequest,
			encodeAuthenticateResponse,
		),
		createUser: gt.NewServer(
			endpoints.CreateUser,
			decodeCreateUserRequest,
			encodeCreateUserResponse,
		),
		getUser: gt.NewServer(
			endpoints.GetUser,
			decodeGetUserRequest,
			encodeGetUserResponse,
		),
		updateUser: gt.NewServer(
			endpoints.UpdateUser,
			decodeUpdateUserRequest,
			encodeUpdateUserResponse,
		),
		deleteUser: gt.NewServer(
			endpoints.DeleteUser,
			decodeDeleteUserRequest,
			encodeDeleteUserResponse,
		),
	}
}

// Authenticate is the gRPCServer method to authenticate
func (s *gRPCServer) Authenticate(ctx context.Context, req *pb.UserAuthRequest) (*pb.UserAuthResponse, error) {
	_, resp, err := s.authenticate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UserAuthResponse), nil
}

// CreateUser is the gRPCServer method to create a user
func (s *gRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateUserResponse), nil
}

// GetUser is the gRPCServer method to get a user
func (s *gRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetUserResponse), nil
}

// UpdateUser is the gRPCServer method to update a user
func (s *gRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UpdateUserResponse), nil
}

// DeleteUser is the gRPCServer method to delete a user
func (s *gRPCServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.DeleteUserResponse), nil
}

func decodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request.(*pb.CreateUserRequest), nil
}

func encodeCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*sharedLib.User)

	return &pb.CreateUserResponse{
		Id:                    resp.ID,
		Name:                  resp.Name,
		Age:                   strconv.Itoa(resp.Age),
		AdditionalInformation: resp.AdditionalInformation,
		Parent:                resp.Parents,
	}, nil
}

func decodeAuthenticateRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request.(*pb.UserAuthRequest), nil
}

func encodeAuthenticateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(string)
	return &pb.UserAuthResponse{
		Message: resp,
	}, nil
}

func decodeGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request.(*pb.GetUserRequest), nil
}

func encodeGetUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(sharedLib.User)

	return &pb.GetUserResponse{
		Id:                    resp.ID,
		Name:                  resp.Name,
		Age:                   strconv.Itoa(resp.Age),
		AdditionalInformation: resp.AdditionalInformation,
		Parent:                resp.Parents,
	}, nil
}

func decodeUpdateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request.(*pb.GetUserRequest), nil
}

func encodeUpdateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(sharedLib.User)

	return &pb.UpdateUserResponse{
		Id:                    resp.ID,
		Name:                  resp.Name,
		Age:                   strconv.Itoa(resp.Age),
		AdditionalInformation: resp.AdditionalInformation,
		Parent:                resp.Parents,
	}, nil
}

func decodeDeleteUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request.(*pb.GetUserRequest), nil
}

func encodeDeleteUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	return &pb.DeleteUserResponse{
		Message: response.(string),
	}, nil
}
