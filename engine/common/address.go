package common

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	rxCityStateZip      *regexp.Regexp
	rxCityStateLooseZip *regexp.Regexp
)

const (
	rxCityStateZipPattern      = ", *([A-Za-z .]{2,}), *([A-Z][A-Z])[, ]*(\\d{5}(?:[-\\s]\\d{4})?)$"
	rxCityStateLooseZipPattern = ", *([A-Za-z .]{2,}), *([A-Z][A-Z])[, ]*(\\d{5}(?:[-\\s]\\d{4})?)?$"
)

func init() {
	if r, err := regexp.Compile(rxCityStateZipPattern); err != nil {
		fmt.Printf("Error compiling rxCityStateZip - %s!!", err.Error())
	} else {
		rxCityStateZip = r
	}
	if r, err := regexp.Compile(rxCityStateLooseZipPattern); err != nil {
		fmt.Printf("Error compiling rxCityStateZip - %s!!", err.Error())
	} else {
		rxCityStateLooseZip = r
	}
}

// ParseAddress accepts a comma delimited address string, and returns the
// address, city, state and zip.  NOTE: this function does very little error checking.
func ParseAddress(fullAddr string, looseZip bool) (addr, city, state, zip string, err error) {
	fail := func(status string) (string, string, string, string, error) {
		return "", "", "", "", fmt.Errorf("%s - %q", status, fullAddr)
	}
	ind, err := matchCityStateZip(fullAddr, looseZip)
	if err != nil {
		return fail(err.Error())
	}
	if len(ind) < 6 {
		return fail("invalid address")
	}
	addr = fullAddr[:ind[0]]
	city = fullAddr[ind[2]:ind[3]]
	state = fullAddr[ind[4]:ind[5]]
	if len(ind) > 6 && ind[6] > 0 {
		zip = strings.Replace(fullAddr[ind[6]:ind[7]], " ", "-", 1)
	} else {
		zip = ""
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
