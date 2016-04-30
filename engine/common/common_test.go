package common

import (
	"fmt"
	"testing"

	// "github.com/davecgh/go-spew/spew"
)

type addressTest struct {
}

func TestParse(t *testing.T) {

	fmt.Printf("\n\n======================================================== TestParse ========================================================\n\n")

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
			a, err := ParseAddress(nt.n, true)
			if (err != nil && nt.ok == true) || (err == nil && nt.ok == false) {
				fmt.Printf("   ERROR - Addr: %q  City: %q  State: %q  Zip: %q\n        err: %s\n", a.Addr, a.City, a.State, a.Zip, err)
			} else {
				fmt.Printf("   PASS  - Addr: %q  City: %q  State: %q  Zip: %q\n        err: %s\n", a.Addr, a.City, a.State, a.Zip, err)
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
			a, err := ParseAddress(nt.n, false)
			if (err != nil && nt.ok == true) || (err == nil && nt.ok == false) {
				fmt.Printf("   ERROR - Addr: %q  City: %q  State: %q  Zip: %q\n        err: %s\n", a.Addr, a.City, a.State, a.Zip, err)
			} else {
				fmt.Printf("   PASS  - Addr: %q  City: %q  State: %q  Zip: %q\n        err: %s\n", a.Addr, a.City, a.State, a.Zip, err)
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

func TestZip(t *testing.T) {

	fmt.Printf("\n\n======================================================== TestZip ========================================================\n\n")

	var tests = []struct {
		n  string // input
		ok bool   // expected result
	}{
		{"95037", true},
		{"99999", true},
		{"00000", true},
		{"12345", true},
		{"1", false},
		{"12", false},
		{"123", false},
		{"1234", false},
		{"123456", false},
		{"1234567", false},
		{"123456789", true},
		{"99999-9999", true},
		{"11111-1111", true},
		{"99999 9999", true},
		{"11111 1111", true},

		{"9999-99999", true},
		{"1111-11111", true},
		{"9999 99999", true},
		{"111 111111", true},

		{"99999-99991", false},
		{"11111-11111", false},
		{"99999-999", false},
		{"11111-111", false},
		{"99999-99", false},
		{"11111-11", false},
		{"99999 999", false},
		{"11111 111", false},
		{"99999 99", false},
		{"11111 11", false},
	}

	for i, nt := range tests {
		fmt.Printf("%2d %-20s  %t\n", i, nt.n, rxZip.MatchString(nt.n))
	}
}

type TestLatLongData struct {
	city  string
	state string
	lat   float64
	lng   float64
}

func TestLatLong(t *testing.T) {

	fmt.Printf("\n\n======================================================== TestLatLong ========================================================\n\n")

	// See notes in Quiver RE: "Data / US Cities"

	{

		var tests = []struct {
			n  TestLatLongData // input
			ok bool            // expected result
		}{
			{TestLatLongData{"Anchorage", "AK", 61.2180556, -149.9002778}, false},
			{TestLatLongData{"Honolulu", "HI", 21.306944399999999, -157.8583333}, false},
			{TestLatLongData{"San Jose", "CA", 37.339385700000001, -121.89495549999999}, true},
			{TestLatLongData{"San Francisco", "CA", 37.773972, -122.431297}, true},
			{TestLatLongData{"New York", "NY", 40.714269100000003, -74.005972900000003}, true},
			{TestLatLongData{"Phoenix", "AZ", 33.448377100000002, -112.0740373}, true},
			{TestLatLongData{"Denver", "CO", 39.739153600000002, -104.9847034}, true},
			{TestLatLongData{"Atlanta", "GA", 33.748995399999998, -84.387982399999999}, true},
			{TestLatLongData{"Chicago", "IL", 41.850033000000003, -87.650052299999999}, true},
		}

		for _, nt := range tests {
			a, e := AddrForLatLng(nt.n.lat, nt.n.lng)
			fmt.Printf("%+v  %s\n", a, e)
			if (nt.ok && e != nil) || (!nt.ok && e == nil) {
				t.Errorf("Failed for: %s - error: %v", nt.n.city, e)
			}
		}
	}
}

func TestAddress(t *testing.T) {

	fmt.Printf("\n\n======================================================== TestAddress ========================================================\n\n")

	{
		fmt.Println("--------------------------- NewAddr --------------------------")
		var tests = []struct {
			n  string // input
			r  []float64
			ok bool // expected result
		}{
			{"322 E Santa Clara St, San Jose, CA, 95112", []float64{37.3391629, -121.8836029}, true},
			{"1322 E Santa Clara St, San Jose, CA", []float64{37.3484843, -121.8645008}, true},
			{"322 E Santa Clara St, Apt 1, San Jose, CA 95112 2222", []float64{37.3391629, -121.8836029}, true},
			{"322 E Santa Clara St, Apt 1, San Jose, CA", []float64{37.3391629, -121.8836029}, true},
			{"322 E Santa Clara St, San Jose, CA, 95112-1111", []float64{37.3391629, -121.8836029}, true},
			{"322 E Santa Clara St, San Jose, CA, 95112-111", []float64{0.0, 0.0}, false},
			{"322 E Santa Clara St, San Jose, CA, 95112-11111", []float64{0.0, 0.0}, false},
			{"10630 S De Anza Blvd, Cupertino, CA 95014", []float64{0.0, 0.0}, true},
		}

		for i, nt := range tests {
			a, err := NewAddr(nt.n, true)
			if ((err != nil) && nt.ok) || ((err == nil) && !nt.ok) {
				t.Errorf("Failed for: %v - error: %v", i+1, err)
			}
			if a.Lat != nt.r[0] || a.Long != nt.r[1] {
				t.Errorf("Lat/lng failed for: %v - error: %v", i+1, err)
			}
			fmt.Printf("%2d %-5t for: %v\n%+v\n\n", i+1, a.Valid, nt.n, a)
		}
	}

	{
		fmt.Printf("\n--------------------------- NewAddrP --------------------------\n\n")
		var tests = []struct {
			addr, city, state, zip string // input
			r                      []float64
			ok                     bool // expected result
		}{
			{"322 E Santa Clara St", "San Jose", "CA", "95112", []float64{37.3391629, -121.8836029}, true},
			{"1322 E Santa Clara St", "San Jose", "CA", "", []float64{37.3484843, -121.8645008}, true},
			{"322 E Santa Clara St, Apt 1", "San Jose", "CA", "95112 2222", []float64{37.3391629, -121.8836029}, true},
			{"322 E Santa Clara St, Apt 1", "San Jose", "CA", "", []float64{37.3391629, -121.8836029}, true},
			{"322 E Santa Clara St", "San Jose", "CA", "95112-1111", []float64{37.3391629, -121.8836029}, true},
			{"322 E Santa Clara St", "San Jose", "CA", "95112-111", []float64{0.0, 0.0}, false},
			{"322 E Santa Clara St", "San Jose", "CA", "95112-11111", []float64{0.0, 0.0}, false},
		}

		for i, nt := range tests {
			a, err := NewAddrP(nt.addr, nt.city, nt.state, nt.zip, true)
			if ((err != nil) && nt.ok) || ((err == nil) && !nt.ok) {
				t.Errorf("Failed for: %v - error: %v", i+1, err)
			}
			if a.Lat != nt.r[0] || a.Long != nt.r[1] {
				t.Errorf("Lat/lng failed for: %v - error: %v", i+1, err)
			}
			fmt.Printf("%2d %-5t for: %s, %s, %s %s\n%+v\n\n", i+1, a.Valid, nt.addr, nt.city, nt.state, nt.zip, a)
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
