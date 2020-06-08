package dictpool

//go:generate msgp

// KV struct so it storages key/value data.
type KV struct {
	Key   []byte
	Value interface{}
}

// Dict dictionary as slice with better performance.
type Dict struct {
	// D slice of KV for storage the data
	D []KV
}

// DictMap dictionary as map.
type DictMap map[string]interface{}
