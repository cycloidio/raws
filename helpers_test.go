package raws

import (
	"reflect"
	"testing"
)

func checkErrors(t *testing.T, name string, index int, err error, expected error) {
	t.Helper()

	if err != nil && !reflect.DeepEqual(err, expected) {
		t.Errorf("%s [%d] - errors: received=%+v | expected=%+v",
			name, index, err.Error(), expected.Error())
	}
}
