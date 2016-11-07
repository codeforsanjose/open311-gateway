package cv_test

import (
	"fmt"

	. "github.com/codeforsanjose/open311-gateway/_background/go/common/cv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Common Suite", func() {

	Describe("Conversion", func() {
		Describe("Create", func() {
			It("create a conversion", func() {
				v := NewConversion()
				Ω(fmt.Sprintf("%T", v)).Should(Equal("cv.Conversion"))
			})

		})

		Describe("Convert", func() {
			DescribeTable("various",
				func(t cnvTest) {
					c := NewConversion()
					result := c.Convert(t.name, t.input)
					Ω(c.IsOK(t.name)).Should(Equal(t.isok), t.String()+c.String())
					Ω(c.Ok()).Should(Equal(t.ok), t.String()+c.String())

					switch result := result.(type) {
					case *float64:
						Ω(*result).Should(Equal(t.expected.(float64)))
					case *int:
						Ω(*result).Should(Equal(t.expected.(int)))
					case *int32:
						Ω(*result).Should(Equal(t.expected.(int32)))
					case *int64:
						Ω(*result).Should(Equal(t.expected.(int64)))
					case *bool:
						Ω(*result).Should(Equal(t.expected.(bool)))
					default:
						Fail("Invalid type")
					}
				},
				Entry("valid lat", cnvTest{
					desc:     `Valid Latitude`,
					input:    "32.9999",
					name:     "Latitude",
					isok:     true,
					ok:       true,
					expected: 32.9999,
				}),
				Entry("zero lat", cnvTest{
					desc:     `Zero Latitude`,
					input:    "0",
					name:     "Latitude",
					isok:     true,
					ok:       true,
					expected: 0.0,
				}),
				Entry("empty lat", cnvTest{
					desc:     `Empty Latitude`,
					input:    "",
					name:     "Latitude",
					isok:     false,
					ok:       false,
					expected: 0.0,
				}),
				Entry("invalid lat", cnvTest{
					desc:     `Invalid Latitude`,
					input:    "32,9999",
					name:     "Latitude",
					isok:     false,
					ok:       false,
					expected: 0.0,
				}),

				Entry("valid lng", cnvTest{
					desc:     `Valid Longitude`,
					input:    "-121.2345",
					name:     "Longitude",
					isok:     true,
					ok:       true,
					expected: -121.2345,
				}),
				Entry("zero lng", cnvTest{
					desc:     `Zero Longitude`,
					input:    "0",
					name:     "Longitude",
					isok:     true,
					ok:       true,
					expected: 0.0,
				}),
				Entry("empty lng", cnvTest{
					desc:     `Empty Longitude`,
					input:    "",
					name:     "Longitude",
					isok:     false,
					ok:       false,
					expected: 0.0,
				}),
				Entry("invalid lng", cnvTest{
					desc:     `Invalid Longitude`,
					input:    "-121,2345",
					name:     "Latitude",
					isok:     false,
					ok:       false,
					expected: 0.0,
				}),

				Entry("valid radius", cnvTest{
					desc:     `Valid Radius`,
					input:    "129",
					name:     "Radius",
					isok:     true,
					ok:       true,
					expected: 129,
				}),
				Entry("zero radius", cnvTest{
					desc:     `Zero Radius`,
					input:    "0",
					name:     "Radius",
					isok:     true,
					ok:       true,
					expected: 0,
				}),
				Entry("empty radius", cnvTest{
					desc:     `Empty Radius`,
					input:    "",
					name:     "Radius",
					isok:     true,
					ok:       true,
					expected: 100,
				}),
				Entry("invalid radius", cnvTest{
					desc:     `Invalid Radius`,
					input:    "1/2 mile",
					name:     "Radius",
					isok:     false,
					ok:       false,
					expected: 0,
				}),

				Entry("valid maxresults", cnvTest{
					desc:     `Valid MaxResults`,
					input:    "23",
					name:     "MaxResults",
					isok:     true,
					ok:       true,
					expected: 23,
				}),
				Entry("zero maxresults", cnvTest{
					desc:     `Zero MaxResults`,
					input:    "0",
					name:     "MaxResults",
					isok:     true,
					ok:       true,
					expected: 0,
				}),
				Entry("empty maxresults", cnvTest{
					desc:     `Empty MaxResults`,
					input:    "",
					name:     "MaxResults",
					isok:     true,
					ok:       true,
					expected: 10,
				}),

				Entry("valid includedetails", cnvTest{
					desc:     `Valid IncludeDetails`,
					input:    "true",
					name:     "IncludeDetails",
					isok:     true,
					ok:       true,
					expected: true,
				}),
				Entry("zero includedetails", cnvTest{
					desc:     `Zero IncludeDetails`,
					input:    "false",
					name:     "IncludeDetails",
					isok:     true,
					ok:       true,
					expected: false,
				}),
				Entry("empty includedetails", cnvTest{
					desc:     `Empty IncludeDetails`,
					input:    "",
					name:     "IncludeDetails",
					isok:     true,
					ok:       true,
					expected: false,
				}),
			)

			Describe("float", func() {
				It("should be a valid float", func() {
					c := NewConversion()
					result := c.Float("Latitude", "18.001")
					Ω(c.IsOK("Latitude")).Should(BeTrue(), c.String())
					Ω(c.Ok()).Should(BeTrue(), c.String())

					Ω(result).Should(Equal(18.001))
				})
			})

			DescribeTable("Float",
				func(t cnvTest) {
					c := NewConversion()
					result := c.Float(t.name, t.input)
					Ω(c.IsOK(t.name)).Should(Equal(t.isok), t.String()+c.String())
					Ω(c.Ok()).Should(Equal(t.ok), t.String()+c.String())

					Ω(result).Should(Equal(t.expected.(float64)))
				},
				Entry("valid lat", cnvTest{
					desc:     `Valid Latitude`,
					input:    "32.9999",
					name:     "Latitude",
					isok:     true,
					ok:       true,
					expected: 32.9999,
				}),
				Entry("zero lat", cnvTest{
					desc:     `Zero Latitude`,
					input:    "0",
					name:     "Latitude",
					isok:     true,
					ok:       true,
					expected: 0.0,
				}),
				Entry("empty lat", cnvTest{
					desc:     `Empty Latitude`,
					input:    "",
					name:     "Latitude",
					isok:     false,
					ok:       false,
					expected: 0.0,
				}),
				Entry("invalid lat", cnvTest{
					desc:     `Invalid Latitude`,
					input:    "32,9999",
					name:     "Latitude",
					isok:     false,
					ok:       false,
					expected: 0.0,
				}),

				Entry("valid lng", cnvTest{
					desc:     `Valid Longitude`,
					input:    "-121.2345",
					name:     "Longitude",
					isok:     true,
					ok:       true,
					expected: -121.2345,
				}),
				Entry("zero lng", cnvTest{
					desc:     `Zero Longitude`,
					input:    "0",
					name:     "Longitude",
					isok:     true,
					ok:       true,
					expected: 0.0,
				}),
				Entry("empty lng", cnvTest{
					desc:     `Empty Longitude`,
					input:    "",
					name:     "Longitude",
					isok:     false,
					ok:       false,
					expected: 0.0,
				}),
				Entry("invalid lng", cnvTest{
					desc:     `Invalid Longitude`,
					input:    "-121,2345",
					name:     "Latitude",
					isok:     false,
					ok:       false,
					expected: 0.0,
				}),
			)

			DescribeTable("Int",
				func(t cnvTest) {
					c := NewConversion()
					result := c.Int(t.name, t.input)
					Ω(c.IsOK(t.name)).Should(Equal(t.isok), t.String()+c.String())
					Ω(c.Ok()).Should(Equal(t.ok), t.String()+c.String())

					Ω(result).Should(Equal(t.expected.(int)))
				},
				Entry("valid radius", cnvTest{
					desc:     `Valid Radius`,
					input:    "129",
					name:     "Radius",
					isok:     true,
					ok:       true,
					expected: 129,
				}),
				Entry("zero radius", cnvTest{
					desc:     `Zero Radius`,
					input:    "0",
					name:     "Radius",
					isok:     true,
					ok:       true,
					expected: 0,
				}),
				Entry("empty radius", cnvTest{
					desc:     `Empty Radius`,
					input:    "",
					name:     "Radius",
					isok:     true,
					ok:       true,
					expected: 100,
				}),
				Entry("invalid radius", cnvTest{
					desc:     `Invalid Radius`,
					input:    "1/2 mile",
					name:     "Radius",
					isok:     false,
					ok:       false,
					expected: 0,
				}),

				Entry("valid maxresults", cnvTest{
					desc:     `Valid MaxResults`,
					input:    "23",
					name:     "MaxResults",
					isok:     true,
					ok:       true,
					expected: 23,
				}),
				Entry("zero maxresults", cnvTest{
					desc:     `Zero MaxResults`,
					input:    "0",
					name:     "MaxResults",
					isok:     true,
					ok:       true,
					expected: 0,
				}),
				Entry("empty maxresults", cnvTest{
					desc:     `Empty MaxResults`,
					input:    "",
					name:     "MaxResults",
					isok:     true,
					ok:       true,
					expected: 10,
				}),
			)

			DescribeTable("Bool",
				func(t cnvTest) {
					c := NewConversion()
					result := c.Bool(t.name, t.input)
					Ω(c.IsOK(t.name)).Should(Equal(t.isok), t.String()+c.String())
					Ω(c.Ok()).Should(Equal(t.ok), t.String()+c.String())

					Ω(result).Should(Equal(t.expected.(bool)))
				},
				Entry("valid includedetails", cnvTest{
					desc:     `Valid IncludeDetails`,
					input:    "true",
					name:     "IncludeDetails",
					isok:     true,
					ok:       true,
					expected: true,
				}),
				Entry("zero includedetails", cnvTest{
					desc:     `Zero IncludeDetails`,
					input:    "false",
					name:     "IncludeDetails",
					isok:     true,
					ok:       true,
					expected: false,
				}),
				Entry("empty includedetails", cnvTest{
					desc:     `Empty IncludeDetails`,
					input:    "",
					name:     "IncludeDetails",
					isok:     true,
					ok:       true,
					expected: false,
				}),
			)
		})
	})
})

// cnvTest is used for more complicated Validation tests
type cnvTest struct {
	desc     string
	input    string
	name     string
	isok     bool
	ok       bool
	expected interface{}
}

func (r cnvTest) String() string {
	return fmt.Sprintf("   desc: %q\n   input: %q\n   IsOk: %v\n   Ok: %v\n   expected: (%T)%#v\n", r.desc, r.input, r.isok, r.ok, r.expected, r.expected)
}
