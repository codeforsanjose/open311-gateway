package mystrings_test

import (
	"fmt"

	. "github.com/codeforsanjose/open311-gateway/_background/go/common/mystrings"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("MyStrings", func() {

	Describe("Single Operations", func() {
		DescribeTable("Trim",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{Trim(", ")},
				desc:     `Trim(", ")`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{Trim(", ")},
				desc:     `Trim(", ")`,
				input:    " , ",
				expected: "",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{Trim(", ")},
				desc:     `Trim(", ")`,
				input:    "99999, ",
				expected: "99999",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{Trim(", ")},
				desc:     `Trim(", ")`,
				input:    " 99999, ",
				expected: "99999",
			}),
		)

		DescribeTable("TrimLeft",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{TrimLeft(", ")},
				desc:     `TrimLeft(", ")`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{TrimLeft(", ")},
				desc:     `TrimLeft(", ")`,
				input:    " , ",
				expected: "",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{TrimLeft(", ")},
				desc:     `TrimLeft(", ")`,
				input:    "99999, ",
				expected: "99999, ",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{TrimLeft(", ")},
				desc:     `TrimLeft(", ")`,
				input:    " , 99999, ",
				expected: "99999, ",
			}),
		)

		DescribeTable("TrimRight",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{TrimRight(", ")},
				desc:     `TrimRight(", ")`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{TrimRight(", ")},
				desc:     `TrimRight(", ")`,
				input:    " , ",
				expected: "",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{TrimRight(", ")},
				desc:     `TrimRight(", ")`,
				input:    "99999, ",
				expected: "99999",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{TrimRight(", ")},
				desc:     `TrimRight(", ")`,
				input:    " , 99999, ",
				expected: " , 99999",
			}),
		)

		DescribeTable("TrimPrefix",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{TrimPrefix("cut this")},
				desc:     `TrimPrefix("cut this")`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{TrimPrefix("cut this")},
				desc:     `TrimPrefix("cut this")`,
				input:    "cut this",
				expected: "",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{TrimPrefix("cut this")},
				desc:     `TrimPrefix("cut this")`,
				input:    "cut thisHi, how are you?",
				expected: "Hi, how are you?",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{TrimPrefix("cut this")},
				desc:     `TrimPrefix("cut this")`,
				input:    "cut thisHi, how are you? cut this",
				expected: "Hi, how are you? cut this",
			}),
		)

		DescribeTable("TrimSuffix",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{TrimSuffix("cut this")},
				desc:     `TrimSuffix("cut this")`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{TrimSuffix("cut this")},
				desc:     `TrimSuffix("cut this")`,
				input:    "cut this",
				expected: "",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{TrimSuffix("cut this")},
				desc:     `TrimSuffix("cut this")`,
				input:    "cut thisHi, how are you?",
				expected: "cut thisHi, how are you?",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{TrimSuffix("cut this")},
				desc:     `TrimSuffix("cut this")`,
				input:    "cut thisHi, how are you? cut this",
				expected: "cut thisHi, how are you? ",
			}),
		)

		DescribeTable("ReplaceOne",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{ReplaceOne("goat", "cat", -1)},
				desc:     `ReplaceOne("goat", "cat", -1)`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{ReplaceOne("goat", "cat", -1)},
				desc:     `ReplaceOne("goat", "cat", -1)`,
				input:    "goat",
				expected: "cat",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{ReplaceOne("goat", "cat", -1)},
				desc:     `ReplaceOne("goat", "cat", -1)`,
				input:    "All goats have claws, and goats have sharp teeth",
				expected: "All cats have claws, and cats have sharp teeth",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{ReplaceOne("goat", "cat", 1)},
				desc:     `ReplaceOne("goat", "cat", -1)`,
				input:    "All goats have claws, and goats have sharp teeth",
				expected: "All cats have claws, and goats have sharp teeth",
			}),
		)

		DescribeTable("DeleteChars",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{DeleteChars("- ")},
				desc:     `DeleteChars("- ")`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{DeleteChars("- ")},
				desc:     `DeleteChars("- ")`,
				input:    "123",
				expected: "123",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{DeleteChars("- ")},
				desc:     `DeleteChars("- ")`,
				input:    "99999-1234",
				expected: "999991234",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{DeleteChars("- ")},
				desc:     `DeleteChars("- ")`,
				input:    "99999 1234",
				expected: "999991234",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{DeleteChars("- ")},
				desc:     `DeleteChars("- ")`,
				input:    "99999- 1234",
				expected: "999991234",
			}),
		)

		DescribeTable("Replace",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{Replace([]string{"-", "*", ",", "."})},
				desc:     `Replace("-", "*", ",", ".")`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{Replace([]string{"-", "*", ",", "."})},
				desc:     `Replace("-", "*", ",", ".")`,
				input:    "-,",
				expected: "*.",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{Replace([]string{"-", "*", ",", "."})},
				desc:     `Replace("-", "*", ",", ".")`,
				input:    "99999-1234",
				expected: "99999*1234",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{Replace([]string{"-", "*", ",", "."})},
				desc:     `Replace("-", "*", ",", ".")`,
				input:    "some-text,  more*text",
				expected: "some*text.  more*text",
			}),
		)
		DescribeTable("TrimSpace",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{TrimSpace},
				desc:     `TrimSpace`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{TrimSpace},
				desc:     `TrimSpace`,
				input:    "          ",
				expected: "",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{TrimSpace},
				desc:     `TrimSpace`,
				input:    "99999 ",
				expected: "99999",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{TrimSpace},
				desc:     `TrimSpace`,
				input:    "     99999  ",
				expected: "99999",
			}),
		)

		DescribeTable("Upper",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{Upper},
				desc:     `Upper`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{Upper},
				desc:     `Upper`,
				input:    "t",
				expected: "T",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{Upper},
				desc:     `Upper`,
				input:    " here's an Upper Case string ",
				expected: " HERE'S AN UPPER CASE STRING ",
			}),
		)

		DescribeTable("Lower",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{Lower},
				desc:     `Lower`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{Lower},
				desc:     `Lower`,
				input:    "T",
				expected: "t",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{Lower},
				desc:     `Lower`,
				input:    " AND here's A Lower Case string ",
				expected: " and here's a lower case string ",
			}),
		)

		DescribeTable("Title",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("empty string", MGTest{
				ops:      []func(string) string{Title},
				desc:     `Title`,
				input:    "",
				expected: "",
			}),
			Entry("fully trimmed string", MGTest{
				ops:      []func(string) string{Title},
				desc:     `Title`,
				input:    "t",
				expected: "T",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{Title},
				desc:     `Title`,
				input:    "here is a title case string",
				expected: "Here Is A Title Case String",
			}),
			Entry("normal string", MGTest{
				ops:      []func(string) string{Title},
				desc:     `Title`,
				input:    "here is a TITLE case string",
				expected: "Here Is A TITLE Case String",
			}),
		)

	})

	Describe("Multiple Operations", func() {
		DescribeTable("Various Ops",
			func(test MGTest) {
				mg, e := NewMacGyver(test.ops...)
				Ω(e).ShouldNot(HaveOccurred(), test.String())

				v := mg.Process(test.input)
				Ω(v).Should(Equal(test.expected), test.String())
			},
			Entry("TrimSpace, Upper on empty string", MGTest{
				ops:      []func(string) string{TrimSpace, Upper},
				desc:     `TrimSpace, Upper`,
				input:    "",
				expected: "",
			}),
			Entry("TrimSpace, Upper", MGTest{
				ops:      []func(string) string{TrimSpace, Upper},
				desc:     `TrimSpace, Upper`,
				input:    "  this is a string  ",
				expected: "THIS IS A STRING",
			}),
			Entry("TrimSpace, Upper, Delete", MGTest{
				ops:      []func(string) string{TrimSpace, Upper, DeleteChars("*")},
				desc:     `TrimSpace, Upper, DeleteChars`,
				input:    "  this is a *string*  ",
				expected: "THIS IS A STRING",
			}),
			Entry("TrimSpace, Upper, Delete, ReplaceOne", MGTest{
				ops:      []func(string) string{TrimSpace, DeleteChars("*"), ReplaceOne(" is ", " was ", 1), Upper},
				desc:     `TrimSpace, Upper, DeleteChars`,
				input:    "  this is a *string*  ",
				expected: "THIS WAS A STRING",
			}),
			Entry("TrimSpace, Upper, Delete, Replace", MGTest{
				ops:      []func(string) string{TrimSpace, DeleteChars("*"), Replace([]string{" is ", " was "}), Upper},
				desc:     `TrimSpace, Upper, DeleteChars`,
				input:    "  this is a *string*  ",
				expected: "THIS WAS A STRING",
			}),
		)

	})
})

// TStruct contains a test and the Ωed expected
type MGTest struct {
	ops      []func(string) string
	desc     string
	input    string
	expected string
}

func (r MGTest) String() string {
	return fmt.Sprintf("   func: %#v\n   desc: %q\n   input: %q\n   expected: %#v\n", r.ops, r.desc, r.input, r.expected)
}
