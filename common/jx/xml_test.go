package jx_test

import (
	"encoding/xml"
	"io/ioutil"
	"strings"

	. "github.com/codeforsanjose/open311-gateway/_background/go/common/jx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

const (
	xmlEncodeResult1 = `<XMLTestT1 AttrString="string" AttrBool="true" AttrInt="999"><ValString>value</ValString><ValInt>-999</ValInt></XMLTestT1>`
	xmlEncodeResult2 = `<?xml version="1.0" encoding="UTF-8"?>
<XMLTestT1 AttrString="string" AttrBool="true" AttrInt="999"><ValString>value</ValString><ValInt>-999</ValInt></XMLTestT1>`
	xmlEncodeResult3 = `<?xml version="1.0" encoding="UTF-8"?>
  <XMLTestT1 AttrString="string" AttrBool="true" AttrInt="999">
      <ValString>value</ValString>
      <ValInt>-999</ValInt>
  </XMLTestT1>`
)

var _ = Describe("XML", func() {
	Describe("EncodeXML", func() {
		Context("Valid Input", func() {
			DescribeTable("test values",
				func(ts interface{}, indent, header bool, expected string) {
					x, e := EncodeXML(ts, indent, header)
					Ω(e).NotTo(HaveOccurred())
					Ω(x.String()).Should(Equal(expected))
				},
				Entry("no indent, header", XMLTestT1{
					AttribS: "string",
					AttribB: true,
					AttribI: 999,
					ValS:    "value",
					ValI:    -999,
				}, false, false, xmlEncodeResult1),
				Entry("no indent, has header", XMLTestT1{
					AttribS: "string",
					AttribB: true,
					AttribI: 999,
					ValS:    "value",
					ValI:    -999,
				}, false, true, xmlEncodeResult2),
				Entry("indent and header", XMLTestT1{
					AttribS: "string",
					AttribB: true,
					AttribI: 999,
					ValS:    "value",
					ValI:    -999,
				}, true, true, xmlEncodeResult3),
				Entry("nil pointer", nil, false, false, ""),
				Entry("string", "some text", false, false, "<string>some text</string>"),
				Entry("int", 999, false, false, "<int>999</int>"),
				Entry("int slice", []int{1, 2, 3}, false, false, "<int>1</int><int>2</int><int>3</int>"),
			)
		})
	})

	Describe("EncodeXMLByte", func() {
		Context("Valid Input", func() {
			DescribeTable("test values",
				func(ts interface{}, indent, header bool, expected string) {
					x, e := EncodeXMLByte(ts, indent, header)
					Ω(e).NotTo(HaveOccurred())
					Ω(string(x)).Should(Equal(expected))
				},
				Entry("no indent, header", XMLTestT1{
					AttribS: "string",
					AttribB: true,
					AttribI: 999,
					ValS:    "value",
					ValI:    -999,
				}, false, false, xmlEncodeResult1),
				Entry("no indent, has header", XMLTestT1{
					AttribS: "string",
					AttribB: true,
					AttribI: 999,
					ValS:    "value",
					ValI:    -999,
				}, false, true, xmlEncodeResult2),
				Entry("indent and header", XMLTestT1{
					AttribS: "string",
					AttribB: true,
					AttribI: 999,
					ValS:    "value",
					ValI:    -999,
				}, true, true, xmlEncodeResult3),
				Entry("nil pointer", nil, false, false, ""),
				Entry("string", "some text", false, false, "<string>some text</string>"),
				Entry("int", 999, false, false, "<int>999</int>"),
				Entry("int slice", []int{1, 2, 3}, false, false, "<int>1</int><int>2</int><int>3</int>"),
			)
		})
	})

	Describe("LoadXMLFile", func() {
		Context("Valid File", func() {
			It("should load and parse the file", func() {
				var t XMLTestT1
				filename := "xml_test_file1.xml"
				e := LoadXMLFile(filename, &t)
				Ω(e).NotTo(HaveOccurred())
				Ω(t.AttribS).Should(Equal("string"))

				b, _ := EncodeXMLByte(t, false, false)
				xc, _ := ioutil.ReadFile(filename)
				Ω(strings.TrimSpace(string(b))).Should(Equal(strings.TrimSpace(string(xc))))

			})
		})
	})

	// Describe("ConcatBuffers", func() {
	// 	Describe("Valid Input", func() {
	// 		DescribeTable("test values",
	// 			func(b1, b2 *bytes.Buffer, expected string) {
	// 				l := int64(b2.Len())
	// 				n, e := ConcatBuffers(b1, b2)
	// 				Ω(e).NotTo(HaveOccurred())
	// 				Ω(b1.String()).Should(Equal(expected))
	// 				Ω(n).Should(Equal(l))
	// 				_, e = ConcatBuffers(b1, b2)
	// 				Ω(e).NotTo(HaveOccurred())
	// 			},
	// 			Entry("First Second", bytes.NewBufferString("First "), bytes.NewBufferString("Second"), "First Second"),
	// 			Entry("First\nSecond", bytes.NewBufferString("First\n"), bytes.NewBufferString("Second"), "First\nSecond"),
	// 			Entry("First", bytes.NewBufferString("First"), new(bytes.Buffer), "First"),
	// 		)
	// 	})
	// })
})

// XMLTestT1 is for testing..
type XMLTestT1 struct {
	XMLName xml.Name `xml:"XMLTestT1"`
	AttribS string   `xml:"AttrString,attr"`
	AttribB bool     `xml:"AttrBool,attr"`
	AttribI int      `xml:"AttrInt,attr"`
	ValS    string   `xml:"ValString"`
	ValI    int      `xml:"ValInt"`
}
