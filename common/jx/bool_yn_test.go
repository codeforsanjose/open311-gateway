package jx_test

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/codeforsanjose/open311-gateway/_background/go/common"
	. "github.com/codeforsanjose/open311-gateway/_background/go/common/jx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("BoolYNType", func() {
	Describe("XML Unmarshal", func() {
		Describe("Valid Input", func() {
			DescribeTable("value:",
				func(test BYNTest) {
					var ts BYNxml
					Ω(ts.loadXML(test.encoded)).Should(Succeed())
					Ω(ts.BYN).Should(Equal(test.x.BYN))
				},
				Entry("Yes", BYNTest{
					encoded: `<?xml version="1.0" encoding="UTF-8" ?><test><belem>Yes</belem></test>`,
					x:       BYNxml{BYN: true},
				}),
				Entry("No", BYNTest{
					encoded: `<?xml version="1.0" encoding="UTF-8" ?><test><belem>No</belem></test>`,
					x:       BYNxml{BYN: false},
				}),
				Entry("xxx", BYNTest{
					encoded: `<?xml version="1.0" encoding="UTF-8" ?><test><belem>xxx</belem></test>`,
					x:       BYNxml{BYN: false},
				}),
			)
		})

	})

	Describe("XML Marshal", func() {
		Describe("Valid Input", func() {
			DescribeTable("value:",
				func(test BYNTest) {
					out, e := xml.Marshal(test.x)
					Ω(e).ShouldNot(HaveOccurred())
					sout := common.ByteToString(out, 0)
					Ω(sout).Should(Equal(test.encoded))
				},
				Entry("true-true", BYNTest{
					encoded: `<test battr="Yes"><belem>Yes</belem></test>`,
					x:       BYNxml{BYNattr: true, BYN: true},
				}),
				Entry("true-false", BYNTest{
					encoded: `<test battr="Yes"><belem>No</belem></test>`,
					x:       BYNxml{BYNattr: true, BYN: false},
				}),
				Entry("false-true", BYNTest{
					encoded: `<test battr="No"><belem>Yes</belem></test>`,
					x:       BYNxml{BYNattr: false, BYN: true},
				}),
				Entry("false-false", BYNTest{
					encoded: `<test battr="No"><belem>No</belem></test>`,
					x:       BYNxml{BYNattr: false, BYN: false},
				}),
			)
		})

	})

	Describe("JSON Unmarshal", func() {
		Describe("Valid Input", func() {
			DescribeTable("value:",
				func(test BYNTest) {
					var ts BYNjson
					Ω(ts.loadJSON(test.encoded)).Should(Succeed())
					Ω(ts.BYN).Should(Equal(test.j.BYN))
				},
				Entry("Yes", BYNTest{
					encoded: `{"belem": "Yes"}`,
					// encoded: `{"test": {"belem": "Yes"}}`,
					j: BYNjson{BYN: true},
				}),
				Entry("No", BYNTest{
					encoded: `{"belem": "No"}`,
					// encoded: `{"test": {"belem": "Yes"}}`,
					j: BYNjson{BYN: false},
				}),
				Entry("xxx", BYNTest{
					encoded: `{"belem": "xxx"}`,
					// encoded: `{"test": {"belem": "Yes"}}`,
					j: BYNjson{BYN: false},
				}),
			)
		})

	})

	Describe("JSON Marshal", func() {
		Describe("Valid Input", func() {
			DescribeTable("value:",
				func(test BYNTest) {
					out, e := json.Marshal(test.j)
					Ω(e).ShouldNot(HaveOccurred())
					sout := common.ByteToString(out, 0)
					Ω(sout).Should(Equal(test.encoded))
				},
				Entry("true", BYNTest{
					encoded: `{"belem":"Yes"}`,
					j:       BYNjson{BYN: true},
				}),
				Entry("false", BYNTest{
					encoded: `{"belem":"No"}`,
					j:       BYNjson{BYN: false},
				}),
			)
		})

	})

})

// BYNxml is for testing unmarshalling of jx.BoolYNType.
type BYNxml struct {
	XMLName xml.Name   `xml:"test"`
	BYNattr BoolYNType `xml:"battr,attr"`
	BYN     BoolYNType `xml:"belem"`
}

func (r *BYNxml) loadXML(input string) error {
	err := xml.Unmarshal([]byte(input), r)
	if err != nil {
		msg := fmt.Sprintf("Unable to unmarshal XML data input - error: %v", err)
		return errors.New(msg)
	}
	return nil
}

// BYNxml is for testing unmarshalling of jx.BoolYNType.
type BYNjson struct {
	BYN BoolYNType `json:"belem"`
}

func (r *BYNjson) loadJSON(input string) error {
	err := json.Unmarshal([]byte(input), r)
	if err != nil {
		msg := fmt.Sprintf("Unable to unmarshal JSON data input - error: %v", err)
		return errors.New(msg)
	}
	return nil
}

type BYNTest struct {
	encoded string
	x       BYNxml
	j       BYNjson
}
