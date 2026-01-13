package rdf

import (
	"fmt"
	"testing"
)

func TestBlankNodeFactoryDefault_FmtConcise(t *testing.T) {
	if _a, _e := fmt.Sprintf("%#+v", bnDefault{v: 1}), "rdf.bnDefault{v:1}"; _a != _e {
		t.Errorf("expected %q, got %q", _e, _a)
	}
}
