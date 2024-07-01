package response

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestResponseWrite(t *testing.T) {
	t.Run("without errors", func(t *testing.T) {
		// Save and restore original write function
		saved := write
		defer func() { write = saved }()
		// Fake write function
		write = func(w http.ResponseWriter, f serializer, r *response) error {
			return nil
		}

		w := httptest.NewRecorder()
		r := new(200, success, "message", "this is a message")

		rw := r.Write(w)

		if rw.Status() != "200 OK" {
			t.Errorf("got response status code %q, want %q", rw.Status(), "200 OK")
		}
	})

	t.Run("with errors", func(t *testing.T) {
		// Save and restore original write function
		saved := write
		defer func() { write = saved }()
		// Fake write function
		write = func(w http.ResponseWriter, f serializer, r *response) error {
			return errors.New("an error")
		}

		w := httptest.NewRecorder()
		r := new(200, success, "message", "this is a message")

		rw := r.Write(w)

		if rw.Status() != "500 Internal Server Error" {
			t.Errorf("got response status code %q, want %q", rw.Status(), "500 Internal Server Error")
		}
	})
}

func TestResponseStatus(t *testing.T) {
	tests := []struct {
		code   int
		status string
	}{
		{200, "200 OK"},
		{201, "201 Created"},
		{404, "404 Not Found"},
		{500, "500 Internal Server Error"},
	}

	for _, tt := range tests {
		n := fmt.Sprintf("Status Code %d", tt.code)
		t.Run(n, func(t *testing.T) {
			rw := new(tt.code, success, "test", "test")
			if rw.Status() != tt.status {
				t.Errorf("got %s, but want %s", rw.Status(), tt.status)
			}
		})
	}
}

func TestResponseBody(t *testing.T) {
	rw := new(200, success, "test", "test")
	if rw.Body() == nil {
		t.Errorf("body response does not exist")
	}
}

func TestResponseHeader(t *testing.T) {
	rw := new(200, success, "test", "test")
	if rw.Header() == nil {
		t.Errorf("body header does not exist")
	}
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
