package repository

import (
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/log"
	"github.com/jumaroar-globant/go-bootcamp/user/config"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	userRepo := NewUserRepository(db, log.NewJSONLogger(os.Stdout))

	user := &User{
		ID:                    "USR123",
		Name:                  "test",
		PasswordHash:          "clave123",
		Age:                   99,
		AdditionalInformation: "not much",
		Parents:               []string{"John Doe", "Jane Doe"},
	}

	sqlString := regexp.QuoteMeta(`INSERT INTO users (id, name, password_hash, age, additional_information) VALUES($1, $2, $3, $4, $5)`)
	mock.ExpectExec(sqlString).WithArgs(user.ID, user.Name, user.PasswordHash, user.Age, user.AdditionalInformation).WillReturnResult(sqlmock.NewResult(0, 1))

	parentSSQLString := regexp.QuoteMeta(`INSERT INTO user_parents (user_id, name) VALUES($1, $2)`)
	mock.ExpectExec(parentSSQLString).WithArgs(user.ID, user.Parents[0]).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(parentSSQLString).WithArgs(user.ID, user.Parents[1]).WillReturnResult(sqlmock.NewResult(0, 1))

	err := userRepo.CreateUser(context.Background(), user)
	c.NoError(err)
}
