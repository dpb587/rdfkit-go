package sparql

import (
	"testing"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

var testingBnode = blanknodeutil.NewStringMapper()

func assertQueryResponseResultBindingEqual(t *testing.T, e, a QueryResponseResultBinding) {
	if e.Name != a.Name {
		t.Fatalf("expected name `%+v`, got `%+v`", e.Name, a.Name)
	} else if e.Term == nil || a.Term == nil {
		t.Fatalf("expected %+v, got %+v", e.Term, a.Term)
	}

	switch eT := e.Term.(type) {
	case rdf.BlankNode:
		aT, ok := a.Term.(rdf.BlankNode)
		if !ok {
			t.Fatalf("expected bnode, got %[1]T(%[1]#+v)", a.Term)
		} else if !aT.TermEquals(eT) {
			t.Fatalf("expected bnode `%+v`, got `%+v`", eT, aT)
		}
	default:
		if !e.Term.TermEquals(a.Term) {
			t.Fatalf("expected %+v, got %+v", e.Term, a.Term)
		}
	}

}
