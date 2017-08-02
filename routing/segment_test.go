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
		wantParams map[string]string
	}{
		{"/products", "GET", "products", "GET", true, map[string]string{}},
		{"/products", "GET", "/products", "GET", true, map[string]string{}},
		{"/products", "GET", "/products/new", "GET", false, map[string]string{}},
		{"/products/new", "GET", "/products/new", "GET", true, map[string]string{}},
		{"/products/:id", "GET", "/products/1", "GET", true, map[string]string{"id": "1"}},
		{"/products/:id", "GET", "/products/1", "POST", false, map[string]string{}},
		{"/products", "GET", "/product", "GET", false, map[string]string{}},
		{"/products", "POST", "/products", "GET", false, map[string]string{}},
		{"/products/:id/edit", "PATCH", "/products/4/edit", "PATCH", true, map[string]string{"id": "4"}},
	}

	for _, test := range tests {
		seg, err := NewSegmentEndpoint(test.uri, test.method)
		if err != nil {
			t.Errorf("NewSegmentEndpoint(%s, %s) error, want <nil>, got %v", test.uri, test.method, err)
			continue
		}

		gotMatch, gotMatches := seg.Matches(test.wantPath, test.wantMethod)
		if test.wantMatch != gotMatches {
			t.Errorf("Segment.Matches(%s, %s) match, want %v, got %v", test.wantPath, test.wantMethod, test.wantMatch, gotMatches)
			continue
		}

		if !gotMatches {
			continue
		}

		if gotMatch.RequestMethod != test.wantMethod {
			t.Errorf("Match.RequestMethod, want %s, got %s", test.wantMethod, gotMatch.RequestMethod)
			continue
		}

		if gotMatch.RequestURI != test.wantPath {
			t.Errorf("Match.RequestURI, want %s, got %s", test.wantPath, gotMatch.RequestURI)
			continue
		}

		for paramName, paramValue := range test.wantParams {
			actualValue, found := gotMatch.Param(paramName)
			if !found || actualValue != paramValue {
				t.Errorf("Match.Param(%s), want (%s, true), got (%s, %v)", paramName, paramValue, actualValue, found)
			}
		}
	}
}
