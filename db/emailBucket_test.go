package db

import (
	"encoding/json"
	"testing"

	"github.com/noctispine/gogmail/gservice"
	util "github.com/noctispine/gogmail/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser() User {
	randomInfos := gservice.OAuthInfos{
		ClientID:     util.RandomString(20),
		ClientSecret: util.RandomString(30),
		RefreshToken: util.RandomString(30),
		AccessToken:  util.RandomString(30),
	}
	return User{
		EmailAddress: util.RandomMail(),
		Infos:        randomInfos,
	}
}

func TestInitEmailBucket(t *testing.T) {
	err := InitEmailBucket()
	require.NoError(t, err)
}

func TestAddUser(t *testing.T) {
	user := createRandomUser()
	err := AddUser(user)
	require.NoError(t, err)

	// try to get pw after add
	_, le := GetInfos(user)
	b, _ := json.Marshal(user.Infos)
	require.Equal(t, len(string(b)), le)
	require.NotZero(t, le)

}

func TestRemoveUserEmail(t *testing.T) {
	user := createRandomUser()
	err := AddUser(user)
	require.NoError(t, err)

	err = RemoveUserEmail(user.EmailAddress)
	require.NoError(t, err)

	// try to get pw after remove
	infos, le := GetInfos(user)
	require.Zero(t, le)
	require.Equal(t, 0, le)
	require.Empty(t, infos)

}

func TestGetInfos(t *testing.T) {
	user := createRandomUser()
	err := AddUser(user)
	require.NoError(t, err)
	b, _ := json.Marshal(user.Infos)

	infos, le := GetInfos(user)
	require.NotZero(t, le)
	require.Equal(t, len(string(b)), le)
	require.NotEmpty(t, b)
	require.Equal(t, string(b), infos)
}

func TestChangeEmailInfos(t *testing.T) {
	user := createRandomUser()
	newUser := createRandomUser()
	b, _ := json.Marshal(newUser.Infos)

	err := AddUser(user)
	require.NoError(t, err)

	err = ChangeEmailInfos(user.EmailAddress, newUser.Infos)
	require.NoError(t, err)

	infos, le := GetInfos(user)
	require.NotZero(t, le)
	require.NotEmpty(t, infos)
	require.Equal(t, len(string(b)), le)
	require.Equal(t, string(b), infos)

}

func TestIterateEmailBucket(t *testing.T) {
	err := IterateEmailBucket()
	require.NoError(t, err)
}

func TestMakeSliceFromUser(t *testing.T) {
	emails, err := MakeSliceFromUser()
	require.NoError(t, err)

	require.NotEmpty(t, emails)
}
