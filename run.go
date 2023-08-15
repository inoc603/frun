package run

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/sync/singleflight"
)

type runner struct {
	backoff     backoff.BackOff
	maxRetries  uint64
	middlewares []Middleware
	ctx         context.Context
	ctxPos      int
}

func newRunner(opts ...Option) *runner {
	res := &runner{
		backoff:    backoff.NewConstantBackOff(0),
		maxRetries: 0,
		ctx:        context.TODO(),
	}
	for _, o := range opts {
		o(res)
	}
	return res
}

type CtxErrFunc func(ctx context.Context) error

type Middleware func(ctx context.Context, f CtxErrFunc) error

func (r *runner) wrap(f CtxErrFunc) func() error {
	wrapped := f
	for _, m := range r.middlewares {
		last := wrapped
		wrapped = func(ctx context.Context) error { return m(ctx, last) }
	}
	return func() error { return wrapped(r.ctx) }
}

func (r *runner) run(f func(context.Context) error) error {
	return backoff.Retry(r.wrap(f), backoff.WithMaxRetries(r.backoff, r.maxRetries))
}

type Option func(r *runner)

func WithMaxRetries(max uint64) Option {
	return func(r *runner) {
		if max > 0 {
			r.maxRetries = max
		}
	}
}

func WithConstantBackoff(d time.Duration) Option {
	return func(r *runner) {
		r.backoff = backoff.NewConstantBackOff(d)
	}
}

func WithExponentialBackOff(initialInterval time.Duration) Option {
	return func(r *runner) {
		// TODO: more customization for backoff
		b := backoff.NewExponentialBackOff()
		b.InitialInterval = initialInterval
		r.backoff = b
	}
}

func WithMiddleware(m Middleware) Option {
	return func(r *runner) {
		r.middlewares = append(r.middlewares, m)
	}
}

func WithWaitGroup(wg *sync.WaitGroup) Option {
	return WithMiddleware(func(ctx context.Context, f CtxErrFunc) error {
		wg.Add(1)
		defer wg.Done()
		return f(ctx)
	})
}

func WithRecover() Option {
	return WithMiddleware(func(ctx context.Context, f CtxErrFunc) (err error) {
		defer func() {
			if re := recover(); re != nil {
				err = fmt.Errorf("panic recovered: %v", err)
			}
		}()
		return f(ctx)
	})
}

func WithSingleflight() Option {
	var g singleflight.Group
	return WithMiddleware(func(ctx context.Context, f CtxErrFunc) error {
		_, err, _ := g.Do("run", func() (interface{}, error) {
			return nil, f(ctx)
		})
		return err
	})
}

//go:generate go run ./internal/gen ./wrap.go
func Wrap(f func(), opts ...Option) func() {
	return Wrap00(f, opts...)
}

func WrapE(f func() error, opts ...Option) func() error {
	return WrapE01(f, opts...)
}

func WrapCE(f func(context.Context) error, opts ...Option) func(context.Context) error {
	return WrapCE11(f, opts...)
}
