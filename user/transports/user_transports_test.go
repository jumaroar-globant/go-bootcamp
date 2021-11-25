package transports

import (
	"context"
	"os"
	"regexp"
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

	sqlString := regexp.QuoteMeta(`SELECT password_hash FROM users WHERE name=?`)
	mock.ExpectQuery(sqlString).WithArgs(username).WillReturnRows(row)

	result, err := grpcServer.Authenticate(context.Background(), req)

	c.Equal("User authenticated!", result.Message)
	c.NoError(err)
}
