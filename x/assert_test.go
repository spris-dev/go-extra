package x_test

import (
	"testing"

	"go.spris.dev/x"
)

func TestAssertPanics(t *testing.T) {
	const message = "something is wrong"

	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("expected a panic, but got none")
		}

		gotMessage, ok := got.(string)
		if !ok || gotMessage != message {
			t.Errorf("expected panic message %q, got %v", message, got)
		}
	}()

	x.Assert(false, message)
}

func TestAssertDoesNotPanic(t *testing.T) {
	defer func() {
		if got := recover(); got != nil {
			t.Errorf("unexpected panic: %v", got)
		}
	}()

	x.Assert(true, "something is wrong")
}
