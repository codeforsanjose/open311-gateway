package structs

import (
	"encoding/json"
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
	// log.Debug("[UnmarshalJSON] parts: %+v\n", parts)
	s.AdpID = parts[0]
	s.AreaID = parts[1]
	s.ProviderID = cnvInt(parts[2])
	s.ID = cnvInt(parts[3])
	// log.Debug("[UnmarshalJSON] AdpID: %#v  AreaID: %#v  ProviderID: %#v  ID: %#v\n", s.AdpID, s.AreaID, s.ProviderID, s.ID)
	return nil
}

// MarshalJSON implements the conversion from the ServiceID struct to the JSON "ID".
func (s ServiceID) MarshalJSON() ([]byte, error) {
	// fmt.Printf("  Marshaling s: %#v\n", s)
	return []byte(fmt.Sprintf("\"%s\"", s.MID())), nil
}

// ==============================================================================================================================
//                                      Report ID
// ==============================================================================================================================

// UnmarshalJSON implements the conversion from the JSON "ID" to the ReportID struct.
func (s *ReportID) UnmarshalJSON(value []byte) error {
	cnvInt := func(x string) int {
		y, _ := strconv.ParseInt(x, 10, 64)
		return int(y)
	}
	parts := strings.Split(strings.Trim(string(value), "\" "), "-")
	// log.Debug("[UnmarshalJSON] parts: %+v\n", parts)
	s.AdpID = parts[0]
	s.AreaID = parts[1]
	s.ProviderID = cnvInt(parts[2])
	s.ID = parts[3]
	// log.Debug("[UnmarshalJSON] AdpID: %#v  AreaID: %#v  ProviderID: %#v  ID: %#v\n", s.AdpID, s.AreaID, s.ProviderID, s.ID)
	return nil
}

// MarshalJSON implements the conversion from the ReportID struct to the JSON "ID".
func (s ReportID) MarshalJSON() ([]byte, error) {
	// fmt.Printf("  Marshaling s: %#v\n", s)
	return []byte(fmt.Sprintf("\"%s\"", s.RID())), nil
}

// ==============================================================================================================================
//                                      NService
// ==============================================================================================================================

// UnmarshalJSON implements the conversion from the JSON "ID" to the ServiceID struct.
func (srv *NService) UnmarshalJSON(value []byte) error {
	type T struct {
		ID            int
		Name          string
		Description   string
		Metadata      bool
		Group         string
		Keywords      []string
		ServiceNotice string `json:"service_notice"`
	}
	var t T
	err := json.Unmarshal(value, &t)
	if err != nil {
		return err
	}
	srv.ID = t.ID
	srv.Name = t.Name
	srv.Description = t.Description
	srv.Metadata = t.Metadata
	srv.ServiceNotice = t.ServiceNotice
	srv.Keywords = t.Keywords
	srv.Group = t.Group
	return nil
}
