package storage

import (
	"github.com/dgraph-io/badger/v4"
	"log"
	"testing"
)

// NOTE: BadgerDB is definitely **NOT** a cache storage engine.
// But we use it as k-v storage in our kick-start project.
func TestBadgerDB(t *testing.T) {
	opts := badger.DefaultOptions("/tmp/badger")
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
