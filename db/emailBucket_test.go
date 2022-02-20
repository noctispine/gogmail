package db

import (
	"testing"

	util "github.com/noctispine/go-email-app/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUserEmail() UserEmail {
	return UserEmail{
		email:    util.RandomMail(),
		password: util.RandomPassword(),
	}
}

func TestInitEmailBucket(t *testing.T) {
	err := testDB.InitEmailBucket()
	require.NoError(t, err)
}

func TestAddUserEmail(t *testing.T) {
	userEmail := createRandomUserEmail()
	err := testDB.AddUserEmail(userEmail)
	require.NoError(t, err)

	// try to get pw after add
	pw, le := testDB.GetPassword(userEmail)
	require.NotZero(t, le)
	require.Equal(t, len(userEmail.password), le)
	require.NotEmpty(t, pw)
	require.Equal(t, userEmail.password, string(pw))
}

func TestRemoveUserEmail(t *testing.T) {
	userEmail := createRandomUserEmail()
	err := testDB.AddUserEmail(userEmail)
	require.NoError(t, err)

	err = testDB.RemoveUserEmail(userEmail)
	require.NoError(t, err)

	// try to get pw after remove
	pw, le := testDB.GetPassword(userEmail)
	require.Zero(t, le)
	require.Equal(t, 0, le)
	require.Empty(t, pw)

}

func TestGetPassword(t *testing.T) {
	userEmail := createRandomUserEmail()
	err := testDB.AddUserEmail(userEmail)
	require.NoError(t, err)

	pw, le := testDB.GetPassword(userEmail)
	require.NotZero(t, le)
	require.Equal(t, len(userEmail.password), le)
	require.NotEmpty(t, pw)
	require.Equal(t, userEmail.password, string(pw))
}

func TestChangeMailPassword(t *testing.T) {
	userEmail := createRandomUserEmail()
	newPassword := util.RandomPassword()

	err := testDB.AddUserEmail(userEmail)
	require.NoError(t, err)

	err = testDB.ChangeMailPassword(userEmail, newPassword)
	require.NoError(t, err)

	pw, le := testDB.GetPassword(userEmail)
	require.NotZero(t, le)
	require.NotEmpty(t, pw)
	require.Equal(t, len(newPassword), le)
	require.Equal(t, newPassword, pw)

}
