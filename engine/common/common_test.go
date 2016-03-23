package common

import (
	"fmt"
	"testing"

	// "github.com/davecgh/go-spew/spew"
)

type addressTest struct {
}

func TestAddress(t *testing.T) {

	fmt.Printf("\n\n======================================================== TestAddress ========================================================\n\n")

	{
		fmt.Printf("\nRunning LOOSE zip parse\n\n")

		var tests = []struct {
			n  string // input
			r  []int
			ok bool // expected result
		}{
			{"322 E Santa Clara St, San Jose, CA, 95112", []int{20, 41, 22, 30, 32, 34, 36, 41}, true},
			{"1322 E Santa Clara St, San Jose, CA", []int{21, 35, 23, 31, 33, 35, -1, -1}, true},
			{"322 E Santa Clara St, Apt 1, San Jose, CA 95112 2222", []int{27, 52, 29, 37, 39, 41, 42, 52}, true},
			{"322 E Santa Clara St, Apt 1, San Jose, CA", []int{27, 41, 29, 37, 39, 41, -1, -1}, true},
			{"322 E Santa Clara St, San Jose, CA, 95112-1111", []int{20, 46, 22, 30, 32, 34, 36, 46}, true},
			{"322 E Santa Clara St, San Jose, CA, 95112-111", []int(nil), false},
			{"322 E Santa Clara St, San Jose, CA, 95112-11111", []int(nil), false},
		}

		for _, nt := range tests {
			fmt.Println(nt.n)
			addr, city, state, zip, err := ParseAddress(nt.n, true)
			if (err != nil && nt.ok == true) || (err == nil && nt.ok == false) {
				fmt.Printf("   ERROR - Addr: %q  City: %q  State: %q  Zip: %q\n        err: %s\n", addr, city, state, zip, err)
			} else {
				fmt.Printf("   PASS  - Addr: %q  City: %q  State: %q  Zip: %q\n        err: %s\n", addr, city, state, zip, err)
			}

			var ind []int
			ind, err = matchCityStateZip(nt.n, true)
			fmt.Printf("     Ind: %#v  err: %s\n\n", ind, err)
			if !testEq(ind, nt.r) {
				t.Errorf("[LOOSE] Index match (%#v <-> %#v) failed for: %s", ind, nt.r, nt.n)
			}
		}
	}

	{
		fmt.Printf("\n\nRunning TIGHT zip parse\n\n")

		var tests = []struct {
			n  string // input
			r  []int
			ok bool // expected result
		}{
			{"322 E Santa Clara St, San Jose, CA, 95112", []int{20, 41, 22, 30, 32, 34, 36, 41}, true},
			{"1322 E Santa Clara St, San Jose, CA", []int(nil), false},
			{"322 E Santa Clara St, Apt 1, San Jose, CA 95112 2222", []int{27, 52, 29, 37, 39, 41, 42, 52}, true},
			{"322 E Santa Clara St, Apt 1, San Jose, CA", []int(nil), false},
			{"322 E Santa Clara St, San Jose, CA, 95112-1111", []int{20, 46, 22, 30, 32, 34, 36, 46}, true},
			{"322 E Santa Clara St, San Jose, CA, 95112-111", []int(nil), false},
			{"322 E Santa Clara St, San Jose, CA, 95112-11111", []int(nil), false},
		}

		for _, nt := range tests {
			fmt.Println(nt.n)
			addr, city, state, zip, err := ParseAddress(nt.n, false)
			if (err != nil && nt.ok == true) || (err == nil && nt.ok == false) {
				fmt.Printf("   ERROR - Addr: %q  City: %q  State: %q  Zip: %q\n        err: %s\n", addr, city, state, zip, err)
			} else {
				fmt.Printf("   PASS  - Addr: %q  City: %q  State: %q  Zip: %q\n        err: %s\n", addr, city, state, zip, err)
			}

			var ind []int
			ind, err = matchCityStateZip(nt.n, false)
			fmt.Printf("     Ind: %#v  err: %s\n\n", ind, err)
			if !testEq(ind, nt.r) {
				t.Errorf("[TIGHT] Index match (%#v <-> %#v) failed for: %s", ind, nt.r, nt.n)
			}
		}
	}
}

func testEq(a, b []int) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
