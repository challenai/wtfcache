package mem

import "sync/atomic"

type MemStat struct {
	// entries count
	keys int64
	// size of all the entries
	sz int64
}

func (ms *MemStat) Count() int64 {
	return atomic.LoadInt64(&ms.keys)
}

func (ms *MemStat) GetSz() int64 {
	return atomic.LoadInt64(&ms.sz)
}

func (ms *MemStat) incr(num int64) {
	atomic.AddInt64(&ms.keys, num)
}

func (ms *MemStat) incrSz(sz int64) {
	atomic.AddInt64(&ms.sz, sz)
}

func (ms *MemStat) decr(num int64) {
	atomic.AddInt64(&ms.keys, -num)
}

func (ms *MemStat) decrSz(sz int64) {
	atomic.AddInt64(&ms.sz, -sz)
}
