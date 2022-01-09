package cache

type Cache interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Del(string) error
}
