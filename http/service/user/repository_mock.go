package userservice

import (
	"context"
	"errors"

	userrepository "github.com/jumaroar-globant/go-bootcamp/http/repository/user"
	"github.com/jumaroar-globant/go-bootcamp/shared"
)

var (
	// ErrForcedFailure is an error for a forced failure
	ErrForcedFailure = errors.New("forced failure")
	// ForceMockFail is a variable to fail a mock
	ForceMockFail = false
)

type repoMock struct {
	userrepository.UserRepository
}

func (m *repoMock) Authenticate(ctx context.Context, username string, password string) (string, error) {
	if ForceMockFail {
		return "", ErrForcedFailure
	}

	return "User Authenticated!", nil
}

func (m *repoMock) CreateUser(ctx context.Context, user *shared.User) (*shared.User, error) {
	if ForceMockFail {
		return nil, ErrForcedFailure
	}

	return &shared.User{
		ID:   "USR123",
		Name: "test",
	}, nil
}

func (m *repoMock) GetUser(ctx context.Context, userID string) (*shared.User, error) {
	if ForceMockFail {
		return nil, ErrForcedFailure
	}

	return &shared.User{
		ID:   userID,
		Name: "test",
	}, nil
}

func (m *repoMock) UpdateUser(ctx context.Context, user *shared.User) (*shared.User, error) {
	if ForceMockFail {
		return nil, ErrForcedFailure
	}

	return &shared.User{
		ID:   "USR123",
		Name: "test",
	}, nil
}

func (m *repoMock) DeleteUser(ctx context.Context, userID string) (string, error) {
	if ForceMockFail {
		return "", ErrForcedFailure
	}

	return "user deleted successfully", nil
}
