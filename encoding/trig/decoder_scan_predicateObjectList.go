package trig

import (
	"unicode"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal"
	"github.com/dpb587/rdfkit-go/encoding/trig/internal/grammar"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

func reader_scan_PredicateObjectList(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_predicateObjectList.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
	}

	switch {
	// case r0 == ';':
	// 	r.commit([]rune{r0})

	// 	return readerStack{}, nil
	case r0.Rune == '<':
		token, err := r.produceIRIREF(r0)
		if err != nil {
			return readerStack{}, grammar.R_predicateObjectList.Err(grammar.R_verb.Err(err))
		}

		resolvedIRI, err := ectx.ResolveIRI(token.Decoded)
		if err != nil {
			return readerStack{}, grammar.R_predicateObjectList.Err(grammar.R_verb.Err(grammar.R_IRIREF.ErrWithTextOffsetRange(err, token.Offsets)))
		}

		ectx.CurPredicate = resolvedIRI
		ectx.CurPredicateLocation = token.Offsets

		r.pushState(ectx, reader_scan_ObjectList_Continue)

		return readerStack{ectx, reader_scan_Object}, nil
	case r0.Rune == 'a':
		r1, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_predicateObjectList.Err(r.newOffsetError(err, r0.AsDecodedRunes(), cursorio.DecodedRunes{}))
		}

		if !unicode.IsSpace(r1.Rune) {
			r.buf.BacktrackRunes(r1)

			token, err := r.producePrefixedName(r0)
			if err != nil {
				return readerStack{}, grammar.R_predicateObjectList.Err(grammar.R_verb.Err(err))
			}

			expanded, ok := ectx.Global.Prefixes.ExpandPrefix(token.NamespaceDecoded, token.LocalDecoded)
			if !ok {
				return readerStack{}, grammar.R_predicateObjectList.Err(grammar.R_verb.Err(grammar.R_PrefixedName.ErrWithTextOffsetRange(iriutil.NewUnknownPrefixError(token.NamespaceDecoded), token.Offsets)))
			}

			ectx.CurPredicate = rdf.IRI(expanded)
			ectx.CurPredicateLocation = token.Offsets
		} else {
			ectx.CurPredicate = rdf.IRI(rdfiri.Type_Property)
			ectx.CurPredicateLocation = r.commitForTextOffsetRange(r0.AsDecodedRunes())

			r.commit(r1.AsDecodedRunes())
		}

		r.pushState(ectx, reader_scan_ObjectList_Continue)

		return readerStack{ectx, reader_scan_Object}, nil
	case r0.Rune != 'a' && (r0.Rune == ':' || internal.IsRune_PN_CHARS_BASE(r0.Rune)):
		token, err := r.producePrefixedName(r0)
		if err != nil {
			return readerStack{}, grammar.R_predicateObjectList.Err(grammar.R_verb.Err(err))
		}

		expanded, ok := ectx.Global.Prefixes.ExpandPrefix(token.NamespaceDecoded, token.LocalDecoded)
		if !ok {
			return readerStack{}, grammar.R_predicateObjectList.Err(grammar.R_verb.Err(grammar.R_PrefixedName.ErrWithTextOffsetRange(iriutil.NewUnknownPrefixError(token.NamespaceDecoded), token.Offsets)))
		}

		ectx.CurPredicate = rdf.IRI(expanded)
		ectx.CurPredicateLocation = token.Offsets

		r.pushState(ectx, reader_scan_ObjectList_Continue)

		return readerStack{ectx, reader_scan_Object}, nil
	}

	r.buf.BacktrackRunes(r0)

	return readerStack{}, nil
}

func reader_scan_PredicateObjectList_Continue(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_predicateObjectList.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
	} else if r0.Rune == ';' {
		r.commit(r0.AsDecodedRunes())

		r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

		return readerStack{ectx, reader_scan_PredicateObjectList}, nil
	}

	r.buf.BacktrackRunes(r0)

	return readerStack{}, nil
}
