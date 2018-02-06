package raws

import (
	"fmt"
	"strings"
)

// Interface defining an API error, this is handy as the region and service are saved
//
// Note: Currently the service is not that necessary, but it could become useful to have as the project evolves and
// start making more complex calls to various endpoints.
type Err interface {
	error
	Region() string
	Service() string
}

// RawsErr type satisfies the standard error interface, thus allowing us to return an error when doing multiple call via
// the go-SDK, even though multiple errors are met; which is why APIErrs save those more specific errors.
type Errs []Err

// NewAPIError returns an Err object filled with the error faced, the region and the service name.
func NewAPIError(region string, service string, e error) Err {
	return &callErr{
		region:  region,
		service: service,
		err:     e,
	}
}

// Error returns a string which summarize how many errors happened and for
// each error, the region, the service and  the error message reported by
// AWS original error.
func (e Errs) Error() string {
	if len(e) == 0 {
		return ""
	}

	var details []string

	for _, err := range e {
		details = append(details, err.Error())
	}

	return fmt.Sprintf("%d error(s) occurred.\n\t%s", len(e), strings.Join(details, "\n\t"))
}

// Error returns a string containing the region, service as well as the original API error message.
func (e *callErr) Error() string {
	return fmt.Sprintf("region: %s, service: %s, Error message: %q",
		e.region,
		e.service,
		e.err.Error())
}

// Returns the region of the error
func (e *callErr) Region() string {
	return e.region
}

// Returns the service name of the error
func (e *callErr) Service() string {
	return e.service
}

type callErr struct {
	Err
	err     error
	region  string
	service string
}
