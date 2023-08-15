package run

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChunk(t *testing.T) {
	cases := []struct {
		size      int
		chunkSize int
		output    [][2]int
	}{
		{
			size:      8,
			chunkSize: 2,
			output: [][2]int{
				{0, 2},
				{2, 4},
				{4, 6},
				{6, 8},
			},
		},
		{
			size:      9,
			chunkSize: 4,
			output: [][2]int{
				{0, 4},
				{4, 8},
				{8, 9},
			},
		},
		{
			size:      0,
			chunkSize: 4,
			output:    [][2]int{},
		},
		{
			size:      -1,
			chunkSize: 4,
			output:    [][2]int{},
		},
		{
			size:      -1,
			chunkSize: -1,
			output:    [][2]int{},
		},
		{
			size:      3,
			chunkSize: 0,
			output: [][2]int{
				{0, 3},
			},
		},
		{
			size:      1,
			chunkSize: 10,
			output: [][2]int{
				{0, 1},
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%d:%d", c.size, c.chunkSize), func(t *testing.T) {
			r := require.New(t)

			chunk := NewChunks(c.size, c.chunkSize)

			var called, total int

			var values []int
			if c.size >= 0 {
				values = make([]int, c.size)
			}

			for chunk.Next() {
				called++
				r.Equal(c.output[called-1][0], chunk.Start())
				r.Equal(c.output[called-1][1], chunk.End())
				total += len(values[chunk.Start():chunk.End()])
			}

			r.Equal(len(c.output), called)
			if c.size >= 0 {
				r.Equal(c.size, total)
			} else {
				r.Equal(0, total)
			}
		})
	}
}
