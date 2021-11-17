package userendpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	userservice "github.com/jumaroar-globant/go-bootcamp/http/service/user"
)

type UserEndpoints struct {
	Authenticate endpoint.Endpoint
}

type AuthenticationRequest struct {
	Username string
	Password string
}

type AuthenticationResponse struct {
	Message string
}

func MakeEndpoints(s userservice.Service) *UserEndpoints {
	return &UserEndpoints{
		Authenticate: makeAuthenticationEndpoint(s),
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
