package nile

// Payload is the structure that maps HTTP request payloads to structures and
// defines the methodology for how their contents get validated.
type Payload interface {
	// Validate ensures that the Payload and its values are structured properly.
	Validate() error
}
