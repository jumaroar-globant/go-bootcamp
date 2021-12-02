package userendpoints

import (
	"context"
	"errors"

	userservice "github.com/jumaroar-globant/go-bootcamp/http/service/user"
	"github.com/jumaroar-globant/go-bootcamp/shared"
)

var (
	errForcedFailure = errors.New("forced failure")
	forceMockFail    = false
)

type serviceMock struct {
	userservice.Service
}

func (m *serviceMock) Authenticate(ctx context.Context, username string, password string) (string, error) {
	if forceMockFail {
		return "", errForcedFailure
	}

	return "User Authenticated!", nil
}

func (m *serviceMock) CreateUser(ctx context.Context, user shared.User) (shared.User, error) {
	if forceMockFail {
		return shared.User{}, errForcedFailure
	}

	return user, nil
}

func (m *serviceMock) GetUser(ctx context.Context, userID string) (shared.User, error) {
	if forceMockFail {
		return shared.User{}, errForcedFailure
	}

	return shared.User{}, nil
}

func (m *serviceMock) UpdateUser(ctx context.Context, user shared.User) (shared.User, error) {
	if forceMockFail {
		return shared.User{}, errForcedFailure
	}

	return user, nil
}

func (m *serviceMock) DeleteUser(ctx context.Context, userID string) (string, error) {
	if forceMockFail {
		return "", errForcedFailure
	}

	return "user deleted successfully", nil
}
