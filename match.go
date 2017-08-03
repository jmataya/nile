package nile

// Match is a structure that contains the data for when a request matches an
// endpoint.
type Match struct {
	Segment    Segment
	Context    *Context
	RequestURI string
}

// NewMatch creates a new Match object.
func NewMatch(segment Segment, path string) *Match {
	return &Match{
		Segment:    segment,
		Context:    &Context{params: map[string]string{}},
		RequestURI: path,
	}
}

// AddParam adds a parameter value to the Match.
func (m *Match) AddParam(key, value string) {
	m.Context.addParam(key, value)
}

// Param gets the value of a param, if it exists.
func (m *Match) Param(key string) (string, bool) {
	return m.Context.Param(key)
}
