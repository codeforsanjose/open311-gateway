package mystrings_test

import (
	"fmt"

	. "github.com/codeforsanjose/open311-gateway/_background/go/common/mystrings"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("MyRegexp", func() {

	DescribeTable("Simple Tests",
		func(test ReTStruct) {
			rx := NewRegex(test.rexp, test.trimset, test.delset)

			Ω(rx.Match(test.input)).Should(Succeed(), test.String())
			Ω(rx.Ok).Should(BeTrue())
			Ω(rx.Named).Should(Equal(test.named), test.String())
		},
		Entry("test1", ReTStruct{
			rexp:    `(?P<first>\d+)\.(?P<second>\d+)`,
			input:   "1.2",
			trimset: "",
			delset:  "",
			match:   false,
			named:   map[string]string{"first": "1", "second": "2"},
			all:     []string{"1.2", "1", "2"},
		}),
		Entry("test2", ReTStruct{
			rexp:    `(?P<first>\d+)\.(?P<second>\d+)`,
			input:   "1234.5678",
			trimset: "",
			delset:  "",
			match:   true,
			named:   map[string]string{"first": "1234", "second": "5678"},
			all:     []string{"1234.5678", "1234", "5678"},
		}),
		Entry("test3", ReTStruct{
			rexp:    `(?P<first>\d+)\.(\d+)`,
			input:   "1234.5678",
			trimset: "",
			delset:  "",
			match:   true,
			named:   map[string]string{"first": "1234"},
			all:     []string{"1234.5678", "1234", "5678"},
		}),
		Entry("test4", ReTStruct{
			rexp:    `(\d+)\.(\d+)`,
			input:   "1234.5678",
			trimset: "",
			delset:  "",
			match:   true,
			named:   map[string]string{},
			all:     []string{"1234.5678", "1234", "5678"},
		}),
	)

	DescribeTable("Fail Tests",
		func(test ReTStruct) {
			rx := NewRegex(test.rexp, test.trimset, test.delset)

			e := rx.Match(test.input)
			Ω(e).Should(MatchError("no matches found"))
			Ω(rx.Ok).Should(BeFalse())
		},
		Entry("test1", ReTStruct{
			rexp:    `(?P<first>\d+)\.(?P<second>\d+)`,
			input:   ".",
			trimset: "",
			delset:  "",
			match:   true,
			named:   nil,
			all:     nil,
		}),
		Entry("test2", ReTStruct{
			rexp:    `(?i)^((?P<name>.*)[ ,]+)?(?P<addr>\d{2,}.*)[ ,]+(?P<city>morgan hill|gilroy|san jose|hollister|san martin)[ ,]*(?P<state>CA)? *(?P<zip>(?:\d\d\d\d\d)-*(?:\d\d\d\d)*)?([ .,-]+(?P<phone>(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}))?([ .,-]+(?P<post>.*))?$`,
			input:   "Morgan Hill, CA",
			trimset: "",
			delset:  "",
			match:   false,
			named:   nil,
			all:     nil,
		}),
	)

	Context("311Gateway", func() {
		DescribeTable("Query Parm Tests",
			func(test ReTStruct) {
				rx := NewRegex(test.rexp, test.trimset, test.delset)

				Ω(rx.Match(test.input)).Should(Succeed(), test.String())
				Ω(rx.Ok).Should(BeTrue())
				Ω(rx.Named).Should(Equal(test.named), test.String())
			},
			Entry("test1", ReTStruct{
				rexp:    `(?i)^attribute\[(?P<name>(\w+))\]`,
				input:   "attribute[234]",
				trimset: "",
				delset:  "",
				match:   false,
				named:   map[string]string{"name": "234"},
				all:     []string{"attribute[234]Ugly graffiti.", "234"},
			}),
			Entry("test2", ReTStruct{
				rexp:    `(?i)^attribute\[(?P<name>(\w+))\]`,
				input:   "attribute[WHISPAWN]",
				trimset: "",
				delset:  "",
				match:   false,
				named:   map[string]string{"name": "WHISPAWN"},
				all:     []string{"attribute[WHISPAWN]", "WHISPAWN"},
			}),
		)
	})

	DescribeTable("Address Tests",
		func(test ReTStruct) {
			rx := NewRegex(test.rexp, test.trimset, test.delset)

			Ω(rx.Match(test.input)).Should(Succeed(), test.String())
			Ω(rx.Ok).Should(BeTrue())
			Ω(rx.Named).Should(Equal(test.named), test.String())
		},
		Entry("name, address and phone", ReTStruct{
			rexp:    `(?i)^((?P<name>.*)[ ,]+)?(?P<addr>\d{2,}.*)[ ,]+(?P<city>morgan hill|gilroy|san jose|hollister|san martin)[ ,]*(?P<state>CA)? *(?P<zip>(?:\d\d\d\d\d)-*(?:\d\d\d\d)*)?([ .,-]+(?P<phone>(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}))?([ .,-]+(?P<post>.*))?$`,
			input:   "James Haskell, 17200 Quail Ct., Morgan Hill, CA 95037  111-222-3333",
			trimset: "",
			delset:  "",
			match:   true,
			named:   map[string]string{"name": "James Haskell,", "addr": "17200 Quail Ct.,", "city": "Morgan Hill", "state": "CA", "zip": "95037", "phone": "111-222-3333", "post": ""},
			all:     []string{"James Haskell, 17200 Quail Ct., Morgan Hill, CA 95037", "James Haskell,", "17200 Quail Ct.,", "Morgan Hill", "CA", "95037"},
		}),
		Entry("name, address, and phone; with strip", ReTStruct{
			rexp:    `(?i)^((?P<name>.*)[ ,]+)?(?P<addr>\d{2,}.*)[ ,]+(?P<city>morgan hill|gilroy|san jose|hollister|san martin)[ ,]*(?P<state>CA)? *(?P<zip>(?:\d\d\d\d\d)-*(?:\d\d\d\d)*)?([ .,-]+(?P<phone>(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}))?([ .,-]+(?P<post>.*))?$`,
			input:   "James Haskell, 17200 Quail Ct., Morgan Hill, CA 95037  111-222-3333",
			trimset: " ,",
			delset:  "",
			match:   true,
			named:   map[string]string{"name": "James Haskell", "addr": "17200 Quail Ct.", "city": "Morgan Hill", "state": "CA", "zip": "95037", "phone": "111-222-3333", "post": ""},
			all:     []string{"James Haskell, 17200 Quail Ct., Morgan Hill, CA 95037", "James Haskell", "17200 Quail Ct.", "Morgan Hill", "CA", "95037"},
		}),
	)

	DescribeTable("Zip Code Tests",
		func(test ReTStruct) {
			rx := NewRegex(test.rexp, test.trimset, test.delset)

			Ω(rx.Match(test.input)).Should(Succeed(), test.String())
			Ω(rx.Ok).Should(BeTrue())
			Ω(rx.Named).Should(Equal(test.named), test.String())
		},
		Entry("5 digit zip", ReTStruct{
			rexp:    `^(?P<zip>\d{5}(?:[-\s]*\d{4})?)$`,
			input:   "95037",
			trimset: "",
			delset:  "",
			match:   true,
			named:   map[string]string{"zip": "95037"},
			all:     []string{"95037", "95037"},
		}),
		Entry("9 digit zip", ReTStruct{
			rexp:    `^(?P<zip>\d{5}(?:[-\s]*\d{4})?)$`,
			input:   "999991234",
			trimset: " ,",
			delset:  "- ",
			match:   true,
			named:   map[string]string{"zip": "999991234"},
			all:     []string{"999991234", "999991234"},
		}),
		Entry("9 digit zip with dash", ReTStruct{
			rexp:    `^(?P<zip>\d{5}(?:[-\s]*\d{4})?)$`,
			input:   "99999-1234",
			trimset: " ,",
			delset:  "- ",
			match:   true,
			named:   map[string]string{"zip": "999991234"},
			all:     []string{"999991234", "999991234"},
		}),
		Entry("9 digit zip", ReTStruct{
			rexp:    `^(?P<zip>\d{5}(?:[-\s]*\d{4})?)$`,
			input:   "99999 1234",
			trimset: " ,",
			delset:  "- ",
			match:   true,
			named:   map[string]string{"zip": "999991234"},
			all:     []string{"999991234", "999991234"},
		}),
	)

	DescribeTable("Zip Code FAIL Tests",
		func(test ReTStruct) {
			rx := NewRegex(test.rexp, test.trimset, test.delset)

			Ω(rx.Match(test.input)).ShouldNot(Succeed(), test.String())
			Ω(rx.Ok).Should(BeFalse())
		},
		Entry("6 digit zip", ReTStruct{
			rexp:    `^(?P<zip>\d{5}(?:[-\s]*\d{4})?)$`,
			input:   "123456",
			trimset: "",
			delset:  "",
			match:   false,
			named:   map[string]string{"zip": "95037"},
			all:     []string{"95037", "95037"},
		}),
		Entry("10 digit zip", ReTStruct{
			rexp:    `^(?P<zip>\d{5}(?:[-\s]*\d{4})?)$`,
			input:   "9999912345",
			trimset: " ,",
			delset:  "- ",
			match:   false,
			named:   map[string]string{"zip": "999991234"},
			all:     []string{"999991234", "999991234"},
		}),
		Entry("9 digit zip with asterisk", ReTStruct{
			rexp:    `^(?P<zip>\d{5}(?:[-\s]*\d{4})?)$`,
			input:   "99999*1234",
			trimset: " ,",
			delset:  "- ",
			match:   false,
			named:   map[string]string{"zip": "999991234"},
			all:     []string{"999991234", "999991234"},
		}),
	)
})

// TStruct contains a test and the Ωed result
type ReTStruct struct {
	rexp    string
	input   string
	trimset string
	delset  string
	match   bool
	named   map[string]string
	all     []string
}

func (r ReTStruct) String() string {
	return fmt.Sprintf("   rexp: %q\n   input: %q\n   trimset: %q\n   match: %t\n   named: %#v\n   all: %#v\n", r.rexp, r.input, r.trimset, r.match, r.named, r.all)
}
