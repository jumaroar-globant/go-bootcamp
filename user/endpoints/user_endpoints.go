package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/jumaroar-globant/go-bootcamp/user/pb"
	"github.com/jumaroar-globant/go-bootcamp/user/service"
)

// UserEndpoints are the user endpoints
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
		return s.CreateUser(ctx, req)
	}
}

func makeAuthenticateEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UserAuthRequest)
		return s.Authenticate(ctx, req)
	}
}

func makeGetUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.GetUserRequest)
		return s.GetUser(ctx, req)
	}
}

func makeUpdateUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UpdateUserRequest)
		return s.UpdateUser(ctx, req)
	}
}

func makeDeleteUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.DeleteUserRequest)
		return s.DeleteUser(ctx, req)
	}
}
