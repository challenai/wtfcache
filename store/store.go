package store

import "cacheme/stat"

type Store struct {
	StoreStat
}

func (s *Store) Get(k string) ([]byte, error) {

}

func (s *Store) Set(k string, v []byte) error {

}

func (s *Store) Del(k string) error {

}

func (s *Store) Info() stat.Stat {

}
