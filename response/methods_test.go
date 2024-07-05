package response

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nukiro/modular/internal/tests"
)

func TestSerialize(t *testing.T) {
	t.Run("with no error", func(t *testing.T) {
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return []byte("Hello World"), nil
		})

		js, err := serialize(f, "body")

		if err != nil {
			t.Errorf("an error was returned, when it is not expected")
		}

		if js == nil {
			t.Errorf("serialize did not return a value")
		}
	})

	// If we do not check for response nil pointer, in case a nil pointer is
	// provided to serialize will panic when f tries to deallocate the body,
	// runtime error: invalid memory address or nil pointer dereference
	t.Run("with a nil response", func(t *testing.T) {
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return nil, nil
		})

		defer func() {
			tests.AssertPanicNilParam(t, recover(), "serialize", "body")
		}()

		serialize(f, nil)
	})

	t.Run("with a nil serializer", func(t *testing.T) {
		defer func() {
			tests.AssertPanicNilParam(t, recover(), "serialize", "serializer")
		}()

		serialize(nil, "body")
	})

	t.Run("when serializer returns an error", func(t *testing.T) {
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return nil, errors.New("error")
		})

		js, err := serialize(f, "body")

		if err == nil {
			t.Errorf("serialize did not return an error")
		}

		if js != nil {
			t.Errorf("a value was returned, when it is not expected")
		}
	})
}

func TestWrite(t *testing.T) {
	t.Run("good response", func(t *testing.T) {
		r := New(200)
		r.Body = "Good Response"
		r.Header.Set("Test Key", "Test Value")

		w := httptest.NewRecorder()
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return []byte("Good Response"), nil
		})

		if err := write(w, f, r); err != nil {
			t.Errorf("write return error: %q", err.Error())
		}

		rw := w.Result()

		if rw.StatusCode != 200 {
			t.Errorf("got status code %d, but want %d", rw.StatusCode, 200)
		}

		defer rw.Body.Close()
		body, err := io.ReadAll(rw.Body)
		if err != nil {
			t.Fatal(err)
		}
		body = bytes.TrimSpace(body)

		if string(body) != "Good Response" {
			t.Errorf("got body %q, but want %q", string(body), "Good Response")
		}

		assertHeader(t, rw, "Content-Type", "application/json")
		assertHeader(t, rw, "Test Key", "Test Value")
	})

	// If we do not check for response nil pointer, in case a nil pointer is
	// provided to write will panic when the range loop tries to
	// deallocate the response header
	// runtime error: invalid memory address or nil pointer dereference
	t.Run("with a nil response", func(t *testing.T) {
		w := httptest.NewRecorder()
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return nil, nil
		})

		defer func() {
			tests.AssertPanicNilParam(t, recover(), "write", "response")
		}()

		write(w, f, nil)
	})

	t.Run("with a nil serializer", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := New(200)

		defer func() {
			tests.AssertPanicNilParam(t, recover(), "write", "serializer")
		}()

		write(w, nil, r)
	})

	t.Run("with a nil response writer", func(t *testing.T) {
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return nil, nil
		})
		r := New(200)

		defer func() {
			tests.AssertPanicNilParam(t, recover(), "write", "response writer")
		}()

		write(nil, f, r)
	})

	t.Run("error serializing the response", func(t *testing.T) {
		r := New(200)
		r.Body = "Good Response"
		r.Header.Set("Test Key", "Test Value")

		w := httptest.NewRecorder()
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return nil, errors.New("an error occurred")
		})

		if err := write(w, f, r); err == nil {
			t.Errorf("write did not return an error")
		}
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
