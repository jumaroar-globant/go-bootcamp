package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/jumaroar-globant/go-bootcamp/user/pb"
	"github.com/jumaroar-globant/go-bootcamp/user/service"
)

var (
	errBadRequest = errors.New("bad request")
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
		req, ok := request.(*pb.CreateUserRequest)
		if !ok {
			return nil, errBadRequest
		}

		return s.CreateUser(ctx, req)
	}
}

func makeAuthenticateEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.UserAuthRequest)
		if !ok {
			return nil, errBadRequest
		}

		return s.Authenticate(ctx, req)
	}
}

func makeGetUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.GetUserRequest)
		if !ok {
			return nil, errBadRequest
		}

		return s.GetUser(ctx, req)
	}
}

func makeUpdateUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.UpdateUserRequest)
		if !ok {
			return nil, errBadRequest
		}

		return s.UpdateUser(ctx, req)
	}
}

func makeDeleteUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.DeleteUserRequest)
		if !ok {
			return nil, errBadRequest
		}

		return s.DeleteUser(ctx, req)
	}
}
