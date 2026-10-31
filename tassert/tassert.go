package tassert

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equal reports an error when want and got are not equal.
func Equal[T comparable](t testing.TB, want, got T) {
	t.Helper()

	if want != got {
		t.Errorf("values are not equal: got %v, want %v", got, want)
	}
}

// EmptyDiff reports an error when want and got differ.
func EmptyDiff(t testing.TB, want, got any, opts ...cmp.Option) {
	t.Helper()

	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Errorf("Diff mismatch (-want +got):\n%s", diff)
	}
}

// True reports an error when cond is false.
func True(t testing.TB, cond bool, format string, args ...any) {
	t.Helper()

	if !cond {
		t.Errorf(format, args...)
	}
}
