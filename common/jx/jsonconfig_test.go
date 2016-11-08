package jx_test

import (
	. "github.com/codeforsanjose/open311-gateway/_background/go/common/jx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestStruct1 struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
}

var _ = Describe("JX", func() {
	Describe("JSONConfig", func() {
		Context("Valid Input", func() {
			It("should load the struct", func() {
				jc := new(JSONConfig)
				ts := new(TestStruct1)
				err := jc.Read(ts, "jsonconfig_testfile1.json")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(ts.Account).Should(Equal("xyz@gmail.com"))
				Ω(ts.Password).Should(Equal("opensesame"))
				Ω(ts.Server).Should(Equal("smtp.gmail.com"))
				Ω(ts.Port).Should(Equal(587))
			})

		})
	})
})
