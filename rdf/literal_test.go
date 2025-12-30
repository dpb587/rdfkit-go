package rdf

import (
	"testing"
)

func TestLiteral_TermEquals_NotType(t *testing.T) {
	v := Literal{
		Datatype:    "String",
		LexicalForm: "1234",
	}.TermEquals(IRI(""))
	if _e, _a := false, v; _e != _a {
		t.Errorf("expected %v, got %v", _e, _a)
	}
}

func TestLiteral_TermEquals_NotDatatype(t *testing.T) {
	v := Literal{
		Datatype:    "String",
		LexicalForm: "1234",
	}.TermEquals(Literal{
		Datatype:    "Integer",
		LexicalForm: "1234",
	})
	if _e, _a := false, v; _e != _a {
		t.Errorf("expected %v, got %v", _e, _a)
	}
}

func TestLiteral_TermEquals_NotLexical(t *testing.T) {
	v := Literal{
		Datatype:    "Integer",
		LexicalForm: "1234",
	}.TermEquals(Literal{
		Datatype:    "Integer",
		LexicalForm: "1234.0",
	})
	if _e, _a := false, v; _e != _a {
		t.Errorf("expected %v, got %v", _e, _a)
	}
}

func TestLiteral_TermEquals_NotQualifiersOmitted(t *testing.T) {
	v := Literal{
		Datatype:    "LangString",
		LexicalForm: "hello",
		Tag: LanguageLiteralTag{
			Language: "en",
		},
	}.TermEquals(Literal{
		Datatype:    "LangString",
		LexicalForm: "hello",
	})
	if _e, _a := false, v; _e != _a {
		t.Errorf("expected %v, got %v", _e, _a)
	}
}

func TestLiteral_TermEquals_NotQualifiersOmittedEmpty(t *testing.T) {
	v := Literal{
		Datatype:    "LangString",
		LexicalForm: "hello",
		Tag: LanguageLiteralTag{
			Language: "",
		},
	}.TermEquals(Literal{
		Datatype:    "LangString",
		LexicalForm: "hello",
	})
	if _e, _a := false, v; _e != _a {
		t.Errorf("expected %v, got %v", _e, _a)
	}
}

func TestLiteral_TermEquals_NotQualifiers(t *testing.T) {
	v := Literal{
		Datatype:    "LangString",
		LexicalForm: "hello",
		Tag: LanguageLiteralTag{
			Language: "en",
		},
	}.TermEquals(Literal{
		Datatype:    "LangString",
		LexicalForm: "hello",
		Tag: LanguageLiteralTag{
			Language: "de",
		},
	})
	if _e, _a := false, v; _e != _a {
		t.Errorf("expected %v, got %v", _e, _a)
	}
}
