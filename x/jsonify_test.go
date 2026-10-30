package x_test

import (
	"testing"

	"go.spris.dev/tassert"
	"go.spris.dev/x"
)

type jsonPerson struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestJsonify(t *testing.T) {
	t.Run("returns indented JSON", func(t *testing.T) {
		got := x.Jsonify(jsonPerson{Name: "Art Vandelay", Age: 42})
		want := "{\n \"name\": \"Art Vandelay\",\n \"age\": 42\n}"

		tassert.EmptyDiff(t, want, got)
	})

	t.Run("panics when value cannot be encoded", func(t *testing.T) {
		defer func() {
			if got := recover(); got == nil {
				t.Error("expected Jsonify to panic")
			}
		}()

		x.Jsonify(make(chan int))
	})
}
