package shared

import (
	"golang.org/x/crypto/bcrypt"

	"crypto/rand"
	"encoding/hex"
)

const (
	_attemptsToReadRandomData = 2
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomData returns secure random data with the given size
func GenerateRandomData(bitSize int) []byte {
	buffer := make([]byte, bitSize)

	for i := 0; i < _attemptsToReadRandomData; i++ {
		_, err := rand.Read(buffer)
		if err == nil {
			break
		}
	}

	return buffer
}

// GenerateRandomHexString returns secure random hex string
func GenerateRandomHexString(bitSize int) string {
	return hex.EncodeToString(GenerateRandomData(bitSize))
}

// GenerateID generates a ID with the given prefix
func GenerateID(prefix string) string {
	return prefix + GenerateRandomHexString(16)
}
