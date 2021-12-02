package userrepository

import (
	"context"
	"testing"

	"github.com/jumaroar-globant/go-bootcamp/shared"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	c := require.New(t)

	repo, err := InitGRPCMock()
	c.Nil(err)

	authResponse, err := repo.Authenticate(context.Background(), "test", "testPassword")
	c.NoError(err)
	c.Equal("User Authenticated!", authResponse)

	forceMockFail = true
	defer func() {
		forceMockFail = false
	}()

	authResponse, err = repo.Authenticate(context.Background(), "test", "testPassWord")
	c.Empty(authResponse)
	c.Error(err, "rpc error: code = Unknown desc = forced failure")
}

func TestCreateUser(t *testing.T) {
	c := require.New(t)

	repo, err := InitGRPCMock()
	c.Nil(err)

	createResponse, err := repo.CreateUser(context.Background(), shared.User{Name: "test"})
	c.NoError(err)
	c.Equal("USR123", createResponse.ID)
	c.Equal("test", createResponse.Name)

	forceBadAge = true
	defer func() {
		forceBadAge = false
	}()

	createResponse, err = repo.CreateUser(context.Background(), shared.User{Name: "test"})
	c.Empty(createResponse)
	c.Error(err, "rpc error: code = Unknown desc = age is not a number")

	forceMockFail = true
	defer func() {
		forceMockFail = false
	}()

	createResponse, err = repo.CreateUser(context.Background(), shared.User{Name: "test"})
	c.Empty(createResponse)
	c.Error(err, "rpc error: code = Unknown desc = forced failure")
}

func TestGetUser(t *testing.T) {
	c := require.New(t)

	repo, err := InitGRPCMock()
	c.Nil(err)

	getResponse, err := repo.GetUser(context.Background(), "USR123")
	c.NoError(err)
	c.Equal("USR123", getResponse.ID)
	c.Equal("test", getResponse.Name)

	forceBadAge = true
	defer func() {
		forceBadAge = false
	}()

	getResponse, err = repo.GetUser(context.Background(), "USR123")
	c.Empty(getResponse)
	c.Error(err, "rpc error: code = Unknown desc = age is not a number")

	forceMockFail = true
	defer func() {
		forceMockFail = false
	}()

	getResponse, err = repo.GetUser(context.Background(), "USR123")
	c.Empty(getResponse)
	c.Error(err, "rpc error: code = Unknown desc = forced failure")
}

func TestUpdateUser(t *testing.T) {
	c := require.New(t)

	repo, err := InitGRPCMock()
	c.Nil(err)

	updateResponse, err := repo.UpdateUser(context.Background(), shared.User{Name: "test"})
	c.NoError(err)
	c.Equal("USR123", updateResponse.ID)
	c.Equal("test", updateResponse.Name)

	forceBadAge = true
	defer func() {
		forceBadAge = false
	}()

	updateResponse, err = repo.UpdateUser(context.Background(), shared.User{Name: "test"})
	c.Empty(updateResponse)
	c.Error(err, "rpc error: code = Unknown desc = age is not a number")

	forceMockFail = true
	defer func() {
		forceMockFail = false
	}()

	updateResponse, err = repo.UpdateUser(context.Background(), shared.User{Name: "test"})
	c.Empty(updateResponse)
	c.Error(err, "rpc error: code = Unknown desc = forced failure")
}

func TestDeleteUser(t *testing.T) {
	c := require.New(t)

	repo, err := InitGRPCMock()
	c.Nil(err)

	deleteResponse, err := repo.DeleteUser(context.Background(), "USR123")
	c.NoError(err)
	c.Equal("user deleted successfully", deleteResponse)

	forceMockFail = true
	defer func() {
		forceMockFail = false
	}()

	deleteResponse, err = repo.DeleteUser(context.Background(), "USR123")
	c.Empty(deleteResponse)
	c.Error(err, "rpc error: code = Unknown desc = forced failure")
}
