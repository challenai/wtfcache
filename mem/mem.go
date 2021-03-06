package mem

import (
	"sync"
)

type MemCache struct {
	sync.Map
	MemStat
}

func NewMemCache() *MemCache {
	return &MemCache{sync.Map{}, MemStat{
		keys: 0,
		sz:   0,
	}}
}

func (mc *MemCache) Get(k string) ([]byte, error) {
	value, ok := mc.Load(k)
	if !ok {
		return []byte(""), nil
	}
	return value.([]byte), nil
}

func (mc *MemCache) Set(k string, v []byte) error {
	v_, exist := mc.LoadOrStore(k, v)
	if !exist {
		mc.MemStat.incr(1)
		mc.MemStat.incrSz(int64(len(k) + len(v)))
		return nil
	}
	mc.MemStat.incrSz(int64(len(v) - len(v_.([]byte))))
	return nil
}

func (mc *MemCache) Del(k string) (bool, error) {
	v, exist := mc.LoadAndDelete(k)
	if exist {
		mc.MemStat.decr(1)
		mc.MemStat.decrSz(int64(len(k) + len(v.([]byte))))
	}
	return exist, nil
}
