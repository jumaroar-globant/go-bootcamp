package service

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/log"
	sharedLib "github.com/jumaroar-globant/go-bootcamp/shared"
	"github.com/jumaroar-globant/go-bootcamp/user/config"
	"github.com/jumaroar-globant/go-bootcamp/user/pb"
	"github.com/jumaroar-globant/go-bootcamp/user/repository"
	"github.com/jumaroar-globant/go-bootcamp/user/shared"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(db, logger), logger)

	username := "testUsername"
	passwordHash, err := shared.HashPassword("testPassword")

	c.NoError(err)

	req := &pb.UserAuthRequest{
		Username: username,
		Password: "testPassword",
	}

	row := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)

	sqlString := regexp.QuoteMeta(repository.PasswordHashQuery)
	mock.ExpectQuery(sqlString).WithArgs(username).WillReturnRows(row)

	message, err := service.Authenticate(context.Background(), req)
	c.Equal("User authenticated!", message)
	c.NoError(err)
}

func TestAuthenticateFails(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(db, logger), logger)

	username := "testUsername"

	req := &pb.UserAuthRequest{
		Username: username,
		Password: "testPassword",
	}

	sqlString := regexp.QuoteMeta(repository.PasswordHashQuery)
	mock.ExpectQuery(sqlString).WithArgs(username).WillReturnError(config.ErrMockFails)

	message, err := service.Authenticate(context.Background(), req)
	c.Equal("", message)
	c.Equal(config.ErrMockFails, err)
}

func TestCreateUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(db, logger), logger)

	user := &pb.CreateUserRequest{
		Name:                  "test",
		Password:              "clave123",
		Age:                   "99",
		AdditionalInformation: "not much",
		Parent:                []string{"John Doe", "Jane Doe"},
	}

	intAge, err := strconv.Atoi(user.Age)
	c.NoError(err)

	sqlString := regexp.QuoteMeta(repository.InsertUserStatement)
	mock.ExpectExec(sqlString).WithArgs(sqlmock.AnyArg(), user.Name, sqlmock.AnyArg(), intAge, user.AdditionalInformation).WillReturnResult(sqlmock.NewResult(0, 1))

	parentSSQLString := regexp.QuoteMeta(repository.InsertParentStatement)
	mock.ExpectExec(parentSSQLString).WithArgs(sqlmock.AnyArg(), user.Parent[0]).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(parentSSQLString).WithArgs(sqlmock.AnyArg(), user.Parent[1]).WillReturnResult(sqlmock.NewResult(0, 1))

	savedUser, err := service.CreateUser(context.Background(), user)
	c.Equal(user.Name, savedUser.Name)
	c.NoError(err)
}

func TestCreateUserValidationsFails(t *testing.T) {
	c := require.New(t)

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(nil, logger), logger)

	user := &pb.CreateUserRequest{}

	savedUser, err := service.CreateUser(context.Background(), user)
	c.Empty(savedUser)
	c.Equal(ErrMissingUserName, err)

	user.Name = "test"

	savedUser, err = service.CreateUser(context.Background(), user)
	c.Empty(savedUser)
	c.Equal(ErrMissingPassword, err)

	user.Password = "testPwd"
	user.Age = "badAge"

	savedUser, err = service.CreateUser(context.Background(), user)
	c.Empty(savedUser)
	c.IsType(&strconv.NumError{}, err)
}

func TestCreateUserDatabaseFails(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(db, logger), logger)

	user := &pb.CreateUserRequest{
		Name:                  "test",
		Password:              "clave123",
		Age:                   "99",
		AdditionalInformation: "not much",
		Parent:                []string{"John Doe", "Jane Doe"},
	}

	intAge, err := strconv.Atoi(user.Age)
	c.NoError(err)

	sqlString := regexp.QuoteMeta(repository.InsertUserStatement)
	mock.ExpectExec(sqlString).WithArgs(sqlmock.AnyArg(), user.Name, sqlmock.AnyArg(), intAge, user.AdditionalInformation).WillReturnError(config.ErrMockFails)

	parentSSQLString := regexp.QuoteMeta(repository.InsertParentStatement)
	mock.ExpectExec(parentSSQLString).WithArgs(sqlmock.AnyArg(), user.Parent[0]).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(parentSSQLString).WithArgs(sqlmock.AnyArg(), user.Parent[1]).WillReturnResult(sqlmock.NewResult(0, 1))

	savedUser, err := service.CreateUser(context.Background(), user)
	c.Empty(savedUser)
	c.Equal(config.ErrMockFails, err)
}

func TestUpdateUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(db, logger), logger)

	user := &pb.UpdateUserRequest{
		Id:                    "USR123",
		Name:                  "test",
		Age:                   "99",
		AdditionalInformation: "not much",
		Parent:                []string{"John Doe", "Jane Doe"},
	}

	intAge, err := strconv.Atoi(user.Age)
	c.NoError(err)

	sqlUpdateString := regexp.QuoteMeta(repository.UpdateUserStatement)

	mock.ExpectExec(sqlUpdateString).WithArgs(user.Name, intAge, user.AdditionalInformation, user.Id).WillReturnResult(sqlmock.NewResult(0, 1))

	sqlDeleteString := regexp.QuoteMeta(repository.DeleteUserParentsStatement)

	mock.ExpectExec(sqlDeleteString).WithArgs(user.Id).WillReturnResult(sqlmock.NewResult(0, 1))

	parentsSQLInsertString := regexp.QuoteMeta(repository.InsertParentStatement)

	mock.ExpectExec(parentsSQLInsertString).WithArgs(user.Id, user.Parent[0]).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(parentsSQLInsertString).WithArgs(user.Id, user.Parent[1]).WillReturnResult(sqlmock.NewResult(0, 1))

	row := sqlmock.NewRows([]string{"id", "name", "age", "additional_information"}).AddRow(user.Id, user.Name, user.Age, user.AdditionalInformation)

	sqlSelectString := regexp.QuoteMeta(repository.UserDataQuery)

	mock.ExpectQuery(sqlSelectString).WithArgs(user.Id).WillReturnRows(row)

	parentSSQLString := regexp.QuoteMeta(repository.UserParentsQuery)

	rows := sqlmock.NewRows([]string{"name"}).AddRow(user.Parent[0]).AddRow(user.Parent[1])

	mock.ExpectQuery(parentSSQLString).WithArgs(user.Id).WillReturnRows(rows)

	savedUser, err := service.UpdateUser(context.Background(), user)
	c.Equal(user.Name, savedUser.Name)
	c.NoError(err)
}

func TestUpdateUserValidationsFails(t *testing.T) {
	c := require.New(t)

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(nil, logger), logger)

	user := &pb.UpdateUserRequest{}

	savedUser, err := service.UpdateUser(context.Background(), user)
	c.Empty(savedUser)
	c.Equal(ErrMissingUserID, err)

	user.Id = "USR123"

	user.Age = "badAge"

	savedUser, err = service.UpdateUser(context.Background(), user)
	c.Empty(savedUser)
	c.IsType(&strconv.NumError{}, err)
}

func TestUpdateUserDbFails(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(db, logger), logger)

	user := &pb.UpdateUserRequest{
		Id:                    "USR123",
		Name:                  "test",
		Age:                   "99",
		AdditionalInformation: "not much",
		Parent:                []string{"John Doe", "Jane Doe"},
	}

	intAge, err := strconv.Atoi(user.Age)
	c.NoError(err)

	sqlUpdateString := regexp.QuoteMeta(repository.UpdateUserStatement)

	mock.ExpectExec(sqlUpdateString).WithArgs(user.Name, intAge, user.AdditionalInformation, user.Id).WillReturnError(config.ErrMockFails)

	savedUser, err := service.UpdateUser(context.Background(), user)
	c.Empty(savedUser)
	c.IsType(config.ErrMockFails, err)
}

func TestGetUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(db, logger), logger)

	req := &pb.GetUserRequest{
		Id: "USR123",
	}

	user := sharedLib.User{
		ID:                    "USR123",
		Name:                  "test",
		Age:                   99,
		AdditionalInformation: "not much",
		Parents:               []string{"John Doe", "Jane Doe"},
	}

	row := sqlmock.NewRows([]string{"id", "name", "age", "additional_information"}).AddRow(user.ID, user.Name, user.Age, user.AdditionalInformation)

	sqlString := regexp.QuoteMeta(repository.UserDataQuery)
	mock.ExpectQuery(sqlString).WithArgs(req.Id).WillReturnRows(row)

	parentSSQLString := regexp.QuoteMeta(repository.UserParentsQuery)

	rows := sqlmock.NewRows([]string{"name"}).AddRow(user.Parents[0]).AddRow(user.Parents[1])

	mock.ExpectQuery(parentSSQLString).WithArgs(req.Id).WillReturnRows(rows)

	foundUser, err := service.GetUser(context.Background(), req)
	c.Equal(foundUser, user)
	c.Nil(err)
}

func TestGetUserFails(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(db, logger), logger)

	req := &pb.GetUserRequest{
		Id: "USR123",
	}

	sqlString := regexp.QuoteMeta(repository.UserDataQuery)
	mock.ExpectQuery(sqlString).WithArgs(req.Id).WillReturnError(config.ErrMockFails)

	foundUser, err := service.GetUser(context.Background(), req)
	c.Empty(foundUser)
	c.Equal(config.ErrMockFails, err)

	req.Id = ""

	foundUser, err = service.GetUser(context.Background(), req)
	c.Empty(foundUser)
	c.Equal(ErrMissingUserID, err)
}

func TestDeleteUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(db, logger), logger)

	req := &pb.DeleteUserRequest{
		Id: "USR123",
	}

	sqlString := regexp.QuoteMeta(repository.DeleteUserParentsStatement)
	mock.ExpectExec(sqlString).WithArgs("USR123").WillReturnResult(sqlmock.NewResult(0, 1))

	parentSSQLString := regexp.QuoteMeta(repository.DeleteUserStatement)
	mock.ExpectExec(parentSSQLString).WithArgs("USR123").WillReturnResult(sqlmock.NewResult(0, 1))

	message, err := service.DeleteUser(context.Background(), req)
	c.Equal(userDeletedString, message)
	c.Nil(err)
}

func TestDeleteUserFails(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	service := NewUserService(repository.NewUserRepository(db, logger), logger)

	req := &pb.DeleteUserRequest{
		Id: "USR123",
	}

	sqlString := regexp.QuoteMeta(repository.DeleteUserParentsStatement)
	mock.ExpectExec(sqlString).WithArgs("USR123").WillReturnError(config.ErrMockFails)

	message, err := service.DeleteUser(context.Background(), req)
	c.Equal("", message)
	c.Equal(config.ErrMockFails, err)

	req.Id = ""

	message, err = service.DeleteUser(context.Background(), req)
	c.Equal("", message)
	c.Equal(ErrMissingUserID, err)
}
