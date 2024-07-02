package response

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nukiro/modular/internal/tests"
)

func TestSerialize(t *testing.T) {
	t.Run("with no error", func(t *testing.T) {
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return []byte("Hello World"), nil
		})
		r := new(200, success, "message", "this is a message")

		js, err := serialize(f, r)

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
			tests.AssertPanicNilParam(t, recover(), "serialize", "response")
		}()

		serialize(f, nil)
	})

	t.Run("with a nil serializer", func(t *testing.T) {
		r := new(200, success, "message", "this is a message")
		defer func() {
			tests.AssertPanicNilParam(t, recover(), "serialize", "serializer")
		}()

		serialize(nil, r)
	})

	t.Run("when serializer returns an error", func(t *testing.T) {
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return nil, errors.New("error")
		})
		r := new(200, "suscess", "message", "this is a message")

		js, err := serialize(f, r)

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
		r := new(200, success, "message", "this is the message")
		r.header.Set("Test Key", "Test Value")
		w := httptest.NewRecorder()
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return []byte("Good response"), nil
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

		if string(body) != "Good response" {
			t.Errorf("got body %q, but want %q", string(body), "Good response")
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
		r := new(200, success, "message", "this is the message")

		defer func() {
			tests.AssertPanicNilParam(t, recover(), "write", "serializer")
		}()

		write(w, nil, r)
	})

	t.Run("with a nil response writer", func(t *testing.T) {
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return nil, nil
		})
		r := new(200, success, "message", "this is the message")

		defer func() {
			tests.AssertPanicNilParam(t, recover(), "write", "response writer")
		}()

		write(nil, f, r)
	})

	t.Run("error serializing the response", func(t *testing.T) {
		w := httptest.NewRecorder()
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return nil, errors.New("an error occurred")
		})
		r := new(200, "success", "message", "this is the message")

		if err := write(w, f, r); err == nil {
			t.Errorf("write did not return an error")
		}
	})
}

func TestWriteError(t *testing.T) {
	t.Run("internal server error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		f := json.MarshalIndent

		rs := writeError(w, f)

		if rs == nil {
			t.Errorf("writeError did not return a response")
		}

		rw := w.Result()

		if rw.StatusCode != 500 {
			t.Errorf("got response status code %d, but want %d", rw.StatusCode, 500)
		}

		defer rw.Body.Close()
		body, err := io.ReadAll(rw.Body)
		if err != nil {
			t.Fatal(err)
		}
		body = bytes.TrimSpace(body)

		if !strings.Contains(string(body), internalServerErrorMsg) {
			t.Errorf("response body:\n%s,\ndoes not contain %q", string(body), internalServerErrorMsg)
		}
	})

	t.Run("error serializing internal server error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return nil, errors.New("an error occurred")
		})

		rs := writeError(w, f)

		if rs == nil {
			t.Errorf("writeError did not return a response")
		}

		rw := w.Result()

		if rw.StatusCode != 500 {
			t.Errorf("got response status code %d, but want %d", rw.StatusCode, 500)
		}
	})

	t.Run("with a nil serializer", func(t *testing.T) {
		w := httptest.NewRecorder()

		defer func() {
			tests.AssertPanicNilParam(t, recover(), "write", "serializer")
		}()

		writeError(w, nil)
	})

	t.Run("with a nil response writer", func(t *testing.T) {
		f := Serializer(func(v any, prefix, indent string) ([]byte, error) {
			return nil, nil
		})

		defer func() {
			tests.AssertPanicNilParam(t, recover(), "write", "response writer")
		}()

		writeError(nil, f)
	})

}
