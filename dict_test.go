package dictpool

import (
	"sync"
	"testing"
)

func Benchmark_DictPooltBytes(b *testing.B) {
	foo := []byte("DictPool")

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		d := AcquireDict()

		d.SetBytes(foo, 99)
		// d.GetBytes(foo)

		ReleaseDict(d)
	}
}

func Benchmark_DictPoolString(b *testing.B) {
	foo := "DictPool"

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		d := AcquireDict()

		d.Set(foo, 99)
		// d.Get(foo)

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

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		d := pool.Get().(dictKV)

		d.data["DictPool"] = 99

		// Reset
		for k := range d.data {
			delete(d.data, k)
		}

		pool.Put(d)
	}
}
