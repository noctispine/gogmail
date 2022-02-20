package db

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type Database struct {
	boltDB *bolt.DB
}

func NewDB(dbName string) (*Database, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Database{
		boltDB: db,
	}, nil
}

func (db *Database) createBucketDB(bucketName string) error {
	err := db.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) updateDB(bucketName, key, value []byte) error {
	err := db.boltDB.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		err = bkt.Put(key, value)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (db *Database) queryDB(bucketName, key []byte) (val []byte, length int) {
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bucketName)
		if bkt == nil {
			return fmt.Errorf("Bucket %q not found!", bucketName)
		}
		val = bkt.Get(key)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return val, len(string(val))

}

func (db *Database) iterateDB(bucketName []byte) {
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=>[%s], value=[%s]\n", k, v)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("failure : %s\n", err)
	}
}

func (db *Database) deleteKey(bucketName, keyName []byte) error {
	err := db.boltDB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucketName)
		err := b.Delete(keyName)

		return err
	})

	if err != nil {
		return err
		// log.Fatalf("failure : %s\n", err)
	}

	return nil
}

func (db *Database) CloseDB() {
	db.boltDB.Close()
}
