package trig

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

func reader_scan_triplesOrGraph_E1(value rdf.GraphNameValue, valueRange *cursorio.TextOffsetRange) scanFunc {
	return func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
		if r0.Rune == '{' {
			ectx.CurGraphName = value
			ectx.CurGraphNameLocation = valueRange

			r.commit(r0.AsDecodedRunes())

			r.pushState(ectx, reader_scan_wrappedGraph_End)

			return readerStack{ectx, reader_scan_triplesBlock}, nil
		}

		r.buf.BacktrackRunes(r0)

		ectx.CurSubject = value.(rdf.SubjectValue)
		ectx.CurSubjectLocation = valueRange

		r.pushState(ectx, reader_scan_triples_End)
		r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

		return readerStack{ectx, reader_scan_PredicateObjectList_Required}, nil
	}
}

func reader_scan_triplesOrGraph_labelOrSubject_IRIREF(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_triples.Err(grammar.R_subject.Err(err))
	}

	token, err := r.produceIRIREF(r0)
	if err != nil {
		return readerStack{}, grammar.R_triples.Err(grammar.R_subject.Err(err))
	}

	resolvedIRI, err := ectx.ResolveIRI(token.Decoded)
	if err != nil {
		return readerStack{}, grammar.R_triples.Err(grammar.R_subject.Err(grammar.R_IRIREF.ErrWithTextOffsetRange(err, token.Offsets)))
	}

	return readerStack{ectx, reader_scan_triplesOrGraph_E1(resolvedIRI, token.Offsets)}, nil
}

func reader_scan_triplesOrGraph_labelOrSubject_PrefixedName(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_triples.Err(grammar.R_subject.Err(grammar.R_PrefixedName.Err(err)))
	}

	token, err := r.producePrefixedName(r0)
	if err != nil {
		return readerStack{}, grammar.R_triples.Err(grammar.R_subject.Err(err))
	}

	expanded, ok := ectx.Global.Prefixes.ExpandPrefix(token.NamespaceDecoded, token.LocalDecoded)
	if !ok {
		return readerStack{}, grammar.R_triples.Err(grammar.R_subject.Err(grammar.R_PrefixedName.ErrWithTextOffsetRange(iriutil.NewUnknownPrefixError(token.NamespaceDecoded), token.Offsets)))
	}

	return readerStack{ectx, reader_scan_triplesOrGraph_E1(rdf.IRI(expanded), token.Offsets)}, nil
}

func reader_scan_triplesOrGraph_labelOrSubject_BlankNode(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_triples.Err(grammar.R_subject.Err(grammar.R_BlankNode.Err(err)))
	}

	token, err := r.produceBlankNode(r0)
	if err != nil {
		return readerStack{}, grammar.R_triples.Err(grammar.R_subject.Err(err))
	}

	return readerStack{ectx, reader_scan_triplesOrGraph_E1(ectx.Global.BlankNodeStringMapper.MapBlankNodeIdentifier(token.Decoded), token.Offsets)}, nil
}
