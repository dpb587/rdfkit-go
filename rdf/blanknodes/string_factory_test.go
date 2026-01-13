package blanknodes

import (
	"fmt"
	"regexp"
	"testing"
)

func TestStringFactory_FmtConcise(t *testing.T) {
	if _a, _e := fmt.Sprintf("%#+v", bnString{v: "hello", s: &bnStringF{}}), regexp.MustCompile(`^blanknodes\.bnString\{v:"hello", s:\(\*blanknodes\.bnStringF\)\([^)]+\)\}$`); !_e.MatchString(_a) {
		t.Errorf("expected %q, got %q", _e, _a)
	}
}
