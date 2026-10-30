package x

import (
	"iter"
	"slices"
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

// Fold reduces the values in the iterator into an accumulator.
func (it IterWrapper[T]) Fold[R any](r R, f func(R, T) R) R {
	for v := range it.seq {
		r = f(r, v)
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
