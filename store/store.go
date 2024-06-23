package store

import (
	"cacheme/cache"
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

type Store struct {
	db *badger.DB
	*cache.Stat
}

func NewStore(storePath string) (*Store, error) {
	var err error
	s := &Store{}
	s.db, err = badger.Open(badger.DefaultOptions(storePath))
	if err != nil {
		log.Fatal(err)
	}
	s.Stat = cache.NewStat()

	var statsRecord []byte
	err = s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(STAT_PERSIST_RECORD))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		statsRecord, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return nil, err
	}
	if len(statsRecord) == 0 {
		return s, nil
	}
	err = s.Stat.Load(statsRecord)
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
	defer s.persistStat()

	var valueSize int64 = 0
	err := s.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		if err != badger.ErrKeyNotFound {
			return err
		}
		valueSize = item.ValueSize()

		err = txn.Set([]byte(k), v)
		return err
	})
	if err != nil {
		return err
	}

	if valueSize == 0 {
		s.Stat.IncrKey(1)
		s.Stat.IncrSz(int64(len(k) + len(v)))
		return nil
	}
	s.Stat.IncrSz(int64(len(v)) - valueSize)
	return nil
}

func (s *Store) Del(k string) (bool, error) {
	defer s.persistStat()

	var valueSize int64 = 0
	err := s.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		if err != badger.ErrKeyNotFound {
			return err
		}
		valueSize = item.ValueSize()

		return txn.Delete([]byte(k))
	})

	if valueSize != 0 {
		s.Stat.DecrKey(1)
		s.Stat.DecrSz(int64(len(k)) + valueSize)
	}
	return err == nil, err
}

func (s *Store) persistStat() error {
	data := s.Stat.Dump()
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(STAT_PERSIST_RECORD), data)
	})
	return err
}
