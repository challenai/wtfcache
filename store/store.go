package store

import (
	"cacheme/stat"
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

type Store struct {
	StoreStat
	db *badger.DB
}

func NewStore(storePath string) (*Store, error) {
	var err error
	s := &Store{}
	s.db, err = badger.Open(badger.DefaultOptions(storePath))
	if err != nil {
		log.Fatal(err)
	}
	return s, nil
}

func (s *Store) Get(k string) ([]byte, error) {
	var val []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		return err
	})
	return val, err
}

func (s *Store) Set(k string, v []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(k), v)
		return err
	})
	return err
}

func (s *Store) Del(k string) (bool, error) {
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(k))
	})
	return err == nil, err
}

func (s *Store) Info() stat.Stat {
	return &s.StoreStat
}
