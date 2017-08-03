package nile

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Router is the basic foundation of the HTTP server.
type Router interface {
	GET(path string, fn HandlerFunc) error
	Start(addr string) error
}

type router struct {
	segments map[string]*segment
}

// New creates a new Router instance.
func New() Router {
	return &router{
		segments: map[string]*segment{},
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

	printLogo(addr)

	return server.ListenAndServe()
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	var match *match
	var hasMatch bool
	for _, segment := range r.segments {
		match, hasMatch = segment.Matches(path)
		if hasMatch {
			break
		}
	}

	if !hasMatch {
		r.writeResponse(resourceNotFound, w)
		return
	}

	endpoint, found := match.Segment.Endpoint(method)
	if !found {
		r.writeResponse(methodNotAllowed, w)
		return
	}

	context := match.Context
	context.setRequest(req)
	handler := endpoint.Handler()

	resp := handler(context)
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
	return r.addRoute(path, http.MethodGet, fn)
}

// POST adds a POST request for the matching path that executes the
// corresponding HandlerFunc upon a match.
func (r *router) POST(path string, fn HandlerFunc) error {
	return r.addRoute(path, http.MethodPost, fn)
}

// PATCH adds a PATCH request for the matching path that executes the
// corresponding HandlerFunc upon a match.
func (r *router) PATCH(path string, fn HandlerFunc) error {
	return r.addRoute(path, http.MethodPatch, fn)
}

// PUT adds a PUT request for the matching path that executes the corresponding
// HandlerFunc upon a match.
func (r *router) PUT(path string, fn HandlerFunc) error {
	return r.addRoute(path, http.MethodPut, fn)
}

// DELETE adds a DELETE request for the matching path that executes the
// corresponding HandlerFunc upon a match.
func (r *router) DELETE(path string, fn HandlerFunc) error {
	return r.addRoute(path, http.MethodDelete, fn)
}

func (r *router) addRoute(path string, method string, handler HandlerFunc) error {
	seg, err := newSegmentEndpoint(path, method, handler)
	if err != nil {
		return err
	}

	existing, found := r.segments[seg.Path]
	if found {
		merged, err := mergeSegments(existing, seg)
		if err != nil {
			return err
		}

		r.segments[seg.Path] = merged
	} else {
		r.segments[seg.Path] = seg
	}

	return nil
}

func printLogo(addr string) {
	const logo = `
      (_) |     
 _ __  _| | ___ 
| '_ \| | |/ _ \
| | | | | |  __/
|_| |_|_|_|\___|
`

	fmt.Println(logo)
	fmt.Println(formatAddress(addr))
}

func formatAddress(addr string) string {
	if string(addr[0]) == ":" {
		// Assume this means we start with the port.
		addr = "http://localhost" + addr
	}

	url, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("Server started on: %s", url.String())
}
