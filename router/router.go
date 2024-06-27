package router

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router interface {
	Mux() http.Handler
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
	Put(string, http.HandlerFunc)
	Patch(string, http.HandlerFunc)
	Delete(string, http.HandlerFunc)
}

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func buildRoute(method, path string, handler http.HandlerFunc) *route {
	return &route{method, BuildPath(path), handler}
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

func panicNilHandler(h http.HandlerFunc) {
	if h == nil {
		panic("handler cannot be nil")
	}
}

func (r *router) Get(path string, handler http.HandlerFunc) {
	panicNilHandler(handler)
	r.routes = append(r.routes, buildRoute(http.MethodGet, path, handler))
}

func (r *router) Post(path string, handler http.HandlerFunc) {
	panicNilHandler(handler)
	r.routes = append(r.routes, buildRoute(http.MethodPost, path, handler))
}

func (r *router) Put(path string, handler http.HandlerFunc) {
	panicNilHandler(handler)
	r.routes = append(r.routes, buildRoute(http.MethodPut, path, handler))
}

func (r *router) Patch(path string, handler http.HandlerFunc) {
	panicNilHandler(handler)
	r.routes = append(r.routes, buildRoute(http.MethodPatch, path, handler))
}

func (r *router) Delete(path string, handler http.HandlerFunc) {
	panicNilHandler(handler)
	r.routes = append(r.routes, buildRoute(http.MethodDelete, path, handler))
}

func BuildPath(path string) string {
	return fmt.Sprintf("/%s", path)
}

func build() *router {
	return &router{
		routes: make([]*route, 0),
	}
}

func New() Router {
	return build()
}
