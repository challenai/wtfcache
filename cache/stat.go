package cache

import (
	"encoding/json"
	"errors"
	"sync/atomic"
)

type Stat struct {
	// entries count
	keys int64
	// size of all the entries
	sz int64
}

func NewStat() *Stat {
	return &Stat{
		keys: 0,
		sz:   0,
	}
}

func (ms *Stat) CountKeys() int64 {
	return atomic.LoadInt64(&ms.keys)
}

func (ms *Stat) GetSz() int64 {
	return atomic.LoadInt64(&ms.sz)
}

func (ms *Stat) IncrKey(num int64) int64 {
	return atomic.AddInt64(&ms.keys, num)
}

func (ms *Stat) IncrSz(sz int64) int64 {
	return atomic.AddInt64(&ms.sz, sz)
}

func (ms *Stat) DecrKey(num int64) int64 {
	return atomic.AddInt64(&ms.keys, -num)
}

func (ms *Stat) DecrSz(sz int64) int64 {
	return atomic.AddInt64(&ms.sz, -sz)
}

func (s *Stat) Dump() []byte {
	result, _ := json.Marshal(map[string]int64{
		"keys": s.CountKeys(),
		"sz":   s.GetSz(),
	})
	return result
}

func (s *Stat) Load(data []byte) error {
	stat := make(map[string]int64)
	err := json.Unmarshal(data, &stat)
	if err != nil {
		return err
	}

	var keys, sz int64
	var exist bool
	keys, exist = stat["keys"]
	if !exist {
		return errors.New("fail to load stat from badger")
	}
	sz, exist = stat["sz"]
	if !exist {
		return errors.New("fail to load stat from badger")
	}

	s.keys = keys
	s.sz = sz

	return nil
}
