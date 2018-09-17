package dictpool

import (
	"sync"
	"testing"
)

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
			return dictKV{
				data: make(map[string]interface{}),
			}
		},
	}

	key := "DictPool"

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		d := pool.Get().(dictKV)

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
