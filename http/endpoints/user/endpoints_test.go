package userendpoints

import (
	"context"
	"testing"

	"github.com/jumaroar-globant/go-bootcamp/shared"
	"github.com/stretchr/testify/require"
)

func TestMakeEndpoints(t *testing.T) {
	c := require.New(t)

	endpoints := MakeEndpoints(&serviceMock{})

	authEndpoint, err := endpoints.Authenticate(context.Background(), AuthenticationRequest{"test", "test"})
	c.NoError(err)
	c.Equal("User Authenticated!", authEndpoint.(AuthenticationResponse).Message)
}

func TestMakeAuthenticateEndpoint(t *testing.T) {
	c := require.New(t)

	service := &serviceMock{}

	endpoint := makeAuthenticationEndpoint(service)

	result, err := endpoint(context.Background(), AuthenticationRequest{"test", "test"})
	c.NoError(err)
	c.Equal("User Authenticated!", result.(AuthenticationResponse).Message)
}

func TestMakeCreateUserEndpoint(t *testing.T) {
	c := require.New(t)

	service := &serviceMock{}

	endpoint := makeCreateUserEndpoint(service)

	result, err := endpoint(context.Background(), shared.User{Name: "test"})
	c.NoError(err)
	c.Equal("test", result.(*shared.User).Name)
}

func TestMakeGetUserEndpoint(t *testing.T) {
	c := require.New(t)

	service := &serviceMock{}

	endpoint := makeGetUserEndpoint(service)

	result, err := endpoint(context.Background(), GetUserRequest{UserID: "USR123"})
	c.NoError(err)
	c.Equal("", result.(*shared.User).Name)
}

func TestMakeUpdateUserEndpoint(t *testing.T) {
	c := require.New(t)

	service := &serviceMock{}

	endpoint := makeUpdateUserEndpoint(service)

	result, err := endpoint(context.Background(), shared.User{Name: "test"})
	c.NoError(err)
	c.Equal("test", result.(*shared.User).Name)
}

func TestMakeDeleteUserEndpoint(t *testing.T) {
	c := require.New(t)

	service := &serviceMock{}

	endpoint := makeDeleteUserEndpoint(service)

	result, err := endpoint(context.Background(), DeleteUserRequest{UserID: "USR123"})
	c.NoError(err)
	c.Equal("user deleted successfully", result.(DeleteUserResponse).Message)
}
