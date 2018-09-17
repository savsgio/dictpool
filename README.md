dictpool
========

[![Go Report Card](https://goreportcard.com/badge/github.com/savsgio/dictpool)](https://goreportcard.com/report/github.com/savsgio/dictpool)

Memory store like `map[string]interface{}` with better performance.

## Benchmarks:
```
Benchmark_DictPool-4            30000000                57.4 ns/op             0 B/op          0 allocs/op
Benchmark_DictPoolBytes-4       20000000                55.8 ns/op             0 B/op          0 allocs/op
Benchmark_DictMap-4             20000000                96.3 ns/op             0 B/op          0 allocs/op
```

*Benchmark with Go 1.11*

## Example:
```go
d := dictpool.AcquireDict()
key := "foo"

d.Set(key, "Hello DictPool")

if d.Has(key){
    fmt.Println(d.Get(key))  // Output: Hello DictPool
}

d.Del(key)

dictpool.ReleaseDict(d)
```
