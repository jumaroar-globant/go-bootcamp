package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/jumaroar-globant/go-bootcamp/user/config"
	"github.com/jumaroar-globant/go-bootcamp/user/shared"
	"github.com/stretchr/testify/require"

	sharedLib "github.com/jumaroar-globant/go-bootcamp/shared"
)

func TestMongoAuthenticate(t *testing.T) {
	c := require.New(t)

	passwordHash, err := shared.HashPassword("testPassword")
	c.Nil(err)

	config.MockedItem = []byte(fmt.Sprintf("{\"password\":\"%s\"}", passwordHash))

	repo := NewMongoUserRepository(config.MockMongoDbHelper{}, log.NewJSONLogger(os.Stdout))

	err = repo.Authenticate(context.Background(), "testUsername", "testPassword")
	c.Nil(err)

	config.MockedItem = nil
	config.ForceNotFound = true

	err = repo.Authenticate(context.Background(), "testUsername", "testPassword")
	c.Equal(ErrUserNotFound, err)

	config.ForceMockFail = true
	config.ForceNotFound = false

	err = repo.Authenticate(context.Background(), "testUsername", "testPassword")
	c.Equal(config.ErrForcedFailure, err)

	config.ForceMockFail = false

	config.MockedItem = []byte(fmt.Sprintf("{\"password\":\"%s\"", passwordHash))

	err = repo.Authenticate(context.Background(), "testUsername", "testPassword")
	c.IsType(&json.SyntaxError{}, err)

	config.MockedItem = []byte(fmt.Sprintf("{\"password\":\"%s\"}", passwordHash))
	err = repo.Authenticate(context.Background(), "testUsername", "badPassword")
	c.Equal(ErrWrongPassword, err)
}

func TestMongoCreateUser(t *testing.T) {
	c := require.New(t)

	repo := NewMongoUserRepository(config.MockMongoDbHelper{}, log.NewJSONLogger(os.Stdout))

	req := sharedLib.User{
		Name: "test",
		Age:  99,
	}

	config.ForceMockFail = true

	err := repo.CreateUser(context.Background(), req)
	c.Equal(config.ErrForcedFailure, err)

	config.ForceMockFail = false
	err = repo.CreateUser(context.Background(), req)
	c.Nil(err)
}

func TestMongoGetUser(t *testing.T) {
	c := require.New(t)

	repo := NewMongoUserRepository(config.MockMongoDbHelper{}, log.NewJSONLogger(os.Stdout))

	config.ForceNotFound = true

	user, err := repo.GetUser(context.Background(), "USR123")
	c.Empty(user)
	c.Equal(ErrUserNotFound, err)

	config.ForceNotFound = false
	config.ForceMockFail = true

	user, err = repo.GetUser(context.Background(), "USR123")
	c.Empty(user)
	c.Equal(config.ErrForcedFailure, err)

	config.MockedItem = []byte("{\"name\":\"test\"}")

	config.ForceMockFail = false

	user, err = repo.GetUser(context.Background(), "USR123")
	c.Nil(err)
	c.Equal("test", user.Name)

	config.ForceDecodeError = true
	defer func() {
		config.ForceDecodeError = false
	}()

	user, err = repo.GetUser(context.Background(), "USR123")
	c.Empty(user)
	c.Equal(config.ErrForcedFailure, err)
}

func TestMongoUpdateUser(t *testing.T) {
	c := require.New(t)

	repo := NewMongoUserRepository(config.MockMongoDbHelper{}, log.NewJSONLogger(os.Stdout))

	config.ForceNotFound = true

	req := sharedLib.User{
		Name: "test",
		Age:  99,
	}

	user, err := repo.UpdateUser(context.Background(), req)
	c.Empty(user)
	c.Equal(ErrUserNotFound, err)

	config.ForceNotFound = false
	config.ForceMockFail = true

	user, err = repo.UpdateUser(context.Background(), req)
	c.Empty(user)
	c.Equal(config.ErrForcedFailure, err)

	config.MockedItem = []byte("{\"name\":\"test\"}")
	config.ForceMockFail = false

	user, err = repo.UpdateUser(context.Background(), req)
	c.Equal("test", user.Name)
	c.Nil(err)

	config.ForceDecodeError = true
	defer func() {
		config.ForceDecodeError = false
	}()

	user, err = repo.UpdateUser(context.Background(), req)
	c.Empty(user)
	c.Equal(config.ErrForcedFailure, err)
}

func TestMongoDeleteUser(t *testing.T) {
	c := require.New(t)

	repo := NewMongoUserRepository(config.MockMongoDbHelper{}, log.NewJSONLogger(os.Stdout))

	config.ForceNotFound = true

	err := repo.DeleteUser(context.Background(), "USR123")
	c.Equal(ErrUserNotFound, err)

	config.ForceNotFound = false
	config.MockedItem = []byte("{\"name\":\"test\"}")

	config.ForceDeleteFail = true
	err = repo.DeleteUser(context.Background(), "USR123")
	c.Equal(config.ErrForcedFailure, err)

	config.ForceDeleteFail = false

	err = repo.DeleteUser(context.Background(), "USR123")
	c.Nil(err)
}
