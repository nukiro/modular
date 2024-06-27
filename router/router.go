package router

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router interface {
	Mux() http.Handler
	Get(string, http.HandlerFunc)
}

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

type router struct {
	routes []*route
}

func (r *router) Mux() http.Handler {
	mux := httprouter.New()

	// Register all routes added by the user.
	for _, route := range r.routes {
		mux.HandlerFunc(route.method, route.path, route.handler)
	}

	return mux
}

func (r *router) Get(path string, handler http.HandlerFunc) {
	r.routes = append(
		r.routes,
		&route{http.MethodGet, r.fullPath(path), handler},
	)
}

func (r *router) fullPath(path string) string {
	return fmt.Sprintf("/%s", path)
}

func New() Router {
	return &router{
		routes: make([]*route, 0),
	}
}
