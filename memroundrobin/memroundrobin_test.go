package memroundrobin

import (
	"context"
	"errors"
	"testing"

	"github.com/josestg/keysched"
)

func TestNew(t *testing.T) {
	values := []keysched.Value[string, string]{
		{KID: "1", Key: "one"},
		{KID: "2", Key: "two"},
	}

	rr := New(values)
	if rr.n != len(values) {
		t.Errorf("n should be equal to the length of values")
	}

	if rr.p.Load() != 0 {
		t.Errorf("p should be 0")
	}

	// mutate the values
	values[0].Key = "ONE"
	if rr.values[0].Key != "one" {
		t.Errorf("should not mutate the internal values")
	}

	// append a new value
	values = append(values, keysched.Value[string, string]{KID: "3", Key: "three"})
	if len(rr.values) != 2 {
		t.Errorf("should not mutate the internal values")
	}
}

func TestScheduler_Next(t *testing.T) {
	t.Run("keys not set", func(t *testing.T) {
		rr := New([]keysched.Value[string, string]{})
		_, err := rr.Next(context.Background())
		if !errors.Is(err, keysched.ErrKeyNotSet) {
			t.Errorf("want error: %v, got error: %v", err, keysched.ErrKeyNotSet)
		}
	})

	t.Run("round robin algorithm", func(t *testing.T) {
		values := []keysched.Value[string, string]{
			{KID: "1", Key: "one"},
			{KID: "2", Key: "two"},
			{KID: "3", Key: "three"},
		}

		schedules := []keysched.Value[string, string]{
			{KID: "1", Key: "one"},
			{KID: "2", Key: "two"},
			{KID: "3", Key: "three"},
			{KID: "1", Key: "one"},
			{KID: "2", Key: "two"},
		}

		rr := New(values)

		for i := range schedules {
			val, err := rr.Next(context.Background())
			if err != nil {
				t.Errorf("expect no error; got %v", err)
			}

			if val != schedules[i] {
				t.Errorf("expect scheduled value is: %+v; got: %+v", schedules[i], val)
			}
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		values := []keysched.Value[string, string]{
			{KID: "1", Key: "one"},
			{KID: "2", Key: "two"},
			{KID: "3", Key: "three"},
		}

		rr := New(values)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // immediately cancel.

		_, err := rr.Next(ctx)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("want error: %v, got error: %v", err, context.Canceled)
		}
	})
}

func TestScheduler_Find(t *testing.T) {
	t.Run("key not found", func(t *testing.T) {
		rr := New([]keysched.Value[string, string]{})
		_, err := rr.Find(context.Background(), "1")
		if !errors.Is(err, keysched.ErrKeyNotFound) {
			t.Errorf("want error: %v, got error: %v", err, keysched.ErrKeyNotSet)
		}
	})

	t.Run("key found", func(t *testing.T) {
		values := []keysched.Value[string, string]{
			{KID: "1", Key: "one"},
			{KID: "2", Key: "two"},
			{KID: "3", Key: "three"},
		}

		rr := New(values)
		val, err := rr.Find(context.Background(), "2")
		if err != nil {
			t.Errorf("expect no error; got %v", err)
		}

		if val != values[1].Key {
			t.Errorf("expect scheduled value is: %+v; got: %+v", values[1], val)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		values := []keysched.Value[string, string]{
			{KID: "1", Key: "one"},
			{KID: "2", Key: "two"},
			{KID: "3", Key: "three"},
		}

		rr := New(values)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // immediately cancel.

		_, err := rr.Find(ctx, "3")
		if !errors.Is(err, context.Canceled) {
			t.Errorf("want error: %v, got error: %v", err, context.Canceled)
		}
	})
}
