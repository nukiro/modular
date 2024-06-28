package response

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	resp := new(200, "success", "message", "this is the message")

	if resp.code != 200 {
		t.Errorf("got code %d, but want %d", resp.code, 200)
	}

	if s, ok := resp.payload["time"]; ok {
		x, ok := s.(int64)
		if !ok {
			t.Errorf("payload time is not in the correct format")
		}

		y := time.Unix(int64(x), 0)
		z := time.Now()

		if y.After(z) {
			t.Errorf("payload time is for the future")
		}
	} else {
		t.Errorf("payload does not contain time key")
	}

	if s, ok := resp.payload["status"]; ok {
		x, ok := s.(string)
		if !ok {
			t.Errorf("payload status is not in the correct format")
		}

		if x != "ok" {
			t.Errorf("got payload status %q, but want %q", x, "ok")
		}
	} else {
		t.Errorf("payload does not contain status key")
	}

	if s, ok := resp.payload["result"]; ok {
		x, ok := s.(string)
		if !ok {
			t.Errorf("payload result is not in the correct format")
		}

		if x != "success" {
			t.Errorf("got payload result %q, but want %q", x, "success")
		}
	} else {
		t.Errorf("payload does not contain result key")
	}

	if s, ok := resp.payload["message"]; ok {
		x, ok := s.(string)
		if !ok {
			t.Errorf("payload data key is not in the correct format")
		}

		if x != "this is the message" {
			t.Errorf("got payload data %q, but want %q", x, "this is the message")
		}
	} else {
		t.Errorf("payload does not contain data key")
	}
}

// func assertPayloadKeyFormat(t testing.TB) {}
