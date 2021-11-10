package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jumaroar-globant/go-bootcamp/user/shared"
)

var (
	Db *sql.DB

	dbDriver   = shared.GetStringEnvVar("DATABASE_DRIVER", "mysql")
	dbUsername = shared.GetStringEnvVar("DATABASE_USERNAME", "root")
	dbPassword = shared.GetStringEnvVar("DATABASE_PASSWORD", "")
	dbIP       = shared.GetStringEnvVar("DATABASE_IP", "127.0.0.1")
	dbPort     = shared.GetStringEnvVar("DATABASE_PORT", "3306")
	dbName     = shared.GetStringEnvVar("DATABASE_NAME", "bootcamp")
)

func init() {
	Db = Connect()
}

func Connect() (db *sql.DB) {
	dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbIP, dbPort, dbName)

	db, error := sql.Open(dbDriver, dbConnString)
	if error != nil {
		log.Fatal("There was an error connecting to the database.", error)
	}

	return db
}
