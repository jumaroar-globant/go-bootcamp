package config

import (
	"database/sql"
	"errors"
	"log"

	"github.com/DATA-DOG/go-sqlmock"
)

var ErrMockFails = errors.New("forced failure")

// NewDatabaseMock is a function to initialize a database mock
func NewDatabaseMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
