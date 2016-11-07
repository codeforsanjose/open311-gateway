package geo

const (
	// Continental US only!
	latMin float64 = 18.0
	latMax float64 = 49.0
	lngMin float64 = -124.6
	lngMax float64 = -62.3
)

// ValidateLatLng validates a lat/lng pair against the Continental US boundaries.
func ValidateLatLng(lat, lng float64) bool {
	if lat >= latMin && lat <= latMax && lng >= lngMin && lng <= lngMax {
		return true
	}
	return false
}
