package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/jumaroar-globant/go-bootcamp/user/pb"
	"github.com/jumaroar-globant/go-bootcamp/user/service"
)

type UserEndpoints struct {
	Authenticate endpoint.Endpoint
	CreateUser   endpoint.Endpoint
	GetUser      endpoint.Endpoint
	UpdateUser   endpoint.Endpoint
	DeleteUser   endpoint.Endpoint
}

// MakeEndpoints func initializes the Endpoint instances
func MakeEndpoints(s service.UserService) UserEndpoints {
	return UserEndpoints{
		Authenticate: makeAuthenticateEndpoint(s),
		CreateUser:   makeCreateUserEndpoint(s),
		GetUser:      makeGetUserEndpoint(s),
		UpdateUser:   makeUpdateUserEndpoint(s),
		DeleteUser:   makeDeleteUserEndpoint(s),
	}
}

func makeCreateUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.CreateUserRequest)
		result, _ := s.CreateUser(ctx, req)
		return result, nil
	}
}

func makeAuthenticateEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UserAuthRequest)
		result, _ := s.Authenticate(ctx, req)
		return result, nil
	}
}

func makeGetUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.GetUserRequest)
		result, _ := s.GetUser(ctx, req)
		return result, nil
	}
}

func makeUpdateUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UpdateUserRequest)
		result, _ := s.UpdateUser(ctx, req)
		return result, nil
	}
}

func makeDeleteUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.DeleteUserRequest)
		result, _ := s.DeleteUser(ctx, req)
		return result, nil
	}
}
