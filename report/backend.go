package report

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"

	"github.com/jmcvetta/napping"
)

// BECreate communicates with the specified backend to create a new Report.
func BECreate(jid int64, report *CreateReport) (*CreateReportResp, error) {
	fmt.Printf("[BECreate] jid: %d\n", jid)
	response := CreateReportResp{
		Message:  "Failed",
		ID:       "",
		AuthorID: "",
	}
	xmlPayload, e := prepOutCS(report)
	if e != nil {
		return &response, e
	}

	payload := bytes.NewBuffer(xmlPayload)
	result := CSReportResp{}
	var reply bytes.Buffer
	err := errors.New("Unknown error from Backend.")
	fmt.Printf("   payload type: %T size: %d\n", payload, payload.Len())
	request := napping.Request{
		Url:                 "http://localhost:5050/api/",
		Method:              "POST",
		Payload:             payload,
		RawPayload:          true,
		Result:              &result,
		CaptureResponseBody: true,
		ResponseBody:        &reply,
		Error:               &err,
	}

	resp, err := napping.Send(&request)
	if err != nil {
		panic(err)
	}
	if resp.Status() == 200 {
		fmt.Printf("  SUCCESS - response: %#v\n", resp)
	}
	fmt.Printf("  response body: %v\n", reply.String())

	return &response, nil
}

// ==============================================================================================================================
//                                      CreateReportResp
// ==============================================================================================================================
func prepOutCS(src *CreateReport) ([]byte, error) {
	requestTypeID, err := strconv.ParseInt(src.TypeID, 10, 64)
	if err != nil {
		fmt.Printf("Unable to parse request type id: %q\n", src.TypeID)
		return nil, fmt.Errorf("Unable to parse request type id: %q", src.TypeID)
	}
	r := CSReport{
		APIAuthKey:        "a01234567890z",
		APIRequestType:    "CreateThreeOneOne",
		APIRequestVersion: "1",
		DeviceType:        src.DeviceType,
		DeviceModel:       src.DeviceModel,
		DeviceID:          src.DeviceID,
		RequestType:       src.Type,
		RequestTypeID:     requestTypeID,
		Latitude:          src.Latitude,
		Longitude:         src.Longitude,
		Description:       src.Description,
		AuthorNameFirst:   src.FirstName,
		AuthorNameLast:    src.LastName,
		AuthorEmail:       src.Email,
		AuthorTelephone:   src.Phone,
		AuthorIsAnonymous: src.IsAnonymous,
	}
	fmt.Printf("  CitySourced payload: %v\n", r)
	output, err := xml.MarshalIndent(r, "  ", "    ")
	s := string(output[:])
	fmt.Printf("  Marshaled: %s\n", s)

	return output, nil
}
