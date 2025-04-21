package ntriples

import (
	"context"
	"fmt"
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

type EncoderOption interface {
	apply(s *EncoderConfig)
	newEncoder(w io.Writer) (*Encoder, error)
}

type Encoder struct {
	w                 io.Writer
	blankNodeStringer blanknodeutil.Stringer
	ascii             bool
}

var _ encoding.GraphEncoder = &Encoder{}

func NewEncoder(w io.Writer, opts ...EncoderOption) (*Encoder, error) {
	compiledOpts := EncoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newEncoder(w)
}

func (w *Encoder) GetContentMetadata() encoding.ContentMetadata {
	return encoding.ContentMetadata{
		FileExt:   ".nt",
		MediaType: "application/n-triples",
		// spec says always utf-8 (even if ascii; application-default)
	}
}

func (w *Encoder) Close() error {
	return nil
}

func (w *Encoder) PutTriple(ctx context.Context, t rdf.Triple) error {
	var err error

	switch s := t.Subject.(type) {
	case rdf.BlankNode:
		_, err = w.w.Write([]byte("_:" + w.blankNodeStringer.GetBlankNodeIdentifier(s)))
		if err != nil {
			return err
		}
	case rdf.IRI:
		_, err = WriteIRI(w.w, s, w.ascii)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("subject: invalid type: %T", s)
	}

	_, err = w.w.Write([]byte(" "))
	if err != nil {
		return err
	}

	switch p := t.Predicate.(type) {
	case rdf.IRI:
		_, err = WriteIRI(w.w, p, w.ascii)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("predicate: invalid type: %T", p)
	}

	_, err = w.w.Write([]byte(" "))
	if err != nil {
		return err
	}

	switch o := t.Object.(type) {
	case rdf.BlankNode:
		_, err = w.w.Write([]byte("_:" + w.blankNodeStringer.GetBlankNodeIdentifier(o)))
		if err != nil {
			return err
		}
	case rdf.IRI:
		_, err = WriteIRI(w.w, o, w.ascii)
		if err != nil {
			return err
		}
	case rdf.Literal:
		_, err = WriteLiteral(w.w, o, w.ascii)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("object: invalid type: %T", o)
	}

	_, err = w.w.Write([]byte(" .\n"))
	if err != nil {
		return err
	}

	return nil
}
