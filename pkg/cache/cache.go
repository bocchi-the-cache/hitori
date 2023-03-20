package cache

import (
	"github.com/bocchi-the-cache/hitori/pkg/config"
	"github.com/dgraph-io/badger/v4"
)

// NOTE: BadgerDB is definitely **NOT** a cache storage engine.
// But we use it as k-v storage in our kick-start project.
// We will replace it with a real cache storage engine once `hitori` is ready.

// Cache could be defined by user
type Cache interface {
	Del(key string) error
	// Get could use buffer/streaming
	Get(key string) ([]byte, error)
	// Set could use buffer/streaming
	Set(key string, value []byte) error
}

var DefaultCache Cache

type StorageCache struct {
	db *badger.DB
}

func NewStorageCache(db *badger.DB) *StorageCache {
	return &StorageCache{db: db}
}

func Init(config *config.Config) error {
	opts := badger.DefaultOptions("/tmp/badger")
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}

	DefaultCache = NewStorageCache(db)
	return nil
}

func (c *StorageCache) Del(key string) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (c *StorageCache) Get(key string) ([]byte, error) {
	var val []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(v []byte) error {
			val = append([]byte{}, v...)
			return nil
		})
	})
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return val, err
}

func (c *StorageCache) Set(key string, value []byte) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

func Del(key string) error {
	return DefaultCache.Del(key)
}

func Get(key string) ([]byte, error) {
	return DefaultCache.Get(key)
}

func Set(key string, value []byte) error {
	return DefaultCache.Set(key, value)
}
