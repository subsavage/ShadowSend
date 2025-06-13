package storage

import (
	"errors"
	"go.etcd.io/bbolt"
	"time"
	"encoding/json"
)

const bucketName = "files"

type FileMetadata struct {
	Filename string
	Data []byte
	Nonce []byte
	Salt []byte
	ExpiresAt time.Time
}

type DB struct {
	conn *bbolt.DB
}

func encodeMetadata(meta FileMetadata) ([]byte, error) {
	return json.Marshal(meta)
}

func decodeMetadata(data []byte) (*FileMetadata, error) {
	var meta FileMetadata
	err := json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

func NewDB(path string) (*DB, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &DB{conn: db}, nil
}

func (db *DB) SaveFile(id string, meta FileMetadata) error {
	return db.conn.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}

		bytes, err := encodeMetadata(meta)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(id), bytes)
	})
}

func (db *DB) GetFile(id string) (*FileMetadata, error) {
	var meta *FileMetadata

	err := db.conn.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.New("bucket not found")
		}

		data := bucket.Get([]byte(id))
		if data == nil {
			return errors.New("file not found")
		}

		var err error
		meta, err = decodeMetadata(data)
		return err
	})

	if err != nil {
		return nil, err
	}

	if time.Now().After(meta.ExpiresAt) {
		return nil, errors.New("file has expired")
	}

	return meta, nil
}