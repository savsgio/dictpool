package dictpool

import (
	"reflect"
	"sync"
	"testing"

	"github.com/savsgio/gotils/bytes"
)

func TestDict_allocKV(t *testing.T) {
	d := AcquireDict()
	allocs := 10

	for i := 0; i < allocs; i++ {
		if got := d.allocKV(); got == nil {
			t.Error("Dict.allocKV() returns nil KV pointer")
		}
	}

	if len(d.D) > allocs {
		t.Errorf("Dict.allocKV() len == %d, want %d", len(d.D), allocs)
	}
}

func TestDict_append(t *testing.T) {
	d := AcquireDict()
	key := "test"
	value := "hello"

	beforeLen := len(d.D)

	d.append(key, value)

	if len(d.D) == beforeLen {
		t.Error("Dict.append() it is not created a new entry")
	}

	if d.Get(key) != value {
		t.Errorf("Dict.append() it is not created a new entry with key '%s' and value '%s'", key, value)
	}
}

func TestDict_Len(t *testing.T) {
	d := AcquireDict()
	d.Set("key", "value")

	if d.Len() != len(d.D) {
		t.Errorf("Dict.Len() == %d, want %d", d.Len(), len(d.D))
	}
}

func TestDict_Swap(t *testing.T) {
	d := AcquireDict()
	k1 := "key1"
	k2 := "key2"
	v1 := "value1"
	v2 := "value2"

	d.Set(k1, v1)
	d.Set(k2, v2)

	d.Swap(0, 1)

	if d.D[0].Key != k2 {
		t.Error("Dict.Swap() not change KV position in []KV")
	}

	if d.D[1].Key != k1 {
		t.Error("Dict.Swap() not change KV position in []KV")
	}
}

func TestDict_Less(t *testing.T) {
	d := AcquireDict()
	k1 := "key1"
	k2 := "key2"
	v1 := "value1"
	v2 := "value2"

	d.Set(k1, v1)
	d.Set(k2, v2)

	want := k1 < k2
	got := d.Less(0, 1)

	if got != want {
		t.Errorf("Dict.Less() == %v, want %v", got, want)
	}
}

func TestDict_Get(t *testing.T) {
	const k, v = "key", "value"

	d := AcquireDict()
	d.Set(k, v)

	val := d.Get(k)
	if val != v {
		t.Errorf("Dict.Get() = '%v', want '%s'", val, v)
	}

	val = d.Get("other")
	if val != nil {
		t.Errorf("Dict.Get() = '%v', want '%v'", val, nil)
	}
}

func TestDict_GetBytes(t *testing.T) {
	const v = "value"

	d := AcquireDict()
	k := []byte("key")

	d.SetBytes(k, v)

	val := d.GetBytes(k)
	if val != v {
		t.Errorf("Dict.GetBytes() = '%v', want '%s'", val, v)
	}

	val = d.GetBytes([]byte("other"))
	if val != nil {
		t.Errorf("Dict.GetBytes() = '%v', want '%v'", val, nil)
	}
}

func TestDict_Set(t *testing.T) {
	const k, v = "key", "value"

	d := AcquireDict()
	d.Set(k, v)

	if !d.Has(k) {
		t.Error("Dict.Set() not set the new key and value")
	}

	newVal := "hello"
	d.Set(k, newVal)

	val := d.Get(k)
	if val != newVal {
		t.Errorf("Dict.Set() has not been updated the value")
	}
}

func TestDict_SetBytes(t *testing.T) {
	const v = "value"

	d := AcquireDict()
	k := []byte("key")

	d.SetBytes(k, v)

	if !d.HasBytes(k) {
		t.Error("Dict.SetBytes() not set the new key and value")
	}

	newVal := []byte("hello")
	d.SetBytes(k, newVal)

	val := d.GetBytes(k)
	if string(val.([]byte)) != string(newVal) { // nolint:forcetypeassert
		t.Errorf("Dict.Set() has not been updated the value")
	}
}

func TestDict_Del(t *testing.T) {
	const k, v = "key", "value"

	d := AcquireDict()
	d.Set(k, v)
	d.Del(k)

	if d.Has(k) {
		t.Errorf("Dict.Del() not delete the key '%s'", k)
	}
}

func TestDict_DelBytes(t *testing.T) {
	const v = "value"

	d := AcquireDict()
	k := []byte("key")

	d.SetBytes(k, v)
	d.Set("dsad12", v)
	d.Set("dsad123124", v)
	d.Set("dsad234234545", v)
	d.DelBytes(k)

	if d.HasBytes(k) {
		t.Errorf("Dict.DelBytes() not delete the key '%s'", string(k))
	}
}

func TestDict_Has(t *testing.T) {
	const k, v = "key", "value"

	d := AcquireDict()
	d.Set(k, v)

	if got := d.Has(k); !got {
		t.Errorf("Dict.Has() = '%v', want '%v'", got, true)
	}
}

func TestDict_HasBytes(t *testing.T) {
	d := AcquireDict()
	k := []byte("key")
	v := "value"

	d.SetBytes(k, v)

	if got := d.HasBytes(k); !got {
		t.Errorf("Dict.HasBytes() = '%v', want '%v'", got, true)
	}
}

