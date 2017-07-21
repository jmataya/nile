package nile

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type route struct {
	RegexPath         string
	ParameterizedPath string
	MethodHandlers    map[string]*methodHandler
}

func newRoute(method string, path string, handler HandlerFunc) (*route, error) {
	r := route{
		RegexPath:         makeRegexPath(path),
		ParameterizedPath: path,
		MethodHandlers:    map[string]*methodHandler{},
	}

	if err := r.AddMethod(method, handler); err != nil {
		return nil, err
	}

	return &r, nil
}

func makeRegexPath(path string) string {
	if path == "" {
		return ""
	}

	if string(path[0]) != "/" {
		path = "/" + path
	}

	if string(path[len(path)-1]) != "/" {
		path += "/*"
	}

	path += "$"

	return path
}

func (r *route) AddMethod(method string, handler HandlerFunc) error {
	if _, ok := r.MethodHandlers[method]; ok {
		return fmt.Errorf("Path %s already contains method %s", r.ParameterizedPath, method)
	}

	mh := methodHandler{
		Method:  method,
		Handler: handler,
	}

	r.MethodHandlers[method] = &mh
	return nil
}

func (r *route) Compare(other *route) int {
	return strings.Compare(r.ParameterizedPath, other.ParameterizedPath)
}

func (r *route) MatchPath(path string) (bool, error) {
	return regexp.MatchString(r.RegexPath, path)
}

func (r *route) Merge(other *route) error {
	if r.ParameterizedPath != other.ParameterizedPath {
		return errors.New("Routes can't be merged because paths are different")
	}

	for _, mh := range other.MethodHandlers {
		if _, ok := r.MethodHandlers[mh.Method]; ok {
			return fmt.Errorf("Destination route already has method %s", mh.Method)
		}

		r.MethodHandlers[mh.Method] = mh
	}

	return nil
}

type methodHandler struct {
	Method  string
	Handler HandlerFunc
}
