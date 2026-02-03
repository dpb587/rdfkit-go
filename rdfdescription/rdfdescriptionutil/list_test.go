package rdfdescriptionutil_test

import (
	"testing"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/rdfdescription/rdfdescriptionutil"
)

func TestNewObjectValueListStatement_EmptyList(t *testing.T) {
	stmt := rdfdescriptionutil.NewObjectValueListStatement(
		rdf.IRI("http://example.com/hasItems"),
	)

	objStmt, ok := stmt.(rdfdescription.ObjectStatement)
	if !ok {
		t.Fatalf("expected ObjectStatement, but got %T", stmt)
	}

	if _a, _e := objStmt.Predicate, rdf.IRI("http://example.com/hasItems"); _a != _e {
		t.Fatalf("expected %v, but got %v", _e, _a)
	} else if _a, _e := objStmt.Object, rdfiri.Nil_List; _a != _e {
		t.Fatalf("expected %v, but got %v", _e, _a)
	}
}

func TestNewObjectValueListStatement_SingleValue(t *testing.T) {
	stmt := rdfdescriptionutil.NewObjectValueListStatement(
		rdf.IRI("http://example.com/hasItems"),
		rdf.IRI("http://example.com/item1"),
	)

	anonStmt, ok := stmt.(rdfdescription.AnonResourceStatement)
	if !ok {
		t.Fatalf("expected AnonResourceStatement, but got %T", stmt)
	}

	if _a, _e := anonStmt.Predicate, rdf.IRI("http://example.com/hasItems"); _a != _e {
		t.Fatalf("expected %v, but got %v", _e, _a)
	} else if _a, _e := len(anonStmt.AnonResource.Statements), 2; _a != _e {
		t.Fatalf("expected %v statements, but got %v", _e, _a)
	}

	{
		firstStmt, ok := anonStmt.AnonResource.Statements[0].(rdfdescription.ObjectStatement)
		if !ok {
			t.Fatalf("expected ObjectStatement for rdf:first, but got %T", anonStmt.AnonResource.Statements[0])
		}

		if _a, _e := firstStmt.Predicate, rdfiri.First_Property; _a != _e {
			t.Fatalf("expected %v, but got %v", _e, _a)
		} else if _a, _e := firstStmt.Object, rdf.IRI("http://example.com/item1"); _a != _e {
			t.Fatalf("expected %v, but got %v", _e, _a)
		}
	}

	{
		restStmt, ok := anonStmt.AnonResource.Statements[1].(rdfdescription.ObjectStatement)
		if !ok {
			t.Fatalf("expected ObjectStatement for rdf:rest, but got %T", anonStmt.AnonResource.Statements[1])
		}

		if _a, _e := restStmt.Predicate, rdfiri.Rest_Property; _a != _e {
			t.Fatalf("expected %v, but got %v", _e, _a)
		} else if _a, _e := restStmt.Object, rdfiri.Nil_List; _a != _e {
			t.Fatalf("expected %v, but got %v", _e, _a)
		}
	}
}

func TestNewObjectValueListStatement_ThreeValues(t *testing.T) {
	stmt := rdfdescriptionutil.NewObjectValueListStatement(
		rdf.IRI("http://example.com/hasItems"),
		rdf.IRI("http://example.com/item1"),
		rdf.Literal{
			LexicalForm: "item2",
			Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#string"),
		},
		rdf.IRI("http://example.com/item3"),
	)

	anonStmt, ok := stmt.(rdfdescription.AnonResourceStatement)
	if !ok {
		t.Fatalf("expected AnonResourceStatement, but got %T", stmt)
	}

	if _a, _e := anonStmt.Predicate, rdf.IRI("http://example.com/hasItems"); _a != _e {
		t.Fatalf("expected %v, but got %v", _e, _a)
	} else if _a, _e := len(anonStmt.AnonResource.Statements), 2; _a != _e {
		t.Fatalf("expected %v statements in node1, but got %v", _e, _a)
	}

	firstStmt1, ok := anonStmt.AnonResource.Statements[0].(rdfdescription.ObjectStatement)
	if !ok {
		t.Fatalf("expected ObjectStatement for node1 rdf:first, but got %T", anonStmt.AnonResource.Statements[0])
	}

	if _a, _e := firstStmt1.Predicate, rdfiri.First_Property; _a != _e {
		t.Fatalf("expected %v, but got %v", _e, _a)
	} else if _a, _e := firstStmt1.Object, rdf.IRI("http://example.com/item1"); _a != _e {
		t.Fatalf("expected %v, but got %v", _e, _a)
	}

	restStmt1, ok := anonStmt.AnonResource.Statements[1].(rdfdescription.AnonResourceStatement)
	if !ok {
		t.Fatalf("expected AnonResourceStatement for node1 rdf:rest, but got %T", anonStmt.AnonResource.Statements[1])
	}

	if _a, _e := restStmt1.Predicate, rdfiri.Rest_Property; _a != _e {
		t.Fatalf("expected %v, but got %v", _e, _a)
	} else if _a, _e := len(restStmt1.AnonResource.Statements), 2; _a != _e {
		t.Fatalf("expected %v statements in node2, but got %v", _e, _a)
	}

	{
		firstStmt2, ok := restStmt1.AnonResource.Statements[0].(rdfdescription.ObjectStatement)
		if !ok {
			t.Fatalf("expected ObjectStatement for node2 rdf:first, but got %T", restStmt1.AnonResource.Statements[0])
		}

		if _a, _e := firstStmt2.Predicate, rdfiri.First_Property; _a != _e {
			t.Fatalf("expected %v, but got %v", _e, _a)
		}

		literal2, ok := firstStmt2.Object.(rdf.Literal)
		if !ok {
			t.Fatalf("expected Literal for node2 object, but got %T", firstStmt2.Object)
		}

		if _a, _e := literal2.LexicalForm, "item2"; _a != _e {
			t.Fatalf("expected %v, but got %v", _e, _a)
		}
	}

	restStmt2, ok := restStmt1.AnonResource.Statements[1].(rdfdescription.AnonResourceStatement)
	if !ok {
		t.Fatalf("expected AnonResourceStatement for node2 rdf:rest, but got %T", restStmt1.AnonResource.Statements[1])
	}

	if _a, _e := restStmt2.Predicate, rdfiri.Rest_Property; _a != _e {
		t.Fatalf("expected %v, but got %v", _e, _a)
	}

	{
		firstStmt3, ok := restStmt2.AnonResource.Statements[0].(rdfdescription.ObjectStatement)
		if !ok {
			t.Fatalf("expected ObjectStatement for node3 rdf:first, but got %T", restStmt2.AnonResource.Statements[0])
		}

		if _a, _e := firstStmt3.Predicate, rdfiri.First_Property; _a != _e {
			t.Fatalf("expected %v, but got %v", _e, _a)
		} else if _a, _e := firstStmt3.Object, rdf.IRI("http://example.com/item3"); _a != _e {
			t.Fatalf("expected %v, but got %v", _e, _a)
		}
	}

	restStmt3, ok := restStmt2.AnonResource.Statements[1].(rdfdescription.ObjectStatement)
	if !ok {
		t.Fatalf("expected ObjectStatement for node3 rdf:rest, but got %T", restStmt2.AnonResource.Statements[1])
	}

	if _a, _e := restStmt3.Predicate, rdfiri.Rest_Property; _a != _e {
		t.Fatalf("expected %v, but got %v", _e, _a)
	} else if _a, _e := restStmt3.Object, rdfiri.Nil_List; _a != _e {
		t.Fatalf("expected %v, but got %v", _e, _a)
	}
}
