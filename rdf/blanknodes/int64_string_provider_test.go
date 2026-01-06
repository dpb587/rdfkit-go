package blanknodes

import (
	"testing"

	"github.com/dpb587/rdfkit-go/rdf"
)

func TestInt64StringProvider(t *testing.T) {
	subject := NewInt64StringProvider("b%d")

	bn1 := rdf.NewBlankNode()
	bn2 := rdf.NewBlankNode()

	if _a, _e := subject.GetBlankNodeString(bn1), "b0"; _a != _e {
		t.Fatalf("expected %s, got %s", _e, _a)
	} else if _a, _e := subject.GetBlankNodeString(bn2), "b1"; _a != _e {
		t.Fatalf("expected %s, got %s", _e, _a)
	} else if _a, _e := subject.GetBlankNodeString(bn1), "b0"; _a != _e {
		t.Fatalf("expected %s, got %s", _e, _a)
	} else if _a, _e := subject.GetBlankNodeString(rdf.NewBlankNode()), "b2"; _a != _e {
		t.Fatalf("expected %s, got %s", _e, _a)
	}
}
