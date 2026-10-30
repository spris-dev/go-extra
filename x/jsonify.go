package x

import (
	"encoding/json"
)

// Jsonify returns value encoded as indented JSON and panics if encoding fails.
func Jsonify(value any) string {
	return string(Must(json.MarshalIndent(value, "", " ")))
}
