package trig

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

func reader_scan_triples(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_triplesBlock.Err(err)
	}

	switch r0.Rune {
	case '<':
		r.buf.BacktrackRunes(r0)

		return readerStack{ectx, reader_scan_triples_subject_IRIREF}, nil
	case '_':
		// TODO require : since static?
		r.buf.BacktrackRunes(r0)

		return readerStack{ectx, reader_scan_triples_subject_BlankNode}, nil
	case '[':
		blankNode := ectx.Global.BlankNodeStringFactory.NewBlankNode()
		blankNodeRange := r.commitForTextOffsetRange(r0.AsDecodedRunes())

		ectx.CurSubject = blankNode
		ectx.CurSubjectLocation = blankNodeRange

		r.pushState(ectx, reader_scan_PredicateObjectList_Continue)
		r.pushState(ectx, reader_scan_PredicateObjectList)

		r.pushState(ectx, reader_scan_blankNodePropertyList_End)
		r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

		return readerStack{ectx, reader_scan_PredicateObjectList}, nil
	case '(':
		blankNode := ectx.Global.BlankNodeStringFactory.NewBlankNode()
		blankNodeRange := r.commitForTextOffsetRange(r0.AsDecodedRunes())

		fn := scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
			nectx := ectx

			if r0.Rune == ')' {
				nectx.CurSubject = rdfiri.Nil_List

				closeOffsets := r.commitForTextOffsetRange(r0.AsDecodedRunes())

				if closeOffsets != nil {
					nectx.CurSubjectLocation = &cursorio.TextOffsetRange{
						From:  blankNodeRange.From,
						Until: closeOffsets.Until,
					}
				} else {
					nectx.CurSubjectLocation = nil
				}

				r.pushState(nectx, reader_scan_PredicateObjectList_Continue)

				return readerStack{nectx, reader_scan_PredicateObjectList_Required}, nil
			}

			r.buf.BacktrackRunes(r0)

			nectx.CurSubject = blankNode
			nectx.CurSubjectLocation = blankNodeRange

			r.pushState(ectx, reader_scan_PredicateObjectList_Continue)
			r.pushState(nectx, reader_scan_PredicateObjectList_Required)
			fn := scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
				return reader_scan_collection(r, ectx, r0, nectx.CurSubject, nectx.CurSubjectLocation)
			})

			return readerStack{ectx, fn}, nil
		})

		return readerStack{ectx, fn}, nil
	}

	if r0.Rune == ':' || internal.IsRune_PN_CHARS_BASE(r0.Rune) {
		r.buf.BacktrackRunes(r0)

		return readerStack{ectx, reader_scan_triples_subject_PrefixedName}, nil
	}

	return readerStack{}, grammar.R_block.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
}

// TODO rename to generic '.' expect
func reader_scan_triples_End(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_triples.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
	} else if r0.Rune == '.' {
		r.commit(r0.AsDecodedRunes())

		return readerStack{}, nil
	}

	r.buf.BacktrackRunes(r0)

	return readerStack{}, grammar.R_triples.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
}

func reader_scan_triples_subject_IRIREF(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
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

func reader_scan_triples_subject_BlankNode(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
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

func reader_scan_triples_subject_PrefixedName(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
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
