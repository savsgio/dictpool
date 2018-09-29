dictpool
========


[![Build Status](https://travis-ci.org/savsgio/dictpool.svg?branch=master)](https://travis-ci.org/savsgio/dictpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/savsgio/dictpool)](https://goreportcard.com/report/github.com/savsgio/dictpool)
[![GoDoc](https://godoc.org/github.com/savsgio/dictpool?status.svg)](https://godoc.org/github.com/savsgio/dictpool)

Memory store like `map[string]interface{}` with better performance and safe concurrency.

## Benchmarks:
```
Benchmark_DictPool-4            20000000                61.4 ns/op             0 B/op          0 allocs/op
Benchmark_DictPoolBytes-4       20000000                59.2 ns/op             0 B/op          0 allocs/op
Benchmark_DictMap-4             20000000               102 ns/op               0 B/op          0 allocs/op
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
