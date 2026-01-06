package rdfjson

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/rdfjson/rdfjsoncontent"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
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
	prefix            string
	indent            string

	buf map[string]map[string][]any
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
	return rdfjsoncontent.DefaultMetadata
}

func (w *Encoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return rdfjsoncontent.TypeIdentifier
}

func (w *Encoder) Close() error {
	e := json.NewEncoder(w.w)
	e.SetIndent(w.prefix, w.indent)

	return e.Encode(w.buf)
}

func (w *Encoder) AddTriple(ctx context.Context, t rdf.Triple) error {
	var sData string

	switch s := t.Subject.(type) {
	case rdf.BlankNode:
		sData = "_:" + w.blankNodeStringer.GetBlankNodeIdentifier(s)
	case rdf.IRI:
		sData = string(s)
	default:
		return fmt.Errorf("subject: invalid type: %T", s)
	}

	var pData string

	switch p := t.Predicate.(type) {
	case rdf.IRI:
		pData = string(p)
	default:
		return fmt.Errorf("predicate: invalid type: %T", p)
	}

	var oData any

	switch o := t.Object.(type) {
	case rdf.BlankNode:
		oData = map[string]string{
			"type":  "bnode",
			"value": "_:" + w.blankNodeStringer.GetBlankNodeIdentifier(o),
		}
	case rdf.IRI:
		oData = map[string]string{
			"type":  "uri",
			"value": string(o),
		}
	case rdf.Literal:
		if o.Datatype == rdfiri.LangString_Datatype {
			if langTag, ok := o.Tag.(rdf.LanguageLiteralTag); ok {
				oData = map[string]string{
					"type":  "literal",
					"value": o.LexicalForm,
					"lang":  langTag.Language,
				}
			} else {
				oData = map[string]string{
					"type":  "literal",
					"value": o.LexicalForm,
				}
			}
		} else if o.Datatype != xsdiri.String_Datatype {
			oData = map[string]string{
				"type":     "literal",
				"value":    o.LexicalForm,
				"datatype": string(o.Datatype),
			}
		} else {
			oData = map[string]string{
				"type":  "literal",
				"value": o.LexicalForm,
			}
		}
	default:
		return fmt.Errorf("object: invalid type: %T", o)
	}

	if _, ok := w.buf[sData]; !ok {
		w.buf[sData] = map[string][]any{}
	}

	w.buf[sData][pData] = append(w.buf[sData][pData], oData)

	return nil
}
