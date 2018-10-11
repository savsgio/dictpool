package dictpool

import (
	"fmt"
)

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
