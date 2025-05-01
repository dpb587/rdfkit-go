package iriutil

import (
	"testing"

	"github.com/dpb587/rdfkit-go/rdf"
)

func TestPrefixTracker(t *testing.T) {
	s := NewPrefixTracker(PrefixMap{
		"ex":  "http://example.com/",
		"exp": "http://example.com/path/",
	})

	{
		original := rdf.IRI("http://example.com/path/1234")

		prefix, reference, ok := s.CompactPrefix(original)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := prefix, "exp"; _a != _e {
			t.Fatalf("prefix: expected %v, got %v", _e, _a)
		} else if _a, _e := reference, "1234"; _a != _e {
			t.Fatalf("reference: expected %v, got %v", _e, _a)
		}
	}

	usedPrefixes := s.GetUsedPrefixMappings()
	if _a, _e := len(usedPrefixes), 1; _a != _e {
		t.Fatalf("len: expected %v, got %v", _e, _a)
	} else if _a, _e := usedPrefixes[0].Prefix, "exp"; _a != _e {
		t.Fatalf("prefix: expected %v, got %v", _e, _a)
	} else if _a, _e := usedPrefixes[0].Expanded, "http://example.com/path/"; _a != _e {
		t.Fatalf("expanded: expected %v, got %v", _e, _a)
	}
}
