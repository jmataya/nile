package nile

type routeTree struct {
	segments []*segment
}

type segment interface {
	Path() string
	HasParam() bool
	ParamName() (string, error)
	IntValue() (int, error)
	StringValue() (string, error)

	Match(path string) *segment
	Segments() []*segment
}

type routeSegment struct {
	path     string
	segments []*segment
}

func (rs *routeSegment) Match(path string) *segment {
	// Extract the first part of the path and match it against the current segment.
	return nil
}
