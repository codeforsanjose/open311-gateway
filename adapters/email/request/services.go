package request

import (
	"github.com/open311-gateway/adapters/email/data"
	"github.com/open311-gateway/adapters/email/structs"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

// ================================================================================================
//                                      SERVICES
// ================================================================================================

// Services is the RPC container struct for the Services service.  This service
// providers a directory of services (i.e. report categories) available for each
// AreaSourced city.
type Services struct{}

// Area returns a list of services for the specifed city.
func (c *Services) Area(rqst *structs.NServiceRequest, resp *structs.NServicesResponse) error {
	fmt.Println(rqst)

	x, err := data.ServicesArea(rqst.Area)
	if err == nil {
		fmt.Printf("  resp: %p", resp)
		resp.SetIDF(rqst.GetID)
		resp.AdpID = data.AdapterName()
		resp.Message = "OK"
		resp.Services = *x
		fmt.Printf("%s\n", spew.Sdump(resp))
	} else {
		fmt.Printf("[Area]: error: %s\n", err)
	}
	return err
}

// All fills resp with a list of services for the specifed city.
func (c *Services) All(rqst *structs.NServiceRequest, resp *structs.NServicesResponse) error {
	fmt.Println(rqst)

	x, err := data.ServicesAll()
	if err == nil {
		fmt.Printf("  resp: %p", resp)
		resp.SetIDF(rqst.GetID)
		resp.AdpID = data.AdapterName()
		resp.Message = "OK"
		resp.Services = *x
		// time.Sleep(time.Second * 4)
		fmt.Printf("%s\n", spew.Sdump(resp))
	} else {
		fmt.Printf("[All]: error: %s\n", err)
	}
	return err
}
