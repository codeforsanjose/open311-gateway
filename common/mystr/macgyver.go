package mystr

import (
	"fmt"
	"strings"
)

// MacGyver is the Swiss Army Knife of string conditioning.  It iterates through
//  a slice of functions, applying each in order to the input string.
type MacGyver struct {
	ops []func(string) string
}

// NewMacGyver returns a pointer to a new MacGyver instance.  Input is a variadic
// list of string processing functions.
func NewMacGyver(fs ...func(string) string) (*MacGyver, error) {
	if len(fs) == 0 {
		return nil, fmt.Errorf("cannot create new MacGyver - no functions specified")
	}
	var sak MacGyver
	sak.ops = append(sak.ops, fs...)
	return &sak, nil
}

// Process sequentially runs the list of string operations.
func (m *MacGyver) Process(input string) (output string) {
	output = input
	if input == "" {
		return output
	}
	for _, f := range m.ops {
		output = f(output)
	}
	return output
}

// --------------------- CLOSURE FUNCTIONS ---------------------------

// Trim removes all of the characters in cutset from the front and
// tail of the input string.
func Trim(cutset string) func(string) string {
	return func(input string) (output string) {
		output = input
		output = strings.Trim(output, cutset)
		return
	}
}

// TrimLeft returns a function that removes all of the characters in cutset
// from the front of the input string.
func TrimLeft(cutset string) func(string) string {
	return func(input string) string {
		return strings.TrimLeft(input, cutset)
	}
}

// TrimRight returns a function that removes all of the characters in
// cutset from the tail of the input string.
func TrimRight(cutset string) func(string) string {
	return func(input string) string {
		return strings.TrimRight(input, cutset)
	}
}

// TrimPrefix returns a function that removes the prefix string from the
// front of the input string.
func TrimPrefix(prefix string) func(string) string {
	return func(input string) string {
		return strings.TrimPrefix(input, prefix)
	}
}

// TrimSuffix returns a function that removes the suffix string from the
// tail of the input string.
func TrimSuffix(suffix string) func(string) string {
	return func(input string) string {
		return strings.TrimSuffix(input, suffix)
	}
}

// ReplaceOne returns a function that returns a copy of the input string with
// the first n instances of old replaced by new.  NOTE: old and new are
// strings, not characters.  In other words, this function replaces instances
// of the full string old with the full string new.
func ReplaceOne(old, new string, n int) func(string) string {
	return func(input string) string {
		return strings.Replace(input, old, new, n)
	}
}

// DeleteChars returns a function that returns a copy of the input string after
// removing each of the characters in delset.  Example: removing the dash from a
// zipcode match: "99999-1234" -> "999991234".
func DeleteChars(delset string) func(string) string {
	var (
		replacer *strings.Replacer
		rl       []string
	)
	if delset > "" {
		for _, d := range delset {
			rl = append(rl, string(d), "")
		}
		replacer = strings.NewReplacer(rl...)
	}

	return func(input string) string {
		if input == "" || replacer == nil {
			return input
		}
		return replacer.Replace(input)
	}
}

// Replace returns a function that returns a copy of the input string, with
// the specified strings replaced.  replaceset is a list of replacement pairs.
// In other words, any instances of replaceset[0] in the input string are
// replaced by replaceset[1].  replaceset[2] -> replaceset[3], etc.
func Replace(replaceset []string) func(string) string {
	var (
		replacer *strings.Replacer
	)
	if replaceset != nil && (len(replaceset)%2 == 0) {
		replacer = strings.NewReplacer(replaceset...)
	}

	return func(input string) string {
		if input == "" || replacer == nil {
			dprint("skipping\n")
			return input
		}
		return replacer.Replace(input)
	}
}

// --------------------- SIMPLE FUNCTIONS ---------------------------

// TrimSpace removes whitespace from both ends of the input string.
func TrimSpace(input string) string {
	return strings.TrimSpace(input)
}

// Upper returns the input converted to all upper case.
func Upper(input string) string {
	return strings.ToUpper(input)
}

// Lower returns the input converted to all lower case.
func Lower(input string) string {
	return strings.ToLower(input)
}

// Title returns the input converted to title case (the beginning of each
// word is capitalized).
func Title(input string) string {
	return strings.Title(input)
}
