package blanknodeutil

import (
	"testing"

	"github.com/dpb587/rdfkit-go/rdf"
)

func TestMapper_ReproducibleInput(t *testing.T) {
	subject := NewMapper(rdf.DefaultBlankNodeFactory)

	input := rdf.NewBlankNode()
	output := subject.MapBlankNode(input)

	if _a, _e := output, input; input.TermEquals(output) {
		t.Fatalf("expected %v, but got %v", _a, _e)
	}

	output2 := subject.MapBlankNode(input)

	if _a, _e := output2, output; !output.TermEquals(output2) {
		t.Fatalf("expected %v, but got %v", _a, _e)
	}
}
