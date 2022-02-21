package db

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/fatih/color"
)

const DEFAULT_EMAIL_BUCKET_NAME = "my_emails"

type UserEmail struct {
	Email    string
	Password string
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
func AddUserEmail(userEmail UserEmail) error {
	err := updateDB([]byte(DEFAULT_EMAIL_BUCKET_NAME), []byte(userEmail.Email), []byte(userEmail.Password))
	if err != nil {
		return err
	}

	return nil
}

// remove mail-passowrd pair from the bucket
func RemoveUserEmail(userEmail UserEmail) error {
	err := deleteKey([]byte(DEFAULT_EMAIL_BUCKET_NAME), []byte(userEmail.Email))
	if err != nil {
		return err

	}

	return nil
}

// change mail's password with given new password
// actually it removes the pair assassociated with the given email address
// then add a new pair with new Password
func ChangeMailPassword(userEmail UserEmail, newPassword string) error {

	err := RemoveUserEmail(userEmail)
	if err != nil {
		return err
	}

	newUserEmail := UserEmail{
		Email:    userEmail.Email,
		Password: newPassword,
	}

	err = AddUserEmail(newUserEmail)
	if err != nil {
		return err
	}

	return nil
}

//
func GetPassword(userEmail UserEmail) (string, int) {
	val, len := queryDB([]byte(DEFAULT_EMAIL_BUCKET_NAME), []byte(userEmail.Email))
	return string(val), len
}

// for testing purposes
func IterateEmailBucket() error {
	fmt.Println("Added Emails:")
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
