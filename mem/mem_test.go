package mem

import (
	"bytes"
	"cacheme/tests"
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

type Entry struct {
	k string
	v []byte
}

func TestSetKey(t *testing.T) {
	entries, _, _ := mockEntries(100)
	mc := NewMemCache()
	for _, e := range entries {
		go func(k string, v []byte) {
			err := mc.Set(k, v)
			if err != nil {
				t.Fail()
			}
		}(e.k, e.v)
	}
}

func TestGetKey(t *testing.T) {
	entries, _, _ := mockEntries(100)
	mc := NewMemCache()
	for _, e := range entries {
		err := mc.Set(e.k, e.v)
		if err != nil {
			t.Fail()
		}
	}
	for _, e := range entries {
		go func(k string, v []byte) {
			v_, err := mc.Get(k)
			if err != nil || !bytes.Equal(v, v_) {
				t.Fail()
			}
		}(e.k, e.v)
	}
}

func TestDelKey(t *testing.T) {
	entries, _, _ := mockEntries(100)
	mc := NewMemCache()
	for _, e := range entries {
		err := mc.Set(e.k, e.v)
		if err != nil {
			t.Fail()
		}
	}
	wg := sync.WaitGroup{}
	wg.Add(len(entries))
	for _, e := range entries {
		go func(k string, v []byte) {
			exist, err := mc.Del(k)
			if err != nil || !exist {
				t.Fail()
			}
			wg.Done()
		}(e.k, e.v)
	}
	wg.Wait()
	for _, e := range entries {
		exist, err := mc.Del(e.k)
		if err != nil || exist {
			t.Fail()
		}
	}
}

func TestCountKeys(t *testing.T) {
	entries, _, _ := mockEntries(100)
	mc := NewMemCache()
	for _, e := range entries {
		err := mc.Set(e.k, e.v)
		if err != nil {
			t.Fail()
		}
	}
	// set the same keys again
	for _, e := range entries[:60] {
		err := mc.Set(e.k, e.v)
		if err != nil {
			t.Fail()
		}
	}
	if mc.Count() != 100 {
		t.Fail()
	}
}

func TestGetSz(t *testing.T) {
	entries, ksz, vsz := mockEntries(100)
	mc := NewMemCache()
	for _, e := range entries {
		err := mc.Set(e.k, e.v)
		if err != nil {
			t.Fail()
		}
	}
	for _, e := range entries[33:66] {
		err := mc.Set(e.k, e.v)
		if err != nil {
			t.Fail()
		}
	}
	if mc.GetSz() != int64(ksz+vsz) {
		t.Fail()
	}
}

// mock entries and return size of keys and values
func mockEntries(sz int) ([]Entry, int64, int64) {
	if sz <= 0 {
		return []Entry{}, 0, 0
	}
	entries := make([]Entry, sz)
	var ksz int64 = 0
	var vsz int64 = 0
	for i := 0; i < sz; i++ {
		entries[i] = Entry{
			k: fmt.Sprintf("%s%d", tests.RandStringRunes(2+rand.Intn(4)), i),
			// k: "key_" + strconv.Itoa(i),
			v: []byte(tests.RandStringRunes(8 + rand.Intn(32))),
		}
		ksz += int64(len(entries[i].k))
		vsz += int64(len(entries[i].v))
		// println(entries[i].k, string(entries[i].v))
	}
	return entries, ksz, vsz
}
