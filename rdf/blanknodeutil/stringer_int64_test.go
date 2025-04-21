package blanknodeutil

import (
	"testing"

	"github.com/dpb587/rdfkit-go/rdf"
)

func TestStringerInt64(t *testing.T) {
	subject := NewStringerInt64()

	bn1 := rdf.NewBlankNode()
	bn2 := rdf.NewBlankNode()

	if _a, _e := subject.GetBlankNodeIdentifier(bn1), "b0"; _a != _e {
		t.Fatalf("expected %s, got %s", _e, _a)
	} else if _a, _e := subject.GetBlankNodeIdentifier(bn2), "b1"; _a != _e {
		t.Fatalf("expected %s, got %s", _e, _a)
	} else if _a, _e := subject.GetBlankNodeIdentifier(bn1), "b0"; _a != _e {
		t.Fatalf("expected %s, got %s", _e, _a)
	} else if _a, _e := subject.GetBlankNodeIdentifier(rdf.NewBlankNode()), "b2"; _a != _e {
		t.Fatalf("expected %s, got %s", _e, _a)
	}
}
