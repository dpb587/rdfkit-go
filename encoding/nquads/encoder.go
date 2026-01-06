package nquads

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/nquads/nquadscontent"
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
	buf               *bytes.Buffer
}

var _ encoding.QuadsEncoder = &Encoder{}

func NewEncoder(w io.Writer, opts ...EncoderOption) (*Encoder, error) {
	compiledOpts := EncoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newEncoder(w)
}

func (w *Encoder) GetContentMetadata() encoding.ContentMetadata {
	return nquadscontent.DefaultMetadata
}

func (w *Encoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return nquadscontent.TypeIdentifier
}

func (w *Encoder) Close() error {
	return nil
}

func (w *Encoder) AddQuad(ctx context.Context, t rdf.Quad) error {
	var err error

	w.buf.Reset()

	switch s := t.Triple.Subject.(type) {
	case rdf.BlankNode:
		w.buf.Write([]byte{'_', ':'})
		w.buf.Write([]byte(w.blankNodeStringer.GetBlankNodeIdentifier(s)))
	case rdf.IRI:
		WriteIRI(w.buf, s, w.ascii)
	default:
		return fmt.Errorf("subject: invalid type: %T", s)
	}

	w.buf.Write([]byte{' '})

	switch p := t.Triple.Predicate.(type) {
	case rdf.IRI:
		WriteIRI(w.buf, p, w.ascii)
	default:
		return fmt.Errorf("predicate: invalid type: %T", p)
	}

	w.buf.Write([]byte{' '})

	switch o := t.Triple.Object.(type) {
	case rdf.BlankNode:
		w.buf.Write([]byte{'_', ':'})
		w.buf.Write([]byte(w.blankNodeStringer.GetBlankNodeIdentifier(o)))
	case rdf.IRI:
		WriteIRI(w.buf, o, w.ascii)
	case rdf.Literal:
		WriteLiteral(w.buf, o, w.ascii)
	default:
		return fmt.Errorf("object: invalid type: %T", o)
	}

	if t.GraphName != nil {
		w.buf.Write([]byte{' '})

		switch g := t.GraphName.(type) {
		case rdf.BlankNode:
			w.buf.Write([]byte{'_', ':'})
			w.buf.Write([]byte(w.blankNodeStringer.GetBlankNodeIdentifier(g)))
		case rdf.IRI:
			WriteIRI(w.buf, g, w.ascii)
		default:
			return fmt.Errorf("graph: invalid type: %T", g)
		}
	}

	w.buf.Write([]byte{' ', '.', '\n'})

	_, err = w.buf.WriteTo(w.w)
	if err != nil {
		return err
	}

	return nil
}
