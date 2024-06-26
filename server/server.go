package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Configuration struct {
	TLS          bool
	Host         string
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server interface {
	Run() error
	Logger(*slog.Logger)
	Handler(http.Handler)
}

type server struct {
	*http.Server
	*Configuration
	logger *slog.Logger
}

func (s *server) address() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

func (s *server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Call builtin recover function to check
			// it there has been a panic or not.
			defer func() {
				// If recover is called outside the deferred
				// function it will not stop a panicking sequence.
				if err := recover(); err != nil {
					// Close the connection works as a trigger for the Go's
					// HTTP server to automatically close the current connection.
					w.Header().Set("Connection", "close")
					w.WriteHeader(500)
					s.logger.Error("server recover panic", "error", fmt.Sprintf("%s", err))
				}
			}()
			next.ServeHTTP(w, r)
		})
}

func (s *server) Logger(logger *slog.Logger) {
	s.logger = logger
	s.Server.ErrorLog = slog.NewLogLogger(s.logger.Handler(), slog.LevelError)
}

func (s *server) Handler(h http.Handler) {
	s.Server.Handler = s.recoverPanic(h)
}

func (s *server) Run() error {
	if s.logger == nil {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		s.Server.ErrorLog = slog.NewLogLogger(logger.Handler(), slog.LevelError)
	}

	if s.Server.Handler == nil {
		return errors.New("server handler has not been set")
	}

	shutdownError := make(chan error)

	// Start background shutdown routine
	// to safely stop the running application.
	go func() {
		// Channel which carries signal values.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		c := <-quit

		// Clean up when a signal has been caught.
		s.logger.Info("shutting down server", "signal", c.String())
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- s.Server.Shutdown(ctx)
	}()

	// Starting the server.
	s.logger.Info("starting server", "addr", s.address())
	err := s.Server.ListenAndServe()
	// Calling Shutdown() on our server will cause ListenAndServe()
	// to immediately return a server closed error.
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise, we wait to receive the return value from Shutdown().
	// If return value is an error, we know that there was a
	// problem with the gracefull shutdown and we return the error.
	err = <-shutdownError
	if err != nil {
		return err
	}

	s.logger.Info("server stopped", "addr", s.address())

	return nil
}

func New(c *Configuration) Server {
	config := c
	if c == nil {
		config = &Configuration{
			false,
			"localhost",
			"8080",
			time.Minute,
			5 * time.Second,
			10 * time.Second,
		}
	}
	srv := &server{Configuration: config}
	srv.Server = &http.Server{
		Addr:         srv.address(),
		IdleTimeout:  config.IdleTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	return srv
}
