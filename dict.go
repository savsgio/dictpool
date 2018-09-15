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
	// D slice of KV for storage the data
	D []KV
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
	d.D = d.D[:0]
}

func (d *Dict) allocKV() *KV {
	n := len(d.D)

	if cap(d.D) > n {
		d.D = d.D[:n+1]
	} else {
		d.D = append(d.D, KV{})
	}

	return &d.D[n]
}

func (d *Dict) appendArgs(key string, value interface{}) {
	kv := d.allocKV()

	kv.key = append(kv.key[:0], key...)
	kv.value = value
}

func (d *Dict) swap(i, j int) {
	d.D[i], d.D[j] = d.D[j], d.D[i]
}

func (d *Dict) getArgs(key string) *KV {
	n := len(d.D)
	for i := 0; i < n; i++ {
		kv := &d.D[i]
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

func (d *Dict) delArgs(key string) {
	for i, n := 0, len(d.D); i < n; i++ {
		kv := &d.D[i]
		if key == string(kv.key) {
			n--
			if i != n {
				d.swap(i, n)
				i--
			}
			d.D = d.D[:n] // Remove last position
		}
	}
}

func (d *Dict) hasArgs(key string) bool {
	for i, n := 0, len(d.D); i < n; i++ {
		kv := &d.D[i]
		if key == string(kv.key) {
			return true
		}
	}

	return false
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

// Del delete key
func (d *Dict) Del(key string) {
	d.delArgs(key)
}

// DelBytes delete key
func (d *Dict) DelBytes(key []byte) {
	d.delArgs(string(key))
}

// Has check if key exists
func (d *Dict) Has(key string) bool {
	return d.hasArgs(key)
}

// HasBytes check if key exists
func (d *Dict) HasBytes(key []byte) bool {
	return d.hasArgs(string(key))
}
