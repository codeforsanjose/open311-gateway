package mystrings

import (
	"fmt"
	"regexp"
)

// MyRegexp supports regex substring matches, both with named and unnamed groups.
type MyRegexp struct {
	rx    *regexp.Regexp
	mg    *MacGyver
	Ok    bool
	Named map[string]string
	All   []string
}

// NewRegex returns a regex matcher.
func NewRegex(exp, trimset, delset string) *MyRegexp {
	rx := regexp.MustCompile(exp)
	mg, e := NewMacGyver(Trim(trimset), DeleteChars(delset))
	if e != nil {
		return nil
	}
	return &MyRegexp{
		rx:    rx,
		mg:    mg,
		Ok:    false,
		Named: nil,
		All:   nil}
}

// Match executes a match.
func (r *MyRegexp) Match(s string) error {
	r.Named, r.All = nil, nil

	matches := r.rx.FindStringSubmatch(s)
	if matches == nil {
		return fmt.Errorf("no matches found")
	}
	r.Ok = true

	for i, m := range matches {
		matches[i] = r.mg.Process(m)
		// dprint("Replace on: %q -> %q\n", m, matches[i])
	}

	results := make(map[string]string)
	for i, name := range r.rx.SubexpNames() {
		// dprint("%2d  %-12v ", i, fmt.Sprintf("%s:", name))
		if i > 0 && name > "" {
			results[name] = matches[i]
			// dprint("%v\n", results[name])
		} else {
			// dprint("~~\n")
		}
	}
	// dprint("\n")
	r.Named, r.All = results, matches
	return nil
}
