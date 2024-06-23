package mem

import (
	"cacheme/cache"
	"sync"
)

type MemCache struct {
	sync.Map
	*cache.Stat
}

func NewMemCache() *MemCache {
	return &MemCache{
		sync.Map{},
		cache.NewStat(),
	}
}

func (mc *MemCache) Get(k string) ([]byte, error) {
	value, ok := mc.Map.Load(k)
	if !ok {
		return []byte(""), nil
	}
	return value.([]byte), nil
}

func (mc *MemCache) Set(k string, v []byte) error {
	v_, exist := mc.LoadOrStore(k, v)
	if !exist {
		mc.Stat.IncrKey(1)
		mc.Stat.IncrSz(int64(len(k) + len(v)))
		return nil
	}
	mc.Stat.IncrSz(int64(len(v) - len(v_.([]byte))))
	return nil
}

func (mc *MemCache) Del(k string) (bool, error) {
	v, exist := mc.LoadAndDelete(k)
	if exist {
		mc.Stat.DecrKey(1)
		mc.Stat.DecrSz(int64(len(k) + len(v.([]byte))))
	}
	return exist, nil
}
