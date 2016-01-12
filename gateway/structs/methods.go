package structs

import (
	"fmt"
	"strconv"
	"strings"
)

// ==============================================================================================================================
//                                      Service ID
// ==============================================================================================================================

// UnmarshalJSON implements the conversion from the JSON "ID" to the ServiceID struct.
func (s *ServiceID) UnmarshalJSON(value []byte) error {
	cnvInt := func(x string) int {
		y, _ := strconv.ParseInt(x, 10, 64)
		return int(y)
	}
	parts := strings.Split(strings.Trim(string(value), "\" "), "-")
	s.IFID = parts[0]
	s.AreaID = parts[1]
	s.ProviderID = cnvInt(parts[2])
	s.ID = cnvInt(parts[3])
	return nil
}

// MarshalJSON implements the conversion from the ServiceID struct to the JSON "ID".
func (s ServiceID) MarshalJSON() ([]byte, error) {
	fmt.Printf("  Marshaling s: %#v\n", s)
	return []byte(fmt.Sprintf("\"%s\"", s.MID())), nil
}
