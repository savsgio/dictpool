dictpool
========

[![Go Report Card](https://goreportcard.com/badge/github.com/savsgio/dictpool)](https://goreportcard.com/report/github.com/savsgio/dictpool)

Dictionary (Python like) with better performance than the `map[key]value`

```
Benchmark_DictPooltBytes-4      30000000                52.7 ns/op             0 B/op          0 allocs/op
Benchmark_DictPoolString-4      30000000                45.2 ns/op             0 B/op          0 allocs/op
Benchmark_DictMap-4             10000000               144 ns/op               0 B/op          0 allocs/op
```

## Example:
```go
d := AcquireDict()

d.Set("foo", "Hello DictPool")

fmt.Println(d.Get("foo"))  // Output: Hello DictPool

ReleaseDict(d)
```
