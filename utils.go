package dictpool

import (
	"fmt"
	"unsafe"
)

// b2s converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func printDict(data *Dict) {
	for _, kv := range data.D {
		switch kv.Value.(type) {
		case *Dict:
			print(string(kv.Key) + ": ")
			printDict(kv.Value.(*Dict))
		default:
			fmt.Println(string(kv.Key), ": ", kv.Value)
		}
	}
}
