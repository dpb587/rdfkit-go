package rdf

import "testing"

func TestDefaultBlankNodeFactory(t *testing.T) {
	g := DefaultBlankNodeFactory
	bn1 := g.NewBlankNode()
	bn2 := g.NewBlankNode()

	if bn1 == bn2 {
		t.Fatalf("expected different blank nodes, got same: %v", bn1)
	} else if bn1.(blankNode).identifier == bn2.(blankNode).identifier {
		t.Fatalf("expected different blank node identifiers, got same: %v", bn1.(blankNode).identifier)
	} else if bn1.TermEquals(bn2) {
		t.Fatalf("expected different blank nodes, got same: %v", bn1)
	} else if !bn1.TermEquals(bn1) {
		t.Fatalf("expected same blank node, got different: %v", bn1)
	} else if !bn2.TermEquals(bn2) {
		t.Fatalf("expected same blank node, got different: %v", bn2)
	}

	bn := blankNode{
		identifier: bn1.(blankNode).identifier,
	}

	if !bn.TermEquals(bn1) {
		t.Fatalf("expected same blank node, got different: %v", bn)
	}
}
