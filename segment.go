package nile

import (
	"errors"
	"fmt"
	"strings"
)

// segment describes a portion of a URI path. Each Segment represents a subset
// of the path between the forward-slash '/' of a URI. Segments are structured
// into a tree that can be used to construct the full path.
type segment struct {
	// Path is the portion of the URI that is described by this Segment.
	Path       string
	children   map[string]*segment
	childOrder []string
	paramChild *segment
	endpoints  map[string]endpoint
}

// newSegment accepts a path and creates a new Segment.
func newSegment(path string) *segment {
	head, tail := splitPath(path)
	seg := &segment{
		Path:       head,
		children:   map[string]*segment{},
		childOrder: []string{},
		endpoints:  map[string]endpoint{},
	}

	if tail != "" {
		seg.AddChild(newSegment(tail))
	}

	return seg
}

// newSegmentEndpoint creates a Segment and attaches an Endpoint at the leaf
// node.
func newSegmentEndpoint(path string, method string, handler HandlerFunc) (*segment, error) {
	head, tail := splitPath(path)
	seg := &segment{
		Path:      head,
		children:  map[string]*segment{},
		endpoints: map[string]endpoint{},
	}

	if tail != "" {
		child, err := newSegmentEndpoint(tail, method, handler)
		if err != nil {
			return nil, err
		}

		if err := seg.AddChild(child); err != nil {
			return nil, err
		}
	} else {
		endPt, err := newEndpoint(method, handler)
		if err != nil {
			return nil, err
		}

		if err := seg.AddEndpoint(endPt); err != nil {
			return nil, err
		}
	}

	return seg, nil
}

// mergeSegments combines two different segments and returns the combined object.
// It assumes that the root of each segment have identical paths.
func mergeSegments(first, second *segment) (*segment, error) {
	if first.Path != second.Path {
		return nil, errors.New("May only merge segments with the same path")
	}

	merged := newSegment(first.Path)

	addChildren := func(children []*segment) error {
		for _, child := range children {
			if err := merged.AddChild(child); err != nil {
				return err
			}
		}

		return nil
	}

	addEndpoints := func(endpoints []endpoint) error {
		for _, endpoint := range endpoints {
			if err := merged.AddEndpoint(endpoint); err != nil {
				return err
			}
		}

		return nil
	}

	if err := addChildren(first.Children()); err != nil {
		return nil, err
	}
	if err := addChildren(second.Children()); err != nil {
		return nil, err
	}

	if err := addEndpoints(first.Endpoints()); err != nil {
		return nil, err
	}
	if err := addEndpoints(second.Endpoints()); err != nil {
		return nil, err
	}

	return merged, nil
}

// Children gets the list of paths that exist under the current path.
func (s *segment) Children() []*segment {
	children := make([]*segment, len(s.childOrder))
	for idx, childPath := range s.childOrder {
		children[idx] = s.children[childPath]
	}

	if s.paramChild != nil {
		children = append(children, s.paramChild)
	}

	return children
}

// AddChild adds a child path that should exist under the current path.
func (s *segment) AddChild(child *segment) error {
	if isParam(child.Path) {
		if s.paramChild != nil {
			return fmt.Errorf("Segment %s already has a route with a parameter", s.Path)
		}

		s.paramChild = child
		return nil
	}

	if currentChild, exists := s.children[child.Path]; exists {
		merged, err := mergeSegments(currentChild, child)
		if err != nil {
			return err
		}

		s.children[child.Path] = merged
		return nil
	}

	s.children[child.Path] = child

	for idx, childPath := range s.childOrder {
		if strings.Compare(child.Path, childPath) == 1 {
			s.childOrder = append(s.childOrder, "")
			copy(s.childOrder[idx+1:], s.childOrder[idx:])
			s.childOrder[idx] = child.Path

			return nil
		}
	}

	s.childOrder = append(s.childOrder, child.Path)
	return nil
}

// RemoveChild removes a child path from the current path.
func (s *segment) RemoveChild(path string) error {
	if isParam(path) {
		s.paramChild = nil
		return nil
	}

	if _, exists := s.children[path]; !exists {
		return fmt.Errorf("Unable to remove child %s from segment %s: child does not exist", path, s.Path)
	}

	delete(s.children, path)
	return nil
}

// Endpoint gets the Endpoint that matches an HTTP method.
func (s *segment) Endpoint(method string) (endpoint, bool) {
	endpoint, found := s.endpoints[method]
	return endpoint, found
}

// Endpoints gets the list of HTTP endpoints that resolve exactly at this
// path.
func (s *segment) Endpoints() []endpoint {
	endpoints := make([]endpoint, len(s.endpoints))
	idx := 0

	for _, endPt := range s.endpoints {
		endpoints[idx] = endPt
		idx++
	}

	return endpoints
}

// AddEndpoint adds an Endpoint to the Segment.
func (s *segment) AddEndpoint(endPt endpoint) error {
	if _, exists := s.endpoints[endPt.Method()]; exists {
		return fmt.Errorf("Unable to add %s endpoint to segment %s: endpoint already exists", endPt.Method(), s.Path)
	}

	s.endpoints[endPt.Method()] = endPt
	return nil
}

// Matches checks a path against the current Segment's endpoints.
// If a match doesn't exist, it checks against the Segment's children.
func (s *segment) Matches(path string) (*match, bool) {
	head, tail := splitPath(path)

	if head != s.Path && !isParam(s.Path) {
		return nil, false
	}

	if tail == "" {
		match := newMatch(s, path)
		if isParam(s.Path) {
			match.AddParam(s.Path[1:], head)
		}
		return match, true
	}

	// Check the children
	for _, child := range s.children {
		match, matches := child.Matches(tail)
		if matches {
			match.RequestURI = path
			if isParam(s.Path) {
				match.AddParam(s.Path[1:], head)
			}
			return match, matches
		}
	}

	if s.paramChild != nil {
		match, matches := s.paramChild.Matches(tail)
		if matches {
			match.RequestURI = path
			match.AddParam(s.Path[1:], head)
			return match, matches
		}
	}

	return nil, false
}

func isParam(path string) bool {
	if len(path) == 0 {
		return false
	}

	return string(path[0]) == ":"
}

func splitPath(path string) (head string, tail string) {
	if len(path) == 0 {
		return
	}

	if string(path[0]) == "/" {
		path = path[1:]
	}

	lastIdx := len(path) - 1
	if lastIdx >= 0 && string(path[lastIdx]) == "/" {
		path = path[:lastIdx]
	}

	head = strings.Split(path, "/")[0]
	tail = path[len(head):]

	return
}
