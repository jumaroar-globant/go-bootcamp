package userservice

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/jumaroar-globant/go-bootcamp/shared"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	c := require.New(t)

	service := NewService(&repoMock{}, log.NewJSONLogger(os.Stdout))

	result, err := service.Authenticate(context.Background(), "test", "test")
	c.NoError(err)
	c.Equal("User Authenticated!", result)

	forceMockFail = true
	defer func() {
		forceMockFail = false
	}()

	_, err = service.Authenticate(context.Background(), "test", "test")
	c.Equal(errForcedFailure, err)
}

func TestCreateUser(t *testing.T) {
	c := require.New(t)

	service := NewService(&repoMock{}, log.NewJSONLogger(os.Stdout))

	result, err := service.CreateUser(context.Background(), shared.User{Name: "test"})
	c.NoError(err)
	c.Equal("USR123", result.ID)

	forceMockFail = true
	defer func() {
		forceMockFail = false
	}()

	_, err = service.CreateUser(context.Background(), shared.User{Name: "test"})
	c.Equal(errForcedFailure, err)
}

func TestGetUser(t *testing.T) {
	c := require.New(t)

	service := NewService(&repoMock{}, log.NewJSONLogger(os.Stdout))

	result, err := service.GetUser(context.Background(), "USR123")
	c.NoError(err)
	c.Equal("USR123", result.ID)

	forceMockFail = true
	defer func() {
		forceMockFail = false
	}()

	_, err = service.GetUser(context.Background(), "USR123")
	c.Equal(errForcedFailure, err)
}

func TestUpdateuser(t *testing.T) {
	c := require.New(t)

	service := NewService(&repoMock{}, log.NewJSONLogger(os.Stdout))

	result, err := service.UpdateUser(context.Background(), shared.User{Name: "test"})
	c.NoError(err)
	c.Equal("USR123", result.ID)

	forceMockFail = true
	defer func() {
		forceMockFail = false
	}()

	_, err = service.UpdateUser(context.Background(), shared.User{Name: "test"})
	c.Equal(errForcedFailure, err)
}

func TestDeleteUser(t *testing.T) {
	c := require.New(t)

	service := NewService(&repoMock{}, log.NewJSONLogger(os.Stdout))

	result, err := service.DeleteUser(context.Background(), "USR123")
	c.NoError(err)
	c.Equal("user deleted successfully", result)

	forceMockFail = true
	defer func() {
		forceMockFail = false
	}()

	_, err = service.DeleteUser(context.Background(), "USR123")
	c.Equal(errForcedFailure, err)
}
