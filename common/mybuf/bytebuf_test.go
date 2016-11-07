package mybuf_test

import (
	"bytes"

	. "github.com/codeforsanjose/open311-gateway/_background/go/common/mybuf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Buffer", func() {
	Describe("Concat", func() {
		Context("Valid Input", func() {
			DescribeTable("values",
				func(b1, b2 *bytes.Buffer, expected string) {
					l := int64(b2.Len())
					n, e := Concat(b1, b2)
					Ω(e).NotTo(HaveOccurred())
					Ω(b1.String()).Should(Equal(expected))
					Ω(n).Should(Equal(l))
				},
				Entry("First Second", bytes.NewBufferString("First "), bytes.NewBufferString("Second"), "First Second"),
				Entry("First\nSecond", bytes.NewBufferString("First\n"), bytes.NewBufferString("Second"), "First\nSecond"),
				Entry("First", bytes.NewBufferString("First"), new(bytes.Buffer), "First"),
			)
		})
	})

	Describe("Concat / ByteSlice", func() {
		Context("Valid Input", func() {
			DescribeTable("values",
				func(b1, b2 *bytes.Buffer, expected []byte) {
					l := int64(b2.Len())
					n, e := Concat(b1, b2)
					Ω(e).NotTo(HaveOccurred())
					bs := ToBSlice(b1)
					// fmt.Printf("\nbs: (%T)%v\n", bs, bs)
					Ω(bs).Should(Equal(expected))
					Ω(n).Should(Equal(l))
				},
				Entry("First Second", bytes.NewBufferString("First "), bytes.NewBufferString("Second"), []byte{70, 105, 114, 115, 116, 32, 83, 101, 99, 111, 110, 100}),
				Entry("First\nSecond", bytes.NewBufferString("First\n"), bytes.NewBufferString("Second"), []byte{70, 105, 114, 115, 116, 10, 83, 101, 99, 111, 110, 100}),
				Entry("First", bytes.NewBufferString("First"), new(bytes.Buffer), []byte{70, 105, 114, 115, 116}),
			)
		})
	})

	Describe("Copy", func() {
		Context("Valid Input", func() {
			DescribeTable("values",
				func(orig *bytes.Buffer, expected []byte) {
					copy, size, e := Copy(orig)
					Ω(e).NotTo(HaveOccurred())
					slOrig := orig.Bytes()
					slCopy := copy.Bytes()
					// fmt.Printf("\norig:     %v\ncopy:     %v\nexpected: %v\n", slOrig, slCopy, expected)
					Ω(bytes.Compare(slOrig, expected)).Should(Equal(0))
					Ω(bytes.Compare(slOrig, slCopy)).Should(Equal(0))
					Ω(size).Should(Equal(len(expected)))
				},
				Entry("empty buffer", bytes.NewBuffer([]byte{}), []byte{}),
				Entry("test String", bytes.NewBufferString("Test String 01"), []byte{84, 101, 115, 116, 32, 83, 116, 114, 105, 110, 103, 32, 48, 49}),
			)
		})
	})

})
