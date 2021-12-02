package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jumaroar-globant/go-bootcamp/user/shared"
)

var (
	dbDriver   = shared.GetStringEnvVar("DATABASE_DRIVER", "mysql")
	dbUsername = shared.GetStringEnvVar("DATABASE_USERNAME", "root")
	dbPassword = shared.GetStringEnvVar("DATABASE_PASSWORD", "")
	dbIP       = shared.GetStringEnvVar("DATABASE_IP", "127.0.0.1")
	dbPort     = shared.GetStringEnvVar("DATABASE_PORT", "3306")
	dbName     = shared.GetStringEnvVar("DATABASE_NAME", "bootcamp")
)

// Connect is a function to connect to the database
func Connect() (*sql.DB, error) {
	dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbIP, dbPort, dbName)

	db, err := sql.Open(dbDriver, dbConnString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
