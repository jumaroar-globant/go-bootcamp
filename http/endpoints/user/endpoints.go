package userendpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	userservice "github.com/jumaroar-globant/go-bootcamp/http/service/user"
	"github.com/jumaroar-globant/go-bootcamp/shared"
)

//UserEndpoints are the user endpoints
type UserEndpoints struct {
	Authenticate endpoint.Endpoint
	CreateUser   endpoint.Endpoint
	GetUser      endpoint.Endpoint
	UpdateUser   endpoint.Endpoint
	DeleteUser   endpoint.Endpoint
}

//AuthenticationRequest is the authentication request
type AuthenticationRequest struct {
	Username string
	Password string
}

//AuthenticationResponse is the authentication response
type AuthenticationResponse struct {
	Message string
}

//GetUserRequest is the get user request
type GetUserRequest struct {
	UserID string
}

//DeleteUserRequest is the delete user request
type DeleteUserRequest struct {
	UserID string
}

//DeleteUserResponse is the delete user response
type DeleteUserResponse struct {
	Message string
}

//MakeEndpoints creates the user endpoints
func MakeEndpoints(s userservice.Service) *UserEndpoints {
	return &UserEndpoints{
		Authenticate: makeAuthenticationEndpoint(s),
		CreateUser:   makeCreateUserEndpoint(s),
		GetUser:      makeGetUserEndpoint(s),
		UpdateUser:   makeUpdateUserEndpoint(s),
		DeleteUser:   makeDeleteUserEndpoint(s),
	}
}

func makeAuthenticationEndpoint(s userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthenticationRequest)
		message, err := s.Authenticate(ctx, req.Username, req.Password)
		return AuthenticationResponse{
			Message: message,
		}, err
	}
}

func makeCreateUserEndpoint(s userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(shared.User)
		return s.CreateUser(ctx, &req)
	}
}

func makeGetUserEndpoint(s userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserRequest)
		return s.GetUser(ctx, req.UserID)
	}
}

func makeUpdateUserEndpoint(s userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(shared.User)
		return s.UpdateUser(ctx, &req)
	}
}

func makeDeleteUserEndpoint(s userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteUserRequest)
		message, err := s.DeleteUser(ctx, req.UserID)
		return DeleteUserResponse{
			Message: message,
		}, err
	}
}
