package nile

import (
	"encoding/json"
	"net/http"
	"time"
)

// Router is the basic foundation of the HTTP server.
type Router interface {
	GET(path string, fn HandlerFunc) error
	Start(addr string) error
}

type router struct {
	routes []*route
}

// New creates a new Router instance.
func New() Router {
	return &router{
		routes: []*route{},
	}
}

func (r *router) Start(addr string) error {
	server := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var matchingRoute *route
	for _, route := range r.routes {
		match, err := route.MatchPath(req.URL.Path)
		if err != nil {
			r.writeResponse(internalServiceError, w)
			return
		}

		if match {
			matchingRoute = route
			break
		}
	}

	if matchingRoute == nil {
		r.writeResponse(resourceNotFound, w)
		return
	}

	method, ok := matchingRoute.MethodHandlers[req.Method]
	if !ok {
		r.writeResponse(methodNotAllowed, w)
		return
	}

	resp := method.Handler(nil)
	r.writeResponse(resp, w)
}

func (r *router) writeResponse(resp Response, w http.ResponseWriter) {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode())
	w.Write(respBytes)
}

// GET adds a GET request for the matching path that executes the corresponding
// HandlerFunc upon a match.
func (r *router) GET(path string, fn HandlerFunc) error {
	return r.add("GET", path, fn)
}

func (r *router) add(method string, path string, fn HandlerFunc) error {
	newR, err := newRoute(method, path, fn)
	if err != nil {
		return err
	}

	routes, err := addRoute(r.routes, newR)
	if err != nil {
		return err
	}

	r.routes = routes
	return nil
}

func addRoute(routes []*route, r *route) ([]*route, error) {
	if len(routes) == 0 {
		return []*route{r}, nil
	}

	half := len(routes) / 2

	var err error
	var left []*route
	var right []*route

	switch r.Compare(routes[half]) {
	case -1:
		left = routes[:half]
		right, err = addRoute(routes[half:], r)
	case 1:
		left, err = addRoute(routes[:half], r)
		right = routes[half:]
	default:
		if err := routes[half].Merge(r); err != nil {
			return nil, err
		}
		return routes, nil
	}

	if err != nil {
		return nil, err
	}

	return append(left, right...), nil
}
