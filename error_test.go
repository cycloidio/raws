package raws

import (
	"errors"
	"testing"
)

func TestErrs_Error(t *testing.T) {
	tests := []struct {
		name           string
		input          Errs
		expectedOutput string
	}{{name: "empty error",
		input:          Errs{},
		expectedOutput: "",
	},
		{name: "one error",
			input: Errs{
				&callErr{
					err:     errors.New("fail-1"),
					region:  "region-1",
					service: "service-1",
				},
			},
			expectedOutput: "1 error(s) occurred.\n\t" +
				`region: region-1, service: service-1, Error message: "fail-1"`,
		},
		{name: "two errors",
			input: Errs{
				&callErr{
					err:     errors.New("fail-1"),
					region:  "region-1",
					service: "service-1",
				},
				&callErr{
					err:     errors.New("fail-2"),
					region:  "region-2",
					service: "service-2",
				},
			},
			expectedOutput: "2 error(s) occurred.\n\t" +
				"region: region-1, service: service-1, Error message: \"fail-1\"\n\t" +
				`region: region-2, service: service-2, Error message: "fail-2"`,
		}}

	for i, tt := range tests {
		if tt.input.Error() != tt.expectedOutput {
			t.Errorf("%s [%d] - Errs output: received=%+v | expected=%+v",
				tt.name, i, tt.input.Error(), tt.expectedOutput)
		}
	}

}

func TestCallErr_Error(t *testing.T) {
	tests := []struct {
		name           string
		input          callErr
		expectedOutput string
	}{{name: "one error",
		input: callErr{
			err:     errors.New("fail-1"),
			region:  "region-1",
			service: "service-1",
		},
		expectedOutput: `region: region-1, service: service-1, Error message: "fail-1"`,
	},
		{name: "another error",
			input: callErr{
				err:     errors.New("fail-2"),
				region:  "region-2",
				service: "service-2",
			},
			expectedOutput: `region: region-2, service: service-2, Error message: "fail-2"`,
		}}

	for i, tt := range tests {
		if tt.input.Error() != tt.expectedOutput {
			t.Errorf("%s [%d] - callErrs output: received=%+v | expected=%+v",
				tt.name, i, tt.input.Error(), tt.expectedOutput)
		}
	}
}

func TestCallErr_Region(t *testing.T) {
	tests := []struct {
		name           string
		input          callErr
		expectedRegion string
	}{{name: "one error with one region",
		input: callErr{
			err:     errors.New("fail-1"),
			region:  "region-1",
			service: "service-1",
		},
		expectedRegion: "region-1",
	},
		{name: "one error with another region",
			input: callErr{
				err:     errors.New("fail-2"),
				region:  "region-2",
				service: "service-2",
			},
			expectedRegion: "region-2",
		}}

	for i, tt := range tests {
		if tt.input.Region() != tt.expectedRegion {
			t.Errorf("%s [%d] - callErrs region: received=%+v | expected=%+v",
				tt.name, i, tt.input.Region(), tt.expectedRegion)
		}
	}
}

func TestCallErr_Service(t *testing.T) {
	tests := []struct {
		name            string
		input           callErr
		expectedService string
	}{{name: "one error with one service",
		input: callErr{
			err:     errors.New("fail-1"),
			region:  "region-1",
			service: "service-1",
		},
		expectedService: "service-1",
	},
		{name: "one error with another service",
			input: callErr{
				err:     errors.New("fail-2"),
				region:  "region-1",
				service: "service-2",
			},
			expectedService: "service-2",
		}}

	for i, tt := range tests {
		if tt.input.Service() != tt.expectedService {
			t.Errorf("%s [%d] - callErrs service: received=%+v | expected=%+v",
				tt.name, i, tt.input.Service(), tt.expectedService)
		}
	}
}
