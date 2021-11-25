package repository

import (
	"context"
	"database/sql"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/log"

	sharedLib "github.com/jumaroar-globant/go-bootcamp/shared"
	"github.com/jumaroar-globant/go-bootcamp/user/config"
	"github.com/jumaroar-globant/go-bootcamp/user/shared"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	userRepo := NewUserRepository(db, log.NewJSONLogger(os.Stdout))

	username := "testUsername"
	passwordHash, err := shared.HashPassword("testPassword")

	c.NoError(err)

	row := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)

	sqlString := regexp.QuoteMeta(`SELECT password_hash FROM users WHERE name=?`)
	mock.ExpectQuery(sqlString).WithArgs(username).WillReturnRows(row)

	err = userRepo.Authenticate(context.Background(), username, "testPassword")
	c.NoError(err)
}

func TestAuthenticateFails(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	userRepo := NewUserRepository(db, log.NewJSONLogger(os.Stdout))

	username := "testUsername"
	passwordHash, err := shared.HashPassword("testPassword")

	c.NoError(err)

	row := sqlmock.NewRows([]string{"password_hash"}).AddRow(passwordHash)

	sqlString := regexp.QuoteMeta(`SELECT password_hash FROM users WHERE name=?`)
	mock.ExpectQuery(sqlString).WithArgs(username).WillReturnRows(row)

	err = userRepo.Authenticate(context.Background(), username, "testPassWord")
	c.Equal(ErrWrongPassword, err)

	mock.ExpectQuery(sqlString).WithArgs(username).WillReturnError(config.ErrMockFails)

	err = userRepo.Authenticate(context.Background(), username, "testpassWord")
	c.Equal(config.ErrMockFails, err)

	mock.ExpectQuery(sqlString).WithArgs(username).WillReturnError(sql.ErrNoRows)

	err = userRepo.Authenticate(context.Background(), username, "testpassWord")
	c.Equal(ErrUserNotFound, err)
}

func TestCreateUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	userRepo := NewUserRepository(db, log.NewJSONLogger(os.Stdout))

	user := &sharedLib.User{
		ID:                    "USR123",
		Name:                  "test",
		Password:              "clave123",
		Age:                   99,
		AdditionalInformation: "not much",
		Parents:               []string{"John Doe", "Jane Doe"},
	}

	sqlString := regexp.QuoteMeta(`INSERT INTO users (id, name, password_hash, age, additional_information) VALUES(?, ?, ?, ?, ?)`)
	mock.ExpectExec(sqlString).WithArgs(user.ID, user.Name, user.Password, user.Age, user.AdditionalInformation).WillReturnResult(sqlmock.NewResult(0, 1))

	parentSSQLString := regexp.QuoteMeta(`INSERT INTO user_parents (user_id, name) VALUES(?, ?)`)
	mock.ExpectExec(parentSSQLString).WithArgs(user.ID, user.Parents[0]).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(parentSSQLString).WithArgs(user.ID, user.Parents[1]).WillReturnResult(sqlmock.NewResult(0, 1))

	err := userRepo.CreateUser(context.Background(), user)
	c.NoError(err)
}

func TestCreateUserFails(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	userRepo := NewUserRepository(db, log.NewJSONLogger(os.Stdout))

	user := &sharedLib.User{
		ID:                    "USR123",
		Name:                  "test",
		Password:              "clave123",
		Age:                   99,
		AdditionalInformation: "not much",
		Parents:               []string{"John Doe", "Jane Doe"},
	}

	sqlString := regexp.QuoteMeta(`INSERT INTO users (id, name, password_hash, age, additional_information) VALUES(?, ?, ?, ?, ?)`)
	mock.ExpectExec(sqlString).WithArgs(user.ID, user.Name, user.Password, user.Age, user.AdditionalInformation).WillReturnError(config.ErrMockFails)

	err := userRepo.CreateUser(context.Background(), user)
	c.Equal(config.ErrMockFails, err)

	mock.ExpectExec(sqlString).WithArgs(user.ID, user.Name, user.Password, user.Age, user.AdditionalInformation).WillReturnResult(sqlmock.NewResult(0, 1))

	parentSSQLString := regexp.QuoteMeta(`INSERT INTO user_parents (user_id, name) VALUES(?, ?)`)
	mock.ExpectExec(parentSSQLString).WithArgs(user.ID, user.Parents[0]).WillReturnError(config.ErrMockFails)

	err = userRepo.CreateUser(context.Background(), user)
	c.Equal(config.ErrMockFails, err)
}

func TestGetUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	userRepo := NewUserRepository(db, log.NewJSONLogger(os.Stdout))

	user := &sharedLib.User{
		ID:                    "USR123",
		Name:                  "test",
		Age:                   99,
		AdditionalInformation: "not much",
		Parents:               []string{"John Doe", "Jane Doe"},
	}

	row := sqlmock.NewRows([]string{"id", "name", "age", "additional_information"}).AddRow(user.ID, user.Name, user.Age, user.AdditionalInformation)

	sqlString := regexp.QuoteMeta(`SELECT id, name, age, additional_information FROM users WHERE id=?`)
	mock.ExpectQuery(sqlString).WithArgs(user.ID).WillReturnRows(row)

	parentSSQLString := regexp.QuoteMeta(`SELECT name FROM user_parents WHERE user_id=?`)

	rows := sqlmock.NewRows([]string{"name"}).AddRow(user.Parents[0]).AddRow(user.Parents[1])

	mock.ExpectQuery(parentSSQLString).WithArgs(user.ID).WillReturnRows(rows)

	foundUser, err := userRepo.GetUser(context.Background(), user.ID)
	c.Equal(user, foundUser)
	c.NoError(err)
}

func TestGetUserFails(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	userRepo := NewUserRepository(db, log.NewJSONLogger(os.Stdout))

	user := &sharedLib.User{
		ID:                    "USR123",
		Name:                  "test",
		Age:                   99,
		AdditionalInformation: "not much",
		Parents:               []string{"John Doe", "Jane Doe"},
	}

	sqlString := regexp.QuoteMeta(`SELECT id, name, age, additional_information FROM users WHERE id=?`)
	mock.ExpectQuery(sqlString).WithArgs("USR123").WillReturnError(sql.ErrNoRows)

	_, err := userRepo.GetUser(context.Background(), "USR123")
	c.Equal(ErrUserNotFound, err)

	mock.ExpectQuery(sqlString).WithArgs("USR123").WillReturnError(config.ErrMockFails)

	_, err = userRepo.GetUser(context.Background(), "USR123")
	c.Equal(config.ErrMockFails, err)

	row := sqlmock.NewRows([]string{"id", "name", "age", "additional_information"}).AddRow(user.ID, user.Name, user.Age, user.AdditionalInformation)

	mock.ExpectQuery(sqlString).WithArgs(user.ID).WillReturnRows(row)

	parentSSQLString := regexp.QuoteMeta(`SELECT name FROM user_parents WHERE user_id=?`)

	mock.ExpectQuery(parentSSQLString).WithArgs("USR123").WillReturnError(config.ErrMockFails)

	_, err = userRepo.GetUser(context.Background(), "USR123")
	c.Equal(config.ErrMockFails, err)
}

func TestUpdateUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	userRepo := NewUserRepository(db, log.NewJSONLogger(os.Stdout))

	user := &sharedLib.User{
		ID:                    "USR123",
		Name:                  "test",
		Age:                   99,
		AdditionalInformation: "not much",
		Parents:               []string{"John Doe", "Jane Doe"},
	}

	sqlUpdateString := regexp.QuoteMeta(`UPDATE users SET name=?, age=?, additional_information=?  WHERE id = ?`)

	mock.ExpectExec(sqlUpdateString).WithArgs(user.Name, user.Age, user.AdditionalInformation, user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	sqlDeleteString := regexp.QuoteMeta(`DELETE FROM user_parents WHERE user_id=?`)

	mock.ExpectExec(sqlDeleteString).WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	parentsSQLInsertString := regexp.QuoteMeta(`INSERT INTO user_parents (user_id, name) VALUES(?, ?)`)

	mock.ExpectExec(parentsSQLInsertString).WithArgs(user.ID, user.Parents[0]).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(parentsSQLInsertString).WithArgs(user.ID, user.Parents[1]).WillReturnResult(sqlmock.NewResult(0, 1))

	row := sqlmock.NewRows([]string{"id", "name", "age", "additional_information"}).AddRow(user.ID, user.Name, user.Age, user.AdditionalInformation)

	sqlSelectString := regexp.QuoteMeta(`SELECT id, name, age, additional_information FROM users WHERE id=?`)

	mock.ExpectQuery(sqlSelectString).WithArgs(user.ID).WillReturnRows(row)

	parentSSQLString := regexp.QuoteMeta(`SELECT name FROM user_parents WHERE user_id=?`)

	rows := sqlmock.NewRows([]string{"name"}).AddRow(user.Parents[0]).AddRow(user.Parents[1])

	mock.ExpectQuery(parentSSQLString).WithArgs(user.ID).WillReturnRows(rows)

	foundUser, err := userRepo.UpdateUser(context.Background(), user)
	c.Equal(user, foundUser)
	c.NoError(err)
}

