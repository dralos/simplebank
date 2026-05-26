package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.Error(t, err)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}

func TestPasswordHashErrorOnLongPassword(t *testing.T) {
	// create a password that's longer than 72 bytes (bcrypt limit) using an emoji
	// 72 ASCII 'a' bytes plus an emoji
	password := strings.Repeat("a", 72) + "🙂"

	_, err := HashPassword(password)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to hash password")
}
