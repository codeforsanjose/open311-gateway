package geo

import "fmt"

// LatLngForAddr queries Google for the geolocation of an address.  It returns the lat, lng, and
// an error.
func LatLngForAddr(addr string) (float64, float64, error) {
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

// AddrForLatLng queries Google for the geolocation of the input LatitudeV and LongitudeV.
// It returns the full address and an error.
func AddrForLatLng(lat, lng float64) (string, error) {
	loc := Point{lat, lng}
	req := &Request{
		Location: &loc,
		Provider: GOOGLE,
	}
	resp, err := req.Lookup(nil)
	fmt.Printf(">>>Found: %s\n", resp.Found)
	// fmt.Printf(">>>Response:\n%s\n", spew.Sdump(resp))
	// fmt.Printf("+++Response:\n%#v\n", resp.GoogleResponse.Results[0].AddressParts)
	// fmt.Println("---------------------------- Address Parts -------------------")
	// for i, v := range resp.GoogleResponse.Results[0].AddressParts {
	// 	fmt.Printf("%t %2d %#v\n", contains(v.Types, "political") && contains(v.Types, "locality"), i, v)
	//
	// }
	if err != nil || resp.Status != "OK" {
		return "", fmt.Errorf("Unable to determine GeoLoc for %v | %v", lat, lng)
	}
	return resp.Found, nil
}

// AddrForLatLng queries Google for the geolocation of the input LatitudeV and LongitudeV.
// It returns the full address and an error.
func AddrDetailForLatLng(lat, lng float64) (string, error) {
	loc := Point{lat, lng}
	req := &Request{
		Location: &loc,
		Provider: GOOGLE,
	}
	resp, err := req.Lookup(nil)
	fmt.Printf(">>>Found: %s\n", resp.Found)
	// fmt.Printf(">>>Response:\n%s\n", spew.Sdump(resp))
	// fmt.Printf("+++Response:\n%#v\n", resp.GoogleResponse.Results[0].AddressParts)
	// fmt.Println("---------------------------- Address Parts -------------------")
	// for i, v := range resp.GoogleResponse.Results[0].AddressParts {
	// 	fmt.Printf("%t %2d %#v\n", contains(v.Types, "political") && contains(v.Types, "locality"), i, v)
	//
	// }
	if err != nil || resp.Status != "OK" {
		return "", fmt.Errorf("Unable to determine GeoLoc for %v | %v", lat, lng)
	}
	return resp.Found, nil
}

// CityForLatLng queries Google for the geolocation of the input LatitudeV and LongitudeV.
// It returns the city, and
func CityForLatLng(lat, lng float64) (string, error) {
	loc := Point{lat, lng}
	req := &Request{
		Location: &loc,
		Provider: GOOGLE,
	}
	resp, _ := req.Lookup(nil)
	for _, v := range resp.GoogleResponse.Results[0].AddressParts {
		if contains(v.Types, "political") && contains(v.Types, "locality") {
			return v.Name, nil
		}
	}
	return "", fmt.Errorf("Unable to find the city for this location")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
