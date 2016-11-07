package geo

import (
	"fmt"

	"github.com/codeforsanjose/open311-gateway/_background/go/common"
	"github.com/codeforsanjose/open311-gateway/_background/go/common/mystrings"
)

var (
	rxCityStateZip, rxCityStateLooseZip, rxGoogleFound, rxZip *mystrings.MyRegexp
	validZipLen                                               map[int]bool
)

// Address represents a Google Address search.  RawAddr contains the requested
// address.  Ok will be set to true if the address is found.  Found is the
// found full address.
type Address struct {
	RawAddr      string
	Found        string
	streetNumber string
	route        string
	subpremise   string
	Addr         string
	City         string
	County       string
	State        string
	Zip          string
	zipSuffix    string
	Lat          float64
	Lng          float64
	Ok           bool
	Errors       []string
}

// FullAddr returns the address, city, state and zip as one string.
func (a *Address) FullAddr() string {
	return fmt.Sprintf("%s, %s, %s %s", a.Addr, a.City, a.State, a.Zip)
}

// String displays the contents of the GooAddr type.
func (a *Address) String() string {
	ls := new(common.FmtBoxer)
	ls.AddF("GooAddr\n")
	ls.AddF("Ok: %t  RawAddr: %q\n", a.Ok, a.RawAddr)
	ls.AddF("Found: %q\n", a.Found)
	ls.AddF("Street Number: %q   Route: %q  Subpremise: %q\n", a.streetNumber, a.route, a.subpremise)
	ls.AddF("Zip - %q   Suffix: %q\n", a.Zip, a.zipSuffix)
	ls.AddF("   %s\n", a.Addr)
	ls.AddF("   %s, %s %s\n", a.City, a.State, a.Zip)
	ls.AddF("   %s\n", a.County)
	ls.AddF("   %v : %v\n", a.Lat, a.Lng)
	ls.AddF("Errors: %+v\n", a.Errors)
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

// NewAddr returns a new address struct from a full, comma delimited street address.
func NewAddr(fullAddr string) (addr *Address, err error) {
	return GooAddr(fullAddr)
}

// NewAddrP returns a new address struct from the address, city, state, and option zipcode.
func NewAddrP(streetAddr, city, state, zip string) (addr *Address, err error) {
	return GooAddr(fmt.Sprintf("%s, %s, %s %s", streetAddr, city, state, zip))
}

// NewAddrLL returns a new address struct from a latitude and longitude.
func NewAddrLL(lat, lng float64) (addr *Address, err error) {
	return GooAddr(lat, lng)
}

// AddrForLatLng finds the nearest address to the coordinates.
func AddrForLatLng(lat, lng float64) (addr *Address, err error) {
	if !ValidateLatLng(lat, lng) {
		return nil, fmt.Errorf("coordinates must be in the continental US")
	}
	addr, err = GooAddr(lat, lng)
	return
}

// ParseAddress accepts a comma delimited address string, and returns the
// address, city, state and zip.  NOTE: this function does very little error checking.
func ParseAddress(fullAddr string, looseZip bool) (a Address, err error) {
	fail := func(status string) (Address, error) {
		return a, fmt.Errorf("Failed to parse address %q - %s", fullAddr, status)
	}

	var (
		e error
		d map[string]string
	)
	if looseZip {
		e = rxCityStateLooseZip.Match(fullAddr)
		d = rxCityStateLooseZip.Named
	} else {
		e = rxCityStateZip.Match(fullAddr)
		d = rxCityStateZip.Named
	}
	if e != nil {
		return fail(e.Error())
	}

	a.Addr = d["addr"]
	a.City = d["city"]
	a.State = d["state"]
	if z, ok := d["zip"]; ok {
		if len(z) >= 5 {
			e := rxZip.Match(z)
			if e == nil && rxZip.Ok {
				a.Zip = rxZip.Named["zip"]
			}
		}
	}

	if a.Zip == "" && !looseZip {
		return fail("invalid zip code")
	}
	return
}

func init() {
	rxCityStateZip = mystrings.NewRegex(`(?i), *(?P<city>[A-Za-z .]{2,}), *(?P<state>[A-Z][A-Z])[, ]*(?P<zip>\d{5}(?:[-\s]\d{4})?)$`, "", "")
	rxCityStateLooseZip = mystrings.NewRegex(`(?i), *(?P<city>[A-Za-z .]{2,}), *(?P<state>[A-Z][A-Z])[, ]*(?P<zip>\d{5}(?:[-\s]\d{4})?)?$`, "", "")
	rxGoogleFound = mystrings.NewRegex(`(?i)(?P<addr>.*), (?P<city>[A-Za-z .]{2,}), (?P<state>[A-Z][A-Z]) (?P<zip>\d{5}), USA$`, "", "")
	rxZip = mystrings.NewRegex(`^(?P<zip>\d{5}(?:[-\s]*\d{4})?)$`, "", "- ")

	validZipLen = map[int]bool{
		0:  true,
		5:  true,
		9:  true,
		10: true,
	}
}
