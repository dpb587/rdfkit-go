package turtle

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"sort"
	"strings"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/rdfdescription/rdfdescriptionio"
)

type EncoderOption interface {
	apply(s *EncoderConfig)
	newEncoder(w io.Writer) (*Encoder, error)
}

type Encoder struct {
	w                 io.Writer
	base              *iriutil.BaseIRI
	prefixes          *iriutil.PrefixTracker
	blankNodeStringer blanknodeutil.Stringer

	err              error
	buffered         bool
	bufferedSections [][]byte
}

var _ encoding.GraphEncoder = &Encoder{}
var _ rdfdescriptionio.GraphEncoder = &Encoder{}

func NewEncoder(w io.Writer, opts ...EncoderOption) (*Encoder, error) {
	compiledOpts := EncoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newEncoder(w)

}

func (e *Encoder) GetContentMetadata() encoding.ContentMetadata {
	return encoding.ContentMetadata{
		FileExt:   ".ttl",
		MediaType: "text/turtle",
		Charset:   "utf-8",
	}
}

func (w *Encoder) Close() error {
	if w.err != nil {
		if errors.Is(w.err, io.ErrClosedPipe) {
			return nil
		}

		return w.err
	}

	if w.buffered && len(w.bufferedSections) > 0 {
		sort.Slice(w.bufferedSections, func(i, j int) bool {
			return bytes.Compare(w.bufferedSections[i], w.bufferedSections[j]) < 0
		})

		if prefixMappings := w.prefixes.GetUsedPrefixMappings(); w.base != nil || len(prefixMappings) > 0 {
			slices.SortFunc(prefixMappings, iriutil.ComparePrefixMappingByPrefix)

			var baseString string

			if w.base != nil {
				baseString = w.base.String()
			}

			_, err := WriteDocumentHeader(w.w, baseString, prefixMappings)
			if err != nil {
				return err
			}

			_, err = w.w.Write([]byte("\n"))
			if err != nil {
				return fmt.Errorf("write header: %v", err)
			}
		}

		for _, section := range w.bufferedSections {
			_, err := w.w.Write(section)
			if err != nil {
				return fmt.Errorf("write: %v", err)
			}
		}
	}

	w.err = io.ErrClosedPipe

	return nil
}

func (w *Encoder) PutResource(ctx context.Context, r rdfdescription.Resource) error {
	if w.err != nil {
		return w.err
	}

	subject := r.GetResourceSubject()
	statements := r.GetResourceStatements()

	if len(statements) == 0 {
		return nil
	}

	if subject == nil {
		subject = rdf.NewBlankNode()
	}

	if len(statements) == 1 {
		if sO, ok := statements[0].(rdfdescription.ObjectStatement); ok {
			return w.PutTriple(ctx, rdf.Triple{
				Subject:   subject,
				Predicate: sO.Predicate,
				Object:    sO.Object,
			})
		}
	}

	buf := &bytes.Buffer{}

	err := w.writeSubjectValue(buf, subject)
	if err != nil {
		return fmt.Errorf("subject: %v", err)
	}

	err = w.putResourceStatements(ctx, buf, "\t", statements)
	if err != nil {
		return fmt.Errorf("resource: %v", err)
	}

	buf.WriteString(" .\n")

	if w.buffered {
		w.bufferedSections = append(w.bufferedSections, buf.Bytes())
	} else {
		_, err = buf.WriteTo(w.w)
		if err != nil {
			return fmt.Errorf("write: %v", err)
		}
	}

	return nil
}

