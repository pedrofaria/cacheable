package ristretto

import (
	"context"
	"fmt"
)

func ExampleNew() {
	clt, _ := NewCache(&Config[string, []byte]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	driver := New(clt)
	defer driver.Close()

	ctx := context.Background()

	if err := driver.Set(ctx, "key1", []byte("value1"), 0); err != nil {
		fmt.Println(err.Error())
	}

	data, err := driver.Get(ctx, "key1")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(data))

	if err := driver.Del(ctx, "key1"); err != nil {
		panic(err)
	}

	// Output: value1
}
