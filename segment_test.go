package nile

import (
	"log"
	"net/http"
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
		root := newSegment("/")
		for _, toInsert := range test.insertions {
			child := newSegment(toInsert)
			if err := root.AddChild(child); err != nil {
				t.Errorf("Segment.AddChild({path: %s}) error, want <nil>, got %v", child.Path, err)
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
			if wantChild != gotChildren[idx].Path {
				t.Errorf("Test %d: Path of child %d, want %s, got %s", testIdx, idx, wantChild, gotChildren[idx].Path)
			}
		}
	}
}

func TestAdvancedSegmentMatching(t *testing.T) {
	var called bool

	badHandler := func(c Context) Response {
		called = false
		resp := map[string]string{"message": "hello"}
		return NewGenericResponse(http.StatusOK, resp)
	}

	goodHandler := func(c Context) Response {
		called = true
		resp := map[string]string{"message": "hello"}
		return NewGenericResponse(http.StatusOK, resp)
	}

	dynamic, err := newSegmentEndpoint("/products/:id", http.MethodGet, badHandler)
	if err != nil {
		t.Errorf("newSegmentEndpoint(%s, %s, fn), want <nil> err, got %v err", "/products/:id", http.MethodGet, err)
		return
	}

	static, err := newSegmentEndpoint("/products/new", http.MethodGet, goodHandler)
	if err != nil {
		t.Errorf("newSegmentEndpoint(%s, %s, fn), want <nil> err, got %v err", "/products/new", http.MethodGet, err)
		return
	}

	merged, err := mergeSegments(dynamic, static)
	if err != nil {
		t.Errorf("mergeSegments(), want <nil> err, got %v err", err)
		return
	}

	endpoints := merged.Endpoints()
	log.Printf("%+v", endpoints)

	gotMatch, hasMatch := merged.Matches("/products/new")
	if !hasMatch {
		t.Errorf("Segment.Matches(%s), want matches %v, got matches %v", "/products/new", true, hasMatch)
		return
	}

	gotEndpoint, found := gotMatch.Segment.Endpoint("GET")
	if !found {
		t.Errorf("Segment.Endpoint(GET), want found true, got found %v", found)
		return
	}

	gotEndpoint.Handler()(nil)
	if !called {
		t.Error("Expected called to have been executed.")
	}
}

func TestSegmentMatching(t *testing.T) {
	dummyHandler := func(context Context) Response {
		resp := map[string]string{"message": "hello"}
		return NewGenericResponse(http.StatusOK, resp)
	}

	var tests = []struct {
		uri        string
		method     string
		wantPath   string
		wantMatch  bool
		wantParams map[string]string
	}{
		{"/products", "GET", "products", true, map[string]string{}},
		{"/products", "GET", "/products", true, map[string]string{}},
		{"/products", "GET", "/products/new", false, map[string]string{}},
		{"/products/new", "GET", "/products/new", true, map[string]string{}},
		{"/products/:id", "GET", "/products/1", true, map[string]string{"id": "1"}},
		{"/products", "GET", "/product", false, map[string]string{}},
		{"/products/:id/edit", "PATCH", "/products/4/edit", true, map[string]string{"id": "4"}},
	}

	for _, test := range tests {
		seg, err := newSegmentEndpoint(test.uri, test.method, dummyHandler)
		if err != nil {
			t.Errorf("NewSegmentEndpoint(%s, %s) error, want <nil>, got %v", test.uri, test.method, err)
			continue
		}

		gotMatch, gotMatches := seg.Matches(test.wantPath)
		if test.wantMatch != gotMatches {
			t.Errorf("Segment.Matches(%s) match, want %v, got %v", test.wantPath, test.wantMatch, gotMatches)
			continue
		}

		if !gotMatches {
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
