package routing

import "testing"

func TestSegmentMatching(t *testing.T) {
	var tests = []struct {
		uri        string
		method     string
		wantPath   string
		wantMethod string
		wantMatch  bool
	}{
		{"/products", "GET", "products", "GET", true},
		{"/products", "GET", "/products", "GET", true},
		{"/products", "GET", "/products/new", "GET", false},
		{"/products/new", "GET", "/products/new", "GET", true},
		{"/products/:id", "GET", "/products/1", "GET", true},
		{"/products/:id", "GET", "/products/1", "POST", false},
		{"/products", "GET", "/product", "GET", false},
		{"/products", "POST", "/products", "GET", false},
	}

	for _, test := range tests {
		seg, err := NewSegmentEndpoint(test.uri, test.method)
		if err != nil {
			t.Errorf("NewSegmentEndpoint(%s, %s) error, want <nil>, got %v", test.uri, test.method, err)
			continue
		}

		gotEndPt, gotMatch := seg.Matches(test.wantPath, test.wantMethod)
		if test.wantMatch != gotMatch {
			t.Errorf("Segment.Matches(%s, %s) match, want %v, got %v", test.wantPath, test.wantMethod, test.wantMatch, gotMatch)
			continue
		}

		if !gotMatch {
			continue
		}

		if gotEndPt.Method() != test.wantMethod {
			t.Errorf("Endpoint.Method(), want %s, got %s", test.wantMethod, gotEndPt.Method())
		}
	}
}
