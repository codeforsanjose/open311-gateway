package geo_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/codeforsanjose/open311-gateway/_background/go/common/geo"
)

var _ = Describe("Geo Suite", func() {
	Describe("Address", func() {
		Describe("Parsing", func() {
			Context("Full Address", func() {
				DescribeTable("various full addresses",
					func(t AddrTest) {
						addr, e := NewAddr(t.addr)
						Ω(e).ShouldNot(HaveOccurred(), t.String())

						Ω(addr.FullAddr()).Should(Equal(t.fullAddr), t.String())
						Ω(addr.Lat).Should(Equal(t.lat), t.String())
						Ω(addr.Lng).Should(Equal(t.lng), t.String())
					},
					Entry("Aquatic Center / MH", AddrTest{
						desc:     `MH Aquatic Center`,
						addr:     "16200 Condit Road, Morgan Hill, CA",
						fullAddr: "16200 Condit Rd, Morgan Hill, CA 95037",
						lat:      37.1233549,
						lng:      -121.625339,
					}),
					Entry("City Hall", AddrTest{
						desc:     "San Jose City Hall",
						addr:     "200 E Santa Clara St, San Jose",
						fullAddr: "200 E Santa Clara St, San Jose, CA 95113",
						lat:      37.3377501,
						lng:      -121.885961,
					}),
					Entry("CfSJ Meeting", AddrTest{
						desc:     "Tech Museum Meeting",
						addr:     "159 W San Carlos St, San Jose",
						fullAddr: "159 W San Carlos St, San Jose, CA 95113",
						lat:      37.3302218,
						lng:      -121.8900109,
					}),
					Entry("Fresno", AddrTest{
						desc:     "Fresno",
						addr:     "478 W. San Ramon Ave, Unit 100, Fresno, CA",
						fullAddr: "478 W San Ramon Ave #100, Fresno, CA 93704",
						lat:      36.814449,
						lng:      -119.801436,
					}),
				)
			})

			Context("Lat/Lng", func() {
				DescribeTable("lat/lng",
					func(t LLAddrTest) {
						addr, e := NewAddrLL(t.ilat, t.ilng)
						Ω(e).ShouldNot(HaveOccurred(), t.String())

						Ω(addr.FullAddr()).Should(Equal(t.fullAddr), t.String())
						Ω(addr.Lat).Should(Equal(t.lat), t.String())
						Ω(addr.Lng).Should(Equal(t.lng), t.String())
					},
					Entry("Aquatic Center / MH", LLAddrTest{
						desc:     `MH Aquatic Center`,
						ilat:     37.1233549,
						ilng:     -121.625339,
						fullAddr: "16200 Condit Rd, Morgan Hill, CA 95037",
						lat:      37.1233549,
						lng:      -121.625339,
					}),
					Entry("City Hall", LLAddrTest{
						desc:     "San Jose City Hall",
						ilat:     37.3377501,
						ilng:     -121.885961,
						fullAddr: "200 E Santa Clara St, San Jose, CA 95113",
						lat:      37.3377501,
						lng:      -121.885961,
					}),
					Entry("CfSJ Meeting", LLAddrTest{
						desc:     "Tech Museum Meeting",
						ilat:     37.3302237,
						ilng:     -121.8900072,
						fullAddr: "145 W San Carlos St, San Jose, CA 95113",
						lat:      37.33037059999999,
						lng:      -121.8899529,
					}),
				)
			})
		})

	})
})

// TStruct contains a test and the Ωed expected
type AddrTest struct {
	desc     string
	addr     string
	fullAddr string
	lat, lng float64
}

func (r AddrTest) String() string {
	return fmt.Sprintf("   desc: %q\n   addr: %q\n   fullAddr: %#v\n   Lat, Lng: %v, %v\n", r.desc, r.addr, r.fullAddr, r.lat, r.lng)
}

// TStruct contains a test and the Ωed expected
type LLAddrTest struct {
	desc       string
	ilat, ilng float64
	fullAddr   string
	lat, lng   float64
}

func (r LLAddrTest) String() string {
	return fmt.Sprintf("   desc: %q\n   input lat/lng: %v, %v\n   fullAddr: %#v\n   Lat, Lng: %v, %v\n", r.desc, r.ilat, r.ilng, r.fullAddr, r.lat, r.lng)
}
