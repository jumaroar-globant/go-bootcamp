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

	ForceMockFail = true
	defer func() {
		ForceMockFail = false
	}()

	authResponse, err = repo.Authenticate(context.Background(), "test", "testPassWord")
	c.Empty(authResponse)
	c.Error(err, "rpc error: code = Unknown desc = forced failure")
}

func TestCreateUser(t *testing.T) {
	c := require.New(t)

	repo, err := InitGRPCMock()
	c.Nil(err)

	createResponse, err := repo.CreateUser(context.Background(), &shared.User{Name: "test"})
	c.NoError(err)
	c.Equal("USR123", createResponse.ID)
	c.Equal("test", createResponse.Name)

	ForceBadAge = true
	defer func() {
		ForceBadAge = false
	}()

	createResponse, err = repo.CreateUser(context.Background(), &shared.User{Name: "test"})
	c.Nil(createResponse)
	c.Error(err, "rpc error: code = Unknown desc = age is not a number")

	ForceMockFail = true
	defer func() {
		ForceMockFail = false
	}()

	createResponse, err = repo.CreateUser(context.Background(), &shared.User{Name: "test"})
	c.Nil(createResponse)
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

	ForceBadAge = true
	defer func() {
		ForceBadAge = false
	}()

	getResponse, err = repo.GetUser(context.Background(), "USR123")
	c.Nil(getResponse)
	c.Error(err, "rpc error: code = Unknown desc = age is not a number")

	ForceMockFail = true
	defer func() {
		ForceMockFail = false
	}()

	getResponse, err = repo.GetUser(context.Background(), "USR123")
	c.Nil(getResponse)
	c.Error(err, "rpc error: code = Unknown desc = forced failure")
}

func TestUpdateUser(t *testing.T) {
	c := require.New(t)

	repo, err := InitGRPCMock()
	c.Nil(err)

	updateResponse, err := repo.UpdateUser(context.Background(), &shared.User{Name: "test"})
	c.NoError(err)
	c.Equal("USR123", updateResponse.ID)
	c.Equal("test", updateResponse.Name)

	ForceBadAge = true
	defer func() {
		ForceBadAge = false
	}()

	updateResponse, err = repo.UpdateUser(context.Background(), &shared.User{Name: "test"})
	c.Nil(updateResponse)
	c.Error(err, "rpc error: code = Unknown desc = age is not a number")

	ForceMockFail = true
	defer func() {
		ForceMockFail = false
	}()

	updateResponse, err = repo.UpdateUser(context.Background(), &shared.User{Name: "test"})
	c.Nil(updateResponse)
	c.Error(err, "rpc error: code = Unknown desc = forced failure")
}

func TestDeleteUser(t *testing.T) {
	c := require.New(t)

	repo, err := InitGRPCMock()
	c.Nil(err)

	deleteResponse, err := repo.DeleteUser(context.Background(), "USR123")
	c.NoError(err)
	c.Equal("user deleted successfully", deleteResponse)

	ForceMockFail = true
	defer func() {
		ForceMockFail = false
	}()

	deleteResponse, err = repo.DeleteUser(context.Background(), "USR123")
	c.Empty(deleteResponse)
	c.Error(err, "rpc error: code = Unknown desc = forced failure")
}
