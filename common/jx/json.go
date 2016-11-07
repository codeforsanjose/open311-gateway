package jx

import (
	"bytes"
	"encoding/json"
)

// EncodeJSON encodes the input struct as JSON, and returns the result as a string.
func EncodeJSON(source interface{}) (output string, err error) {
	var boutput = new(bytes.Buffer)
	enc := json.NewEncoder(boutput)
	enc.SetIndent(" ", "   ")
	if err := enc.Encode(source); err != nil {
		return "", err
	}

	return boutput.String(), nil
}
