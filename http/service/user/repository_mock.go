package userservice

import (
	"context"
	"errors"

	userrepository "github.com/jumaroar-globant/go-bootcamp/http/repository/user"
	"github.com/jumaroar-globant/go-bootcamp/shared"
)

var (
	errForcedFailure = errors.New("forced failure")
	forceMockFail    = false
)

type repoMock struct {
	userrepository.UserRepository
}

func (m *repoMock) Authenticate(ctx context.Context, username string, password string) (string, error) {
	if forceMockFail {
		return "", errForcedFailure
	}

	return "User Authenticated!", nil
}

func (m *repoMock) CreateUser(ctx context.Context, user shared.User) (shared.User, error) {
	if forceMockFail {
		return shared.User{}, errForcedFailure
	}

	return shared.User{
		ID:   "USR123",
		Name: "test",
	}, nil
}

func (m *repoMock) GetUser(ctx context.Context, userID string) (shared.User, error) {
	if forceMockFail {
		return shared.User{}, errForcedFailure
	}

	return shared.User{
		ID:   userID,
		Name: "test",
	}, nil
}

func (m *repoMock) UpdateUser(ctx context.Context, user shared.User) (shared.User, error) {
	if forceMockFail {
		return shared.User{}, errForcedFailure
	}

	return shared.User{
		ID:   "USR123",
		Name: "test",
	}, nil
}

func (m *repoMock) DeleteUser(ctx context.Context, userID string) (string, error) {
	if forceMockFail {
		return "", errForcedFailure
	}

	return "user deleted successfully", nil
}
