package x

import (
	"context"
	"errors"
	"iter"
	"slices"

	"golang.org/x/sync/errgroup"
)

// IterWrapper provides operations over an iterator sequence.
type IterWrapper[T any] struct {
	seq iter.Seq[T]
}

// Iter returns an iterator over v.
func Iter[T any](v []T) IterWrapper[T] {
	return IterWrapper[T]{slices.Values(v)}
}

// IterSeq returns a wrapper around v.
func IterSeq[T any](v iter.Seq[T]) IterWrapper[T] {
	return IterWrapper[T]{v}
}

// Map transforms each value in the iterator.
func (it IterWrapper[T]) Map[R any](f func(T) R) IterWrapper[R] {
	return IterWrapper[R]{
		func(yield func(R) bool) {
			for v := range it.seq {
				if !yield(f(v)) {
					return
				}
			}
		},
	}
}

// MapIndex transforms each value in the iterator with its index.
func (it IterWrapper[T]) MapIndex[R any](f func(T, int) R) IterWrapper[R] {
	return IterWrapper[R]{
		func(yield func(R) bool) {
			i := -1
			for v := range it.seq {
				i += 1
				if !yield(f(v, i)) {
					return
				}
			}
		},
	}
}

// Filter returns the values that satisfy cond.
func (it IterWrapper[T]) Filter(cond func(T) bool) IterWrapper[T] {
	return IterWrapper[T]{
		func(yield func(T) bool) {
			for v := range it.seq {
				if !cond(v) {
					continue
				}

				if !yield(v) {
					return
				}
			}
		},
	}
}

// FilterIndex returns the values for which cond, given the value and index, is true.
func (it IterWrapper[T]) FilterIndex(cond func(T, int) bool) IterWrapper[T] {
	return IterWrapper[T]{
		func(yield func(T) bool) {
			i := -1
			for v := range it.seq {
				i += 1
				if !cond(v, i) {
					continue
				}

				if !yield(v) {
					return
				}
			}
		},
	}
}

// Fold reduces the values in the iterator into an accumulator.
func (it IterWrapper[T]) Fold[R any](r R, f func(R, T) R) R {
	for v := range it.seq {
		r = f(r, v)
	}

	return r
}

// FoldIndex reduces the values in the iterator into an accumulator with each index.
func (it IterWrapper[T]) FoldIndex[R any](r R, f func(R, T, int) R) R {
	i := -1
	for v := range it.seq {
		i += 1
		r = f(r, v, i)
	}

	return r
}

// Some reports whether any value satisfies cond. It returns false for an empty iterator.
func (it IterWrapper[T]) Some(cond func(T) bool) bool {
	for v := range it.seq {
		if cond(v) {
			return true
		}
	}

	return false
}

// None reports whether no value satisfies cond. It returns true for an empty iterator.
func (it IterWrapper[T]) None(cond func(T) bool) bool {
	return !it.Some(cond)
}

// Every reports whether all values satisfy cond. It returns true for an empty iterator.
func (it IterWrapper[T]) Every(cond func(T) bool) bool {
	for v := range it.seq {
		if !cond(v) {
			return false
		}
	}

	return true
}

// Collect returns the iterator values as a slice.
func (it IterWrapper[T]) Collect() []T {
	return slices.Collect(it.seq)
}

// Go runs f concurrently for each value. Results and errors are in input order; an error does not stop other calls.
func (it IterWrapper[T]) Go[R any](f func(val T) (R, error)) ([]R, []error) {
	input := it.Collect()
	results := make([]R, len(input))
	errors := make([]error, len(input))

	g := new(errgroup.Group)

	for i, val := range input {
		g.Go(func() error {
			res, err := f(val)

			if err != nil {
				errors[i] = err
			} else {
				results[i] = res
			}

			return nil
		})
	}

	g.Wait()

	return results, errors
}

// GoCtx runs f concurrently for each value with ctx and returns the first error.
func (it IterWrapper[T]) GoCtx[R any](ctx context.Context, f func(ctx context.Context, val T) (R, error)) ([]R, error) {
	return it.GoCtxWithLimit(ctx, -1, f)
}

// GoCtxWithLimit runs f concurrently with at most limit calls at once using ctx.
// A negative limit is unlimited; zero is invalid. It returns the first error.
func (it IterWrapper[T]) GoCtxWithLimit[R any](ctx context.Context, limit int, f func(ctx context.Context, val T) (R, error)) ([]R, error) {
	if limit == 0 {
		return nil, errors.New("limit must not be 0")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	input := it.Collect()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	results := make([]R, len(input))

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(limit)

	for i, val := range input {
		g.Go(func() error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			res, err := f(ctx, val)

			if err != nil {
				return err
			}

			results[i] = res

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}
