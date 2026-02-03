package turtle

import (
	"bytes"
	"testing"

	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

func TestEncoder_Buffered_Base(t *testing.T) {
	ctx := t.Context()

	buf := &bytes.Buffer{}
	e, err := NewEncoder(buf, EncoderConfig{}.
		SetBuffered(true).
		SetBase("http://example.com/path/"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e.AddTriple(ctx, rdf.Triple{
		Subject:   rdf.IRI("http://example.com/path/subject"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})

	err = e.Close()

	if _a, _e := buf.String(), `@base <http://example.com/path/> .

<subject> a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Buffered_Prefixes(t *testing.T) {
	ctx := t.Context()

	buf := &bytes.Buffer{}
	e, err := NewEncoder(buf, EncoderConfig{}.
		SetBuffered(true).
		SetPrefixes(iri.PrefixMappingList{
			{
				Prefix:   "ex",
				Expanded: "http://example.com/path/",
			},
			{
				Prefix:   "unused",
				Expanded: "http://example.com/unused/",
			},
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = e.AddTriple(ctx, rdf.Triple{
		Subject:   rdf.IRI("http://example.com/path/subject"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})

	err = e.Close()

	if _a, _e := buf.String(), `@prefix ex: <http://example.com/path/> .

ex:subject a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Unbuffered_Header(t *testing.T) {
	ctx := t.Context()

	buf := &bytes.Buffer{}
	e, err := NewEncoder(buf, EncoderConfig{}.
		SetBase("http://example.com/anywhere/").
		SetPrefixes(iri.PrefixMappingList{
			{
				Prefix:   "ex",
				Expanded: "http://example.com/path/",
			},
			{
				Prefix:   "unused",
				Expanded: "http://example.com/unused/",
			},
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e.AddTriple(ctx, rdf.Triple{
		Subject:   rdf.IRI("http://example.com/path/subject"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})

	err = e.Close()

	if _a, _e := buf.String(), `@base <http://example.com/anywhere/> .
@prefix ex: <http://example.com/path/> .
@prefix unused: <http://example.com/unused/> .

ex:subject a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Prefix_LocalNameEscapes(t *testing.T) {
	ctx := t.Context()

	buf := &bytes.Buffer{}
	e, err := NewEncoder(buf, EncoderConfig{}.
		SetPrefixes(iri.PrefixMappingList{
			{
				Prefix:   "ex",
				Expanded: "http://example.com/path/",
			},
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e.AddTriple(ctx, rdf.Triple{
		Subject:   rdf.IRI("http://example.com/path/nested/subject"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})

	e.Close()

	if _a, _e := buf.String(), `@prefix ex: <http://example.com/path/> .

ex:nested\/subject a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Prefix_UnderscoreAndDashNotEscaped(t *testing.T) {
	ctx := t.Context()

	buf := &bytes.Buffer{}
	e, err := NewEncoder(buf, EncoderConfig{}.
		SetPrefixes(iri.PrefixMappingList{
			{
				Prefix:   "ex",
				Expanded: "http://example.com/",
			},
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e.AddTriple(ctx, rdf.Triple{
		Subject:   rdf.IRI("http://example.com/hello_world"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})

	e.AddTriple(ctx, rdf.Triple{
		Subject:   rdf.IRI("http://example.com/hello-world"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})

	e.Close()

	if _a, _e := buf.String(), `@prefix ex: <http://example.com/> .

ex:hello_world a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
ex:hello-world a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Prefix_SpecialCharsNeedEscape(t *testing.T) {
	ctx := t.Context()

	buf := &bytes.Buffer{}
	e, err := NewEncoder(buf, EncoderConfig{}.
		SetPrefixes(iri.PrefixMappingList{
			{
				Prefix:   "ex",
				Expanded: "http://example.com/",
			},
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e.AddTriple(ctx, rdf.Triple{
		Subject:   rdf.IRI("http://example.com/hello/world"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})

	e.Close()

	if _a, _e := buf.String(), `@prefix ex: <http://example.com/> .

ex:hello\/world a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Prefix_DotPositionDependent(t *testing.T) {
	ctx := t.Context()

	buf := &bytes.Buffer{}
	e, err := NewEncoder(buf, EncoderConfig{}.
		SetPrefixes(iri.PrefixMappingList{
			{
				Prefix:   "ex",
				Expanded: "http://example.com/",
			},
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e.AddTriple(ctx, rdf.Triple{
		Subject:   rdf.IRI("http://example.com/.hidden"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})

	e.AddTriple(ctx, rdf.Triple{
		Subject:   rdf.IRI("http://example.com/hello.world"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})

	e.AddTriple(ctx, rdf.Triple{
		Subject:   rdf.IRI("http://example.com/test."),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})

	e.Close()

	if _a, _e := buf.String(), `@prefix ex: <http://example.com/> .

ex:\.hidden a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
ex:hello.world a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
ex:test\. a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Prefix_BackslashPercentEncoded(t *testing.T) {
	ctx := t.Context()

	buf := &bytes.Buffer{}
	e, err := NewEncoder(buf, EncoderConfig{}.
		SetPrefixes(iri.PrefixMappingList{
			{
				Prefix:   "ex",
				Expanded: "http://example.com/",
			},
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e.AddTriple(ctx, rdf.Triple{
		Subject:   rdf.IRI("http://example.com/back\\slash"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})

	e.Close()

	// Backslash is the escape character itself, so it must be percent-encoded
	if _a, _e := buf.String(), `@prefix ex: <http://example.com/> .

ex:back%5Cslash a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Resources_AnonSinglePropertyWS(t *testing.T) {
	ctx := t.Context()

	buf := &bytes.Buffer{}
	e, err := NewEncoder(buf, EncoderConfig{}.
		SetPrefixes(iri.PrefixMappingList{
			{
				Prefix:   "schema",
				Expanded: "http://schema.org/",
			},
			{
				Prefix:   "sh",
				Expanded: "http://www.w3.org/ns/shacl#",
			},
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resource := rdfdescription.SubjectResource{
		Subject: rdf.IRI("http://schema.org/weight_Property"),
		Statements: rdfdescription.StatementList{
			rdfdescription.AnonResourceStatement{
				Predicate: rdf.IRI("http://www.w3.org/ns/shacl#or"),
				AnonResource: rdfdescription.AnonResource{
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/ns/shacl#class"),
							Object:    rdf.IRI("http://schema.org/Mass"),
						},
					},
				},
			},
		},
	}

	err = e.AddResource(ctx, resource)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e.Close()

	// The closing bracket should have a space before it
	if _a, _e := buf.String(), `@prefix schema: <http://schema.org/> .
@prefix sh: <http://www.w3.org/ns/shacl#> .

schema:weight_Property sh:or [ sh:class schema:Mass ] .
`; _a != _e {
		t.Fatalf("expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Resources_AnonMultiplePropertiesWS(t *testing.T) {
	ctx := t.Context()

	buf := &bytes.Buffer{}
	e, err := NewEncoder(buf, EncoderConfig{}.
		SetPrefixes(iri.PrefixMappingList{
			{
				Prefix:   "schema",
				Expanded: "http://schema.org/",
			},
			{
				Prefix:   "sh",
				Expanded: "http://www.w3.org/ns/shacl#",
			},
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resource := rdfdescription.SubjectResource{
		Subject: rdf.IRI("http://schema.org/weight_Property"),
		Statements: rdfdescription.StatementList{
			rdfdescription.AnonResourceStatement{
				Predicate: rdf.IRI("http://www.w3.org/ns/shacl#property"),
				AnonResource: rdfdescription.AnonResource{
					Statements: rdfdescription.StatementList{
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/ns/shacl#class"),
							Object:    rdf.IRI("http://schema.org/Mass"),
						},
						rdfdescription.ObjectStatement{
							Predicate: rdf.IRI("http://www.w3.org/ns/shacl#name"),
							Object:    rdf.Literal{LexicalForm: "Weight", Datatype: rdfiri.LangString_Datatype, Tag: rdf.LanguageLiteralTag{Language: "en"}},
						},
					},
				},
			},
		},
	}

	err = e.AddResource(ctx, resource)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e.Close()

	// When multiline, the closing bracket should be on its own line with no space before it
	if _a, _e := buf.String(), `@prefix schema: <http://schema.org/> .
@prefix sh: <http://www.w3.org/ns/shacl#> .

schema:weight_Property sh:property [
	sh:class schema:Mass ;
	sh:name "Weight"@en
] .
`; _a != _e {
		t.Fatalf("expected %q, got %q", _e, _a)
	}
}
