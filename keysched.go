package keysched

import "context"

type errno int

// Error returns the error text.
func (e errno) Error() string { return _errorTexts[e] }

// List of errors.
const (
	ErrKeyNotSet errno = iota
	ErrKeyNotFound
)

var _errorTexts = []string{
	ErrKeyNotSet:   "keysched: key not set",
	ErrKeyNotFound: "keysched: key not found",
}

// Value represents the scheduled item.
type Value[ID comparable, Key any] struct {
	KID ID
	Key Key
}

// Scheduler knows how to manage key schedule for each request.
type Scheduler[ID comparable, K any] interface {
	// Next returns the next scheduled item.
	Next(ctx context.Context) (Value[ID, K], error)

	// Find finds the key by its id.
	Find(ctx context.Context, id ID) (K, error)
}
