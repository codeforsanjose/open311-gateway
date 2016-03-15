package request

import (
	"fmt"
	"testing"

	"Gateway311/adapters/email/data"
	"Gateway311/adapters/email/logs"
	"Gateway311/adapters/email/structs"
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

func TestCreate(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestCreate] =============================\n\n")

	rqsts := []*structs.NCreateRequest{
		&structs.NCreateRequest{
			NRequestCommon: structs.NRequestCommon{
				ID:    structs.NID{1000, 1},
				Route: structs.NRoute{"EM1", "CU", 1},
				Rtype: structs.NRTCreate,
			},
			MID:         structs.ServiceID{"EM1", "CU", 1, 10},
			Type:        "Gang Activity",
			DeviceType:  "~~~DeviceType~~~",
			DeviceModel: "~~~DeviceModel~~~",
			DeviceID:    "~~~DeviceID~~~",
			Latitude:    40.00,
			Longitude:   -100.00,
			Address:     "Address1",
			Area:        "Cupertino",
			State:       "CA",
			Zip:         "99999",
			FirstName:   "James",
			LastName:    "Haskell",
			Email:       "jameskhaskell@gmail.com",
			Phone:       "4084084008",
			IsAnonymous: false,
			Description: "There are scary guys outside!",
		},
		&structs.NCreateRequest{
			NRequestCommon: structs.NRequestCommon{
				ID:    structs.NID{2000, 2},
				Route: structs.NRoute{"EM1", "CU", 2},
				Rtype: structs.NRTCreate,
			},
			MID:         structs.ServiceID{"EM1", "CU", 2, 30},
			Type:        "Illegal Dumping / Trash",
			DeviceType:  "~~~DeviceType~~~",
			DeviceModel: "~~~DeviceModel~~~",
			DeviceID:    "~~~DeviceID~~~",
			Latitude:    40.00,
			Longitude:   -100.00,
			Address:     "Address1",
			Area:        "Cupertino",
			State:       "CA",
			Zip:         "99999",
			FirstName:   "James",
			LastName:    "Haskell",
			Email:       "jameskhaskell@gmail.com",
			Phone:       "4084084008",
			IsAnonymous: false,
			Description: "There's an old couch on my sidewalk in Cupertino!",
		},
		&structs.NCreateRequest{
			NRequestCommon: structs.NRequestCommon{
				ID:    structs.NID{3000, 3},
				Route: structs.NRoute{"EM1", "SUN", 1},
				Rtype: structs.NRTCreate,
			},
			MID:         structs.ServiceID{"EM1", "SUN", 1, 10},
			Type:        "Gang Activity",
			DeviceType:  "~~~DeviceType~~~",
			DeviceModel: "~~~DeviceModel~~~",
			DeviceID:    "~~~DeviceID~~~",
			Latitude:    40.00,
			Longitude:   -100.00,
			Address:     "Address1",
			Area:        "Cupertino",
			State:       "CA",
			Zip:         "99999",
			FirstName:   "James",
			LastName:    "Haskell",
			Email:       "jameskhaskell@gmail.com",
			Phone:       "4084084008",
			IsAnonymous: false,
			Description: "Gang bangers are everywhere!",
		},
	}

	resp := new(structs.NCreateResponse)

	rpt := new(Report)
	for _, rqst := range rqsts {
		err := rpt.Create(rqst, resp)
		if err != nil {
			t.Errorf("Create failed - %s", err)
		}
	}
}
