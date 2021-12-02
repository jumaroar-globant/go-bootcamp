package transports

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/log"
	"github.com/jumaroar-globant/go-bootcamp/user/config"
	"github.com/jumaroar-globant/go-bootcamp/user/endpoints"
	"github.com/jumaroar-globant/go-bootcamp/user/pb"
	"github.com/jumaroar-globant/go-bootcamp/user/repository"
	"github.com/jumaroar-globant/go-bootcamp/user/service"

	sharedLib "github.com/jumaroar-globant/go-bootcamp/user/shared"

	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	svc := service.NewUserService(repository.NewUserRepository(db, logger), logger)

	userEndpoints := endpoints.MakeEndpoints(svc)

	grpcServer := NewGRPCServer(userEndpoints, log.NewJSONLogger(os.Stdout))

	username := "testUsername"

	passwordHash, err := sharedLib.HashPassword("testPassword")
	c.NoError(err)

	req := &pb.UserAuthRequest{
		Username: username,
		Password: "testPassword",
	}

	row := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)

	sqlString := regexp.QuoteMeta(repository.PasswordHashQuery)
	mock.ExpectQuery(sqlString).WithArgs(username).WillReturnRows(row)

	result, err := grpcServer.Authenticate(context.Background(), req)

	c.Equal("User authenticated!", result.Message)
	c.NoError(err)

	mock.ExpectQuery(sqlString).WithArgs(username).WillReturnError(config.ErrMockFails)

	_, err = grpcServer.Authenticate(context.Background(), req)
	c.Equal(config.ErrMockFails, err)
}

func TestCreateUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	svc := service.NewUserService(repository.NewUserRepository(db, logger), logger)

	userEndpoints := endpoints.MakeEndpoints(svc)

	grpcServer := NewGRPCServer(userEndpoints, log.NewJSONLogger(os.Stdout))

	req := &pb.CreateUserRequest{
		Name:                  "test",
		Password:              "clave123",
		Age:                   "99",
		AdditionalInformation: "not much",
		Parent:                []string{"John Doe", "Jane Doe"},
	}

	intAge, err := strconv.Atoi(req.Age)
	c.NoError(err)

	sqlString := regexp.QuoteMeta(repository.InsertUserStatement)
	mock.ExpectExec(sqlString).WithArgs(sqlmock.AnyArg(), req.Name, sqlmock.AnyArg(), intAge, req.AdditionalInformation).WillReturnResult(sqlmock.NewResult(0, 1))

	parentSSQLString := regexp.QuoteMeta(repository.InsertParentStatement)
	mock.ExpectExec(parentSSQLString).WithArgs(sqlmock.AnyArg(), req.Parent[0]).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(parentSSQLString).WithArgs(sqlmock.AnyArg(), req.Parent[1]).WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := grpcServer.CreateUser(context.Background(), req)

	c.Equal(req.Name, result.Name)
	c.NoError(err)

	req.Name = ""
	_, err = grpcServer.CreateUser(context.Background(), req)
	c.Equal(service.ErrMissingUserName, err)
}

func TestGetUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	svc := service.NewUserService(repository.NewUserRepository(db, logger), logger)

	userEndpoints := endpoints.MakeEndpoints(svc)

	grpcServer := NewGRPCServer(userEndpoints, log.NewJSONLogger(os.Stdout))

	req := &pb.GetUserRequest{
		Id: "USR123",
	}

	user := &pb.GetUserResponse{
		Id:                    "USR123",
		Name:                  "test",
		Age:                   "99",
		AdditionalInformation: "not much",
		Parent:                []string{"John Doe", "Jane Doe"},
	}

	intAge, err := strconv.Atoi(user.Age)
	c.NoError(err)

	row := sqlmock.NewRows([]string{"id", "name", "age", "additional_information"}).AddRow(user.Id, user.Name, intAge, user.AdditionalInformation)

	sqlString := regexp.QuoteMeta(`SELECT id, name, age, additional_information FROM users WHERE id=?`)
	mock.ExpectQuery(sqlString).WithArgs(req.Id).WillReturnRows(row)

	parentSSQLString := regexp.QuoteMeta(`SELECT name FROM user_parents WHERE user_id=?`)

	rows := sqlmock.NewRows([]string{"name"}).AddRow(user.Parent[0]).AddRow(user.Parent[1])

	mock.ExpectQuery(parentSSQLString).WithArgs(req.Id).WillReturnRows(rows)

	result, err := grpcServer.GetUser(context.Background(), req)

	c.Equal(user, result)
	c.NoError(err)

	req.Id = ""

	_, err = grpcServer.GetUser(context.Background(), req)
	c.Equal(service.ErrMissingUserID, err)
}

func TestUpdateUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	svc := service.NewUserService(repository.NewUserRepository(db, logger), logger)

	userEndpoints := endpoints.MakeEndpoints(svc)

	grpcServer := NewGRPCServer(userEndpoints, log.NewJSONLogger(os.Stdout))

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

	result, err := grpcServer.UpdateUser(context.Background(), user)

	c.Equal(user.Name, result.Name)
	c.NoError(err)

	user.Id = ""

	_, err = grpcServer.UpdateUser(context.Background(), user)
	c.Equal(service.ErrMissingUserID, err)
}

func TestDeleteUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	svc := service.NewUserService(repository.NewUserRepository(db, logger), logger)

	userEndpoints := endpoints.MakeEndpoints(svc)

	grpcServer := NewGRPCServer(userEndpoints, log.NewJSONLogger(os.Stdout))

	req := &pb.DeleteUserRequest{
		Id: "USR123",
	}

	sqlString := regexp.QuoteMeta(repository.DeleteUserParentsStatement)
	mock.ExpectExec(sqlString).WithArgs("USR123").WillReturnResult(sqlmock.NewResult(0, 1))

	parentSSQLString := regexp.QuoteMeta(repository.DeleteUserStatement)
	mock.ExpectExec(parentSSQLString).WithArgs("USR123").WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := grpcServer.DeleteUser(context.Background(), req)

	c.Equal("user deleted successfully", result.Message)
	c.NoError(err)

	req.Id = ""

	_, err = grpcServer.DeleteUser(context.Background(), req)
	c.Equal(service.ErrMissingUserID, err)
}
