package jx

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// XJFloat64 is a float64 type providing safer unmarshalling of XML or JSON.
// It will return 0.0 if the value in the XML or JSON is missing or empty.
// JSON will accept either a number or string value.
// If the JSON/XML value is non-numeric, and/or if it cannot be converted to a
// numeric value using the standard strconv function, it will return an error.
type XJFloat64 float64

// UnmarshalXML replaces the standard XML unmarshalling method.
// It will return 0.0 if the value in the XML is missing or empty.
func (r *XJFloat64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	value, err := decodeElement(d, start)
	if err != nil {
		return err
	}
	return r.set(value)
}

func (r *XJFloat64) UnmarshalText(text []byte) error {
	return r.set(string(text))
}

// UnmarshalJSON replaces the standard JSON unmarshalling method.
// It will return 0.0 if the value in the JSON is missing or empty.
func (r *XJFloat64) UnmarshalJSON(value []byte) error {
	return r.set(string(value))
}

func (r *XJFloat64) set(value string) error {
	x, err := parseVal(value, rtFloat64)
	*r = XJFloat64(x.(float64))
	return err
}

// XJInt is a int type providing safer unmarshalling of XML or JSON.
// It will return 0 if the value in the XML or JSON is missing or empty.
// If the JSON/XML value is non-numeric, it will return an error.
// JSON will accept either a number or string value.
// If the JSON/XML value is non-numeric, and/or if it cannot be converted to a
// numeric value using the standard strconv function, it will return an error.
type XJInt int

// UnmarshalXML replaces the standard XML unmarshalling method.  It will return 0.0 if the value in the XML is missing or empty.
func (r *XJInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	value, err := decodeElement(d, start)
	if err != nil {
		return err
	}
	return r.set(value)
}

// UnmarshalJSON replaces the standard JSON unmarshalling method.  It will return 0.0 if the value in the JSON is missing or empty.
func (r *XJInt) UnmarshalJSON(value []byte) error {
	return r.set(string(value))
}

func (r *XJInt) set(value string) error {
	x, err := parseVal(value, rtInt)
	*r = XJInt(x.(int))
	return err
}

type nType int

const (
	rtFloat64 nType = iota
	rtInt
)

func trimNumber(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	}
	switch r {
	case '"', '\'':
		return true
	}
	return false
}

func parseVal(value string, outtype nType) (interface{}, error) {
	cvalue := strings.TrimFunc(value, trimNumber)

	if cvalue == "" {
		cvalue = "0"
	}

	switch outtype {
	case rtFloat64:
		var x float64
		x, err := strconv.ParseFloat(cvalue, 64)
		return x, err
	case rtInt:
		var x int
		x, err := strconv.Atoi(cvalue)
		return x, err
	default:
		return nil, fmt.Errorf("parseVal received an invalid outtype: %v", outtype)
	}
}

func decodeElement(d *xml.Decoder, start xml.StartElement) (value string, err error) {
	err = d.DecodeElement(&value, &start)
	return value, err
}