func TestUpdateUserFails(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	userRepo := NewUserRepository(db, log.NewJSONLogger(os.Stdout))

	user := &sharedLib.User{
		ID:                    "USR123",
		Name:                  "test",
		Age:                   99,
		AdditionalInformation: "not much",
		Parents:               []string{"John Doe", "Jane Doe"},
	}

	sqlUpdateString := regexp.QuoteMeta(`UPDATE users SET name=?, age=?, additional_information=?  WHERE id = ?`)

	mock.ExpectExec(sqlUpdateString).WithArgs(user.Name, user.Age, user.AdditionalInformation, user.ID).WillReturnError(config.ErrMockFails)

	_, err := userRepo.UpdateUser(context.Background(), user)
	c.Equal(config.ErrMockFails, err)

	mock.ExpectExec(sqlUpdateString).WithArgs(user.Name, user.Age, user.AdditionalInformation, user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	sqlDeleteString := regexp.QuoteMeta(`DELETE FROM user_parents WHERE user_id=?`)

	mock.ExpectExec(sqlDeleteString).WithArgs(user.ID).WillReturnError(config.ErrMockFails)

	_, err = userRepo.UpdateUser(context.Background(), user)
	c.Equal(config.ErrMockFails, err)

	mock.ExpectExec(sqlUpdateString).WithArgs(user.Name, user.Age, user.AdditionalInformation, user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(sqlDeleteString).WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	parentsSQLInsertString := regexp.QuoteMeta(`INSERT INTO user_parents (user_id, name) VALUES(?, ?)`)

	mock.ExpectExec(parentsSQLInsertString).WithArgs(user.ID, user.Parents[0]).WillReturnError(config.ErrMockFails)
	_, err = userRepo.UpdateUser(context.Background(), user)
	c.Equal(config.ErrMockFails, err)
}

func TestDeleteUser(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	userRepo := NewUserRepository(db, log.NewJSONLogger(os.Stdout))

	sqlString := regexp.QuoteMeta(`DELETE FROM user_parents WHERE user_id=?`)
	mock.ExpectExec(sqlString).WithArgs("USR123").WillReturnResult(sqlmock.NewResult(0, 1))

	parentSSQLString := regexp.QuoteMeta(`DELETE FROM users WHERE id=?`)
	mock.ExpectExec(parentSSQLString).WithArgs("USR123").WillReturnResult(sqlmock.NewResult(0, 1))

	err := userRepo.DeleteUser(context.Background(), "USR123")
	c.NoError(err)
}

func TestDeleteUserFails(t *testing.T) {
	c := require.New(t)

	db, mock := config.NewDatabaseMock()

	userRepo := NewUserRepository(db, log.NewJSONLogger(os.Stdout))

	sqlString := regexp.QuoteMeta(`DELETE FROM user_parents WHERE user_id=?`)
	mock.ExpectExec(sqlString).WithArgs("USR123").WillReturnError(config.ErrMockFails)

	err := userRepo.DeleteUser(context.Background(), "USR123")
	c.Equal(config.ErrMockFails, err)

	mock.ExpectExec(sqlString).WithArgs("USR123").WillReturnResult(sqlmock.NewResult(0, 1))

	parentSSQLString := regexp.QuoteMeta(`DELETE FROM users WHERE id=?`)
	mock.ExpectExec(parentSSQLString).WithArgs("USR123").WillReturnError(config.ErrMockFails)

	err = userRepo.DeleteUser(context.Background(), "USR123")
	c.Equal(config.ErrMockFails, err)
}
