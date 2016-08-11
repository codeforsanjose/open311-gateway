package create

import (
	"fmt"
	"testing"

	"github.com/open311-gateway/adapters/email/data"
	"github.com/open311-gateway/adapters/email/structs"

	"github.com/davecgh/go-spew/spew"
)

var Debug = true

func init() {

	fmt.Println("Reading config...")
	if err := data.Init("../data/config.json"); err != nil {
		fmt.Printf("Init() failed: %s", err)
	}
}

type testResultS struct {
	input string
	isOK  bool
}

func isOK(e error) bool {
	if e == nil {
		return false
	}
	return true
}

func TestCreateEmail(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestCreateEmail] =============================\n\n")

	var tests = []struct {
		n    structs.NRoute // input
		isOK bool           // expected result
	}{
		{structs.NRoute{"EM1", "CU", 1}, true},
		{structs.NRoute{"EM1", "CU", 2}, true},
		{structs.NRoute{"EM1", "SUN", 1}, true},
	}

	c := Request{
		RequestType:       "Graffiti",
		RequestTypeID:     10,
		ImageURL:          "",
		Latitude:          100.0,
		Longitude:         -200.0,
		Description:       "There's a bunch of Bricks in the Wall over here",
		AuthorNameFirst:   "James",
		AuthorNameLast:    "Haskell",
		AuthorEmail:       "jameskhaskell@gmail.com",
		AuthorTelephone:   "4084084008",
		AuthorIsAnonymous: false,
	}

	for _, tr := range tests {
		prov, err := data.RouteProvider(tr.n)
		if err != nil {
			t.Errorf("RouteProvider failed - %s", err)
		}
		c.Sender = prov.Email
		resp, err := c.Process()
		if err != nil {
			t.Errorf("Process failed - %s", err)
		}
		fmt.Print(spew.Sdump(resp))
	}
}
