package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	c := require.New(t)

	db, err := Connect()
	c.NotNil(db)
	c.Nil(err)
}

func TestConnectFails(t *testing.T) {
	c := require.New(t)

	dbDriver = ""
	db, err := Connect()
	c.Nil(db)
	c.EqualError(err, "sql: unknown driver \"\" (forgotten import?)")
}
