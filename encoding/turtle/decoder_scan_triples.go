package turtle

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

func reader_scan_Triples_End(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_triples.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
	} else if r0.Rune == '.' {
		r.commit(r0.AsDecodedRunes())

		return readerStack{}, nil
	}

	return readerStack{}, grammar.R_triples.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
}

func reader_scan_Triples_Subject_IRIREF(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
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

	ectx.CurSubject = resolvedIRI
	ectx.CurSubjectLocation = token.Offsets

	r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

	return readerStack{ectx, reader_scan_PredicateObjectList_Required}, nil
}

func reader_scan_Triples_Subject_PrefixedName(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
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

	ectx.CurSubject = rdf.IRI(expanded)
	ectx.CurSubjectLocation = token.Offsets

	r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

	return readerStack{ectx, reader_scan_PredicateObjectList_Required}, nil
}

func reader_scan_Triples_Subject_BlankNode(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_triples.Err(grammar.R_subject.Err(grammar.R_BlankNode.Err(err)))
	}

	token, err := r.produceBlankNode(r0)
	if err != nil {
		return readerStack{}, grammar.R_triples.Err(grammar.R_subject.Err(err))
	}

	ectx.CurSubject = ectx.Global.BlankNodeStringFactory.NewStringBlankNode(token.Decoded)
	ectx.CurSubjectLocation = token.Offsets

	r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

	return readerStack{ectx, reader_scan_PredicateObjectList_Required}, nil
}
