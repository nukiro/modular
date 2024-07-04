package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nukiro/modular/internal/tests"
)

func TestCheckCode(t *testing.T) {
	tests := []struct {
		code int
		want tests.Expect
		msg  string
	}{
		{200, tests.ExpectNil, "did panic"},
		{-1, tests.ExpectNotNil, "did not panic"},
		{900, tests.ExpectNotNil, "did not panic"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%d status code", tt.code)
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); tt.want(r) {
					t.Errorf("check status code %d %s", tt.code, tt.msg)
				}
			}()

			checkCode(tt.code)
		})
	}

	t.Run("panic message", func(t *testing.T) {
		defer func() {
			r := recover()
			s := fmt.Sprint(r)
			if s != "response code 0 is unknown" {
				t.Errorf("got %q panic, but want %q", s, "response code 0 is unknown")
			}
		}()

		checkCode(0)
	})
}

func TestNew(t *testing.T) {
	// Faking check code
	saved := checkCode
	defer func() { checkCode = saved }()
	checkCode = func(c int) {}

	r := New(100101)

	if len(r.Header) != 0 {
		t.Errorf("headers was not empty")
	}

	if r.Code != 100101 {
		t.Errorf("got code %d, but want %d", r.Code, 200)
	}
}

func TestJSON(t *testing.T) {
	t.Run("without errors", func(t *testing.T) {
		// Save and restore original write function
		saved := write
		defer func() { write = saved }()
		// Fake write function
		write = func(w http.ResponseWriter, f Serializer, r *Response) error {
			return nil
		}

		w := httptest.NewRecorder()
		r := New(200)

		rw := r.JSON(w, "body")

		if rw.Code != 200 {
			t.Errorf("response status code %q, want %q", rw.Code, 200)
		}

		if rw.Body != "body" {
			t.Errorf("response body %q, want %q", rw.Body, "body")
		}
	})

	t.Run("with errors", func(t *testing.T) {
		// Save and restore original write function
		saved := write
		defer func() { write = saved }()
		// Fake write function
		write = func(w http.ResponseWriter, f Serializer, r *Response) error {
			return errors.New("an error")
		}

		w := httptest.NewRecorder()
		r := New(200)

		defer func() {
			tests.AssertPanic(t, recover(), "JSON", "an error")

		}()

		r.JSON(w, json.MarshalIndent)
	})
}
