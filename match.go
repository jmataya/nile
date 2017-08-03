package nile

// match is a structure that contains the data for when a request matches an
// endpoint.
type match struct {
	Segment    *segment
	Context    *context
	RequestURI string
}

// newMatch creates a new Match object.
func newMatch(seg *segment, path string) *match {
	return &match{
		Segment:    seg,
		Context:    &context{params: map[string]string{}},
		RequestURI: path,
	}
}

// AddParam adds a parameter value to the Match.
func (m *match) AddParam(key, value string) {
	m.Context.addParam(key, value)
}

// Param gets the value of a param, if it exists.
func (m *match) Param(key string) (string, bool) {
	return m.Context.Param(key)
}
