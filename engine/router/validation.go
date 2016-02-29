package router

import "Gateway311/engine/structs"

// "github.com/davecgh/go-spew/spew"

// ValidateRID verifies the ReportID is routable.
func ValidateRID(rid structs.ReportID) bool {
	_, adpOK := adapters.getAdapter(rid.AdpID)
	_, areaOK := adapters.getAreaAdapters(rid.AreaID)
	log.Debug("adpOK: %t  areaOK: %t")
	if adpOK != nil || areaOK != nil {
		return false
	}
	return true
}
