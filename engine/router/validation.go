package router

import "github.com/open311-gateway/engine/structs"

// "github.com/davecgh/go-spew/spew"

// ValidateRID verifies the ReportID is routable.
func ValidateRID(rid structs.ReportID) bool {
	_, adpOK := adapters.getAdapter(rid.AdpID)
	_, areaOK := adapters.getAreaAdapters(rid.AreaID)
	// log.Debug("adpOK: %s  areaOK: %s", adpOK, areaOK)
	if adpOK != nil || areaOK != nil {
		return false
	}
	return true
}
