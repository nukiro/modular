package server

import (
	"fmt"
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

func TestNew(t *testing.T) {
	t.Run("with a configuration", func(t *testing.T) {
		// Create a server configuration
		config := &Configuration{
			Host:         "127.0.0.1",
			Port:         1234,
			IdleTimeout:  10 * time.Minute,
			ReadTimeout:  time.Minute,
			WriteTimeout: 2 * time.Minute,
		}

		srv := new(config)
		assertServer(t, srv, config)
	})

	t.Run("with a default configuration", func(t *testing.T) {
		srv := new(nil)
		assertServer(t, srv, configuration)
	})
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

func assertServer(t testing.TB, srv *server, c *Configuration) {
	t.Helper()

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)

	if srv.Server.Addr != addr {
		t.Errorf("Addr is %q, but want %q", srv.Server.Addr, addr)
	}
	if srv.Server.IdleTimeout != c.IdleTimeout {
		t.Errorf(
			"Idle Timeout is %q, but want %q",
			srv.Server.IdleTimeout, c.IdleTimeout,
		)
	}
	if srv.Server.ReadTimeout != c.ReadTimeout {
		t.Errorf(
			"Read Timeout is %q, but want %q",
			srv.Server.ReadTimeout, c.ReadTimeout,
		)
	}
	if srv.Server.WriteTimeout != c.WriteTimeout {
		t.Errorf(
			"Write Timeout is %q, but want %q",
			srv.Server.WriteTimeout, c.WriteTimeout,
		)
	}
}
