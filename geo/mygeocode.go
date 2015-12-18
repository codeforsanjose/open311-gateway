package geo

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

// GetLatLng queries Google for the geolocation of an address.  It returns the lat, lng, and
// an error.
func GetLatLng(addr string) (float64, float64, error) {
	req := &Request{
		Address:  addr,
		Provider: GOOGLE,
	}
	resp, err := req.Lookup(nil)
	if err != nil || resp.Status != "OK" {
		return 0.0, 0.0, fmt.Errorf("Unable to determine GeoLoc for %q", addr)
	}
	p := resp.GoogleResponse.Results[0].Geometry.Location
	return p.Lat, p.Lng, nil
}

// GetAddress queries Google for the geolocation of the input latitude and longitude.
// It returns the full address and an error.
func GetAddress(lat, lng float64) (string, error) {
	loc := Point{lat, lng}
	req := &Request{
		Location: &loc,
		Provider: GOOGLE,
	}
	resp, err := req.Lookup(nil)
	fmt.Printf(">>>Found: %s\n", resp.Found)
	fmt.Printf(">>>City: %s\n", getCity(resp))
	fmt.Printf(">>>Response:\n%s\n", spew.Sdump(resp))
	fmt.Printf("+++Response:\n%#v\n", resp.GoogleResponse.Results[0].AddressParts)
	fmt.Println("---------------------------- Address Parts -------------------")
	for i, v := range resp.GoogleResponse.Results[0].AddressParts {
		fmt.Printf("%t %2d %#v\n", contains(v.Types, "political") && contains(v.Types, "locality"), i, v)

	}
	if err != nil || resp.Status != "OK" {
		return "", fmt.Errorf("Unable to determine GeoLoc for %v | %v", lat, lng)
	}
	return resp.Found, nil
}

func getCity(resp *Response) string {
	for _, v := range resp.GoogleResponse.Results[0].AddressParts {
		if contains(v.Types, "political") && contains(v.Types, "locality") {
			return v.Name
		}
	}
	return ""
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
