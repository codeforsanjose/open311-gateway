package cv_test

import (
	"fmt"

	. "github.com/codeforsanjose/open311-gateway/_background/go/common/cv"
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Common Suite", func() {

	Describe("Validation", func() {
		Describe("Create", func() {
			It("create a validation", func() {
				v := NewValidation()
				Ω(fmt.Sprintf("%T", v)).Should(Equal("cv.Validation"))
			})

			It("create and set a validation", func() {
				v := NewValidation()
				v.Set("item1", "item 1 is OK.", true)
				v.Set("item2", "item 2 is OK.", true)
				Ω(v.IsOK("item1")).Should(BeTrue())
				Ω(v.IsOK("item2")).Should(BeTrue())
				// fmt.Printf(v.Error())
				// fmt.Printf(v.String())
				Ω(v.Ok()).Should(BeTrue())
			})

			It("create and set a not OK validation", func() {
				v := NewValidation()
				v.Set("item1", "item 1 is OK.", true)
				v.Set("item2", "item 2 is broken!", false)
				Ω(v.IsOK("item1")).Should(BeTrue())
				Ω(v.IsOK("item2")).Should(BeFalse())
				// fmt.Printf("\n%s", v.Error())
				// fmt.Printf(v.String())
				Ω(v.Ok()).Should(BeFalse())
			})
		})
	})
})
