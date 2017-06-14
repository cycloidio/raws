package core

import (
	"fmt"
)

type Err interface {
	error
	Region() string
	Service() string
}

type RawsErr struct {
	APIErrs []callErr
}

func NewAPIError(e error, region string, service string) Err {
	return &callErr{
		err:     e,
		region:  region,
		service: service,
	}
}

func (r *RawsErr) AppendError(regionErr string, serviceErr string, originErr error) {
	if originErr != nil {
		r.APIErrs = append(r.APIErrs, callErr{
			region:  regionErr,
			service: serviceErr,
			err:     originErr,
		})
	}
}

func (r RawsErr) Error() string {
	var output [][]string

	for _, callErr := range r.APIErrs {
		output = append(output, []string{callErr.region, callErr.service})
	}
	return fmt.Sprintf("%d error(s) occured: %s", len(r.APIErrs), output)
}

func (e *callErr) Error() string {
	return fmt.Sprintf("%s: error while using '%s' service - %s",
		e.region,
		e.service,
		e.err.Error())
}

func (e *callErr) Region() string {
	return e.region
}

func (e *callErr) Service() string {
	return e.service
}

type callErr struct {
	err     error
	region  string
	service string
}
