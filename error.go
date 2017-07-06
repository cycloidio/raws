package raws

import (
	"fmt"
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

// Error rerturns a string which summarize how many errors happened, in which regions and for which services.
func (e Errs) Error() string {
	var output [][]string

	for _, err := range e {
		output = append(output, []string{err.Region(), err.Service()})
	}
	return fmt.Sprintf("%d error(s) occured: %s", len(e), output)
}

// Error returns a string containing the region, service as well as the original API error message.
func (e *callErr) Error() string {
	return fmt.Sprintf("%s: error while using '%s' service - %s",
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
