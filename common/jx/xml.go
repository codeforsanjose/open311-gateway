package jx

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/codeforsanjose/open311-gateway/_background/go/common/mybuf"
)

// EncodeXML encodes the input struct as XML, and returns the result as a pointer
// *bytes.Buffer.
func EncodeXML(source interface{}, indent, header bool) (*bytes.Buffer, error) {
	var bout = new(bytes.Buffer)
	enc := xml.NewEncoder(bout)
	if indent {
		enc.Indent("  ", "    ")
	}
	if err := enc.Encode(source); err != nil {
		return nil, err
	}

	if header {
		b := bytes.NewBufferString(xml.Header)
		_, e := mybuf.Concat(b, bout)
		if e != nil {
			return nil, e
		}
		return b, nil
	}
	return bout, nil
}

// EncodeXMLByte encodes the input struct as XML, and returns the result as a slice of bytes.
func EncodeXMLByte(source interface{}, indent, header bool) ([]byte, error) {
	buf, err := EncodeXML(source, indent, header)
	if err != nil {
		return nil, err
	}
	return mybuf.ToBSlice(buf), nil
}

// EncodeXMLString encodes the input struct as XML, and returns the result as a slice of bytes.
func EncodeXMLString(source interface{}, indent, header bool) (string, error) {
	buf, err := EncodeXML(source, indent, header)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// LoadXMLFile loads the specified file into the target struct.
func LoadXMLFile(file string, target interface{}) error {

	xmlFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer xmlFile.Close()

	if err := xml.NewDecoder(xmlFile).Decode(target); err != nil {
		return fmt.Errorf("unable to load xml: %s", err)
	}

	return nil
}
