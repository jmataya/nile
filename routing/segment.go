package routing

import (
	"errors"
	"fmt"
	"strings"
)

// Segment describes a portion of a URI path. Each Segment represents a subset
// of the path between the forward-slash '/' of a URI. Segments are structured
// into a tree that can be used to construct the full path.
type Segment interface {
	// Path gets the portion of the URI that is described by this Segment.
	Path() string

	// Children gets the list of paths that exist under the current path.
	Children() []Segment

	// AddChild adds a child path that should exist under the current path.
	AddChild(child Segment) error

	// RemoveChild removes a child path from the current path.
	RemoveChild(path string) error

	// Endpoints gets the list of HTTP endpoints that resolve exactly at this
	// path.
	Endpoints() []Endpoint

	// AddEndpoint adds an Endpoint to the Segment.
	AddEndpoint(endPt Endpoint) error

	// Parent gets the Segment that precedes the current Segment in the path. The
	// second parameter of the tuple will return true if a parent exists. A lack
	// of parent exists that this Segment is at the root.
	Parent() (Segment, bool)

	// SetParent sets the Segment that preceds the current Segment in the path.
	SetParent(parent Segment)

	// Matches checks a path and method against the current Segment's endpoints.
	// If a match doesn't exist, it checks against the Segment's children.
	Matches(path, method string) (Endpoint, bool)
}

// NewSegment accepts a path and creates a new Segment.
func NewSegment(path string) Segment {
	head, tail := splitPath(path)
	seg := &segment{
		path:       head,
		parent:     nil,
		children:   map[string]Segment{},
		childOrder: []string{},
		endpoints:  map[string]Endpoint{},
	}

	if tail != "" {
		seg.AddChild(NewSegment(tail))
	}

	return seg
}

// NewSegmentEndpoint creates a Segment and attaches an Endpoint at the leaf
// node.
func NewSegmentEndpoint(path, method string) (Segment, error) {
	head, tail := splitPath(path)
	seg := &segment{
		path:      head,
		parent:    nil,
		children:  map[string]Segment{},
		endpoints: map[string]Endpoint{},
	}

	if tail != "" {
		child, err := NewSegmentEndpoint(tail, method)
		if err != nil {
			return nil, err
		}

		if err := seg.AddChild(child); err != nil {
			return nil, err
		}
	} else {
		endPt, err := NewEndpoint(method)
		if err != nil {
			return nil, err
		}

		if err := seg.AddEndpoint(endPt); err != nil {
			return nil, err
		}
	}

	return seg, nil
}

type segment struct {
	path       string
	parent     Segment
	children   map[string]Segment
	childOrder []string
	endpoints  map[string]Endpoint
}

func (s *segment) Path() string {
	return s.path
}

func (s *segment) Parent() (Segment, bool) {
	hasParent := s.parent != nil
	return s.parent, hasParent
}

func (s *segment) SetParent(parent Segment) {
	s.parent = parent
}

func (s *segment) Children() []Segment {
	children := make([]Segment, len(s.childOrder))
	for idx, childPath := range s.childOrder {
		children[idx] = s.children[childPath]
	}
	return children
}

func (s *segment) AddChild(child Segment) error {
	if currentChild, exists := s.children[child.Path()]; exists {
		merged, err := mergeSegments(currentChild, child)
		if err != nil {
			return err
		}

		s.children[child.Path()] = merged
		return nil
	}

	s.children[child.Path()] = child

	for idx, childPath := range s.childOrder {
		if strings.Compare(child.Path(), childPath) == 1 {
			s.childOrder = append(s.childOrder, "")
			copy(s.childOrder[idx+1:], s.childOrder[idx:])
			s.childOrder[idx] = child.Path()

			return nil
		}
	}

	s.childOrder = append(s.childOrder, child.Path())
	return nil
}

func mergeSegments(first, second Segment) (Segment, error) {
	if first.Path() != second.Path() {
		return nil, errors.New("May only merge segments with the same path")
	}

	merged := NewSegment(first.Path())

	addChildren := func(children []Segment) error {
		for _, child := range children {
			if err := merged.AddChild(child); err != nil {
				return err
			}
		}

		return nil
	}

	addEndpoints := func(endpoints []Endpoint) error {
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

func (s *segment) RemoveChild(path string) error {
	if _, exists := s.children[path]; !exists {
		return fmt.Errorf("Unable to remove child %s from segment %s: child does not exist", path, s.path)
	}

	delete(s.children, path)
	return nil
}

func (s *segment) Endpoints() []Endpoint {
	endpoints := make([]Endpoint, len(s.endpoints))
	idx := 0

	for _, endPt := range s.endpoints {
		endpoints[idx] = endPt
		idx++
	}

	return endpoints
}

func (s *segment) AddEndpoint(endPt Endpoint) error {
	if _, exists := s.endpoints[endPt.Method()]; exists {
		return fmt.Errorf("Unable to add %s endpoint to segment %s: endpoint already exists", endPt.Method(), s.path)
	}

	s.endpoints[endPt.Method()] = endPt
	return nil
}

func (s *segment) RemoveEndpoint(method string) error {
	if _, exists := s.endpoints[method]; !exists {
		return fmt.Errorf("Unable to remove %s endpoint from segment %s: endpoint does not exist", method, s.path)
	}

	delete(s.endpoints, method)
	return nil
}

func (s *segment) Matches(path, method string) (Endpoint, bool) {
	head, tail := splitPath(path)

	if head != s.path && !isParam(s.path) {
		return nil, false
	}

	if tail == "" {
		// Check the endpoints
		endPt, matches := s.endpoints[method]
		return endPt, matches
	}

	// Check the children
	for _, child := range s.children {
		endPt, matches := child.Matches(tail, method)
		if matches {
			return endPt, matches
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
