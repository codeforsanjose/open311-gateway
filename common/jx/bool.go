package jx

import (
	"encoding/xml"
	"strings"

	"github.com/codeforsanjose/open311-gateway/_background/go/common"
)

const (
	jsonTrimChars = `"' `

	sYES = "YES"
	sYes = "Yes"
	qYes = `"Yes"`

	sNo = "No"
	qNo = `"No"`

	sTRUE = "TRUE"
	sTrue = "True"
	qTrue = `"True"`

	sFalse = "False"
	qFalse = `"False"`
)

// ==============================================================================================================================
//                                      BoolYN Type
// ==============================================================================================================================

// BoolYNType is a boolean type which supports JSON and XML marshaling/unmarshalling
// to "sYes" / "sNo" values.
type BoolYNType bool

// MarshalJSON returns a JSON representation of the boolean.
func (b BoolYNType) MarshalJSON() ([]byte, error) {
	if b {
		return []byte(qYes), nil
	}
	return []byte(qNo), nil
}

// UnmarshalJSON unmarshals the JSON representation ("sYes" or "sNo") into a boolean
// value.
func (b *BoolYNType) UnmarshalJSON(data []byte) error {
	// sNOTE: the incoming JSON data INCLUDES the quotation marks!
	switch strings.Trim(strings.ToUpper(common.ByteToString(data, 0)), jsonTrimChars) {
	case sYES, sTRUE:
		*b = true
	default:
		*b = false
	}
	return nil
}

// MarshalXML returns an XML representation of the boolean.
func (b BoolYNType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s := sNo
	if b {
		s = sYes
	}
	return e.EncodeElement(s, start)
}

// MarshalXMLAttr returns an XML representation of the boolean.
func (b BoolYNType) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if b {
		return xml.Attr{Name: name, Value: sYes}, nil
	}
	return xml.Attr{Name: name, Value: sNo}, nil
}

// UnmarshalXML unmarshals the XML representation ("sYes" or "sNo") into a boolean
// value.
func (b *BoolYNType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	_ = d.DecodeElement(&v, &start)
	switch strings.ToUpper(v) {
	case sYES, sTRUE:
		*b = true
	default:
		*b = false
	}
	return nil
}

// ==============================================================================================================================
//                                      BoolTF Type
// ==============================================================================================================================

// BoolTFType is a boolean type which supports JSON and XML marshaling/unmarshalling
// to "sYes" / "sNo" values.
type BoolTFType bool

// MarshalJSON returns a JSON representation of the boolean.
func (b BoolTFType) MarshalJSON() ([]byte, error) {
	if b {
		return []byte(qTrue), nil
	}
	return []byte(qFalse), nil
}

// UnmarshalJSON unmarshals the JSON representation ("sYes" or "sNo") into a boolean
// value.
func (b *BoolTFType) UnmarshalJSON(data []byte) error {
	// sNOTE: the incoming JSON data INCLUDES the quotation marks!
	switch strings.Trim(strings.ToUpper(common.ByteToString(data, 0)), jsonTrimChars) {
	case sYES, sTRUE:
		*b = true
	default:
		*b = false
	}
	return nil
}

// MarshalXML returns an XML representation of the boolean.
func (b BoolTFType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s := sFalse
	if b {
		s = sTrue
	}
	return e.EncodeElement(s, start)
}

// MarshalXMLAttr returns an XML representation of the boolean.
func (b BoolTFType) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if b {
		return xml.Attr{Name: name, Value: sTrue}, nil
	}
	return xml.Attr{Name: name, Value: sFalse}, nil
}

// UnmarshalXML unmarshals the XML representation ("sYes" or "sNo") into a boolean
// value.
func (b *BoolTFType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	_ = d.DecodeElement(&v, &start)
	switch strings.ToUpper(v) {
	case sYES, sTRUE:
		*b = true
	default:
		*b = false
	}
	return nil
}

// ==============================================================================================================================
//                                      MISC
// ==============================================================================================================================

// BoolToStringTF converts a boolean value to "True" or "False" (note capitilization!).
func BoolToStringTF(v bool) string {
	if v {
		return "True"
	}
	return "False"
}
