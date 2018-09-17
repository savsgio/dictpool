package dictpool

import (
	"sync"
)

// KV struct so it storages key/value data
type KV struct {
	Key   []byte
	Value interface{}
}

// DictMap dictionary as map
type DictMap map[string]interface{}

// Dict dictionary as slice with better performance
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

	kv.Key = append(kv.Key[:0], key...)
	kv.Value = value
}

func (d *Dict) swap(i, j int) {
	d.D[i], d.D[j] = d.D[j], d.D[i]
}

func (d *Dict) getArgs(key string) *KV {
	n := len(d.D)
	for i := 0; i < n; i++ {
		kv := &d.D[i]
		if key == b2s(kv.Key) {
			return kv
		}
	}

	return nil
}

func (d *Dict) setArgs(key string, value interface{}) {
	kv := d.getArgs(key)
	if kv != nil {
		kv.Value = value
		return
	}

	d.appendArgs(key, value)
}

func (d *Dict) delArgs(key string) {
	for i, n := 0, len(d.D); i < n; i++ {
		kv := &d.D[i]
		if key == b2s(kv.Key) {
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
		if key == b2s(kv.Key) {
			return true
		}
	}

	return false
}

// Get get data from key
func (d *Dict) Get(key string) interface{} {
	kv := d.getArgs(key)
	if kv != nil {
		return kv.Value
	}

	return nil
}

// GetBytes get data from key
func (d *Dict) GetBytes(key []byte) interface{} {
	kv := d.getArgs(b2s(key))
	if kv != nil {
		return kv.Value
	}

	return nil
}

// Set set new key
func (d *Dict) Set(key string, value interface{}) {
	d.setArgs(key, value)
}

// SetBytes set new key
func (d *Dict) SetBytes(key []byte, value interface{}) {
	d.setArgs(b2s(key), value)
}

// Del delete key
func (d *Dict) Del(key string) {
	d.delArgs(key)
}

// DelBytes delete key
func (d *Dict) DelBytes(key []byte) {
	d.delArgs(b2s(key))
}

// Has check if key exists
func (d *Dict) Has(key string) bool {
	return d.hasArgs(key)
}

// HasBytes check if key exists
func (d *Dict) HasBytes(key []byte) bool {
	return d.hasArgs(b2s(key))
}

// Map convert to map
func (d *Dict) Map(dst DictMap) {
	for _, kv := range d.D {
		switch kv.Value.(type) {
		case *Dict:
			subDst := make(DictMap)
			kv.Value.(*Dict).Map(subDst)
			dst[b2s(kv.Key)] = subDst
		default:
			dst[b2s(kv.Key)] = kv.Value
		}
	}
}

// Parse convert map to Dict
func (d *Dict) Parse(src DictMap) {
	for k, v := range src {
		switch v.(type) {
		case map[string]interface{}:
			subDict := new(Dict)
			subDict.Parse(v.(map[string]interface{}))

			d.Set(k, subDict)
		default:
			d.Set(k, v)
		}
	}
}
