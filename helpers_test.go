package raws

import (
	"reflect"
	"testing"
)

func checkErrors(t *testing.T, name string, index int, err error, expected error) {
	t.Helper()

	if err == nil && expected != nil {
		if index >= 0 {
			t.Errorf("%s [%d] - errors: received=nil | expected=%v", name, index, expected)
		} else {
			t.Errorf("%s - errors: received=nil | expected=%v", name, expected)
		}

		return
	}

	if !reflect.DeepEqual(err, expected) {
		if index >= 0 {
			t.Errorf("%s [%d] - errors: received=%+v | expected=%+v", name, index, err, expected)
		} else {
			t.Errorf("%s - errors: received=%+v | expected=%+v", name, index, expected)
		}
	}
}

func checkError(t *testing.T, err error, expected error) {
	t.Helper()
	checkErrors(t, t.Name(), -1, err, expected)
}
