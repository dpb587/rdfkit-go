package turtle

import (
	"bytes"
	"testing"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
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
		SetPrefixes(iriutil.NewPrefixMap(iriutil.PrefixMapping{
			Prefix:   "ex",
			Expanded: "http://example.com/path/",
		}, iriutil.PrefixMapping{
			Prefix:   "unused",
			Expanded: "http://example.com/unused/",
		})),
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
		SetPrefixes(iriutil.NewPrefixMap(iriutil.PrefixMapping{
			Prefix:   "ex",
			Expanded: "http://example.com/path/",
		}, iriutil.PrefixMapping{
			Prefix:   "unused",
			Expanded: "http://example.com/unused/",
		})),
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
		SetPrefixes(iriutil.NewPrefixMap(iriutil.PrefixMapping{
			Prefix:   "ex",
			Expanded: "http://example.com/path/",
		})),
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
		SetPrefixes(iriutil.NewPrefixMap(iriutil.PrefixMapping{
			Prefix:   "ex",
			Expanded: "http://example.com/",
		})),
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
		SetPrefixes(iriutil.NewPrefixMap(iriutil.PrefixMapping{
			Prefix:   "ex",
			Expanded: "http://example.com/",
		})),
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
		SetPrefixes(iriutil.NewPrefixMap(iriutil.PrefixMapping{
			Prefix:   "ex",
			Expanded: "http://example.com/",
		})),
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
		SetPrefixes(iriutil.NewPrefixMap(iriutil.PrefixMapping{
			Prefix:   "ex",
			Expanded: "http://example.com/",
		})),
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
