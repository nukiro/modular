package router

import (
	"net/http"
	"testing"

	"github.com/nukiro/modular/internal/tests"
)

func TestFullPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"", "/"},
		{"articles", "/articles"},
		{"articles/:id", "/articles/:id"},
		{"articles/:id/comments", "/articles/:id/comments"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := BuildPath(tt.path)

			if got != tt.want {
				t.Errorf("got %q, but want %q", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	t.Run("new route", func(t *testing.T) {
		router := build()

		path := "articles"
		handler := func(w http.ResponseWriter, r *http.Request) {}
		router.Get(path, handler)

		assertRoutes(t, router.routes, "GET", "/articles")
	})

	t.Run("nil route handler", func(t *testing.T) {
		defer func() {
			tests.AssertPanicNilParam(t, recover(), "Get", "handler")
		}()

		router := build()
		router.Get("", nil)
	})
}

func TestPost(t *testing.T) {
	t.Run("new route", func(t *testing.T) {
		router := build()

		path := "articles"
		handler := func(w http.ResponseWriter, r *http.Request) {}
		router.Post(path, handler)

		assertRoutes(t, router.routes, "POST", "/articles")
	})

	t.Run("nil route handler", func(t *testing.T) {
		defer func() {
			tests.AssertPanicNilParam(t, recover(), "Post", "handler")
		}()

		router := build()
		router.Post("", nil)
	})
}

func TestPut(t *testing.T) {
	t.Run("new route", func(t *testing.T) {
		router := build()

		path := "articles"
		handler := func(w http.ResponseWriter, r *http.Request) {}
		router.Put(path, handler)

		assertRoutes(t, router.routes, "PUT", "/articles")
	})

	t.Run("nil route handler", func(t *testing.T) {
		defer func() {
			tests.AssertPanicNilParam(t, recover(), "Put", "handler")
		}()

		router := build()
		router.Put("", nil)
	})
}

func TestPatch(t *testing.T) {
	t.Run("new route", func(t *testing.T) {
		router := build()

		path := "articles"
		handler := func(w http.ResponseWriter, r *http.Request) {}
		router.Patch(path, handler)

		assertRoutes(t, router.routes, "PATCH", "/articles")
	})

	t.Run("nil route handler", func(t *testing.T) {
		defer func() {
			tests.AssertPanicNilParam(t, recover(), "Patch", "handler")
		}()

		router := build()
		router.Patch("", nil)
	})
}

func TestDelete(t *testing.T) {
	t.Run("new route", func(t *testing.T) {
		router := build()

		path := "articles"
		handler := func(w http.ResponseWriter, r *http.Request) {}
		router.Delete(path, handler)

		assertRoutes(t, router.routes, "DELETE", "/articles")
	})

	t.Run("nil route handler", func(t *testing.T) {
		defer func() {
			tests.AssertPanicNilParam(t, recover(), "Delete", "handler")
		}()

		router := build()
		router.Delete("", nil)
	})
}

func assertRoutes(t testing.TB, routes []*route, method, path string) {
	t.Helper()

	if len(routes) != 1 {
		t.Errorf("more than one route was added")
	}

	got := routes[0]

	if got.method != method {
		t.Errorf("got %q, but want %s method", got.method, method)
	}

	if got.path != path {
		t.Errorf("got %q, but want %q", got.path, path)
	}

	if got.handler == nil {
		t.Errorf("handler cannot be nil")
	}
}
