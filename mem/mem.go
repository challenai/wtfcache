package mem

import "cacheme/stat"

type MemCache struct {
}

func (mc *MemCache) Get(k string) ([]byte, error) {

}

func (mc *MemCache) Set(k string, v []byte) error {

}

func (mc *MemCache) Del(k string) error {

}

func (mc *MemCache) Info() stat.Stat {

}
