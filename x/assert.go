package x

import "fmt"

// AssertFunc panics when cond is false.
func AssertFunc(cond bool, message func() string) {
	if !cond {
		panic(fmt.Sprintf("assertion failed: %v", message()))
	}
}

// Assert panics when cond is false.
func Assert(cond bool, message string) {
	if !cond {
		panic(message)
	}
}
