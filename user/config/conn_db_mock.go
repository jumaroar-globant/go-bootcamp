package config

import (
	"database/sql"
	"log"

	"github.com/DATA-DOG/go-sqlmock"
)

// NewDatabaseMock is a function to initialize a database mock
func NewDatabaseMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