func (w *Encoder) putResourceStatements(ctx context.Context, buf *bytes.Buffer, indent string, statements rdfdescription.StatementList) error {
	statementsByPredicate := statements.GroupByPredicate()

	var predicateList rdf.PredicateValueList

	{
		raw := statementsByPredicate.GetPredicateList()
		slices.SortFunc(raw, func(a, b rdf.PredicateValue) int {
			return strings.Compare(string(a.(rdf.IRI)), string(b.(rdf.IRI)))
		})

		for _, p := range raw {
			if p == rdfiri.Type_Property {
				predicateList = append([]rdf.PredicateValue{p}, predicateList...)
			} else {
				predicateList = append(predicateList, p)
			}
		}
	}

	for pIdx, p := range predicateList {
		if pIdx > 0 {
			buf.WriteString(" ;")
		}

		if len(predicateList) > 1 {
			buf.WriteString("\n" + indent)
		} else {
			buf.WriteString(" ")
		}

		pIRI, ok := p.(rdf.IRI)
		if !ok {
			return fmt.Errorf("predicate: invalid type: %T", p)
		}

		if pIRI == rdfiri.Type_Property {
			buf.WriteString("a")
		} else {
			w.writeIRI(buf, pIRI)
		}

		pStatements := statementsByPredicate[p]

		for statementIdx, statement := range pStatements {
			indent := indent

			if statementIdx > 0 {
				buf.WriteString(" ,")
			}

			if len(pStatements) > 1 {
				indent += "\t"
				buf.WriteString("\n" + indent)
			} else {
				buf.WriteString(" ")
			}

			switch statementT := statement.(type) {
			case rdfdescription.ObjectStatement:
				w.writeObjectValue(buf, statementT.Object)
			case rdfdescription.AnonResourceStatement:
				if len(statementT.AnonResource.Statements) == 0 {
					buf.WriteString("[]")
				} else {
					buf.WriteString("[")
					// buf.WriteString("[\n")
					// buf.WriteString(indent + "\t")

					err := w.putResourceStatements(ctx, buf, indent+"\t", statementT.AnonResource.Statements)
					if err != nil {
						return fmt.Errorf("resource: %v", err)
					}

					buf.WriteString("\n" + indent + "]")
				}
			default:
				return fmt.Errorf("object: invalid type: %T", statement)
			}
		}
	}

	return nil
}

func (w *Encoder) PutTriple(ctx context.Context, t rdf.Triple) error {
	if w.err != nil {
		return w.err
	}

	buf := &bytes.Buffer{}

	err := w.writeSubjectValue(buf, t.Subject)
	if err != nil {
		return fmt.Errorf("subject: %v", err)
	}

	buf.WriteString(" ")

	switch p := t.Predicate.(type) {
	case rdf.IRI:
		if p == rdfiri.Type_Property {
			buf.WriteString("a")
		} else {
			w.writeIRI(buf, p)
		}
	default:
		return fmt.Errorf("predicate: invalid type: %T", p)
	}

	buf.WriteString(" ")

	w.writeObjectValue(buf, t.Object)

	buf.WriteString(" .\n")

	if w.buffered {
		w.bufferedSections = append(w.bufferedSections, buf.Bytes())
	} else {
		_, err = buf.WriteTo(w.w)
		if err != nil {
			return fmt.Errorf("write: %v", err)
		}
	}

	return nil
}

func (w *Encoder) writeSubjectValue(buf *bytes.Buffer, v rdf.SubjectValue) error {
	switch s := v.(type) {
	case rdf.BlankNode:
		buf.WriteString("_:" + w.blankNodeStringer.GetBlankNodeIdentifier(s))

		return nil
	case rdf.IRI:
		w.writeIRI(buf, s)

		return nil
	}

	return fmt.Errorf("invalid type: %T", v)
}

func (w *Encoder) writeIRI(buffered *bytes.Buffer, v rdf.IRI) {
	prefix, suffix, ok := w.prefixes.CompactPrefix(v)
	if ok {
		buffered.WriteString(prefix + ":" + format_PN_LOCAL(suffix))

		return
	}

	if w.base != nil {
		if relativized, ok := w.base.RelativizeIRI(v); ok {
			fmt.Fprintf(buffered, "<%s>", formatIRI(relativized, false))

			return
		}
	}

	fmt.Fprintf(buffered, "<%s>", formatIRI(string(v), false))
}

func (e *Encoder) writeObjectValue(w *bytes.Buffer, v rdf.ObjectValue) error {
	var literal rdf.Literal

	switch o := v.(type) {
	case rdf.BlankNode:
		label := e.blankNodeStringer.GetBlankNodeIdentifier(o)

		w.WriteString("_:" + label)

		return nil
	case rdf.IRI:
		e.writeIRI(w, o)

		return nil
	case rdf.Literal:
		literal = o

		switch literal.Datatype {
		case xsdiri.Boolean_Datatype, xsdiri.Decimal_Datatype, xsdiri.Double_Datatype, xsdiri.Integer_Datatype, xsdiri.Long_Datatype:
			w.Write([]byte(literal.LexicalForm))

			return nil
		case rdfiri.LangString_Datatype:
			w.WriteString(formatLiteralLexicalForm(literal.LexicalForm, false))

			fmt.Fprintf(w, "@%s", literal.Tags[rdf.LanguageLiteralTag])

			return nil
		}

		w.WriteString(formatLiteralLexicalForm(literal.LexicalForm, false))

		if literal.Datatype != xsdiri.String_Datatype {
			w.WriteString("^^")
			e.writeIRI(w, literal.Datatype)
		}

		return nil
	}

	return fmt.Errorf("invalid type: %T", v)
}
