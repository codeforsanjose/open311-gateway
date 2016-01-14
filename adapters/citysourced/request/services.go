package request

import (
	"Gateway311/adapters/citysourced/data"
	"Gateway311/adapters/citysourced/structs"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

// ================================================================================================
//                                      SERVICES
// ================================================================================================

// Service is the RPC container struct for the Services service.  This service
// providers a directory of services (i.e. report categories) available for each
// CitySourced city.
type Service struct{}

// City returns a list of services for the specifed city.
func (c *Service) City(rqst *structs.NServiceRequest, resp *structs.NServicesResponse) error {
	fmt.Println(rqst)

	x, err := data.ServicesCity(rqst.City)
	if err == nil {
		fmt.Printf("  resp: %p", resp)
		resp.Message = "OK"
		resp.Services = *x
		fmt.Printf(" --> %p\n", resp)
		fmt.Printf("      %s\n", spew.Sdump(resp))
	} else {
		fmt.Printf("[City]: error: %s\n", err)
	}
	return err
}

// All fills resp with a list of services for the specifed city.
func (c *Service) All(rqst *structs.NServiceRequest, resp *structs.NServicesResponse) error {
	fmt.Println(rqst)

	x, err := data.ServicesAll()
	if err == nil {
		fmt.Printf("  resp: %p", resp)
		resp.Message = "OK"
		resp.Services = *x
		fmt.Printf(" --> %p\n", resp)
		fmt.Printf("      %s\n", spew.Sdump(resp))
	} else {
		fmt.Printf("[All]: error: %s\n", err)
	}
	return err
}
