package structs

import (
	"fmt"
	"strconv"
	"strings"
)

// =======================================================================================
//                                      RID
// =======================================================================================

const emInvalidRid = "Invalid RID: %q"

// ------------------------------- ReportID -------------------------------

// NewRID creates a ReportID by concatenating a Route (NRoute string) with a message
// response ID.
func NewRID(route NRoute, id string) ReportID {
	return ReportID{
		NRoute: route,
		ID:     id,
	}
}

// ReportID adds routing information to a ReportID returned by a call to a
// Service Provider.
type ReportID struct {
	NRoute
	ID string
}

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

// RIDFromString converts a reportID string to a new ReportID struct.
func RIDFromString(rids string) (ReportID, NRoute, error) {
	if rids == "" {
		return ReportID{}, NRoute{}, fmt.Errorf("empty RID: %q", rids)
	}
	adpID, areaID, providerID, reportID, err := SplitRID(rids)
	if err != nil {
		return ReportID{}, NRoute{}, fmt.Errorf(emInvalidRid, rids)
	}
	nr := NRoute{
		AdpID:      adpID,
		AreaID:     areaID,
		ProviderID: providerID,
	}

	return ReportID{
		NRoute: nr,
		ID:     fmt.Sprintf("%v", reportID),
	}, nr, nil
}

// SplitRID breaks down an RID, and returns all subfields.
func SplitRID(rid string) (string, string, int, int, error) {
	adpID, areaID, provID, id, err := SplitRMID(rid)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf(emInvalidRid, rid)
	}
	return adpID, areaID, provID, id, nil
}

// RidAdpID breaks down a RID, and returns the AdpID.
func RidAdpID(rid string) (string, error) {
	adpID, _, _, _, err := SplitRMID(rid)
	if err != nil {
		return "", fmt.Errorf(emInvalidRid, rid)
	}
	return adpID, nil
}

// RidAreaID breaks down a RID, and returns the AreaID.
func RidAreaID(rid string) (string, error) {
	_, areaID, _, _, err := SplitRMID(rid)
	if err != nil {
		return "", fmt.Errorf(emInvalidRid, rid)
	}
	return areaID, nil
}

// RidProviderID breaks down an RID, and returns the ProviderID.
func RidProviderID(rid string) (int, error) {
	_, _, provID, _, err := SplitRMID(rid)
	if err != nil {
		return 0, fmt.Errorf(emInvalidRid, rid)
	}
	return provID, nil
}

// RidID breaks down an RID, and returns the Report ID.
func RidID(rid string) (int, error) {
	_, _, _, id, err := SplitRMID(rid)
	if err != nil {
		return 0, fmt.Errorf(emInvalidRid, rid)
	}
	return id, nil
}

// =======================================================================================
//                                      STRINGS
// =======================================================================================

// RID creates the Master ID string for the Service.
func (r ReportID) RID() string {
	if r.AdpID == "" && r.AreaID == "" && r.ProviderID == 0 && r.ID == "" {
		return ""
	}
	return fmt.Sprintf("%s-%s-%d-%s", r.AdpID, r.AreaID, r.ProviderID, r.ID)
}

// Display the string represenation of a ReportID.
func (r ReportID) String() string {
	return fmt.Sprintf("%s-%s", r.NRoute, r.ID)
}
