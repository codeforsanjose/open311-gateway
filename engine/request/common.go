package request

import (
	"Gateway311/engine/common"
	"fmt"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

// ==============================================================================================================================
//                                      COMMON
// ==============================================================================================================================
type cType struct {
	self cIface
	cRType
}

func (r *cType) load(p cIface, rqstID int64, rqst *rest.Request) error {
	r.self = p
	r.id = rqstID

	if err := rqst.DecodeJsonPayload(r.self); err != nil {
		if err.Error() != "JSON payload is empty" {
			return fmt.Errorf("Unable to process request - %s", err)
		}
	}
	if err := r.self.parseQP(rqst); err != nil {
		return fmt.Errorf("Unable to process request - %s", err)
	}

	if err := r.self.validate(); err != nil {
		return fmt.Errorf("Unable to process request - %s", err)
	}

	return nil
}

type cIface interface {
	parseQP(r *rest.Request) error
	validate() error
}

type cRType struct {
	id int64
}

func (r *cRType) SetID(id int64) {
	r.id = id
}

func (r *cRType) GetID() int64 {
	return r.id
}

// ==============================================================================================================================
//                                      VALIDATION
// ==============================================================================================================================
const (
	// 48 contiguous states.
	latMin float64 = 18.0
	latMax float64 = 49.0
	lngMin float64 = -124.6
	lngMax float64 = -62.3
)

func validateLatLng(lat, lng float64) bool {
	if lat >= latMin && lat <= latMax && lng >= lngMin && lng <= lngMax {
		return true
	}
	return false
}

// ------------------ Validation System

// ValidationDetail is a simple method for compiling validation results.
type ValidationDetail struct {
	ok     bool
	result string
}

// Validation is a simple method for compiling validation results.
type Validation map[string]*ValidationDetail

func newValidation() Validation {
	return make(map[string]*ValidationDetail)
}

// Set creates a validation as ok (true) or not (false).
func (r Validation) Set(item, result string, isOK bool) {
	v, ok := r[item]
	if ok {
		v.ok = isOK
		if result > "" {
			v.result = result
		}
	} else {
		r[item] = &ValidationDetail{
			ok:     isOK,
			result: result,
		}
	}
}

// IsOK returns the state of the requested Validation.  If the Validation has
// not been set, it will return FALSE.
func (r Validation) IsOK(item string) bool {
	v, ok := r[item]
	if !ok {
		return false
	}
	return v.ok
}

// Ok scans all validations - if all are true (i.e. they passed that validation
// test), then it returns true.
func (r Validation) Ok() bool {
	for _, v := range r {
		if !v.ok {
			return false
		}
	}
	return true
}

// String returns a string representation of the validation entries.
func (r Validation) String() string {
	ls := new(common.LogString)
	ls.AddF("Validation\n")
	ls.AddS("-Item-         -Valid-  -Reason-\n")
	for k, v := range r {
		ls.AddF("%-15s %-5t  %-90.90s\n", k, v.ok, v.result)
	}
	return ls.Box(110)
}

// Error is a standard error interface, returning a string listing any failed
// validations.
func (r Validation) Error() string {
	validMsg := ""
	for k, v := range r {
		if !v.ok {
			if validMsg == "" {
				validMsg = k
			} else {
				validMsg = validMsg + ", " + k
			}
		}
	}
	if validMsg != "" {
		return fmt.Sprintf("errors: %s", validMsg)
	}
	return ""
}

// ==============================================================================================================================
//                                      CONVERSION
// ==============================================================================================================================

// ----------------------  Conversion Parameters
type vparm struct {
	vtype    string
	required bool
	dflt     string
}

// vparms is initialized in init below.
var vparms map[string]vparm

// ----------------------  Convert

func newConversion() conversion {
	var c conversion
	c.Validation = make(map[string]*ValidationDetail)
	return c
}

type conversion struct {
	Validation
}

func (r *conversion) float(name, val string) float64 {
	return *r.convert(name, val).(*float64)
}

func (r *conversion) int(name, val string) int {
	return *r.convert(name, val).(*int)
}

func (r *conversion) bool(name, val string) bool {
	return *r.convert(name, val).(*bool)
}

func (r conversion) String() string {
	return r.Validation.String()
}

func (r *conversion) convert(name, val string) interface{} {
	log.Debug("Converting %q val: %q", name, val)
	vp, vpOK := vparms[name]
	fail := func(item, msg string) interface{} {
		log.Debug("FAIL: %s", msg)
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
		return fail(name, fmt.Sprintf("unknown conversion type: %q", vp.vtype))
	}
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
