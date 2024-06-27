package server

import (
	"testing"
	"time"
)

func TestServerAddress(t *testing.T) {
	// From a server default configuration
	config := configuration
	// Set server host and port
	config.Host = "127.0.0.1"
	config.Port = "4000"

	want := "127.0.0.1:4000"

	// Generate a new server
	srv := server{Configuration: configuration}

	got := srv.address()

	if srv.address() != want {
		t.Errorf("got %q, but want %q", got, want)
	}
}

func TestBuildServer(t *testing.T) {
	// Create a server configuration
	config := &Configuration{
		Host:         "127.0.0.1",
		Port:         "1234",
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  time.Minute,
		WriteTimeout: 2 * time.Minute,
	}

	srv := build(config)
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
