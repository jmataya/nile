package routing

import (
	"testing"
)

func TestSegmentChildren(t *testing.T) {
	var tests = []struct {
		insertions   []string
		wantChildren []string
	}{
		{
			insertions:   []string{"abc", "def"},
			wantChildren: []string{"def", "abc"},
		},
		{
			insertions:   []string{"def", "abc"},
			wantChildren: []string{"def", "abc"},
		},
		{
			insertions:   []string{"/products/abc", "/products/def"},
			wantChildren: []string{"products"},
		},
		{
			insertions:   []string{"/products/def", "/products/abc"},
			wantChildren: []string{"products"},
		},
	}

	var errAdding bool
	for testIdx, test := range tests {
		root := NewSegment("/")
		for _, toInsert := range test.insertions {
			child := NewSegment(toInsert)
			if err := root.AddChild(child); err != nil {
				t.Errorf("Segment.AddChild({path: %s}) error, want <nil>, got %v", child.Path(), err)
				errAdding = true
			}
		}

		if errAdding {
			continue
		}

		gotChildren := root.Children()
		if len(gotChildren) != len(test.wantChildren) {
			t.Errorf("len(Segment.Children()), want %d, got %d", len(test.wantChildren), len(gotChildren))
			continue
		}

		for idx, wantChild := range test.wantChildren {
			if wantChild != gotChildren[idx].Path() {
				t.Errorf("Test %d: Path of child %d, want %s, got %s", testIdx, idx, wantChild, gotChildren[idx].Path())
			}
		}
	}
}

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
