package raws

import (
	"errors"
	"testing"
)

func TestErrors_Error(t *testing.T) {
	tests := []struct {
		name           string
		input          Errors
		expectedOutput string
	}{{name: "empty error",
		input:          Errors{},
		expectedOutput: "",
	},
		{name: "one error",
			input: Errors{
				Error{
					err:     errors.New("fail-1"),
					region:  "region-1",
					service: "service-1",
				},
			},
			expectedOutput: "1 error(s) occurred.\n\t" +
				`region: region-1, service: service-1, Error message: "fail-1"`,
		},
		{name: "two errors",
			input: Errors{
				Error{
					err:     errors.New("fail-1"),
					region:  "region-1",
					service: "service-1",
				},
				Error{
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

func TestError_Error(t *testing.T) {
	tests := []struct {
		name           string
		input          Error
		expectedOutput string
	}{{name: "one error",
		input: Error{
			err:     errors.New("fail-1"),
			region:  "region-1",
			service: "service-1",
		},
		expectedOutput: `region: region-1, service: service-1, Error message: "fail-1"`,
	},
		{name: "another error",
			input: Error{
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

func TestError_Region(t *testing.T) {
	tests := []struct {
		name           string
		input          Error
		expectedRegion string
	}{{name: "one error with one region",
		input: Error{
			err:     errors.New("fail-1"),
			region:  "region-1",
			service: "service-1",
		},
		expectedRegion: "region-1",
	},
		{name: "one error with another region",
			input: Error{
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

func TestError_Service(t *testing.T) {
	tests := []struct {
		name            string
		input           Error
		expectedService string
	}{{name: "one error with one service",
		input: Error{
			err:     errors.New("fail-1"),
			region:  "region-1",
			service: "service-1",
		},
		expectedService: "service-1",
	},
		{name: "one error with another service",
			input: Error{
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

func TestErrorIn(t *testing.T) {
	t.Run("when the input error is nil", func(t *testing.T) {
		var err = ErrorIn("some-region", nil)
		if err != nil {
			t.Errorf("unexpected result. received=%+v | expected=nil", err)
		}
	})

	t.Run("when the input error is of the type Error", func(t *testing.T) {
		var inputErr = NewError("some-region", "some-service", errors.New("some error"))

		var err = ErrorIn("other-region", inputErr)
		if err != nil {
			t.Errorf("unexpected result. received=%+v | expected=nil", err)
		}

		err = ErrorIn("some-region", inputErr)
		if err == nil || err != inputErr {
			t.Errorf("unexpected result. received=%+v | expected=%+v", err, inputErr)
		}
	})

	t.Run("when the input error is of the type *Error", func(t *testing.T) {
		var inputErr *Error

		{
			var err = NewError("some-region", "some-service", errors.New("some error"))
			inputErr = &err
		}

		var err = ErrorIn("other-region", inputErr)
		if err != nil {
			t.Errorf("unexpected result. received=%+v | expected=nil", err)
		}

		err = ErrorIn("some-region", inputErr)
		if err == nil || err != inputErr {
			t.Errorf("unexpected result. received=%+v | expected=%+v", err, inputErr)
		}

		inputErr = nil
		err = ErrorIn("other-region", inputErr)
		if err != nil {
			t.Errorf("unexpected result. received=%+v | expected=nil", err)
		}
	})

	t.Run("when the input error is of the type Errors", func(t *testing.T) {
		var inputErr = Errors{
			NewError("region-1", "service-1", errors.New("some error")),
			NewError("region-2", "service-2", errors.New("some error")),
			NewError("region-3", "service-3", errors.New("some error")),
		}

		for _, ie := range inputErr {
			var err = ErrorIn(ie.Region(), inputErr)
			if err == nil {
				t.Errorf("unexpected result. received=nil | expected=%+v", Errors{ie})
				continue
			}

			var errs, ok = err.(Errors)
			if !ok {
				t.Errorf("unexpected return type. received=%T | expected=%T", errs, Errors{})
				continue
			}

			if len(errs) != 1 || errs[0] != ie {
				t.Errorf("unexpected result. received=%+v | expected=%+v", errs, Errors{ie})
			}
		}

		inputErr = append(inputErr, NewError("region-3", "service-3.1", errors.New("another error")))
		var err = ErrorIn("region-3", inputErr)
		{
			if err == nil {
				t.Errorf("unexpected result. received=nil | expected=%+v", Errors{inputErr[2], inputErr[3]})
				goto nextTest
			}

			var errs, ok = err.(Errors)
			if !ok {
				t.Errorf("unexpected return type. received=%T | expected=%T", errs, Errors{})
				goto nextTest
			}

			if len(errs) != 2 || errs[0] != inputErr[2] || errs[1] != inputErr[3] {
				t.Errorf(
					"unexpected result. received=%+v | expected=%+v",
					errs,
					Errors{inputErr[2], inputErr[3]},
				)
			}
		}

	nextTest:
		err = ErrorIn("region-5", inputErr)
		if err != nil {
			t.Errorf("unexpected result. received=%+v | expected=nil", err)
		}

		err = ErrorIn("region-1", Errors{})
		if err != nil {
			t.Errorf("unexpected result. received=%+v | expected=nil", err)
		}

		inputErr = nil
		err = ErrorIn("region-1", inputErr)
		if err != nil {
			t.Errorf("unexpected result. received=%+v | expected=nil", err)
		}
	})

	t.Run("when the input error is not of the type Error nor Errors", func(t *testing.T) {
		var inputErr = errors.New("some error")

		var err = ErrorIn("some-region", inputErr)
		if err != nil {
			t.Errorf("unexpected result. received=%+v | expected=nil", err)
		}

		inputErr = nil
		err = ErrorIn("some-region", inputErr)
		if err != nil {
			t.Errorf("unexpected result. received=%+v | expected=nil", err)
		}
	})
}
