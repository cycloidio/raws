package core

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
type RawsErr struct {
	APIErrs []callErr
}

// NewAPIError returns an Err object filled with the error faced, the region and the service name.
func NewAPIError(e error, region string, service string) Err {
	return &callErr{
		err:     e,
		region:  region,
		service: service,
	}
}

// AppendError is to add errors to the list in an easy way. If there was indeed an error, it is added to the list of
// APIErrs of the RawsErr struct.
func (r *RawsErr) AppendError(regionErr string, serviceErr string, originErr error) {
	if originErr != nil {
		r.APIErrs = append(r.APIErrs, callErr{
			region:  regionErr,
			service: serviceErr,
			err:     originErr,
		})
	}
}

// Error rerturns a string which summarize how many errors happened, in which regions and for which services.
func (r RawsErr) Error() string {
	var output [][]string

	for _, callErr := range r.APIErrs {
		output = append(output, []string{callErr.region, callErr.service})
	}
	return fmt.Sprintf("%d error(s) occured: %s", len(r.APIErrs), output)
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
	err     error
	region  string
	service string
}
