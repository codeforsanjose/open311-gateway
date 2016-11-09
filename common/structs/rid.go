package structs

import "fmt"

// =======================================================================================
//                                      RID
// =======================================================================================

const emInvalidRid = "Invalid RID: %q"

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

// NRouteFromString converts a reportID string to a new NRoute struct.
func NRouteFromString(rids string) (NRoute, error) {
	adpID, areaID, providerID, _, err := SplitRID(rids)
	if err != nil {
		return NRoute{}, fmt.Errorf("invalid RID: %q", rids)
	}
	return NRoute{
		AdpID:      adpID,
		AreaID:     areaID,
		ProviderID: providerID,
	}, nil

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
