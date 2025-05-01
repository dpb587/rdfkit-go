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

func TestPrefixMap_NewPrefixMap(t *testing.T) {
	s0 := PrefixMap{
		"ex": "http://example.com/",
	}

	originalCom := rdf.IRI("http://example.com/path")
	originalOrg := rdf.IRI("http://example.org/path")

	{
		prefix, reference, ok := s0.CompactIRI(originalCom)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := prefix, "ex"; _a != _e {
			t.Fatalf("prefix: expected %v, got %v", _e, _a)
		} else if _a, _e := reference, "path"; _a != _e {
			t.Fatalf("reference: expected %v, got %v", _e, _a)
		}

		_, _, ok = s0.CompactIRI(originalOrg)
		if _a, _e := ok, false; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		}
	}

	s1 := s0.NewPrefixMap(PrefixMapping{
		Prefix:   "ex",
		Expanded: "http://example.org/",
	})

	{ // s0 unchanged
		prefix, reference, ok := s0.CompactIRI(originalCom)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := prefix, "ex"; _a != _e {
			t.Fatalf("prefix: expected %v, got %v", _e, _a)
		} else if _a, _e := reference, "path"; _a != _e {
			t.Fatalf("reference: expected %v, got %v", _e, _a)
		}

		_, _, ok = s0.CompactIRI(originalOrg)
		if _a, _e := ok, false; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		}
	}

	{
		prefix, reference, ok := s1.CompactIRI(originalOrg)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := prefix, "ex"; _a != _e {
			t.Fatalf("prefix: expected %v, got %v", _e, _a)
		} else if _a, _e := reference, "path"; _a != _e {
			t.Fatalf("reference: expected %v, got %v", _e, _a)
		}

		_, _, ok = s1.CompactIRI(originalCom)
		if _a, _e := ok, false; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		}
	}
}
