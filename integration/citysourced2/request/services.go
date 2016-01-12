package request

import (
	"Gateway311/integration/citysourced/data"
	"Gateway311/integration/citysourced/structs"
	"_sketches/spew"
	"fmt"
)

// ================================================================================================
//                                      SERVICES
// ================================================================================================

// Service is the RPC container struct for the Services service.  This service
// providers a directory of services (i.e. report categories) available for each
// CitySourced city.
type Service struct{}

// ServicesForCity fills resp with a list of services for the specifed city.
func (c *Service) ServicesForCity(rqst *structs.NServiceRequest, resp *structs.NServicesResponse) error {
	fmt.Println(rqst)

	x, err := data.ServicesForCity(rqst.City)
	if err == nil {
		fmt.Printf("  resp: %p", resp)
		resp.Message = "OK"
		resp.Services = *x
		fmt.Printf(" --> %p\n", resp)
		fmt.Printf("      %s\n", spew.Sdump(resp))
	} else {
		fmt.Printf("[ServicesForCity]: error: %s\n", err)
	}
	return err
}
