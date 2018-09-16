package dictpool

import (
	"sync"
	"testing"
)

func Benchmark_DictPoolSet(b *testing.B) {
	foo := "DictPool"

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		d := AcquireDict()
		d.Set(foo, 99)
		ReleaseDict(d)
	}
}

func Benchmark_DictPooltSetBytes(b *testing.B) {
	foo := []byte("DictPool")

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		d := AcquireDict()
		d.SetBytes(foo, 99)
		ReleaseDict(d)
	}
}

func Benchmark_DictPoolGet(b *testing.B) {
	foo := "DictPool"

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		d := AcquireDict()
		d.Set(foo, 99)
		d.Get(foo)
		ReleaseDict(d)
	}
}

func Benchmark_DictPoolGetBytes(b *testing.B) {
	foo := []byte("DictPool")

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		d := AcquireDict()
		d.SetBytes(foo, 99)
		d.GetBytes(foo)
		ReleaseDict(d)
	}
}

func Benchmark_DictMapSet(b *testing.B) {
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

		// Reset
		for k := range d.data {
			delete(d.data, k)
		}

		pool.Put(d)
	}
}

func Benchmark_DictMapGet(b *testing.B) {
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

func Benchmark_Marshal(b *testing.B) {
	d1 := AcquireDict()
	d2 := AcquireDict()

	d1.Set("Foo", "Bar")
	d2.Set("DictPool", d1)

	buff := []byte(nil)

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		buff, _ = d2.Marshal()

		buff = buff[:0]
	}
	b.StopTimer()

	ReleaseDict(d1)
	ReleaseDict(d2)
}

func Benchmark_Unmarshal(b *testing.B) {
	jsonStr := []byte("{\"DictPool\":{\"Foo\":\"Bar\"}}")
	d := AcquireDict()

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		d.Unmarshal(jsonStr)
	}
	b.StopTimer()

	ReleaseDict(d)
}

// func printDict(data *Dict) {
// 	for _, kv := range data.D {
// 		switch kv.Value.(type) {
// 		case *Dict:
// 			print(string(kv.Key) + ": ")
// 			printDict(kv.Value.(*Dict))
// 		default:
// 			fmt.Println(string(kv.Key), ": ", kv.Value)
// 		}
// 	}
// }
