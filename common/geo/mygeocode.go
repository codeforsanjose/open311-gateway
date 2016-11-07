package geo

import "fmt"

const googleRespOK = "OK"

// GooAddr runs the Google Geocode API on the input address, and retuns a pointer
// to an Address.
func GooAddr(input ...interface{}) (*Address, error) {
	var req *Request

	ga := Address{
		Ok:     false,
		Errors: make([]string, 0, 0),
	}

	// Create the request
	switch inp := input[0].(type) {
	case string:
		ga.RawAddr = inp
		req = &Request{
			Address:  inp,
			Provider: GOOGLE,
		}
	case float64, float32:
		loc := Point{input[0].(float64), input[1].(float64)}
		req = &Request{
			Location: &loc,
			Provider: GOOGLE,
		}
	default:
		return nil, fmt.Errorf("invalid input type")
	}

	// Send the request to Google
	resp, err := req.Lookup(nil)
	if err != nil || resp.Status != googleRespOK {
		return nil, fmt.Errorf("Unable to determine GeoLoc for %q", input)
	}

	// Parse the Google results back into Address
	ga.Found = resp.Found
	err = ga.unpackGResponse(resp.GoogleResponse)
	// fmt.Printf("%s\n", ga.String())
	return &ga, nil
}

func (a *Address) unpackGResponse(resp *GoogleResponse) error {
	loadAddr := func(ir int) error {
		for _, ap := range resp.Results[ir].AddressParts {
			for _, aptype := range ap.Types {
				switch aptype {
				case "street_number":
					a.streetNumber = ap.Name
				case "route":
					a.route = ap.ShortName
				case "subpremise":
					a.subpremise = ap.Name
				case "locality":
					a.City = ap.Name
				case "administrative_area_level_2":
					a.County = ap.ShortName
				case "administrative_area_level_1":
					a.State = ap.ShortName
				case "postal_code":
					a.Zip = ap.Name
				case "postal_code_suffix":
					a.zipSuffix = ap.Name
				}
			}
		}
		a.Lat = resp.Results[ir].Geometry.Location.Lat
		a.Lng = resp.Results[ir].Geometry.Location.Lng

		a.Addr = a.streetNumber + " " + a.route
		if a.subpremise > "" {
			a.Addr += " #" + a.subpremise
		}
		return nil
	}

	// Parse the Google results back into Address
	var err error
	for ir, result := range resp.Results {
		for _, apart := range result.AddressParts {
			for _, t := range apart.Types {
				if t == "street_number" {
					err = loadAddr(ir)
					break
				}
			}
		}
	}

	if err != nil {
		a.Ok = true
		return err
	}
	return nil
}

// GooLatLngForAddr queries Google for the geolocation of an address.  It returns the lat, lng, and
// an error.
func GooLatLngForAddr(addr string) (float64, float64, error) {
	ga, e := GooAddr(addr)
	if e != nil {
		return 0, 0, e
	}
	return ga.Lat, ga.Lng, nil
}

// GooAddrForLatLng queries Google for the geolocation of the input Latitude and Longitude.
// It returns the full address and an error.
func GooAddrForLatLng(lat, lng float64) (string, error) {
	ga, e := GooAddr(lat, lng)
	if e != nil {
		return "", e
	}
	return ga.FullAddr(), nil
}

// GooCityForLatLng queries Google for the geolocation of the input Latitude and Longitude.
// It returns the city.
func GooCityForLatLng(lat, lng float64) (string, error) {
	ga, e := GooAddr(lat, lng)
	if e != nil {
		return "", fmt.Errorf("unable to find the city for this location - %s", e)
	}
	return ga.City, nil
}
