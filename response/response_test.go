package response

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("new response", func(t *testing.T) {
		r := new(200, success, "message", "this is the message")

		if r.headers != nil {
			t.Errorf("headers was not empty")
		}

		if r.code != 200 {
			t.Errorf("got code %d, but want %d", r.code, 200)
		}

		if s, ok := r.payload["time"]; ok {
			x := assertPayloadKeyFormat[int64](t, "time", s)

			y := time.Unix(int64(x), 0)
			z := time.Now()

			if y.After(z) {
				t.Errorf("payload time is for the future")
			}
		} else {
			t.Errorf("payload does not contain time key")
		}

		if s, ok := r.payload["status"]; ok {
			x := assertPayloadKeyFormat[string](t, "status", s)

			if x != "ok" {
				t.Errorf("got payload status %q, but want %q", x, "ok")
			}
		} else {
			t.Errorf("payload does not contain status key")
		}

		if s, ok := r.payload["result"]; ok {
			x := assertPayloadKeyFormat[result](t, "result", s)

			if x != "success" {
				t.Errorf("got payload result %q, but want %q", x, "success")
			}
		} else {
			t.Errorf("payload does not contain result key")
		}

		if s, ok := r.payload["message"]; ok {
			x := assertPayloadKeyFormat[string](t, "data", s)

			if x != "this is the message" {
				t.Errorf("got payload data %q, but want %q", x, "this is the message")
			}
		} else {
			t.Errorf("payload does not contain data key")
		}
	})

	t.Run("new response with an empty key", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("new did not panic")
			}
		}()

		new(200, "success", "", "this is the message")
	})

	t.Run("status code is not defined", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("new did not panic")
			}
		}()

		new(0, "success", "message", "this is the message")
	})
}

func assertPayloadKeyFormat[T string | int64 | result](t testing.TB, name string, value any) T {
	t.Helper()
	x, ok := value.(T)
	if !ok {
		t.Errorf("payload %s key is not in the correct format", name)
	}
	return x
}
