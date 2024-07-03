package server

import (
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/nukiro/modular/internal/tests"
)

func TestAddress(t *testing.T) {
	// From a server default configuration
	config := *configuration
	// Set server host and port
	config.Host = "127.0.0.1"
	config.Port = 4000

	want := "127.0.0.1:4000"

	// Generate a new server
	srv := server{config: &config}

	got := srv.address()
	if got != want {
		t.Errorf("got %q, but want %q", got, want)
	}
}

func TestBuild(t *testing.T) {
	// Create a server configuration
	config := &Configuration{
		Host:         "127.0.0.1",
		Port:         1234,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  time.Minute,
		WriteTimeout: 2 * time.Minute,
	}

	srv := new(config)
	addr := "127.0.0.1:1234"

	if srv.Server.Addr != addr {
		t.Errorf("Addr is %q, but want %q", srv.Server.Addr, addr)
	}
	if srv.Server.IdleTimeout != config.IdleTimeout {
		t.Errorf(
			"Idle Timeout is %q, but want %q",
			srv.Server.IdleTimeout, config.IdleTimeout,
		)
	}
	if srv.Server.ReadTimeout != config.ReadTimeout {
		t.Errorf(
			"Read Timeout is %q, but want %q",
			srv.Server.ReadTimeout, config.ReadTimeout,
		)
	}
	if srv.Server.WriteTimeout != config.WriteTimeout {
		t.Errorf(
			"Write Timeout is %q, but want %q",
			srv.Server.WriteTimeout, config.WriteTimeout,
		)
	}
}

func TestLogger(t *testing.T) {
	t.Run("new logger", func(t *testing.T) {
		// Generate a new logger
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

		// Build a server and set the logger
		srv := new(nil)
		srv.Logger(logger)

		if srv.logger == nil {
			t.Errorf("logger was not set")
		}
		if srv.Server.ErrorLog == nil {
			t.Errorf("server error log was not set")
		}
	})

	t.Run("nil pointer logger", func(t *testing.T) {
		defer func() {
			tests.AssertPanicNilParam(t, recover(), "Logger", "logger")
		}()

		// Build a server and set the logger
		srv := New()
		srv.Logger(nil)
	})
}

func TestHandler(t *testing.T) {
	t.Run("new http handler", func(t *testing.T) {
		// Generate a new mux
		mux := http.DefaultServeMux

		// Build a server and set the handler
		srv := new(nil)
		srv.Handler(mux)

		if srv.Server.Handler == nil {
			t.Errorf("server handler was not set")
		}
	})

	t.Run("nil pointer handler", func(t *testing.T) {
		defer func() {
			tests.AssertPanicNilParam(t, recover(), "Handler", "handler")
		}()

		srv := New()
		srv.Handler(nil)
	})
}

func TestRun(t *testing.T) {
	t.Run("nil server handler", func(t *testing.T) {
		srv := New()
		err := srv.Run()

		if err == nil {
			t.Errorf("run did not return an error")
		}
	})
}
