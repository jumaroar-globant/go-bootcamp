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

//MakeEndpoints creates the user endpoints
func MakeEndpoints(s userservice.Service) *UserEndpoints {
	return &UserEndpoints{
		Authenticate: makeAuthenticationEndpoint(s),
		CreateUser:   makeCreateUserEndpoint(s),
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
