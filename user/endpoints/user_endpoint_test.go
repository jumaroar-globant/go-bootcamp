package endpoints

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/log"
	"github.com/jumaroar-globant/go-bootcamp/shared"
	"github.com/jumaroar-globant/go-bootcamp/user/config"
	"github.com/jumaroar-globant/go-bootcamp/user/pb"
	"github.com/jumaroar-globant/go-bootcamp/user/repository"
	"github.com/jumaroar-globant/go-bootcamp/user/service"

	sharedLib "github.com/jumaroar-globant/go-bootcamp/user/shared"

	"github.com/stretchr/testify/require"
)

func TestMakeCreateUserEndpoint(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	svc := service.NewUserService(repository.NewUserRepository(db, logger), logger)

	createUserEndpoint := makeCreateUserEndpoint(svc)

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

	result, err := createUserEndpoint(context.Background(), req)

	c.Equal(req.Name, result.(shared.User).Name)
	c.NoError(err)

	_, err = createUserEndpoint(context.Background(), "bad request")
	c.Equal(errBadRequest, err)
}

func TestMakeAuthenticateEndpoint(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	svc := service.NewUserService(repository.NewUserRepository(db, logger), logger)

	authenticatendpoint := makeAuthenticateEndpoint(svc)

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

	result, err := authenticatendpoint(context.Background(), req)

	c.Equal("User authenticated!", result.(string))
	c.NoError(err)

	_, err = authenticatendpoint(context.Background(), "bad request")
	c.Equal(errBadRequest, err)
}

func TestMakeGetUserEndpoint(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	svc := service.NewUserService(repository.NewUserRepository(db, logger), logger)

	getuserendpoint := makeGetUserEndpoint(svc)

	req := &pb.GetUserRequest{
		Id: "USR123",
	}

	user := shared.User{
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

	result, err := getuserendpoint(context.Background(), req)

	c.Equal(user, result.(shared.User))
	c.NoError(err)

	_, err = getuserendpoint(context.Background(), "bad request")
	c.Equal(errBadRequest, err)
}

func TestMakeUpdateUserEndpoint(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	svc := service.NewUserService(repository.NewUserRepository(db, logger), logger)

	updateendpoint := makeUpdateUserEndpoint(svc)

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

	result, err := updateendpoint(context.Background(), user)

	c.Equal(user.Name, result.(shared.User).Name)
	c.NoError(err)

	_, err = updateendpoint(context.Background(), "bad request")
	c.Equal(errBadRequest, err)
}

func TestMakeDeleteEndpoint(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	svc := service.NewUserService(repository.NewUserRepository(db, logger), logger)

	deletendpoint := makeDeleteUserEndpoint(svc)

	req := &pb.DeleteUserRequest{
		Id: "USR123",
	}

	sqlString := regexp.QuoteMeta(repository.DeleteUserParentsStatement)
	mock.ExpectExec(sqlString).WithArgs("USR123").WillReturnResult(sqlmock.NewResult(0, 1))

	parentSSQLString := regexp.QuoteMeta(repository.DeleteUserStatement)
	mock.ExpectExec(parentSSQLString).WithArgs("USR123").WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := deletendpoint(context.Background(), req)

	c.Equal("user deleted successfully", result.(string))
	c.NoError(err)

	_, err = deletendpoint(context.Background(), "bad request")
	c.Equal(errBadRequest, err)
}

func TestMakeEndpoints(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	logger := log.NewJSONLogger(os.Stdout)

	svc := service.NewUserService(repository.NewUserRepository(db, logger), logger)

	endpoints := MakeEndpoints(svc)
	c.IsType(UserEndpoints{}, endpoints)

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

	result, err := endpoints.Authenticate(context.Background(), req)

	c.Equal("User authenticated!", result.(string))
	c.NoError(err)
}
