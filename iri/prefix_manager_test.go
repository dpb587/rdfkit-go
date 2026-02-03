package iri

import (
	"testing"
)

func TestPrefixMap_LongestMatch(t *testing.T) {
	s := NewPrefixManager(PrefixMappingList{
		PrefixMapping{
			Prefix:   "ex",
			Expanded: "http://example.com/",
		},
		PrefixMapping{
			Prefix:   "exp",
			Expanded: "http://example.com/path/",
		},
	})

	{
		original := "http://example.com/path/1234"

		pr, ok := s.CompactPrefix(original)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := pr.Prefix, "exp"; _a != _e {
			t.Fatalf("prefix: expected %v, got %v", _e, _a)
		} else if _a, _e := pr.Reference, "1234"; _a != _e {
			t.Fatalf("reference: expected %v, got %v", _e, _a)
		}

		expanded, ok := s.ExpandPrefix(pr)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := expanded, original; _a != _e {
			t.Fatalf("expanded: expected %v, got %v", _e, _a)
		}
	}

	{
		original := "http://example.com/1234"

		pr, ok := s.CompactPrefix(original)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := pr.Prefix, "ex"; _a != _e {
			t.Fatalf("prefix: expected %v, got %v", _e, _a)
		} else if _a, _e := pr.Reference, "1234"; _a != _e {
			t.Fatalf("reference: expected %v, got %v", _e, _a)
		}

		expanded, ok := s.ExpandPrefix(pr)
		if _a, _e := ok, true; _a != _e {
			t.Fatalf("ok: expected %v, got %v", _e, _a)
		} else if _a, _e := expanded, original; _a != _e {
			t.Fatalf("expanded: expected %v, got %v", _e, _a)
		}
	}
}
