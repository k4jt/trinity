package store

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

const (
	BucketUsers = "users"
)

type Store struct {
	db *bolt.DB
}

func Open() *Store {
	// Open the trinity.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("trinity.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketUsers))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return &Store{db: db}
}

func (s *Store) Close() {
	s.db.Close()
	log.Info("close db connection")
}

// CreateUser saves u to the store. The new user ID is set on u once the data is persisted.
func (s *Store) CreateUser(u *User) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte(BucketUsers))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writable.
		// That can't happen in an Update() call so I ignore the error check.
		u.ID, _ = b.NextSequence()

		// Marshal user data into bytes.
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(u.ID), buf)
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func (s *Store) GetUser(id uint64) (u User, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketUsers))
		if bucket == nil {
			return fmt.Errorf("Bucket '%s' not found!", BucketUsers)
		}

		val := bucket.Get(itob(id))

		if err := json.Unmarshal(val, &u); err != nil {
			return fmt.Errorf("Error unmarshal user: %v", err)
		}

		return nil
	})

	return
}

func (s *Store) GetAllUsers() (users []User, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketUsers))
		if bucket == nil {
			return fmt.Errorf("Bucket '%s' not found!", BucketUsers)
		}

		err := bucket.ForEach(func(k, val []byte) error {
			var u User
			if err := json.Unmarshal(val, &u); err != nil {
				return fmt.Errorf("Error unmarshal user: %v", err)
			}
			users = append(users, u)
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error itterating bucket: %v", err)
		}

		return nil
	})

	return
}
