package run

type Chunks struct {
	size         int
	chunkSize    int
	currentStart int
	currentEnd   int
}

// NewChunks returns a Chunks instance that helps with splitting a slice into chunks.
func NewChunks(size, chunkSize int) *Chunks {
	if size < 0 {
		size = 0
	}

	if chunkSize <= 0 {
		if size > 0 {
			chunkSize = size
		} else {
			chunkSize = 1
		}
	}

	return &Chunks{size: size, chunkSize: chunkSize}
}

// Next returns whether there are remaining chunks.
// Typicallly this should be used with a for loop.
//
//	for c.Next() {
//		list[c.Start():c.End()]
//	}
func (c *Chunks) Next() bool {
	if c.size == 0 {
		return false
	}

	if c.currentEnd == 0 {
		c.currentEnd = c.currentStart + c.chunkSize
		if c.currentEnd >= c.size {
			c.currentEnd = c.size
		}
		return true
	}

	c.currentStart += c.chunkSize

	if c.currentStart >= c.size {
		c.currentStart = c.size
	}

	c.currentEnd = c.currentStart + c.chunkSize
	if c.currentEnd >= c.size {
		c.currentEnd = c.size
	}

	return c.currentStart < c.size
}

// Start returns the start of the current chunk.
func (c *Chunks) Start() int {
	return c.currentStart
}

// End returns the end of the current chunk.
// End() is meant to be used as the second parameter in a slice expression.
func (c *Chunks) End() int {
	return c.currentEnd
}
