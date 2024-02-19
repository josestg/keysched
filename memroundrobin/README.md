# In-Memory Round-Robin Key Scheduler

## Usage

```go
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/josestg/keysched"
	"github.com/josestg/keysched/memroundrobin"
)

func main() {
	rr := memroundrobin.New(genRandRSAKeys(3))
	for i := 0; i < 5; i++ {
		key, err := rr.Next(context.Background())
		fmt.Printf("req_id: %d, next: %+v, error: %v\n", i, key, err)
	}

	// or find key by its id.
	key, err := rr.Find(context.Background(), 1)
	fmt.Printf("found: %+v, error: %v\n", key, err)
}

func genRandRSAKeys(n int) []keysched.Value[int, *rsa.PrivateKey] {
	values := make([]keysched.Value[int, *rsa.PrivateKey], 0, n)
	for i := 0; i < n; i++ {
		pk, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}

		val := keysched.Value[int, *rsa.PrivateKey]{KID: i, Key: pk}
		values = append(values, val)
	}
	return values
}
```