package db

import (
	"testing"

	util "github.com/noctispine/go-email-app/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUserEmail() UserEmail {
	return UserEmail{
		Email:    util.RandomMail(),
		Password: util.RandomPassword(),
	}
}

func TestInitEmailBucket(t *testing.T) {
	err := InitEmailBucket()
	require.NoError(t, err)
}

func TestAddUserEmail(t *testing.T) {
	userEmail := createRandomUserEmail()
	err := AddUserEmail(userEmail)
	require.NoError(t, err)

	// try to get pw after add
	pw, le := GetPassword(userEmail)
	require.NotZero(t, le)
	require.Equal(t, len(userEmail.Password), le)
	require.NotEmpty(t, pw)
	require.Equal(t, userEmail.Password, string(pw))
}

func TestRemoveUserEmail(t *testing.T) {
	userEmail := createRandomUserEmail()
	err := AddUserEmail(userEmail)
	require.NoError(t, err)

	err = RemoveUserEmail(userEmail.Email)
	require.NoError(t, err)

	// try to get pw after remove
	pw, le := GetPassword(userEmail)
	require.Zero(t, le)
	require.Equal(t, 0, le)
	require.Empty(t, pw)

}

func TestGetPassword(t *testing.T) {
	userEmail := createRandomUserEmail()
	err := AddUserEmail(userEmail)
	require.NoError(t, err)

	pw, le := GetPassword(userEmail)
	require.NotZero(t, le)
	require.Equal(t, len(userEmail.Password), le)
	require.NotEmpty(t, pw)
	require.Equal(t, userEmail.Password, string(pw))
}

func TestChangeMailPassword(t *testing.T) {
	userEmail := createRandomUserEmail()
	newPassword := util.RandomPassword()

	err := AddUserEmail(userEmail)
	require.NoError(t, err)

	err = ChangeMailPassword(userEmail, newPassword)
	require.NoError(t, err)

	pw, le := GetPassword(userEmail)
	require.NotZero(t, le)
	require.NotEmpty(t, pw)
	require.Equal(t, len(newPassword), le)
	require.Equal(t, newPassword, pw)

}

func TestIterateEmailBucket(t *testing.T) {
	err := IterateEmailBucket()
	require.NoError(t, err)
}

func TestMakeSliceFromEmailBucket(t *testing.T) {
	emails, err := MakeSliceFromEmailBucket()
	require.NoError(t, err)

	require.NotEmpty(t, emails)
}
