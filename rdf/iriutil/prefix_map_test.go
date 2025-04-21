package iriutil

import (
	"testing"

	"github.com/dpb587/rdfkit-go/rdf"
)

func TestPrefixMap_LongestMatch(t *testing.T) {
	s := PrefixMap{
		"ex":  "http://example.com/",
		"exp": "http://example.com/path/",
	}

	{
		original := rdf.IRI("http://example.com/path/1234")

		prefix, reference, ok := s.CompactIRI(original)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := prefix, "exp"; _a != _e {
			t.Fatalf("prefix: expected %v, got %v", _e, _a)
		} else if _a, _e := reference, "1234"; _a != _e {
			t.Fatalf("reference: expected %v, got %v", _e, _a)
		}

		expanded, ok := s.ExpandIRI(prefix, reference)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := expanded, original; _a != _e {
			t.Fatalf("expanded: expected %v, got %v", _e, _a)
		}
	}

	{
		original := rdf.IRI("http://example.com/1234")

		prefix, reference, ok := s.CompactIRI(original)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := prefix, "ex"; _a != _e {
			t.Fatalf("prefix: expected %v, got %v", _e, _a)
		} else if _a, _e := reference, "1234"; _a != _e {
			t.Fatalf("reference: expected %v, got %v", _e, _a)
		}

		expanded, ok := s.ExpandIRI(prefix, reference)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := expanded, original; _a != _e {
			t.Fatalf("expanded: expected %v, got %v", _e, _a)
		}
	}
}
