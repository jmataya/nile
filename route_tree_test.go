package nile

import "testing"

func TestSegmentParam(t *testing.T) {
	var tests = []struct {
		path         string
		wantHasParam bool
		wantPathName string
	}{
		{path: "hello", wantHasParam: false},
		{path: "", wantHasParam: false},
		{path: ":id", wantHasParam: true, wantPathName: "id"},
	}

	for _, test := range tests {
		seg := &routeSegment{path: test.path}
		gotHasParam := seg.HasParam()

		if test.wantHasParam != gotHasParam {
			t.Errorf("Segment.HasParam() == %v, want %v", gotHasParam, test.wantHasParam)
			continue
		}

		if gotHasParam {
			gotPathName, err := seg.ParamName()
			if err != nil {
				t.Errorf("Segment.ParamName() returned error %v", err)
				continue
			}

			if test.wantPathName != gotPathName {
				t.Errorf("Segment.ParamName() == %s, want %s", gotPathName, test.wantPathName)
			}
		}
	}
}

func TestNewSegment(t *testing.T) {
	var tests = []struct {
		path             string
		wantPath         string
		wantSegmentPaths []string
	}{
		{"/products", "products", []string{""}},
		{"/", "", []string{}},
		{"", "", []string{}},
		{"/products/new", "products", []string{"new"}},
		{"/products/:id", "products", []string{":id"}},
	}

	for _, test := range tests {
		seg := newSegment(test.path)

		gotPath := seg.Path()
		if gotPath != test.wantPath {
			t.Errorf("Segment.Path() == %s, want %s", gotPath, test.wantPath)
		}

		gotSegLen := len(seg.Segments())
		wantSegLen := len(test.wantSegmentPaths)

		if gotSegLen != wantSegLen {
			t.Errorf("Length of Segment.Path() == %d, want %d", gotSegLen, wantSegLen)
			continue
		}
	}
}

func TestSegmentMatch(t *testing.T) {
	var tests = []struct {
		path      string
		tryPath   string
		wantMatch bool
	}{
		{"/products", "/products", true},
		{"/products/", "/products", true},
		{"/products", "/products/", true},
		{"/products/", "/products/", true},
		{"/", "/", true},
		{"/", "", true},
		{"", "/", true},
		{"", "", true},
		{"/products", "/products/new", false},
		{"/products/:id", "/products/1", true},
		{"/products/new", "/products/1", false},
		{"/products/:id", "/products", false},
		{"/products/:id", "/products/", false},
	}

	for _, test := range tests {
		seg := newSegment(test.path)
		_, gotMatch := seg.Match(test.tryPath)

		if gotMatch != test.wantMatch {
			t.Errorf("Segment#Match(%s) for %s returned %v, want %v", test.tryPath, test.path, gotMatch, test.wantMatch)
		}
	}
}
