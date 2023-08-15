package run

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	t.Run("no_retry", func(t *testing.T) {
		count := 0

		r := require.New(t)

		f := WrapE(
			func() error {
				count++
				return fmt.Errorf("no")
			},
			WithMaxRetries(0),
			WithConstantBackoff(time.Millisecond),
		)

		r.ErrorContains(f(), "no")
		r.Equal(count, 1)
	})

	t.Run("retry", func(t *testing.T) {
		count := 0

		r := require.New(t)

		f := WrapE(
			func() error {
				count++
				return fmt.Errorf("no")
			},
			WithMaxRetries(3),
			WithConstantBackoff(time.Millisecond),
		)

		r.ErrorContains(f(), "no")
		r.Equal(count, 4)
	})
}
