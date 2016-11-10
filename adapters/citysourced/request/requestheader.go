package request

import (
	"github.com/codeforsanjose/open311-gateway/common"
	"github.com/codeforsanjose/open311-gateway/common/jx"
)

var stdHeader Header

// Header encapsulates all of the CitySourced API and App information that
// is common to all requests.
type Header struct {
	jx.JSONConfig
	APIAuthKey        string `xml:"ApiAuthKey,omitempty"`
	APIRequestVersion string `xml:"ApiRequestVersion,omitempty"`
	AppKey            string `xml:"AppKey,omitempty"`
	AppType           string `xml:"AppType,omitempty"`
	AppVersion        string `xml:"AppVersion,omitempty"`
	Locale            string `xml:"Locale,omitempty"`
}

// String returns a string representation of a Report.
func (r Header) String() string {
	ls := new(common.FmtBoxer)
	ls.AddS("Header\n")
	ls.AddF("API AuthKey: **********  Version: %s\n", r.APIRequestVersion)
	ls.AddF("App Key: %s   Type: %s   Ver: %s\n", r.AppKey, r.AppType, r.AppVersion)
	ls.AddF("Locale: %s\n", r.Locale)
	return ls.BoxC(80)
}
