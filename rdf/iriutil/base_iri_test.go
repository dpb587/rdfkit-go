package iriutil

import (
	"testing"

	"github.com/dpb587/rdfkit-go/rdf"
)

func TestBaseIRI_RelativizeLiteral(t *testing.T) {
	tests := []struct {
		Base                string
		Input               rdf.IRI
		ExpectedRelativized string
		ExpectedOK          bool
	}{
		{
			Base:  "subpath",
			Input: "https://example.com/",
		},
		{
			Base:  "subpath",
			Input: "other",
		},
		{ // not supporting path behaviors for opaque uris
			Base:  "subpath",
			Input: "subpath/other",
		},
		{
			Base:  "subpath",
			Input: "/other",
		},
		{
			Base:                "subpath",
			Input:               "subpath?query#fragment",
			ExpectedRelativized: "?query#fragment",
			ExpectedOK:          true,
		},
		{
			Base:                "subpath",
			Input:               "subpath?query",
			ExpectedRelativized: "?query",
			ExpectedOK:          true,
		},
		{
			Base:                "subpath",
			Input:               "subpath#fragment",
			ExpectedRelativized: "#fragment",
			ExpectedOK:          true,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.Input), func(t *testing.T) {
			rb, err := ParseBaseIRI(tt.Base)
			if err != nil {
				t.Fatal("setup failed")
			}

			relativized, ok := rb.RelativizeIRI(tt.Input)
			if _a, _e := ok, tt.ExpectedOK; _a != _e {
				t.Fatalf("ok: expected %v, got %v", _e, _a)
			} else if _a, _e := relativized, tt.ExpectedRelativized; _a != _e {
				t.Fatalf("relativized: expected %v, got %v", _e, _a)
			}

			// not doing fuzz-parse here since Go URL parser roots relative paths (subpath => /subpath)
			// probably should support in ParsedIRI?
		})
	}
}

func TestBaseIRI_Relativize(t *testing.T) {
	tests := []struct {
		Base                string
		Input               rdf.IRI
		ExpectedRelativized string
		ExpectedOK          bool
	}{
		{
			Base:  "http://example.com/path/subpath",
			Input: "https://example.com/",
		},
		{
			Base:                "http://example.com/path/subpath",
			Input:               "http://example.com/path/subpath?query#fragment",
			ExpectedRelativized: "?query#fragment",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath",
			Input:               "http://example.com/path/subpath?query",
			ExpectedRelativized: "?query",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath",
			Input:               "http://example.com/path/other#fragment",
			ExpectedRelativized: "other#fragment",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath",
			Input:               "http://example.com/path/other",
			ExpectedRelativized: "other",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath",
			Input:               "http://example.com/other",
			ExpectedRelativized: "/other",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath",
			Input:               "http://example.com/path/subpath/other",
			ExpectedRelativized: "subpath/other",
			ExpectedOK:          true,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.Input), func(t *testing.T) {
			rb, err := ParseBaseIRI(tt.Base)
			if err != nil {
				t.Fatal("setup failed")
			}

			relativized, ok := rb.RelativizeIRI(tt.Input)
			if _a, _e := ok, tt.ExpectedOK; _a != _e {
				t.Fatalf("ok: expected %v, got %v", _e, _a)
			} else if _a, _e := relativized, tt.ExpectedRelativized; _a != _e {
				t.Fatalf("relativized: expected %v, got %v", _e, _a)
			}

			if !ok {
				return
			}

			parsed, err := rb.Parse(relativized)
			if err != nil {
				t.Fatalf("fuzz-parse: unexpected error: %v", err)
			} else if _a, _e := rdf.IRI(parsed.String()), tt.Input; _a != _e {
				t.Fatalf("fuzz-parse: parsed: expected %v, got %v", _e, _a)
			}
		})
	}
}

func TestBaseIRI_RelativizeWithQueryFragment(t *testing.T) {
	tests := []struct {
		Base                string
		Input               rdf.IRI
		ExpectedRelativized string
		ExpectedOK          bool
	}{
		{
			Base:  "http://example.com/path/subpath?query#fragment",
			Input: "https://example.com/",
		},
		{
			Base:                "http://example.com/path/subpath?query#fragment",
			Input:               "http://example.com/path/subpath?query#fragment",
			ExpectedRelativized: "",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath?query#fragment",
			Input:               "http://example.com/path/subpath?query",
			ExpectedRelativized: "?query",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath?query#fragment",
			Input:               "http://example.com/path/subpath#fragment",
			ExpectedRelativized: "subpath#fragment",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath?query#fragment",
			Input:               "http://example.com/path/other",
			ExpectedRelativized: "other",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath?query#fragment",
			Input:               "http://example.com/other",
			ExpectedRelativized: "/other",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath?query#fragment",
			Input:               "http://example.com/path/subpath/other",
			ExpectedRelativized: "subpath/other",
			ExpectedOK:          true,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.Input), func(t *testing.T) {
			rb, err := ParseBaseIRI(tt.Base)
			if err != nil {
				t.Fatal("setup failed")
			}

			relativized, ok := rb.RelativizeIRI(tt.Input)
			if _a, _e := ok, tt.ExpectedOK; _a != _e {
				t.Fatalf("ok: expected %v, got %v", _e, _a)
			} else if _a, _e := relativized, tt.ExpectedRelativized; _a != _e {
				t.Fatalf("relativized: expected %v, got %v", _e, _a)
			}

			if !ok {
				return
			}

			parsed, err := rb.Parse(relativized)
			if err != nil {
				t.Fatalf("fuzz-parse: unexpected error: %v", err)
			} else if _a, _e := rdf.IRI(parsed.String()), tt.Input; _a != _e {
				t.Fatalf("fuzz-parse: parsed: expected %v, got %v", _e, _a)
			}
		})
	}
}

