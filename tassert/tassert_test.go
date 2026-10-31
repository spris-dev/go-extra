package tassert_test

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"go.spris.dev/tassert"
)

func TestEqual(t *testing.T) {
	tassert.Equal(t, 123, 123)
}

func TestEmptyDiff(t *testing.T) {
	tassert.EmptyDiff(t, []int{1, 2, 3}, []int{1, 2, 3})
	tassert.EmptyDiff(t, []int{}, []int(nil), cmpopts.EquateEmpty())
}

func TestTrue(t *testing.T) {
	tassert.True(t, true, "expected condition to be true")
}
