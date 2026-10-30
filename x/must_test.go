package x_test

import (
	"strconv"
	"testing"

	"go.spris.dev/tassert"
	"go.spris.dev/x"
)

func TestMust(t *testing.T) {
	t.Run("returns value when there is no error", func(t *testing.T) {
		got := x.Must(strconv.Atoi("123"))

		tassert.Equal(t, 123, got)
	})

	t.Run("panics when there is an error", func(t *testing.T) {
		defer func() {
			if got := recover(); got == nil {
				t.Error("expected Must to panic")
			}
		}()

		x.Must(strconv.Atoi("not a number"))
	})
}

func TestMustInt(t *testing.T) {
	t.Run("returns parsed integer", func(t *testing.T) {
		got := x.MustInt("123")

		tassert.Equal(t, 123, got)
	})

	t.Run("panics when parsing fails", func(t *testing.T) {
		defer func() {
			if got := recover(); got == nil {
				t.Error("expected MustInt to panic")
			}
		}()

		x.MustInt("not a number")
	})
}