func TestBaseIRI_RelativizeWithQuery(t *testing.T) {
	tests := []struct {
		Base                string
		Input               rdf.IRI
		ExpectedRelativized string
		ExpectedOK          bool
	}{
		{
			Base:  "http://example.com/path/subpath?query",
			Input: "https://example.com/",
		},
		{
			Base:                "http://example.com/path/subpath?query",
			Input:               "http://example.com/path/subpath?query#fragment",
			ExpectedRelativized: "#fragment",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath?query",
			Input:               "http://example.com/path/subpath?query",
			ExpectedRelativized: "",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath?query",
			Input:               "http://example.com/path/subpath#fragment",
			ExpectedRelativized: "subpath#fragment",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath?query",
			Input:               "http://example.com/path/other",
			ExpectedRelativized: "other",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath?query",
			Input:               "http://example.com/other",
			ExpectedRelativized: "/other",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath?query",
			Input:               "http://example.com/path/subpath/other",
			ExpectedRelativized: "subpath/other",
			ExpectedOK:          true,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.Input), func(t *testing.T) {
			rb, err := ParseBaseIRI(tt.Base)
			if err != nil {
				t.Fatal("setup failed")
			}

			relativized, ok := rb.RelativizeIRI(tt.Input)
			if _a, _e := ok, tt.ExpectedOK; _a != _e {
				t.Fatalf("ok: expected %v, got %v", _e, _a)
			} else if _a, _e := relativized, tt.ExpectedRelativized; _a != _e {
				t.Fatalf("relativized: expected %v, got %v", _e, _a)
			}

			if !ok {
				return
			}

			parsed, err := rb.Parse(relativized)
			if err != nil {
				t.Fatalf("fuzz-parse: unexpected error: %v", err)
			} else if _a, _e := rdf.IRI(parsed.String()), tt.Input; _a != _e {
				t.Fatalf("fuzz-parse: parsed: expected %v, got %v", _e, _a)
			}
		})
	}
}

func TestBaseIRI_RelativizeWithFragment(t *testing.T) {
	tests := []struct {
		Base                string
		Input               rdf.IRI
		ExpectedRelativized string
		ExpectedOK          bool
	}{
		{
			Base:  "http://example.com/path/subpath#fragment",
			Input: "https://example.com/",
		},
		{
			Base:                "http://example.com/path/subpath#fragment",
			Input:               "http://example.com/path/subpath#fragment",
			ExpectedRelativized: "",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath#fragment",
			Input:               "http://example.com/path/other",
			ExpectedRelativized: "other",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath#fragment",
			Input:               "http://example.com/other",
			ExpectedRelativized: "/other",
			ExpectedOK:          true,
		},
		{
			Base:                "http://example.com/path/subpath#fragment",
			Input:               "http://example.com/path/subpath/other",
			ExpectedRelativized: "subpath/other",
			ExpectedOK:          true,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.Input), func(t *testing.T) {
			rb, err := ParseBaseIRI(tt.Base)
			if err != nil {
				t.Fatal("setup failed")
			}

			relativized, ok := rb.RelativizeIRI(tt.Input)
			if _a, _e := ok, tt.ExpectedOK; _a != _e {
				t.Fatalf("ok: expected %v, got %v", _e, _a)
			} else if _a, _e := relativized, tt.ExpectedRelativized; _a != _e {
				t.Fatalf("relativized: expected %v, got %v", _e, _a)
			}

			if !ok {
				return
			}

			parsed, err := rb.Parse(relativized)
			if err != nil {
				t.Fatalf("fuzz-parse: unexpected error: %v", err)
			} else if _a, _e := rdf.IRI(parsed.String()), tt.Input; _a != _e {
				t.Fatalf("fuzz-parse: parsed: expected %v, got %v", _e, _a)
			}
		})
	}
}
