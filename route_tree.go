package nile

import (
	"fmt"
	"strings"
)

type routeTree struct {
	segments []segment
}

type segment interface {
	Path() string
	HasParam() bool
	ParamName() (string, error)

	Match(path string) (*routeMatch, bool)
	Segments() []segment
}

func newSegment(path string) segment {
	head, tail := splitPath(path)
	segments := []segment{}

	if head != "" {
		segments = []segment{newSegment(tail)}
	}

	return &routeSegment{
		path:     head,
		segments: segments,
	}
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

type routeSegment struct {
	path     string
	segments []segment
}

func (rs *routeSegment) Path() string {
	return rs.path
}

func (rs *routeSegment) HasParam() bool {
	if len(rs.path) > 0 && string(rs.path[0]) == ":" {
		return true
	}

	return false
}

func (rs *routeSegment) ParamName() (string, error) {
	if rs.HasParam() {
		return rs.path[1:], nil
	}

	return "", fmt.Errorf("Segment %s has not param", rs.path)
}

func (rs *routeSegment) Segments() []segment {
	return rs.segments
}

func (rs *routeSegment) Match(path string) (*routeMatch, bool) {
	rm := &routeMatch{
		URI:    path,
		Params: []*routeParam{},
	}

	hasMatch := match(path, rs, rm)
	if hasMatch {
		return rm, true
	}

	return nil, false
}

func match(path string, seg segment, routeMatch *routeMatch) bool {
	headMatches := false

	// Split the head from the tail of the path to match.
	head, tail := splitPath(path)

	if head == seg.Path() {
		headMatches = true
	} else if seg.HasParam() && len(head) > 0 {
		headMatches = true
	}

	if !headMatches {
		return false
	}

	if tail == "" && len(seg.Segments()) == 0 {
		return true
	}

	for _, childSeg := range seg.Segments() {
		hasMatch := match(tail, childSeg, routeMatch)
		if hasMatch {
			return true
		}
	}

	return false
}
