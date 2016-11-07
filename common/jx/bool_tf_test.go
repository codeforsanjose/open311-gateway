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

var _ = Describe("BoolTFType", func() {
	Describe("XML Unmarshal", func() {
		Describe("Valid Input", func() {
			DescribeTable("value:",
				func(test BTFTest) {
					var ts BTFxml
					Ω(ts.loadXML(test.encoded)).Should(Succeed())
					Ω(ts.BTF).Should(Equal(test.x.BTF))
				},
				Entry("True", BTFTest{
					encoded: `<?xml version="1.0" encoding="UTF-8" ?><test><belem>True</belem></test>`,
					x:       BTFxml{BTF: true},
				}),
				Entry("False", BTFTest{
					encoded: `<?xml version="1.0" encoding="UTF-8" ?><test><belem>False</belem></test>`,
					x:       BTFxml{BTF: false},
				}),
				Entry("xxx", BTFTest{
					encoded: `<?xml version="1.0" encoding="UTF-8" ?><test><belem>xxx</belem></test>`,
					x:       BTFxml{BTF: false},
				}),
			)
		})

	})

	Describe("XML Marshal", func() {
		Describe("Valid Input", func() {
			DescribeTable("value:",
				func(test BTFTest) {
					out, e := xml.Marshal(test.x)
					Ω(e).ShouldNot(HaveOccurred())
					sout := common.ByteToString(out, 0)
					Ω(sout).Should(Equal(test.encoded))
				},
				Entry("true-true", BTFTest{
					encoded: `<test battr="True"><belem>True</belem></test>`,
					x:       BTFxml{BTFattr: true, BTF: true},
				}),
				Entry("true-false", BTFTest{
					encoded: `<test battr="True"><belem>False</belem></test>`,
					x:       BTFxml{BTFattr: true, BTF: false},
				}),
				Entry("false-true", BTFTest{
					encoded: `<test battr="False"><belem>True</belem></test>`,
					x:       BTFxml{BTFattr: false, BTF: true},
				}),
				Entry("false-false", BTFTest{
					encoded: `<test battr="False"><belem>False</belem></test>`,
					x:       BTFxml{BTFattr: false, BTF: false},
				}),
			)
		})

	})

	Describe("JSON Unmarshal", func() {
		Describe("Valid Input", func() {
			DescribeTable("value:",
				func(test BTFTest) {
					var ts BTFjson
					Ω(ts.loadJSON(test.encoded)).Should(Succeed())
					Ω(ts.BTF).Should(Equal(test.j.BTF))
				},
				Entry("True", BTFTest{
					encoded: `{"belem": "True"}`,
					// encoded: `{"test": {"belem": "True"}}`,
					j: BTFjson{BTF: true},
				}),
				Entry("False", BTFTest{
					encoded: `{"belem": "False"}`,
					// encoded: `{"test": {"belem": "True"}}`,
					j: BTFjson{BTF: false},
				}),
				Entry("xxx", BTFTest{
					encoded: `{"belem": "xxx"}`,
					// encoded: `{"test": {"belem": "True"}}`,
					j: BTFjson{BTF: false},
				}),
			)
		})

	})

	Describe("JSON Marshal", func() {
		Describe("Valid Input", func() {
			DescribeTable("value:",
				func(test BTFTest) {
					out, e := json.Marshal(test.j)
					Ω(e).ShouldNot(HaveOccurred())
					sout := common.ByteToString(out, 0)
					Ω(sout).Should(Equal(test.encoded))
				},
				Entry("true", BTFTest{
					encoded: `{"belem":"True"}`,
					j:       BTFjson{BTF: true},
				}),
				Entry("false", BTFTest{
					encoded: `{"belem":"False"}`,
					j:       BTFjson{BTF: false},
				}),
			)
		})

	})

	Describe("BoolToStringTF", func() {
		Describe("valid input", func() {
			It("true", func() {
				Ω(BoolToStringTF(true)).Should(Equal("True"))
			})

			It("false", func() {
				Ω(BoolToStringTF(false)).Should(Equal("False"))
			})

		})
	})
})

// BTFxml is for testing unmarshalling of jx.BoolTFType.
type BTFxml struct {
	XMLName xml.Name   `xml:"test"`
	BTFattr BoolTFType `xml:"battr,attr"`
	BTF     BoolTFType `xml:"belem"`
}

func (r *BTFxml) loadXML(input string) error {
	err := xml.Unmarshal([]byte(input), r)
	if err != nil {
		msg := fmt.Sprintf("Unable to unmarshal XML data input - error: %v", err)
		return errors.New(msg)
	}
	return nil
}

// BTFxml is for testing unmarshalling of jx.BoolTFType.
type BTFjson struct {
	BTF BoolTFType `json:"belem"`
}

func (r *BTFjson) loadJSON(input string) error {
	err := json.Unmarshal([]byte(input), r)
	if err != nil {
		msg := fmt.Sprintf("Unable to unmarshal JSON data input - error: %v", err)
		return errors.New(msg)
	}
	return nil
}

type BTFTest struct {
	encoded string
	x       BTFxml
	j       BTFjson
}
