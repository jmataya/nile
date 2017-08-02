package routing

// Match is a structure that contains the data for when a request matches an
// endpoint.
type Match struct {
	Endpoint      Endpoint
	RequestURI    string
	RequestMethod string
	params        map[string]string
}

// NewMatch creates a new Match object.
func NewMatch(endPt Endpoint, path string, method string) *Match {
	return &Match{
		Endpoint:      endPt,
		RequestURI:    path,
		RequestMethod: method,
		params:        map[string]string{},
	}
}

// AddParam adds a parameter value to the Match.
func (m *Match) AddParam(key, value string) {
	m.params[key] = value
}

// Param gets the value of a param, if it exists.
func (m *Match) Param(key string) (string, bool) {
	value, ok := m.params[key]
	return value, ok
}
