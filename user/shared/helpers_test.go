package shared

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	c := require.New(t)

	hash, err := HashPassword("")
	c.IsType("", hash)
	c.Nil(err)
}

func TestCheckPasswordHash(t *testing.T) {
	c := require.New(t)

	hash, err := HashPassword("test")
	c.Nil(err)

	isValid := CheckPasswordHash("test", hash)
	c.True(isValid)
}

func TestGenerateRandomData(t *testing.T) {
	c := require.New(t)

	c.Equal(8, len(GenerateRandomData(8)))
}

func TestGenerateRandomHexString(t *testing.T) {
	c := require.New(t)

	c.Equal(16, len(GenerateRandomHexString(8)))
}

func TestGenerateID(t *testing.T) {
	c := require.New(t)

	c.Equal(34, len(GenerateID("PR")))
}
