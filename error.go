package raws

import (
	"fmt"
	"strings"
)

// Error is a type which satisfied the standard error interface, but provides
// context over an error that the AWS SDK can originate.
type Error struct {
	err     error
	region  string
	service string
}

// NewError creates an Error object for the specific AWS region, service and
// containing the original error returned by the AWS SDK.
func NewError(region string, service string, e error) Error {
	return Error{
		region:  region,
		service: service,
		err:     e,
	}
}

// Error satisfies the error interface and returns a string containing the
// region, service and the message of the  AWS SDK error.
func (e Error) Error() string {
	return fmt.Sprintf("region: %s, service: %s, Error message: %q",
		e.region,
		e.service,
		e.err.Error())
}

// Region returns the region of the error.
func (e Error) Region() string {
	return e.region
}

// Service Returns the service name of the error.
// NOTE, currently the service is not that necessary, but it could become useful
// to have as the project evolves and start making more complex calls to various
// endpoints.
func (e Error) Service() string {
	return e.service
}

// Errors type satisfies the standard error interface, thus allowing us to return
// an error when doing multiple call via the Go AWS SDK, even though multiple errors
// are met.
type Errors []Error

// Error returns a string which summarize how many errors happened and for
// each error, the region, the service and  the error message reported by
// AWS original error.
func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var details []string

	for _, err := range e {
		details = append(details, err.Error())
	}

	return fmt.Sprintf("%d error(s) occurred.\n\t%s", len(e), strings.Join(details, "\n\t"))
}
