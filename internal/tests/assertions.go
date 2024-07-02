package tests

import (
	"fmt"
	"testing"
)

func AssertPanic(t testing.TB, r any, method, msg string) {
	t.Helper()
	if r == nil {
		t.Errorf("%s did not panic", method)
	}
	if r != nil {
		s := fmt.Sprint(r)
		if s != msg {
			t.Errorf("got %q panic error, but want %q", s, msg)
		}
	}
}

func AssertPanicNilParam(t testing.TB, r any, method, param string) {
	AssertPanic(t, r, method, fmt.Sprintf("%s param cannot be nil", param))
}

func AssertPanicEmptyParam(t testing.TB, r any, method, param string) {
	AssertPanic(t, r, method, fmt.Sprintf("%s param cannot be empty", param))
}
