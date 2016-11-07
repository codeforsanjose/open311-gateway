package jx_test

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"

	. "github.com/codeforsanjose/open311-gateway/_background/go/common/jx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unmarshal-Safe Numbers", func() {
	Describe("Simple Initialization", func() {
		Describe("XJFloat64", func() {
			DescribeTable("value:",
				func(val float64) {
					cval := XJFloat64(val)
					Ω(cval).Should(Equal(XJFloat64(val)))
				},
				Entry("0", 0.0),
				Entry("-1", -1.0),
				Entry("0.0001", 0.0001),
				Entry("999999999.0", 999999999.0),
				Entry("-999999999.0", -999999999.0),
			)
		})
		Describe("XJInt", func() {
			DescribeTable("value:",
				func(val int) {
					cval := XJInt(val)
					Ω(cval).Should(Equal(XJInt(val)))
				},
				Entry("0", 0),
				Entry("-1", -1),
				Entry("999999999", 999999999),
				Entry("-999999999", -999999999),
			)
		})
	})

	Describe("JSON Unmarshal", func() {
		Describe("Valid Input", func() {
			DescribeTable("value:",
				func(test JSONTest) {
					var ts TestStruct
					Ω(ts.loadJSON(test.input)).Should(Succeed())
					Ω(ts.FloatVal).Should(Equal(test.fval))
					Ω(ts.IntVal).Should(Equal(test.ival))
				},
				Entry("positive values", JSONTest{
					input: `{"test": {"floatVal": 123.456, "intVal": 567}}`,
					fval:  XJFloat64(123.456),
					ival:  XJInt(567),
				}),
				Entry("zero values", JSONTest{
					input: `{"test": {"floatVal": 0, "intVal": 0}}`,
					fval:  XJFloat64(0),
					ival:  XJInt(0),
				}),
				Entry("negative values", JSONTest{
					input: `{"test": {"floatVal": -123.456, "intVal": -567}}`,
					fval:  XJFloat64(-123.456),
					ival:  XJInt(-567),
				}),
				Entry("quoted non-zero values", JSONTest{
					input: `{"test": {"floatVal": "123.456", "intVal": "567"}}`,
					fval:  XJFloat64(123.456),
					ival:  XJInt(567),
				}),
				Entry("quoted zero values", JSONTest{
					input: `{"test": {"floatVal": "0", "intVal": "0"}}`,
					fval:  XJFloat64(0),
					ival:  XJInt(0),
				}),
				Entry("empty values", JSONTest{
					input: `{"test": {"floatVal": "", "intVal": ""}}`,
					fval:  XJFloat64(0),
					ival:  XJInt(0),
				}),
			)
		})

		Describe("Invalid Input ", func() {
			DescribeTable("value should fail:",
				func(test JSONTest) {
					var ts TestStruct
					e := ts.loadJSON(test.input)
					Ω(e).Should(HaveOccurred())
				},
				Entry("values containing characters", JSONTest{
					input: `{"test": {"floatVal": "123.456x", "intVal": "a567"}}`,
				}),
				Entry("values containing punctuation", JSONTest{
					input: `{"test": {"floatVal": "123.456.", "intVal": "567,"}}`,
				}),
			)
		})
	})

	Describe("XML Unmarshal", func() {
		Describe("Valid Input", func() {
			DescribeTable("value:",
				func(test XMLTest) {
					var ts TestStruct
					Ω(ts.loadXML(test.input)).Should(Succeed())
					Ω(ts.FloatVal).Should(Equal(test.fval))
					Ω(ts.IntVal).Should(Equal(test.ival))
				},
				Entry("positive values", XMLTest{
					input: `<?xml version="1.0" encoding="UTF-8" ?><test><floatVal>123.456</floatVal><intVal>567</intVal></test>`,
					fval:  XJFloat64(123.456),
					ival:  XJInt(567),
				}),
				Entry("negative values", XMLTest{
					input: `<?xml version="1.0" encoding="UTF-8" ?><test><floatVal>-123.456</floatVal><intVal>-567</intVal></test>`,
					fval:  XJFloat64(-123.456),
					ival:  XJInt(-567),
				}),
				Entry("zero values", XMLTest{
					input: `<?xml version="1.0" encoding="UTF-8" ?><test><floatVal>0</floatVal><intVal>0</intVal></test>`,
					fval:  XJFloat64(0),
					ival:  XJInt(0),
				}),
				Entry("quoted positive values", XMLTest{
					input: `<?xml version="1.0" encoding="UTF-8" ?><test><floatVal>"123.456"</floatVal><intVal>"567"</intVal></test>`,
					fval:  XJFloat64(123.456),
					ival:  XJInt(567),
				}),
				Entry("quoted negative values", XMLTest{
					input: `<?xml version="1.0" encoding="UTF-8" ?><test><floatVal>"-123.456"</floatVal><intVal>"-567"</intVal></test>`,
					fval:  XJFloat64(-123.456),
					ival:  XJInt(-567),
				}),
				Entry("quoted zero values", XMLTest{
					input: `<?xml version="1.0" encoding="UTF-8" ?><test><floatVal>"0"</floatVal><intVal>"0"</intVal></test>`,
					fval:  XJFloat64(0),
					ival:  XJInt(0),
				}),
				Entry("empty values", XMLTest{
					input: `<?xml version="1.0" encoding="UTF-8" ?><test><floatVal></floatVal><intVal></intVal></test>`,
					fval:  XJFloat64(0),
					ival:  XJInt(0),
				}),
			)

		})

		Describe("Invalid Input ", func() {
			DescribeTable("value should fail:",
				func(test XMLTest) {
					var ts TestStruct
					e := ts.loadXML(test.input)
					Ω(e).Should(HaveOccurred())
				},
				Entry("values containing characters", XMLTest{
					input: `<?xml version="1.0" encoding="UTF-8" ?><test><floatVal>123.456x</floatVal><intVal>a567</intVal></test>`,
				}),
				Entry("values containing punctuation", XMLTest{
					input: `<?xml version="1.0" encoding="UTF-8" ?><test><floatVal>[123.456]</floatVal><intVal>5,67</intVal></test>`,
				}),
			)
		})

	})
})

// TestStruct is...
type TestStruct struct {
	XMLName xml.Name `xml:"test"`
	TStruct `json:"test" xml:"test"`
	// T       TStruct  `json:"test"`
}

func (r *TestStruct) loadJSON(input string) error {
	err := json.Unmarshal([]byte(input), r)
	if err != nil {
		msg := fmt.Sprintf("Unable to unmarshal JSON data input - error: %v", err)
		return errors.New(msg)
	}
	return nil
}

func (r *TestStruct) loadXML(input string) error {
	err := xml.Unmarshal([]byte(input), r)
	if err != nil {
		msg := fmt.Sprintf("Unable to unmarshal XML data input - error: %v", err)
		return errors.New(msg)
	}
	return nil
}

// TStruct is for testing unmarshalling of common.XJFloat64 and common.XJInt types
type TStruct struct {
	FloatVal XJFloat64 `json:"floatVal" xml:"floatVal"`
	IntVal   XJInt     `json:"intVal" xml:"intVal"`
}

type JSONTest struct {
	input string
	fval  XJFloat64
	ival  XJInt
}

type XMLTest struct {
	input string
	fval  XJFloat64
	ival  XJInt
}
