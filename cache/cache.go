package cache

import "cacheme/stat"

type Cache interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Del(string) error
	Info() stat.Stat
}
