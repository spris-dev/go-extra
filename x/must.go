package x

import "strconv"

// Must returns val or panics if err is non-nil.
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}

// MustInt parses val as an integer or panics on error.
func MustInt(val string) int {
	return Must(strconv.Atoi(val))
}
