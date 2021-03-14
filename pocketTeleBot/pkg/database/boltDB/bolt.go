package boltDB

import (
	"errors"
	"github.com/boltdb/bolt"
	"pocketTeleBot/pkg/database"
	"strconv"
)

type TokenDB struct {
	db *bolt.DB
}

func NewTokenDB(db *bolt.DB) *TokenDB {
	return &TokenDB{db: db}
}

func (r *TokenDB) Save(chatID int64, token string, bucket database.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatID), []byte(token))
	})
}

func (r *TokenDB) Get(chatID int64, bucket database.Bucket) (string, error) {
	var token string
	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(intToBytes(chatID))
		token = string(data)
		return nil
	})
	if err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token not found")
	}
	return token, nil
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
