package memroundrobin

import (
	"context"
	"sync/atomic"

	"github.com/josestg/keysched"
)

// Scheduler is the round-robin scheduler that uses in-memory storage.
type Scheduler[ID comparable, K any] struct {
	n      int
	p      atomic.Int32
	values []keysched.Value[ID, K] // assume the values are read-only.
}

var _ keysched.Scheduler[string, any] = (*Scheduler[string, any])(nil)

// New returns a new round-robin scheduler.
func New[ID comparable, K any](values []keysched.Value[ID, K]) Scheduler[ID, K] {
	readonlyCopy := make([]keysched.Value[ID, K], len(values))
	copy(readonlyCopy, values)
	return Scheduler[ID, K]{
		n:      len(values),
		p:      atomic.Int32{},
		values: readonlyCopy,
	}
}

// Next returns the next scheduled item.
func (s *Scheduler[ID, K]) Next(ctx context.Context) (keysched.Value[ID, K], error) {
	select {
	case <-ctx.Done():
		return keysched.Value[ID, K]{}, ctx.Err()
	default:
	}

	if s.n == 0 {
		return keysched.Value[ID, K]{}, keysched.ErrKeyNotSet
	}
	idx := int(s.p.Load())
	v := s.values[idx]
	s.p.Store(int32((idx + 1) % s.n))
	return v, nil
}

// Find finds the key by its id.
func (s *Scheduler[ID, K]) Find(ctx context.Context, id ID) (key K, err error) {
	select {
	case <-ctx.Done():
		return key, ctx.Err()
	default:
	}
	for _, v := range s.values {
		if v.KID == id {
			return v.Key, nil
		}
	}
	return key, keysched.ErrKeyNotFound
}
