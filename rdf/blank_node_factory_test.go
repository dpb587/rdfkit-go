package rdf

import (
	"fmt"
	"regexp"
	"testing"
)

func TestBlankNodeFactory_FmtConcise(t *testing.T) {
	if _a, _e := fmt.Sprintf("%#+v", bn{v: 1, s: &bnF{}}), regexp.MustCompile(`^rdf\.bn\{v:1, s:\(\*rdf\.bnF\)\([^)]+\)\}$`); !_e.MatchString(_a) {
		t.Errorf("expected %q, got %q", _e, _a)
	}
}
