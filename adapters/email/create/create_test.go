package create

import (
	"fmt"
	"testing"

	"Gateway311/adapters/email/data"
	"Gateway311/adapters/email/logs"
)

var Debug = true

func init() {
	logs.Init(Debug)

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

	c := Request{
		To:                "jameskhaskell@gmail.com",
		From:              "Code for San Jose",
		Subject:           "New Report",
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

	err := c.SendEmail([]string{"jameskhaskell@gmail.com"})
	if err != nil {
		t.Errorf(err.Error())
	}

}
