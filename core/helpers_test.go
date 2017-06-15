package core

import (
	"testing"
	"reflect"
)

func checkErrors(t *testing.T, name string, index int, err error, expected error) {
	if err != nil && !reflect.DeepEqual(err, expected) {
		t.Errorf("%s [%d] - errors: received=%+v | expected=%+v",
			name, index, err.Error(), expected.Error())
	}
}
