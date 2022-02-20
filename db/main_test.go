package db

import (
	"log"
	"testing"

	util "github.com/noctispine/go-email-app/utils"
	"github.com/stretchr/testify/require"
)

var testDB *Database

const test_bucket = "test_bucket"

func generateRandomKeyValuePair() (string, string) {
	randomKey := util.RandomString(util.RandomInt(4, 10))
	randomVal := util.RandomString(util.RandomInt(5, 20))
	return randomKey, randomVal
}

func TestMain(m *testing.M) {
	var err error
	testDB, err = NewDB("test")
	if err != nil {
		log.Fatal("db connection is failed:", err)
	}
	m.Run()
}

func TestCreateBucket(t *testing.T) {
	err := testDB.createBucketDB(test_bucket)
	require.NoError(t, err)
}

func TestUpdateDB(t *testing.T) {
	randomKey, randomVal := generateRandomKeyValuePair()
	err := testDB.updateDB([]byte(test_bucket), []byte(randomKey), []byte(randomVal))

	require.NoError(t, err)

	// get and test it
	val, length := testDB.queryDB([]byte(test_bucket), []byte(randomKey))

	require.Equal(t, randomVal, string(val))
	require.NotZero(t, len(randomVal), length)
	require.Equal(t, len(randomVal), length)
}

func TestDeleteKey(t *testing.T) {
	randomKey, randomVal := generateRandomKeyValuePair()
	err := testDB.updateDB([]byte(test_bucket), []byte(randomKey), []byte(randomVal))
	require.NoError(t, err)

	// delete
	err = testDB.deleteKey([]byte(test_bucket), []byte(randomKey))
	require.NoError(t, err)

	// try to get deleted key-val pair
	val, length := testDB.queryDB([]byte(test_bucket), []byte(randomKey))
	require.Empty(t, val)
	require.Zero(t, length)
}

func TestQueryDB(t *testing.T) {
	randomKey, randomVal := generateRandomKeyValuePair()
	err := testDB.updateDB([]byte(test_bucket), []byte(randomKey), []byte(randomVal))
	require.NoError(t, err)

	val, length := testDB.queryDB([]byte(test_bucket), []byte(randomKey))
	require.NotEmpty(t, val)
	require.NotZero(t, length)
	require.NotEmpty(t, val)
	require.Equal(t, randomVal, string(val))

}
