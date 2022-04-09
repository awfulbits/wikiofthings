package database

import (
	"encoding/binary"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type DB struct {
	Bolt *bolt.DB
}

func Start(buckets ...string) (db *DB, err error) {
	db = &DB{}
	db.Bolt, err = bolt.Open("wikiofthings.db", 0600, nil)
	for i := range buckets {
		db.Bolt.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(buckets[i]))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})
	}
	return
}

func (db *DB) Stop() error {
	return db.Bolt.Close()
}

func (db *DB) CreateID(bucket string) (id int) {
	db.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		uid, _ := b.NextSequence()
		id = int(uid)
		return nil
	})
	return id
}

func (db *DB) Create(bucket string, key []byte, value []byte) (err error) {
	return db.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if value := b.Get(key); value != nil {
			err = fmt.Errorf("error creating page: storage key %s already exists", key)
			return err
		}
		err = b.Put(key, value)
		return err
	})
}

func (db *DB) Read(bucket string, key []byte) (value []byte, err error) {
	db.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		value = b.Get(key)
		return nil
	})
	if len(value) <= 0 {
		err = fmt.Errorf("key-value not found in bucket")
	}
	return
}

func (db *DB) Update(bucket string, key []byte, value []byte) {
	db.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(key, value)
		return err
	})
}

func (db *DB) Delete(bucket string, key []byte) {}

func (db *DB) IdToBytes(id int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	return b
}

func (db *DB) BytesToId(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
