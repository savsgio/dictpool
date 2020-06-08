package dictpool

import (
	"sync"
	"testing"
)

func TestReleaseDict(t *testing.T) {
	d := AcquireDict()

	d.Set("key", "value")

	ReleaseDict(d)

	if len(d.D) > 0 {
		t.Error("ReleaseDict() not reset the Dict")
	}
}

func TestDict_Reset(t *testing.T) {
	d := AcquireDict()
	d.Set("Test", true)
	d.Reset()

	if len(d.D) > 0 {
		t.Error("Reset() the length of Dict is not 0")
	}
}

func Test_allocKV(t *testing.T) {
	d := AcquireDict()

	newD, kv := allocKV(d.D)

	if len(d.D) == len(newD) {
		t.Error("allocKV() it is not created a new entry")
	}

	if kv == nil {
		t.Error("allocKV() it returns a nil KV pointer")
	}
}

func Test_appendArgs(t *testing.T) {
	d := AcquireDict()
	key := "test"
	value := "hello"

	newD := appendArgs(d.D, key, value)

	if len(d.D) == len(newD) {
		t.Error("appendArgs() it is not created a new entry")
	}

	d.D = newD
	if d.Get(key) != value {
		t.Errorf("appendArgs() it is not created a new entry with key '%s' and value '%s'", key, value)
	}
}

func Test_swap(t *testing.T) {
	d := AcquireDict()
	k1 := "key1"
	k2 := "key2"
	v1 := "value1"
	v2 := "value2"

	d.Set(k1, v1)
	d.Set(k2, v2)

	newD := swap(d.D, 0, 1)

	if string(newD[0].Key) != k2 {
		t.Error("swap() not change KV position in []KV")
	}

	if string(newD[1].Key) != k1 {
		t.Error("swap() not change KV position in []KV")
	}
}

func TestDict_Get(t *testing.T) {
	const k, v = "key", "value"

	d := AcquireDict()
	d.Set(k, v)

	val := d.Get(k)
	if val != v {
		t.Errorf("Get() = '%v', want '%s'", val, v)
	}
}

func TestDict_GetBytes(t *testing.T) {
	const v = "value"

	d := AcquireDict()
	k := []byte("key")

	d.SetBytes(k, v)

	val := d.GetBytes(k)
	if val != v {
		t.Errorf("GetBytes() = '%v', want '%s'", val, v)
	}
}

func TestDict_Set(t *testing.T) {
	const k, v = "key", "value"

	d := AcquireDict()
	d.Set(k, v)

	if !d.Has(k) {
		t.Error("Set() not set the new key and value")
	}
}

func TestDict_SetBytes(t *testing.T) {
	const v = "value"

	d := AcquireDict()
	k := []byte("key")

	d.SetBytes(k, v)

	if !d.HasBytes(k) {
		t.Error("SetBytes() not set the new key and value")
	}
}

func TestDict_Del(t *testing.T) {
	const k, v = "key", "value"

	d := AcquireDict()
	d.Set(k, v)
	d.Del(k)

	if d.Has(k) {
		t.Errorf("Del() not delete the key '%s'", k)
	}
}

func TestDict_DelBytes(t *testing.T) {
	const v = "value"

	d := AcquireDict()
	k := []byte("key")

	d.SetBytes(k, v)
	d.DelBytes(k)

	if d.HasBytes(k) {
		t.Errorf("DelBytes() not delete the key '%s'", string(k))
	}
}

func TestDict_Has(t *testing.T) {
	const k, v = "key", "value"

	d := AcquireDict()
	d.Set(k, v)

	if got := d.Has(k); !got {
		t.Errorf("Has() = '%v', want '%v'", got, true)
	}
}

func TestDict_HasBytes(t *testing.T) {
	d := AcquireDict()
	k := []byte("key")
	v := "value"

	d.SetBytes(k, v)

	if got := d.HasBytes(k); !got {
		t.Errorf("HasBytes() = '%v', want '%v'", got, true)
	}
}

func TestDict_Map(t *testing.T) {
	const k, v = "key", "value"

	m := make(map[string]interface{})

	d := AcquireDict()
	d.Set(k, v)
	d.Map(m)

	if mv, ok := m[k]; !ok {
		t.Errorf("Map() the key '%v' is not set into the map", k)
	} else if mv != v {
		t.Errorf("Map() the value of key '%v' in map is '%v', want '%v'", k, mv, v)
	}
}

func TestDict_Parse(t *testing.T) {
	d := AcquireDict()
	k := "key"
	v := "value"

	m := make(map[string]interface{})
	m[k] = v
	d.Parse(m)

	if !d.Has(k) {
		t.Errorf("Parse() the key '%v' is not set into the Dict", k)
	} else if d.Get(k) != v {
		t.Errorf("Parse() the value of key '%v' in Dict is '%v', want '%v'", k, d.Get(k), v)
	}
}

func Benchmark_DictPool(b *testing.B) {
	foo := "DictPool"

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		d := AcquireDict()
		d.Set(foo, 99)
		d.Get(foo)
		ReleaseDict(d)
	}
}

func Benchmark_DictPoolBytes(b *testing.B) {
	foo := []byte("DictPool")

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		d := AcquireDict()
		d.SetBytes(foo, 99)
		d.GetBytes(foo)
		ReleaseDict(d)
	}
}

func Benchmark_DictMap(b *testing.B) {
	type dictKV struct {
		data map[string]interface{}
	}

	pool := sync.Pool{
		New: func() interface{} {
			return &dictKV{
				data: make(map[string]interface{}),
			}
		},
	}

	key := "DictPool"

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		d := pool.Get().(*dictKV)

		d.data[key] = 99

		if _, ok := d.data[key]; ok {
			// Reset
			for k := range d.data {
				delete(d.data, k)
			}
		}

		pool.Put(d)
	}
}

func Benchmark_Map(b *testing.B) {
	d1 := AcquireDict()

	d1.Set("Foo", "Bar")
	d1.Set("Foo2", "Bar2")

	m := make(DictMap)

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		d1.Map(m)
	}
	b.StopTimer()

	ReleaseDict(d1)
}

func Benchmark_Parse(b *testing.B) {
	m := map[string]interface{}{
		"Hola":  true,
		"Adios": false,
	}

	d := AcquireDict()

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		d.Parse(m)
	}
	b.StopTimer()

	ReleaseDict(d)
}
