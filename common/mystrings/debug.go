package mystrings

import "fmt"

const debug bool = true

func dprint(format string, a ...interface{}) (int, error) {
	if debug {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}
