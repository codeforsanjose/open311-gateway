package common

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/open311-gateway/engine/geo"
)

var (
	rxCityStateZip      *regexp.Regexp
	rxCityStateLooseZip *regexp.Regexp
	rxGoogleFound       *regexp.Regexp
	rxZip               *regexp.Regexp
)

var validZipLen map[int]bool

const (
	rxCityStateZipPattern      = `, *([A-Za-z .]{2,}), *([A-Z][A-Z])[, ]*(\d{5}(?:[-\s]\d{4})?)$`
	rxCityStateLooseZipPattern = `, *([A-Za-z .]{2,}), *([A-Z][A-Z])[, ]*(\d{5}(?:[-\s]\d{4})?)?$`
	rxGoogleAddressPattern     = `(.*), ([A-Za-z .]{2,}), ([A-Z][A-Z]) (\d{5}), USA$`
	rxZipPattern               = `^\d{5}(?:[-\s]\d{4})?$`
)

// Address represents a standard US street address.
type Address struct {
	Addr      string
	City      string
	State     string
	Zip       string
	Lat, Long float64
	Valid     bool
}

// String displays the contents of the CreateRequest type.
func (r Address) String() string {
	ls := new(LogString)
	ls.AddF("Address\n")
	ls.AddF("Valid: %t\n", r.Valid)
	ls.AddF("   %s\n", r.Addr)
	ls.AddF("   %s, %s  %s\n", r.City, r.State, r.Zip)
	ls.AddF("Lat: %v  Long: %v\n", r.Lat, r.Long)
	return ls.Box(60)
}

// Validate does simple range checking to make sure all address fields are populated.
// If zip is false, then the zipcode will not be included in the validation.
func (r *Address) Validate(addr, zip bool) bool {
	if (addr && len(r.Addr) == 0) || (len(r.City) == 0) || (len(r.State) != 2) || (zip && validZipLen[len(r.Zip)] == true) {
		return false
	}
	return true
}

// FullAddr returns the standard comma delimited full addres string, like:
// "{{addr}}, {{city}}, {{state}} {{zip}}"
func (r *Address) FullAddr() string {
	return fmt.Sprintf("%s, %s, %s %s", r.Addr, r.City, r.State, r.Zip)
}

// NewAddr returns a new address struct from a full, comma delimited street address.
func NewAddr(fullAddr string, looseZip bool) (addr Address, err error) {
	if len(fullAddr) == 0 {
		err = fmt.Errorf("full address is empty")
		return
	}

	addr, err = ParseAddress(fullAddr, true)
	if err != nil {
		return
	}
	fmt.Printf("NewAddr:\n%s", addr.String())

	addr.Lat, addr.Long, err = geo.LatLngForAddr(addr.FullAddr())
	if err == nil {
		addr.Valid = true
	}
	return
}

// NewAddrP returns a new address struct from the address, city, state, and option zipcode.
func NewAddrP(addr1, city, state, zip string, looseZip bool) (addr Address, err error) {
	addr.Addr = addr1
	addr.City = city
	addr.State = state
	addr.Zip = zip

	if !addr.Validate(true, false) {
		err = fmt.Errorf("invalid street address")
		return
	}
	if len(addr.Zip) > 0 {
		if !rxZip.MatchString(addr.Zip) {
			err = fmt.Errorf("invalid zip code")
			return
		}
	}

	addr.Lat, addr.Long, err = geo.LatLngForAddr(addr.FullAddr())
	if err == nil {
		addr.Valid = true
	}
	return
}

// ParseAddress accepts a comma delimited address string, and returns the
// address, city, state and zip.  NOTE: this function does very little error checking.
func ParseAddress(fullAddr string, looseZip bool) (a Address, err error) {
	fail := func(status string) (Address, error) {
		return a, fmt.Errorf("%s - %q", status, fullAddr)
	}
	ind, e := matchCityStateZip(fullAddr, looseZip)
	if e != nil {
		return fail(e.Error())
	}
	if len(ind) < 6 {
		return fail("invalid address")
	}
	a.Addr = fullAddr[:ind[0]]
	a.City = fullAddr[ind[2]:ind[3]]
	a.State = fullAddr[ind[4]:ind[5]]
	if len(ind) > 6 && ind[6] > 0 {
		if rxZip.MatchString(fullAddr[ind[6]:ind[7]]) {
			a.Zip = strings.Replace(fullAddr[ind[6]:ind[7]], " ", "-", 1)
		}
	}
	if a.Zip == "" && !looseZip {
		return fail("invalid zip code")
	}
	return
}

func matchCityStateZip(fullAddr string, looseZip bool) (ind []int, err error) {
	if looseZip {
		ind = rxCityStateLooseZip.FindStringSubmatchIndex(fullAddr)
	} else {
		ind = rxCityStateZip.FindStringSubmatchIndex(fullAddr)
	}
	if len(ind) == 0 {
		return ind, fmt.Errorf("invalid city, state or zip")
	}
	return
}

// AddrForLatLng finds the nearest address to the coordinates.
func AddrForLatLng(lat, lng float64) (a Address, err error) {
	fail := func(ers string, e error) (Address, error) {
		if e != nil && len(ers) > 0 {
			return a, fmt.Errorf("unable to find coordinates: %v, %v - %s - %s", lat, lng, ers, e.Error())
		} else if e != nil && len(ers) == 0 {
			return a, fmt.Errorf("unable to find coordinates: %v, %v - %s", lat, lng, e.Error())
		} else {
			return a, fmt.Errorf("unable to find coordinates: %v, %v - %s", lat, lng, ers)
		}
	}

	if !ValidateLatLng(lat, lng) {
		return fail("coordinates must be in the continental US", nil)
	}

	fullAddr, e := geo.AddrForLatLng(lat, lng)
	if e != nil {
		return fail("", e)
	}
	ind := rxGoogleFound.FindStringSubmatchIndex(fullAddr)
	if len(ind) < 10 {
		return fail("address not found", nil)
	}

	a.Addr = fullAddr[ind[2]:ind[3]]
	a.City = fullAddr[ind[4]:ind[5]]
	a.State = fullAddr[ind[6]:ind[7]]
	a.Zip = fullAddr[ind[8]:ind[9]]
	a.Lat = lat
	a.Long = lng
	a.Valid = true

	return
}

func init() {
	rxCityStateZip = regexp.MustCompile(rxCityStateZipPattern)
	rxCityStateLooseZip = regexp.MustCompile(rxCityStateLooseZipPattern)
	rxGoogleFound = regexp.MustCompile(rxGoogleAddressPattern)
	rxZip = regexp.MustCompile(rxZipPattern)

	validZipLen = map[int]bool{
		0:  true,
		5:  true,
		9:  true,
		10: true,
	}
}
