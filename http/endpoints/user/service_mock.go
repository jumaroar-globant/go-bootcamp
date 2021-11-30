package userendpoints

import (
	"context"

	userservice "github.com/jumaroar-globant/go-bootcamp/http/service/user"
	"github.com/jumaroar-globant/go-bootcamp/shared"
)

type serviceMock struct {
	userservice.Service
}

func (m *serviceMock) Authenticate(ctx context.Context, username string, password string) (string, error) {
	return "User Authenticated!", nil
}

func (m *serviceMock) CreateUser(ctx context.Context, user *shared.User) (*shared.User, error) {
	return user, nil
}

func (m *serviceMock) GetUser(ctx context.Context, userID string) (*shared.User, error) {
	return &shared.User{}, nil
}

func (m *serviceMock) UpdateUser(ctx context.Context, user *shared.User) (*shared.User, error) {
	return user, nil
}

func (m *serviceMock) DeleteUser(ctx context.Context, userID string) (string, error) {
	return "user deleted successfully", nil
}
