package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	hash, _ := HashPassword(password)

	// Correct password
	assert.True(t, CheckPassword(password, hash))

	// Wrong password
	assert.False(t, CheckPassword("wrongpassword", hash))
}

func TestCheckPasswordEmpty(t *testing.T) {
	hash, _ := HashPassword("test123")

	// Empty password should fail
	assert.False(t, CheckPassword("", hash))
}