func TestDict_Reset(t *testing.T) {
	d := AcquireDict()
	d.Set("Test", true)
	d.Reset()

	if len(d.D) > 0 {
		t.Error("Dict.Reset() the length of Dict is not 0")
	}
}

func TestDict_Map(t *testing.T) {
	const k, v, k2 = "key", "value", "subkey"

	const subK, subV = "subK", "subV"

	m := make(DictMap)

	want := DictMap{
		k: v,
		k2: DictMap{
			subK: subV,
		},
	}

	d1 := AcquireDict()
	d2 := AcquireDict()

	d2.Set(subK, subV)

	d1.Set(k, v)
	d1.Set(k2, d2)
	d1.Map(m)

	if !reflect.DeepEqual(m, want) {
		t.Errorf("Dict.Map() == %v, want %v", m, want)
	}
}

func isEqual(d *Dict, dm map[string]interface{}) bool {
	for k, v := range dm {
		val := d.Get(k)

		if sv, ok := v.(map[string]interface{}); ok {
			return isEqual(val.(*Dict), sv) // nolint:forcetypeassert
		}

		if val != v {
			return false
		}
	}

	return true
}

func TestDict_Parse(t *testing.T) {
	const k, v, k2 = "key", "value", "subkey"

	const subK, subV = "subK", "subV"

	d1 := AcquireDict()
	d2 := AcquireDict()
	d3 := AcquireDict()

	d3.Set(subK, subV)
	d2.Set(k, v)
	d2.Set(k2, d3)

	m := make(DictMap)
	m[k] = v
	m[k2] = map[string]interface{}{subK: subV}
	d1.Parse(m)

	if len(m) == 0 || !isEqual(d1, m) {
		t.Errorf("Dict.Parse() == %v, want %v", d1, d2)
	}
}

func genKeys(tb testing.TB, size int) []string {
	tb.Helper()

	keys := []string{}

	for i := 0; i < size; i++ {
		keys = append(keys, string(bytes.Rand(make([]byte, 10))))
	}

	return keys
}

func benchmarkGet(b *testing.B, d *Dict, items int) {
	b.Helper()

	keys := genKeys(b, items)
	want := 12345 % items

	for i := range keys {
		d.Set(keys[i], i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.Get(keys[want])
	}
}

func Benchmark_Get(b *testing.B) {
	d := AcquireDict()

	benchmarkGet(b, d, 10)
}

func Benchmark_GetBinary(b *testing.B) {
	d := AcquireDict()
	d.BinarySearch = true

	benchmarkGet(b, d, 10)
}

func Benchmark_GetBigHeap(b *testing.B) {
	d := AcquireDict()

	benchmarkGet(b, d, 1000)
}

func Benchmark_GetBinaryBigHeap(b *testing.B) {
	d := AcquireDict()
	d.BinarySearch = true

	benchmarkGet(b, d, 1000)
}

func benchmarkSet(b *testing.B, d *Dict, items int) {
	b.Helper()

	keys := genKeys(b, items)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.Set(keys[i%items], i)
	}
}

func Benchmark_Set(b *testing.B) {
	d := AcquireDict()

	benchmarkSet(b, d, 10)
}

func Benchmark_SetBinary(b *testing.B) {
	d := AcquireDict()
	d.BinarySearch = true

	benchmarkSet(b, d, 10)
}

func Benchmark_SetBigHeap(b *testing.B) {
	d := AcquireDict()

	benchmarkSet(b, d, 1000)
}

func Benchmark_SetBinaryBigHeap(b *testing.B) {
	d := AcquireDict()
	d.BinarySearch = true

	benchmarkSet(b, d, 1000)
}

func benchmarkDel(b *testing.B, d *Dict, items int) {
	b.Helper()

	keys := genKeys(b, items)
	for i := range keys {
		d.Set(keys[i], i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.Del(keys[i%items])
	}
}

func Benchmark_Del(b *testing.B) {
	d := AcquireDict()

	benchmarkDel(b, d, 10)
}

func Benchmark_DelBinary(b *testing.B) {
	d := AcquireDict()
	d.BinarySearch = true

	benchmarkDel(b, d, 10)
}

func Benchmark_DelBigHeap(b *testing.B) {
	d := AcquireDict()

	benchmarkDel(b, d, 1000)
}

func Benchmark_DelBinaryBigHeap(b *testing.B) {
	d := AcquireDict()
	d.BinarySearch = true

	benchmarkDel(b, d, 1000)
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

func BenchmarkDict(b *testing.B) {
	keys := genKeys(b, 100)

	u := AcquireDict()
	// u.BinarySearch = true

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j, key := range keys {
			u.Set(key, j)
		}

		for _, key := range keys {
			_ = u.Get(key)
		}

		u.Reset()
	}
}

func BenchmarkStdMap(b *testing.B) {
	keys := genKeys(b, 100)
	u := make(map[string]interface{})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j, key := range keys {
			u[key] = j
		}

		for _, key := range keys {
			_ = u[key]
		}

		u = make(map[string]interface{})
	}
}

func BenchmarkSyncMap(b *testing.B) {
	keys := genKeys(b, 100)
	u := new(sync.Map)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j, key := range keys {
			u.Store(key, j)
		}

		for _, key := range keys {
			_, _ = u.Load(key)
		}

		u.Range(func(key, _ interface{}) bool {
			u.Delete(key)

			return true
		})
	}
}
