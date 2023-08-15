package run_test

import (
	"fmt"

	"audit/minos/internal/run"
)

func ExampleNewChunks() {
	values := make([]int, 100)
	for i := range values {
		values[i] = i
	}

	for c := run.NewChunks(len(values), 10); c.Next(); {
		chunk := values[c.Start():c.End()]
		fmt.Println(len(chunk), chunk)
	}
}
