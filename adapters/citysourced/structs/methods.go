package structs

import "encoding/json"

// ------------------------------- JSON -------------------------------

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
