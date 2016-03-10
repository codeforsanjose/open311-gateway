package common

import (
	"fmt"
	"strconv"
)

// NewConversion returns a conversion set.  A Conversion set will log all errors captured
// during one or more conversions.  All possible conversions are specifed below in init().
func NewConversion() Conversion {
	var c Conversion
	c.Validation = make(map[string]*ValidationDetail)
	return c
}

// ==============================================================================================================================
//                                      INIT
// ==============================================================================================================================

func init() {
	vparms = make(map[string]vparm)

	vparms["Latitude"] = vparm{"float", true, "0.0"}
	vparms["Longitude"] = vparm{"float", true, "0.0"}
	vparms["IsAnonymous"] = vparm{"bool", false, "true"}
	vparms["Radius"] = vparm{"int", false, "100"}
	vparms["MaxResults"] = vparm{"int", false, "10"}

	vparms["IncludeDetails"] = vparm{"bool", false, "false"}
	vparms["IncludeComments"] = vparm{"bool", false, "false"}
	vparms["IncludeVotes"] = vparm{"bool", false, "false"}
}

// ==============================================================================================================================
//                                      Conversion
// ==============================================================================================================================

// vparm represents a
type vparm struct {
	vtype    string
	required bool
	dflt     string
}

// vparms is initialized in init below.
var vparms map[string]vparm

// Conversion represents a Conversion.  It is composed of a Validation.
type Conversion struct {
	Validation
}

// Float validates a float target as per the vparms for the specified name.
func (r *Conversion) Float(name, val string) float64 {
	return *r.convert(name, val).(*float64)
}

// Int validates a integer target as per the vparms for the specified name.
func (r *Conversion) Int(name, val string) int {
	return *r.convert(name, val).(*int)
}

// Bool validates a boolean target as per the vparms for the specified name.
func (r *Conversion) Bool(name, val string) bool {
	return *r.convert(name, val).(*bool)
}

func (r Conversion) String() string {
	return r.Validation.String()
}

func (r *Conversion) convert(name, val string) interface{} {
	// log.Debug("Converting %q val: %q", name, val)
	vp, vpOK := vparms[name]
	fail := func(item, msg string) interface{} {
		r.Validation[item] = &ValidationDetail{ok: false, result: msg}
		if !vpOK {
			return nil
		}
		switch vp.vtype {
		case "float":
			retval := 0.0
			return &retval
		case "int":
			retval := 0
			return &retval
		case "bool":
			retval := false
			return &retval
		default:
			return nil
		}
	}
	if !vpOK {
		return fail(name, "unknown field")
	}

	if val == "" {
		if vp.required {
			return fail(name, "empty")
		}
		val = vp.dflt
	}
	switch vp.vtype {
	case "float":
		out, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fail(name, err.Error())
		}
		r.Validation[name] = &ValidationDetail{ok: true}
		return &out
	case "int":
		out, err := strconv.Atoi(val)
		if err != nil {
			return fail(name, err.Error())
		}
		r.Validation[name] = &ValidationDetail{ok: true}
		return &out
	case "bool":
		out, err := strconv.ParseBool(val)
		if err != nil {
			return fail(name, err.Error())
			// return fail(name, fmt.Sprintf("%s trying to convert: [%s] to float\n", err, val))
		}
		r.Validation[name] = &ValidationDetail{ok: true}
		return &out
	default:
		return fail(name, fmt.Sprintf("unknown Conversion type: %q", vp.vtype))
	}
}
