package dictpool

import (
	"sync"
)

// KV struct so it storages key/value data
type KV struct {
	key   []byte
	value interface{}
}

// Dict struct for imitate map[key]value
type Dict struct {
	data []KV
}

var defaultPool = sync.Pool{
	New: func() interface{} {
		return new(Dict)
	},
}

// AcquireDict acquire new dict
func AcquireDict() *Dict {
	return defaultPool.Get().(*Dict)
}

// ReleaseDict release dict
func ReleaseDict(d *Dict) {
	d.Reset()
	defaultPool.Put(d)
}

// Reset reset dict
func (d *Dict) Reset() {
	d.data = d.data[:0]
}

func (d *Dict) allocKV() *KV {
	n := len(d.data)

	if cap(d.data) > n {
		d.data = d.data[:n+1]
	} else {
		d.data = append(d.data, KV{})
	}

	return &d.data[n]
}

func (d *Dict) appendArgs(key string, value interface{}) {
	kv := d.allocKV()

	kv.key = append(kv.key[:0], key...)
	kv.value = value
}

func (d *Dict) getArgs(key string) *KV {
	n := len(d.data)
	for i := 0; i < n; i++ {
		kv := &d.data[i]
		if key == string(kv.key) {
			return kv
		}
	}

	return nil
}

func (d *Dict) setArgs(key string, value interface{}) {
	kv := d.getArgs(key)
	if kv != nil {
		kv.value = value
		return
	}

	d.appendArgs(key, value)
}

// Get get data from key
func (d *Dict) Get(key string) interface{} {
	kv := d.getArgs(key)
	if kv != nil {
		return kv.value
	}

	return nil
}

// GetBytes get data from key
func (d *Dict) GetBytes(key []byte) interface{} {
	kv := d.getArgs(string(key))
	if kv != nil {
		return kv.value
	}

	return nil
}

// Set set new key
func (d *Dict) Set(key string, value interface{}) {
	d.setArgs(key, value)
}

// SetBytes set new key
func (d *Dict) SetBytes(key []byte, value interface{}) {
	d.setArgs(string(key), value)
}
