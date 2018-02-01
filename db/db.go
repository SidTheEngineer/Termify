package db

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

const (
	dbName = "SpotfiyAuth.db"
)

var db *bolt.DB

// StartDB opens and starts a new BoltDB to be connected to.
func StartDB() {
	newDb, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 5 * time.Second})

	if err != nil {
		log.Fatal(err)
	}

	db = newDb
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("auth"))
		return nil
	})
}
