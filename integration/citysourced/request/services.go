package request

import (
	"encoding/json"
	"fmt"
)

// ================================================================================================
//                                      SERVICES
// ================================================================================================

// Service is the RPC container struct for the Services service.  This service
// providers a directory of services (i.e. report categories) available for each
// CitySourced city.
type Service struct{}

// Run mashals and sends the Service request to the proper back-end, and returns
// the response in Native format.
func (c *Service) Run(rqst *NServiceRequest, resp *NServices) error {
	fmt.Printf("resp: %p\n", resp)
	fmt.Println(rqst)

	// irqst, err := c.makeI(rqst)
	// r, err := irqst.Process()
	// r.makeN(resp)
	fmt.Printf("  --> resp: %p\n%s\n", resp, *resp)
	return nil
}

// ------------------------------- SID -------------------------------

// UnmarshalJSON implements the conversion from the JSON "ID" to the ServiceID struct.
func (srv *NService) UnmarshalJSON(value []byte) error {
	type T struct {
		ID         int
		Name       string
		Categories []string `json:"catg"`
	}
	var t T
	err := json.Unmarshal(value, &t)
	if err != nil {
		return err
	}
	srv.ID = t.ID
	srv.Name = t.Name
	srv.Categories = t.Categories
	return nil
}
