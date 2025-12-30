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
	"github.com/dpb587/rdfkit-go/rdfdescription/rdfdescriptionutil"
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

var _ encoding.TriplesEncoder = &Encoder{}
var _ rdfdescriptionutil.ResourceEncoder = &Encoder{}

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

func (w *Encoder) AddResource(ctx context.Context, r rdfdescription.Resource) error {
	if w.err != nil {
		return w.err
	}

	subject := r.GetResourceSubject()
	statements := r.GetResourceStatements()

	if len(statements) == 0 {
		return nil
	}

	buf := &bytes.Buffer{}

	if subject == nil {
		buf.WriteString("[]")
	} else {
		err := w.writeSubjectValue(buf, subject)
		if err != nil {
			return fmt.Errorf("subject: %v", err)
		}
	}

	_, err := w.putResourceStatements(ctx, buf, "", statements)
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

func (w *Encoder) putResourceStatements(ctx context.Context, buf *bytes.Buffer, linePrefix string, statements rdfdescription.StatementList) (bool, error) {
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

	var multiline bool

	if len(predicateList) > 1 {
		multiline = true
		linePrefix += "\t"
	}

	for pIdx, p := range predicateList {
		if pIdx > 0 {
			buf.WriteString(" ;")
		}

		if multiline {
			buf.WriteString("\n" + linePrefix)
		} else {
			buf.WriteString(" ")
		}

		pIRI, ok := p.(rdf.IRI)
		if !ok {
			return false, fmt.Errorf("predicate: invalid type: %T", p)
		}

		if pIRI == rdfiri.Type_Property {
			buf.WriteString("a")
		} else {
			w.writeIRI(buf, pIRI)
		}

		pStatements := statementsByPredicate[p]
		pMultiline := false
		pLinePrefix := linePrefix

		if len(pStatements) > 1 {
			pMultiline = true
			pLinePrefix += "\t"
		}

		for statementIdx, statement := range pStatements {
			if statementIdx > 0 {
				buf.WriteString(" ,")
			}

			if pMultiline {
				buf.WriteString("\n" + pLinePrefix)
			} else {
				buf.WriteString(" ")
			}

			mm, err := w.writeResourceStatement(ctx, buf, pLinePrefix, statement)
			if err != nil {
				return false, fmt.Errorf("statement: %v", err)
			} else if mm {
				multiline = true
			}
		}

		multiline = multiline || pMultiline
	}

	return multiline, nil
}

func (w *Encoder) writeResourceStatement(ctx context.Context, buf *bytes.Buffer, linePrefix string, statement rdfdescription.Statement) (bool, error) {
	var multiline bool

	switch statementT := statement.(type) {
	case rdfdescription.ObjectStatement:
		w.writeObjectValue(buf, statementT.Object)
	case rdfdescription.AnonResourceStatement:
		if len(statementT.AnonResource.Statements) == 0 {
			buf.WriteString("[]")
		} else if entries, ok := w.tryResourceCompactList(statementT.AnonResource); ok {
			mm, err := w.writeResourceCompactList(ctx, buf, linePrefix, entries)
			if err != nil {
				return false, fmt.Errorf("list: %v", err)
			} else if mm {
				multiline = true
			}
		} else {
			buf.WriteString("[")

			mm, err := w.putResourceStatements(ctx, buf, linePrefix, statementT.AnonResource.Statements)
			if err != nil {
				return false, fmt.Errorf("resource: %v", err)
			} else if mm {
				multiline = true

				buf.WriteString("\n" + linePrefix + "]")
			} else {
				buf.WriteString("]")
			}
		}
	default:
		return false, fmt.Errorf("object: invalid type: %T", statement)
	}

	return multiline, nil
}

func (w *Encoder) AddTriple(ctx context.Context, t rdf.Triple) error {
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

			if langTag, ok := literal.Tag.(rdf.LanguageLiteralTag); ok {
				fmt.Fprintf(w, "@%s", langTag.Language)
			}

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

func (e *Encoder) tryResourceCompactList(resource rdfdescription.AnonResource) (rdfdescription.StatementList, bool) {
	statements := resource.GetResourceStatements()
	if len(statements) == 0 {
		return nil, false
	}

	var entries rdfdescription.StatementList
	nextStatements := resource.GetResourceStatements()

	for {
		statementsByPredicate := nextStatements.GroupByPredicate()

		var hasFirst rdfdescription.Statement
		var hasRest rdfdescription.Statement

		for predicate, statements := range statementsByPredicate {
			switch predicate {
			case rdfiri.Type_Property:
				if len(statements) != 1 {
					return nil, false
				}

				s0, ok := statements[0].(rdfdescription.ObjectStatement)
				if !ok {
					return nil, false
				} else if s0.Object != rdfiri.List_Class {
					return nil, false
				}
			case rdfiri.First_Property:
				if len(statements) != 1 {
					return nil, false
				}

				hasFirst = statements[0]
			case rdfiri.Rest_Property:
				if len(statements) != 1 {
					return nil, false
				}

				hasRest = statements[0]
			default:
				return nil, false
			}
		}

		if hasFirst == nil || hasRest == nil {
			return nil, false
		}

		entries = append(entries, hasFirst)

		switch restStmt := hasRest.(type) {
		case rdfdescription.ObjectStatement:
			switch oT := restStmt.Object.(type) {
			case rdf.IRI:
				if oT == rdfiri.Nil_List {
					return entries, true
				}
			}

			return nil, false
		case rdfdescription.AnonResourceStatement:
			nextStatements = restStmt.AnonResource.GetResourceStatements()
		default:
			panic(fmt.Errorf("invalid type: %T", restStmt))
		}
	}
}

func (e *Encoder) writeResourceCompactList(ctx context.Context, buf *bytes.Buffer, linePrefix string, entries rdfdescription.StatementList) (bool, error) {
	if len(entries) == 0 {
		buf.WriteString("()")

		return false, nil
	}

	var multiline bool

	buf.WriteString("(")

	if len(entries) > 0 {
		multiline = true

		itemLinePrefix := linePrefix + "\t"

		for _, statement := range entries {
			buf.WriteString("\n" + itemLinePrefix)

			_, err := e.writeResourceStatement(ctx, buf, itemLinePrefix, statement)
			if err != nil {
				return false, fmt.Errorf("statement: %v", err)
			}
		}

		buf.WriteString("\n" + linePrefix)
	}

	buf.WriteString(")")

	return multiline, nil
}
