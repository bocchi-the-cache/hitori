package cache

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"log"
	"testing"
)

// NOTE: BadgerDB is definitely **NOT** a cache storage engine.
// But we use it as k-v storage in our kick-start project.
// We will replace it with a real cache storage engine once `hitori` is ready.
func TestBadgerDB(t *testing.T) {
	opts := badger.DefaultOptions("/tmp/badger")
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Write
	t.Logf("Write to badger: %v", 42)
	err = db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte("answer"), []byte("42"))
		err := txn.SetEntry(e)
		return err
	})

	// Read
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("answer"))
		if err != nil {
			return err
		}

		var valNot, valCopy []byte
		err = item.Value(func(val []byte) error {
			// This func with val would only be called if item.Value encounters no error.

			// Accessing val here is valid.
			fmt.Printf("The answer is: %s\n", val)

			// Copying or parsing val is valid.
			valCopy = append([]byte{}, val...)

			// Assigning val slice to another variable is NOT OK.
			valNot = val // Do not do this.
			return nil
		})
		if err != nil {
			return err
		}

		// DO NOT access val here. It is the most common cause of bugs.
		fmt.Printf("NEVER do this. %s\n", valNot)

		// You must copy it to use it outside item.Value(...).
		t.Logf("The answer is: %s\n", valCopy)

		// Alternatively, you could also use item.ValueCopy().
		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		fmt.Printf("(using item.ValueCopy) The answer is: %s\n", valCopy)

		return nil
	})

}

func TestCachePublicFunction(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Write to cache: %v", "foo:bar")
	err = Set("foo", []byte("bar"))

	t.Logf("Read from cache: %v", "foo:bar")
	val, err := Get("foo")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Read from cache, value: %v", string(val))
	if string(val) != "bar" {
		t.Fatal("value is not equal to 'bar'")
	}

	t.Logf("Del from cache: %v", "foo:bar")
	err = Del("foo")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Del a non-exist key from cache: %v", "foo:bar")
	err = Del("foo")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Read a non-exist key from cache: %v", "foo:bar")
	val, err = Get("foo")
	// if not found
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Fatal("value should be nil")
	}
}
