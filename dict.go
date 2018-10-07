package dictpool

import (
	"sync"
)

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

func allocKV(data []KV) ([]KV, *KV) {
	n := len(data)

	if cap(data) > n {
		data = data[:n+1]
	} else {
		data = append(data, KV{})
	}

	return data, &data[n]
}

func appendArgs(data []KV, key string, value interface{}) []KV {
	data, kv := allocKV(data)

	kv.Key = append(kv.Key[:0], key...)
	kv.Value = value

	return data
}

func swap(data []KV, i, j int) []KV {
	data[i], data[j] = data[j], data[i]

	return data
}

func getArgs(data []KV, key string) *KV {
	n := len(data)
	for i := 0; i < n; i++ {
		kv := &data[i]
		if key == b2s(kv.Key) {
			return kv
		}
	}

	return nil
}

func setArgs(data []KV, key string, value interface{}) []KV {
	kv := getArgs(data, key)
	if kv != nil {
		kv.Value = value
		return data
	}

	return appendArgs(data, key, value)
}

func delArgs(data []KV, key string) []KV {
	for i, n := 0, len(data); i < n; i++ {
		kv := &data[i]
		if key == b2s(kv.Key) {
			n--
			if i != n {
				swap(data, i, n)
				i--
			}
			data = data[:n] // Remove last position
		}
	}

	return data
}

func hasArgs(data []KV, key string) bool {
	for i, n := 0, len(data); i < n; i++ {
		kv := &data[i]
		if key == b2s(kv.Key) {
			return true
		}
	}

	return false
}

// Get get data from key
func (d *Dict) Get(key string) interface{} {
	kv := getArgs(d.D, key)
	if kv != nil {
		return kv.Value
	}

	return nil
}

// GetBytes get data from key
func (d *Dict) GetBytes(key []byte) interface{} {
	kv := getArgs(d.D, b2s(key))
	if kv != nil {
		return kv.Value
	}

	return nil
}

// Set set new key
func (d *Dict) Set(key string, value interface{}) {
	d.D = setArgs(d.D, key, value)
}

// SetBytes set new key
func (d *Dict) SetBytes(key []byte, value interface{}) {
	d.D = setArgs(d.D, b2s(key), value)
}

// Del delete key
func (d *Dict) Del(key string) {
	d.D = delArgs(d.D, key)
}

// DelBytes delete key
func (d *Dict) DelBytes(key []byte) {
	d.D = delArgs(d.D, b2s(key))
}

// Has check if key exists
func (d *Dict) Has(key string) bool {
	return hasArgs(d.D, key)
}

// HasBytes check if key exists
func (d *Dict) HasBytes(key []byte) bool {
	return hasArgs(d.D, b2s(key))
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
