package raws

import (
	"reflect"
	"testing"
)

func checkErrors(t *testing.T, name string, index int, err error, expected error) {
	t.Helper()

	if err == nil && expected != nil {
		t.Errorf("%s [%d] - errors: received=nil | expected=%v", name, index, expected)
		return
	}

	if !reflect.DeepEqual(err, expected) {
		t.Errorf("%s [%d] - errors: received=%+v | expected=%+v", name, index, err, expected)
	}
}
