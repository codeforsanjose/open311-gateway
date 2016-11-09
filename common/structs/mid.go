package structs

import (
	"fmt"
	"strconv"
	"strings"
)

// =======================================================================================
//                                      MID
// =======================================================================================

// SplitRMID breaks down an MID or RID, and returns all subfields.
func SplitRMID(mid string) (string, string, int, int, error) {
	fail := func() (string, string, int, int, error) {
		return "", "", 0, 0, fmt.Errorf("Invalid RMID: %q", mid)
	}
	parts := strings.Split(mid, "-")
	if len(parts) != 4 {
		fail()
	}
	pid, err := strconv.Atoi(parts[2])
	if err != nil {
		fail()
	}
	id, err := strconv.Atoi(parts[3])
	if err != nil {
		fail()
	}
	return parts[0], parts[1], pid, id, nil
}

const emInvalidMid = "Invalid MID: %q"

// SplitMID breaks down an MID, and returns all subfields.
func SplitMID(mid string) (string, string, int, int, error) {
	adpID, areaID, provID, id, err := SplitRMID(mid)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf(emInvalidMid, mid)
	}
	return adpID, areaID, provID, id, nil
}

// MidAdpID breaks down a MID, and returns the AdpID.
func MidAdpID(mid string) (string, error) {
	adpID, _, _, _, err := SplitRMID(mid)
	if err != nil {
		return "", fmt.Errorf(emInvalidMid, mid)
	}
	return adpID, nil
}

// MidAreaID breaks down a MID, and returns the AreaID.
func MidAreaID(mid string) (string, error) {
	_, areaID, _, _, err := SplitRMID(mid)
	if err != nil {
		return "", fmt.Errorf(emInvalidMid, mid)
	}
	return areaID, nil
}

// MidProviderID breaks down an MID, and returns the ProviderID.
func MidProviderID(mid string) (int, error) {
	_, _, provID, _, err := SplitRMID(mid)
	if err != nil {
		return 0, fmt.Errorf(emInvalidMid, mid)
	}
	return provID, nil
}

// MidID breaks down an MID, and returns the Service ID.
func MidID(mid string) (int, error) {
	_, _, _, id, err := SplitRMID(mid)
	if err != nil {
		return 0, fmt.Errorf(emInvalidMid, mid)
	}
	return id, nil
}

// =======================================================================================
//                                      STRINGS
// =======================================================================================

// MID creates the Master ID string for the Service.
func (r ServiceID) MID() string {
	if r.AdpID == "" && r.AreaID == "" && r.ProviderID == 0 && r.ID == 0 {
		return ""
	}
	return fmt.Sprintf("%s-%s-%d-%d", r.AdpID, r.AreaID, r.ProviderID, r.ID)
}
