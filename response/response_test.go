package response

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/nukiro/modular/internal/tests"
)

func TestCheckStatusCode(t *testing.T) {
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

			checkStatusCode(tt.code)
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

		checkStatusCode(0)
	})
}

func TestNew(t *testing.T) {
	t.Run("good response", func(t *testing.T) {
		r := new(200, success, "message", "this is the message")

		if len(r.header) != 0 {
			t.Errorf("headers was not empty")
		}

		if r.statusCode != 200 {
			t.Errorf("got code %d, but want %d", r.statusCode, 200)
		}

		if s, ok := r.body["time"]; ok {
			x := assertPayloadKeyFormat[int64](t, "time", s)

			y := time.Unix(int64(x), 0)
			z := time.Now()

			if y.After(z) {
				t.Errorf("payload time is for the future")
			}
		} else {
			t.Errorf("payload does not contain time key")
		}

		if s, ok := r.body["status"]; ok {
			x := assertPayloadKeyFormat[string](t, "status", s)

			if x != "ok" {
				t.Errorf("got payload status %q, but want %q", x, "ok")
			}
		} else {
			t.Errorf("payload does not contain status key")
		}

		if s, ok := r.body["result"]; ok {
			x := assertPayloadKeyFormat[result](t, "result", s)

			if x != "success" {
				t.Errorf("got payload result %q, but want %q", x, "success")
			}
		} else {
			t.Errorf("payload does not contain result key")
		}

		if s, ok := r.body["message"]; ok {
			x := assertPayloadKeyFormat[string](t, "data", s)

			if x != "this is the message" {
				t.Errorf("got payload data %q, but want %q", x, "this is the message")
			}
		} else {
			t.Errorf("payload does not contain data key")
		}
	})

	t.Run("empty key response param", func(t *testing.T) {
		defer func() {
			tests.AssertPanicEmptyParam(t, recover(), "new", "response key")
		}()

		new(200, "success", "", "this is the message")
	})

	t.Run("empty data response param", func(t *testing.T) {
		defer func() {
			tests.AssertPanicEmptyParam(t, recover(), "new", "response data")
		}()

		new(200, "success", "message", "")
	})

	t.Run("not defined status code response param", func(t *testing.T) {
		defer func() {
			tests.AssertPanic(t, recover(), "new", "response code 0 is unknown")
		}()

		new(0, "success", "message", "this is the message")
	})
}

func assertHeader(t testing.TB, rw *http.Response, key, value string) {
	t.Helper()
	if v := rw.Header.Get(key); v != "" {
		if v != value {
			t.Errorf("got %q header %q, but want %q", key, v, value)
		}
	} else {
		t.Errorf("response does not contain %q header key", key)
	}
}

func assertPayloadKeyFormat[T string | int64 | result](t testing.TB, name string, value any) T {
	t.Helper()
	x, ok := value.(T)
	if !ok {
		t.Errorf("payload %s key is not in the correct format", name)
	}
	return x
}
