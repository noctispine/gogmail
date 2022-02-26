package db

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/fatih/color"
	"github.com/noctispine/gogmail/gservice"
)

const DEFAULT_EMAIL_BUCKET_NAME = "my_emails"

type User struct {
	EmailAddress string
	Infos        gservice.OAuthInfos
}

// initialize email bucket
// which contains user's emails
func InitEmailBucket() error {
	err := createBucketDB(DEFAULT_EMAIL_BUCKET_NAME)
	if err != nil {
		return err
	}

	return nil
}

// add mail-password pair to the bucket
func AddUser(user User) error {
	encoded, err := json.Marshal(user.Infos)

	if err != nil {
		return err
	}

	err = updateDB([]byte(DEFAULT_EMAIL_BUCKET_NAME), []byte(user.EmailAddress), []byte(encoded))
	if err != nil {
		return err
	}

	return nil
}

// remove mail-passowrd pair from the bucket
func RemoveUserEmail(key string) error {
	err := deleteKey([]byte(DEFAULT_EMAIL_BUCKET_NAME), []byte(key))
	if err != nil {
		return err

	}

	return nil
}

// change mail's password with given new password
// actually it removes the pair assassociated with the given email address
// then add a new pair with new Password
func ChangeEmailInfos(emailAddress string, newInfos gservice.OAuthInfos) error {

	err := RemoveUserEmail(emailAddress)
	if err != nil {
		return err
	}

	newUserEmail := User{
		EmailAddress: emailAddress,
		Infos:        newInfos,
	}

	err = AddUser(newUserEmail)
	if err != nil {
		return err
	}

	return nil
}

//
func GetInfos(user User) (string, int) {
	val, len := queryDB([]byte(DEFAULT_EMAIL_BUCKET_NAME), []byte(user.EmailAddress))
	return string(val), len
}

// for testing purposes
func IterateEmailBucket() error {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DEFAULT_EMAIL_BUCKET_NAME))

		c := b.Cursor()

		i := 0
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if i%2 == 1 {
				color.Set(color.ReverseVideo)
			}
			fmt.Printf("%d: %s : %s\n", i, k, v)

			color.Unset()
			i++
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func MakeSliceFromUser() ([]User, error) {
	// it should be at least 1 len because we add quit option
	// no matter there is emails or not
	emails := make([]User, 1, 10)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DEFAULT_EMAIL_BUCKET_NAME))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var userInfos gservice.OAuthInfos
			err := json.Unmarshal(v, &userInfos)

			if err != nil {
				return err
			}

			user := User{
				EmailAddress: string(k),
				Infos:        userInfos,
			}

			emails = append(emails, user)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return emails, nil
}
