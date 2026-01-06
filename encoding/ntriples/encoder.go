package ntriples

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/ntriples/ntriplescontent"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type EncoderOption interface {
	apply(s *EncoderConfig)
	newEncoder(w io.Writer) (*Encoder, error)
}

type Encoder struct {
	w                io.Writer
	bnStringProvider blanknodes.StringProvider
	ascii            bool
	buf              *bytes.Buffer
}

var _ encoding.TriplesEncoder = &Encoder{}

func NewEncoder(w io.Writer, opts ...EncoderOption) (*Encoder, error) {
	compiledOpts := EncoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newEncoder(w)
}

func (w *Encoder) GetContentMetadata() encoding.ContentMetadata {
	return ntriplescontent.DefaultMetadata
}

func (w *Encoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return ntriplescontent.TypeIdentifier
}

func (w *Encoder) Close() error {
	return nil
}

func (w *Encoder) AddTriple(ctx context.Context, t rdf.Triple) error {
	var err error

	w.buf.Reset()

	switch s := t.Subject.(type) {
	case rdf.BlankNode:
		w.buf.Write([]byte("_:" + w.bnStringProvider.GetBlankNodeString(s)))
	case rdf.IRI:
		WriteIRI(w.buf, s, w.ascii)
	default:
		return fmt.Errorf("subject: invalid type: %T", s)
	}

	w.buf.Write([]byte(" "))

	switch p := t.Predicate.(type) {
	case rdf.IRI:
		WriteIRI(w.buf, p, w.ascii)
	default:
		return fmt.Errorf("predicate: invalid type: %T", p)
	}

	w.buf.Write([]byte(" "))

	switch o := t.Object.(type) {
	case rdf.BlankNode:
		w.buf.Write([]byte("_:" + w.bnStringProvider.GetBlankNodeString(o)))
	case rdf.IRI:
		WriteIRI(w.buf, o, w.ascii)
	case rdf.Literal:
		WriteLiteral(w.buf, o, w.ascii)
	default:
		return fmt.Errorf("object: invalid type: %T", o)
	}

	w.buf.Write([]byte(" .\n"))

	_, err = w.buf.WriteTo(w.w)
	if err != nil {
		return err
	}

	return nil
}
