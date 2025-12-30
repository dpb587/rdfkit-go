package turtle

import (
	"bytes"
	"context"
	"testing"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

func TestEncoder_Buffered_Base(t *testing.T) {
	buf := &bytes.Buffer{}
	e, err := NewEncoder(buf, EncoderConfig{}.
		SetBuffered(true).
		SetBase("http://example.com/path/"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = e.AddTriple(context.Background(), rdf.Triple{
		Subject:   rdf.IRI("http://example.com/path/subject"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _a, _e := len(buf.Bytes()), 0; _a != _e {
		t.Fatalf("buffer length: expected %v, got %v", _e, _a)
	}

	err = e.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _a, _e := buf.String(), `@base <http://example.com/path/> .

<subject> a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("buffer content: expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Buffered_Prefixes(t *testing.T) {
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

	err = e.AddTriple(context.Background(), rdf.Triple{
		Subject:   rdf.IRI("http://example.com/path/subject"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _a, _e := len(buf.Bytes()), 0; _a != _e {
		t.Fatalf("buffer length: expected %v, got %v", _e, _a)
	}

	err = e.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _a, _e := buf.String(), `@prefix ex: <http://example.com/path/> .

ex:subject a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("buffer content: expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Unbuffered_Header(t *testing.T) {
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

	err = e.AddTriple(context.Background(), rdf.Triple{
		Subject:   rdf.IRI("http://example.com/path/subject"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _a, _e := len(buf.Bytes()), 0; _a == _e {
		t.Fatalf("buffer length: expected >%v, got %v", _e, _a)
	}

	err = e.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _a, _e := buf.String(), `@base <http://example.com/anywhere/> .
@prefix ex: <http://example.com/path/> .
@prefix unused: <http://example.com/unused/> .

ex:subject a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("buffer content: expected %q, got %q", _e, _a)
	}
}

func TestEncoder_Prefix_LocalNameEscapes(t *testing.T) {
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

	err = e.AddTriple(context.Background(), rdf.Triple{
		Subject:   rdf.IRI("http://example.com/path/nested/subject"),
		Predicate: rdfiri.Type_Property,
		Object:    rdfiri.Property_Class,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _a, _e := len(buf.Bytes()), 0; _a == _e {
		t.Fatalf("buffer length: expected >%v, got %v", _e, _a)
	}

	err = e.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _a, _e := buf.String(), `@prefix ex: <http://example.com/path/> .

ex:nested\/subject a <http://www.w3.org/1999/02/22-rdf-syntax-ns#Property> .
`; _a != _e {
		t.Fatalf("buffer content: expected %q, got %q", _e, _a)
	}
}
