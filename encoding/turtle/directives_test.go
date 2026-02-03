package turtle

import (
	"bytes"
	"testing"

	"github.com/dpb587/rdfkit-go/iri"
)

func TestDirectiveMode_AtIsZero(t *testing.T) {
	var mode DirectiveMode
	if mode != DirectiveMode_At {
		t.Fatalf("DirectiveMode_At should be the zero value, got %v", mode)
	}
}

func TestWriteDirectives_At(t *testing.T) {
	buf := &bytes.Buffer{}

	_, err := WriteDirectives(buf, WriteDirectivesOptions{
		Base:     "http://example.com/path/",
		BaseMode: DirectiveMode_At,
		Prefixes: iri.PrefixMappingList{
			{Prefix: "ex", Expanded: "http://example.com/ns/"},
		},
		PrefixMode: DirectiveMode_At,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `@base <http://example.com/path/> .
@prefix ex: <http://example.com/ns/> .
`
	if buf.String() != expected {
		t.Fatalf("buffer content: expected %q, got %q", expected, buf.String())
	}
}

func TestWriteDirectives_SPARQL(t *testing.T) {
	buf := &bytes.Buffer{}

	_, err := WriteDirectives(buf, WriteDirectivesOptions{
		Base:     "http://example.com/path/",
		BaseMode: DirectiveMode_SPARQL,
		Prefixes: iri.PrefixMappingList{
			{Prefix: "ex", Expanded: "http://example.com/ns/"},
		},
		PrefixMode: DirectiveMode_SPARQL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `BASE <http://example.com/path/>
PREFIX ex: <http://example.com/ns/>
`
	if buf.String() != expected {
		t.Fatalf("buffer content: expected %q, got %q", expected, buf.String())
	}
}

func TestWriteDirectives_Disabled(t *testing.T) {
	buf := &bytes.Buffer{}

	_, err := WriteDirectives(buf, WriteDirectivesOptions{
		Base:     "http://example.com/path/",
		BaseMode: DirectiveMode_Disabled,
		Prefixes: iri.PrefixMappingList{
			{Prefix: "ex", Expanded: "http://example.com/ns/"},
		},
		PrefixMode: DirectiveMode_Disabled,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() != 0 {
		t.Fatalf("buffer should be empty, got %q", buf.String())
	}
}

func TestWriteDirectives_Mixed_BaseDisabled_PrefixAt(t *testing.T) {
	buf := &bytes.Buffer{}

	_, err := WriteDirectives(buf, WriteDirectivesOptions{
		Base:     "http://example.com/path/",
		BaseMode: DirectiveMode_Disabled,
		Prefixes: iri.PrefixMappingList{
			{Prefix: "ex", Expanded: "http://example.com/ns/"},
		},
		PrefixMode: DirectiveMode_At,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `@prefix ex: <http://example.com/ns/> .
`
	if buf.String() != expected {
		t.Fatalf("buffer content: expected %q, got %q", expected, buf.String())
	}
}

func TestWriteDirectives_Mixed_BaseAt_PrefixSPARQL(t *testing.T) {
	buf := &bytes.Buffer{}

	_, err := WriteDirectives(buf, WriteDirectivesOptions{
		Base:     "http://example.com/path/",
		BaseMode: DirectiveMode_At,
		Prefixes: iri.PrefixMappingList{
			{Prefix: "ex", Expanded: "http://example.com/ns/"},
		},
		PrefixMode: DirectiveMode_SPARQL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `@base <http://example.com/path/> .
PREFIX ex: <http://example.com/ns/>
`
	if buf.String() != expected {
		t.Fatalf("buffer content: expected %q, got %q", expected, buf.String())
	}
}

func TestWriteDirectives_Mixed_BaseSPARQL_PrefixDisabled(t *testing.T) {
	buf := &bytes.Buffer{}

	_, err := WriteDirectives(buf, WriteDirectivesOptions{
		Base:     "http://example.com/path/",
		BaseMode: DirectiveMode_SPARQL,
		Prefixes: iri.PrefixMappingList{
			{Prefix: "ex", Expanded: "http://example.com/ns/"},
		},
		PrefixMode: DirectiveMode_Disabled,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `BASE <http://example.com/path/>
`
	if buf.String() != expected {
		t.Fatalf("buffer content: expected %q, got %q", expected, buf.String())
	}
}

func TestWriteDirectives_EmptyBase(t *testing.T) {
	buf := &bytes.Buffer{}

	_, err := WriteDirectives(buf, WriteDirectivesOptions{
		Base:     "",
		BaseMode: DirectiveMode_At,
		Prefixes: iri.PrefixMappingList{
			{Prefix: "ex", Expanded: "http://example.com/ns/"},
		},
		PrefixMode: DirectiveMode_At,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `@prefix ex: <http://example.com/ns/> .
`
	if buf.String() != expected {
		t.Fatalf("buffer content: expected %q, got %q", expected, buf.String())
	}
}

func TestWriteDirectives_EmptyPrefixes(t *testing.T) {
	buf := &bytes.Buffer{}

	_, err := WriteDirectives(buf, WriteDirectivesOptions{
		Base:       "http://example.com/path/",
		BaseMode:   DirectiveMode_At,
		Prefixes:   iri.PrefixMappingList{},
		PrefixMode: DirectiveMode_At,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `@base <http://example.com/path/> .
`
	if buf.String() != expected {
		t.Fatalf("buffer content: expected %q, got %q", expected, buf.String())
	}
}

func TestWriteDirectives_MultiplePrefixes(t *testing.T) {
	buf := &bytes.Buffer{}

	_, err := WriteDirectives(buf, WriteDirectivesOptions{
		Base:     "http://example.com/path/",
		BaseMode: DirectiveMode_At,
		Prefixes: iri.PrefixMappingList{
			{Prefix: "ex", Expanded: "http://example.com/ns/"},
			{Prefix: "rdf", Expanded: "http://www.w3.org/1999/02/22-rdf-syntax-ns#"},
		},
		PrefixMode: DirectiveMode_At,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `@base <http://example.com/path/> .
@prefix ex: <http://example.com/ns/> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
`
	if buf.String() != expected {
		t.Fatalf("buffer content: expected %q, got %q", expected, buf.String())
	}
}
