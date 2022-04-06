package database

import (
	bolt "go.etcd.io/bbolt"
)

type DB struct {
	Bolt *bolt.DB
}

func Start() (db *DB, err error) {
	db = &DB{}
	db.Bolt, err = bolt.Open("wikiofthings.db", 0600, nil)
	return
}

func (db *DB) Get(bucket string, key string) (value []byte) {
	db.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		value = b.Get([]byte(key))
		return nil
	})
	return
}

func (db *DB) Stop() error {
	return db.Bolt.Close()
}
