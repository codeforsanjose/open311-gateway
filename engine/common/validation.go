package common

import (
	"fmt"
)

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

func ValidateLatLng(lat, lng float64) bool {
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

func NewValidation() Validation {
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
	ls := new(LogString)
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
