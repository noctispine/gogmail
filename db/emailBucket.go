package db

const DEFAULT_EMAIL_BUCKET_NAME = "my_emails"

type UserEmail struct {
	email    string
	password string
}

// initialize email bucket
// which contains user's emails
func (db *Database) InitEmailBucket() error {
	err := db.createBucketDB(DEFAULT_EMAIL_BUCKET_NAME)
	if err != nil {
		return err
	}

	return nil
}

// add mail-password pair to the bucket
func (db *Database) AddUserEmail(userEmail UserEmail) error {
	err := db.updateDB([]byte(DEFAULT_EMAIL_BUCKET_NAME), []byte(userEmail.email), []byte(userEmail.password))
	if err != nil {
		return err
	}

	return nil
}

// remove mail-passowrd pair from the bucket
func (db *Database) RemoveUserEmail(userEmail UserEmail) error {
	err := db.deleteKey([]byte(DEFAULT_EMAIL_BUCKET_NAME), []byte(userEmail.email))
	if err != nil {
		return err

	}

	return nil
}

// change mail's password with given new password
// actually it removes the pair assassociated with the given email address
// then add a new pair with new password
func (db *Database) ChangeMailPassword(userEmail UserEmail, newPassword string) error {

	err := db.RemoveUserEmail(userEmail)
	if err != nil {
		return err
	}

	newUserEmail := UserEmail{
		email:    userEmail.email,
		password: newPassword,
	}

	err = db.AddUserEmail(newUserEmail)
	if err != nil {
		return err
	}

	return nil
}

//
func (db *Database) GetPassword(userEmail UserEmail) (string, int) {
	val, len := db.queryDB([]byte(DEFAULT_EMAIL_BUCKET_NAME), []byte(userEmail.email))
	return string(val), len
}

// for testing purposes
func (db *Database) IterateEmailBucket() {
	db.iterateDB([]byte(DEFAULT_EMAIL_BUCKET_NAME))
}
