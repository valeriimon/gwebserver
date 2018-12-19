package mainservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

// Database : wrapper for bolt database to add custom methods
type Database struct {
	ref *bolt.DB
}

// InitDatabase : initialize database
func InitDatabase() error {
	var database *Database
	database, err := database.openSession()
	if err != nil {
		return err
	}
	defer database.ref.Close()

	if err := database.ref.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DbMainBucket))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Open database file stream
func (db *Database) openSession() (*Database, error) {
	boltDb, err := bolt.Open("database.db", 0600, nil)
	db = &Database{boltDb}
	return db, err
}

// Save data to database
func (db *Database) Save(bucket, key string, newData []byte) {
	db, err := db.openSession()
	if err != nil {
		panic(err)
	}
	defer db.ref.Close()

	if err := db.ref.Update(func(tx *bolt.Tx) error {
		mainBucket := tx.Bucket([]byte(DbMainBucket))

		b, err := mainBucket.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		if err = b.Put([]byte(key), newData); err != nil {
			return err
		}
		fmt.Println("Saved succesfully")
		return nil
	}); err != nil {
		panic(err.Error())
	}
}

// GetData : get data from database;
// If key length equels to zero returns whole bucket records
func (db *Database) GetData(bucket, key string) (results []interface{}, err error) {
	db, err = db.openSession()
	if err != nil {
		return nil, err
	}
	defer db.ref.Close()

	if err := db.ref.View(func(tx *bolt.Tx) error {
		var result interface{}

		b := tx.Bucket([]byte(DbMainBucket)).Bucket([]byte(bucket))
		if b == nil {
			return errors.New("Bucket Not Found")
		}

		if len(key) == 0 {
			// TODO: Change forEach to for{} method
			b.ForEach(func(k, v []byte) error {
				if err := json.Unmarshal(v, &result); err != nil {
					return err
				}

				results = append(results, result)
				return nil
			})

			return nil
		}

		buf := b.Get([]byte(key))
		if buf == nil {
			return errors.New("Not Found")
		}

		if err := json.Unmarshal(buf, &result); err != nil {
			return err
		}

		results = append(results, result)

		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

// Get is simple method to get sample from database via simple transaction
func (db *Database) Get(bucket, key string) (interface{}, error) {
	var result interface{}
	db, err := db.openSession()
	if err != nil {
		return nil, err
	}
	defer db.ref.Close()

	tx, err := db.ref.Begin(false)
	if err != nil {
		return nil, err
	}

	b := tx.Bucket([]byte(DbMainBucket)).Bucket([]byte(bucket))
	if b == nil {
		return nil, errors.New("Bucket Not Found")
	}

	buf := b.Get([]byte(key))
	if buf == nil {
		return nil, nil
	}

	if err := json.Unmarshal(buf, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Auth : simple users authentication via email and password
func (db *Database) Auth(bucket, email, password string) (*User, error) {
	var result []byte
	var user User

	db, err := db.openSession()
	if err != nil {
		panic(err)
	}
	defer db.ref.Close()

	if err := db.ref.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(DbMainBucket)).Bucket([]byte(bucket))
		c := b.Cursor()

		email := []byte(fmt.Sprintf(`"email":"%v"`, email))
		pass := []byte(fmt.Sprintf(`"password":"%v"`, password))

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Contains(v, email) && bytes.Contains(v, pass) {
				result = v
				return nil
			}
		}

		return errors.New("Auth failed")

	}); err != nil {
		return nil, err
	}

	err = json.Unmarshal(result, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
