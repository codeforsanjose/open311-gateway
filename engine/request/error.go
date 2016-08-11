package request

import "github.com/open311-gateway/engine/common"

// ErrorsResponseJ represents an error response.  The error response contains one or more errors.
type ErrorsResponseJ []*ErrorResponseJ

// ErrorResponseJ represents an in individual error in an error response.
type ErrorResponseJ struct {
	Code        int    `json:"code" xml:"code"`
	Description string `json:"description" xml:"description"`
}

func newErrorsResponseJ() ErrorsResponseJ {
	return ErrorsResponseJ{}
}

// errorJ creates a new error instance
func (r ErrorsResponseJ) errorJ(code int, descr string) ErrorsResponseJ {
	r = append(r, &ErrorResponseJ{
		Code:        code,
		Description: descr,
	})
	return r
}

// String displays the contents of the CreateRequest type.
func (r ErrorsResponseJ) String() string {
	ls := new(common.LogString)
	ls.AddF("ErrorResponse\n")
	for _, v := range r {
		ls.AddF("%-5v  %v\n", v.Code, v.Description)
	}
	return ls.Box(80)
}
