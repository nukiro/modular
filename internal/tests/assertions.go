package tests

import (
	"fmt"
	"testing"
)

func AssertPanicNilParam(t testing.TB, r any, method, param string) {
	want := fmt.Sprintf("%s param cannot be nil", param)
	if r == nil {
		t.Errorf("%s did not panic", method)
	}
	if r != nil {
		s := fmt.Sprint(r)
		if s != want {
			t.Errorf("got %q panic, but want %q", s, want)
		}
	}
}
