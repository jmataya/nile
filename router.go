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
	segments map[string]Segment
}

// New creates a new Router instance.
func New() Router {
	return &router{
		segments: map[string]Segment{},
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
	path := req.URL.Path
	method := req.Method

	var match *Match
	var hasMatch bool
	for _, segment := range r.segments {
		match, hasMatch = segment.Matches(path)
		if hasMatch {
			break
		}
	}

	if !hasMatch {
		r.writeResponse(ResourceNotFound, w)
		return
	}

	endpoint, found := match.Segment.Endpoint(method)
	if !found {
		r.writeResponse(MethodNotAllowed, w)
		return
	}

	context := match.Context
	context.Request = req
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
	seg, err := NewSegmentEndpoint(path, http.MethodGet, fn)
	if err != nil {
		return err
	}

	existing, found := r.segments[seg.Path()]
	if found {
		merged, err := MergeSegments(existing, seg)
		if err != nil {
			return err
		}

		r.segments[seg.Path()] = merged
	} else {
		r.segments[seg.Path()] = seg
	}

	return nil
}
