package turtle

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

func reader_scan_Object(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_object.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))
	}

	switch {
	case r0.Rune == '<':
		token, err := r.produceIRIREF(r0)
		if err != nil {
			return readerStack{}, grammar.R_object.Err(err)
		}

		resolvedIRI, err := ectx.ResolveIRI(token.Decoded)
		if err != nil {
			return readerStack{}, grammar.R_object.Err(grammar.R_IRIREF.ErrWithTextOffsetRange(err, token.Offsets))
		}

		return r.emit(&statement{
			triple: rdf.Triple{
				Subject:   ectx.CurSubject,
				Predicate: ectx.CurPredicate,
				Object:    resolvedIRI,
			},
			offsets: r.buildTextOffsets(
				encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
				encoding.PredicateStatementOffsets, ectx.CurPredicateLocation,
				encoding.ObjectStatementOffsets, token.Offsets,
			),
		})
	case r0.Rune == '_':
		token, err := r.produceBlankNode(r0)
		if err != nil {
			return readerStack{}, grammar.R_object.Err(err)
		}

		blankNode := ectx.Global.BlankNodeStringMapper.MapBlankNodeIdentifier(token.Decoded)

		return r.emit(&statement{
			triple: rdf.Triple{
				Subject:   ectx.CurSubject,
				Predicate: ectx.CurPredicate,
				Object:    blankNode,
			},
			offsets: r.buildTextOffsets(
				encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
				encoding.PredicateStatementOffsets, ectx.CurPredicateLocation,
				encoding.ObjectStatementOffsets, token.Offsets,
			),
		})
	case r0.Rune == '(':
		cursor := r.commitForTextOffsetRange(r0.AsDecodedRunes())

		fn := scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
			if err != nil {
				return readerStack{}, grammar.R_object.Err(grammar.R_collection.Err(err))
			}

			return reader_scan_collection(r, ectx, r0, ectx.Global.BlankNodeFactory.NewBlankNode(), cursor)
		})

		return readerStack{ectx, fn}, nil
	case r0.Rune == '[':
		blankNode := ectx.Global.BlankNodeFactory.NewBlankNode()
		blankNodeRange := r.commitForTextOffsetRange(r0.AsDecodedRunes())

		nectx := ectx
		nectx.CurSubject = blankNode
		nectx.CurSubjectLocation = blankNodeRange
		nectx.CurPredicate = nil
		nectx.CurPredicateLocation = nil

		r.pushState(nectx, reader_scan_blankNodePropertyList_End)
		r.pushState(nectx, reader_scan_PredicateObjectList_Continue)
		r.pushState(nectx, reader_scan_PredicateObjectList)

		return r.emit(&statement{
			triple: rdf.Triple{
				Subject:   ectx.CurSubject,
				Predicate: ectx.CurPredicate,
				Object:    blankNode,
			},
			offsets: r.buildTextOffsets(
				encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
				encoding.PredicateStatementOffsets, ectx.CurPredicateLocation,
				encoding.ObjectStatementOffsets, blankNodeRange,
			),
		})
	case r0.Rune == '"', r0.Rune == '\'':
		token, err := r.produceString(r0)
		if err != nil {
			return readerStack{}, grammar.R_object.Err(grammar.R_literal.Err(grammar.R_RDFLiteral.Err(err)))
		}

		literal := rdf.Literal{
			Datatype:    xsdiri.String_Datatype,
			LexicalForm: token.Decoded,
		}

		{
			r0, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_object.Err(grammar.R_literal.Err(grammar.R_RDFLiteral.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{}))))
			}

			switch r0.Rune {
			case '@':
				langtagToken, err := r.produceLANGTAG(r0)
				if err != nil {
					return readerStack{}, grammar.R_object.Err(grammar.R_literal.Err(grammar.R_RDFLiteral.Err(err)))
				}

				literal.Datatype = rdfiri.LangString_Datatype
				literal.Tags = map[rdf.LiteralTag]string{
					rdf.LanguageLiteralTag: langtagToken.Decoded,
				}
			case '^':
				r1, err := r.buf.NextRune()
				if err != nil {
					return readerStack{}, grammar.R_object.Err(grammar.R_literal.Err(grammar.R_RDFLiteral.Err(r.newOffsetError(err, r0.AsDecodedRunes(), cursorio.DecodedRunes{}))))
				} else if r1.Rune != '^' {
					return readerStack{}, grammar.R_object.Err(grammar.R_literal.Err(grammar.R_RDFLiteral.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, r0.AsDecodedRunes(), r1.AsDecodedRunes()))))
				}

				r2, err := r.buf.NextRune()
				if err != nil {
					return readerStack{}, grammar.R_object.Err(grammar.R_literal.Err(grammar.R_RDFLiteral.Err(r.newOffsetError(err, cursorio.NewDecodedRunes(r0, r1), cursorio.DecodedRunes{}))))
				}

				r.commit(cursorio.NewDecodedRunes(r0, r1))

				if r2.Rune == '<' {
					datatypeToken, err := r.produceIRIREF(r2)
					if err != nil {
						return readerStack{}, grammar.R_object.Err(grammar.R_literal.Err(grammar.R_RDFLiteral.Err(err)))
					}

					resolvedIRI, err := ectx.ResolveIRI(datatypeToken.Decoded)
					if err != nil {
						return readerStack{}, grammar.R_object.Err(grammar.R_literal.Err(grammar.R_RDFLiteral.Err(grammar.R_IRIREF.ErrWithTextOffsetRange(err, datatypeToken.Offsets))))
					}

					literal.Datatype = resolvedIRI
				} else {
					datatypeToken, err := r.producePrefixedName(r2)
					if err != nil {
						return readerStack{}, grammar.R_object.Err(grammar.R_literal.Err(grammar.R_RDFLiteral.Err(err)))
					}

					expanded, ok := ectx.Global.Prefixes.ExpandPrefix(datatypeToken.NamespaceDecoded, datatypeToken.LocalDecoded)
					if !ok {
						return readerStack{}, grammar.R_object.Err(grammar.R_literal.Err(grammar.R_RDFLiteral.Err(grammar.R_PrefixedName.ErrWithTextOffsetRange(iriutil.NewUnknownPrefixError(datatypeToken.NamespaceDecoded), datatypeToken.Offsets))))
					}

					literal.Datatype = rdf.IRI(expanded)
				}
			default:
				r.buf.BacktrackRunes(r0)
			}
		}

		return r.emit(&statement{
			triple: rdf.Triple{
				Subject:   ectx.CurSubject,
				Predicate: ectx.CurPredicate,
				Object:    literal,
			},
			offsets: r.buildTextOffsets(
				encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
				encoding.PredicateStatementOffsets, ectx.CurPredicateLocation,
				encoding.ObjectStatementOffsets, token.Offsets,
			),
		})
	case r0.Rune == '+', r0.Rune == '-', '0' <= r0.Rune && r0.Rune <= '9', r0.Rune == '.':

		if r0.Rune == '.' {
			r1, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, err // TODO EOF check
			} else if r1.Rune < '0' || r1.Rune > '9' {
				r.buf.BacktrackRunes(r0, r1)

				return readerStack{}, nil
			}

			r.buf.BacktrackRunes(r1)
		}

		token, err := r.produceNumericLiteral(r0)
		if err != nil {
			return readerStack{}, grammar.R_object.Err(grammar.R_literal.Err(grammar.R_NumericLiteral.Err(err)))
		}

		literal := rdf.Literal{
			LexicalForm: token.Decoded,
		}

		switch token.GrammarRule {
		case grammar.R_INTEGER:
			literal.Datatype = xsdiri.Integer_Datatype
		case grammar.R_DECIMAL:
			literal.Datatype = xsdiri.Decimal_Datatype
		case grammar.R_DOUBLE:
			literal.Datatype = xsdiri.Double_Datatype
		default:
			panic("unreachable")
		}

		return r.emit(&statement{
			triple: rdf.Triple{
				Subject:   ectx.CurSubject,
				Predicate: ectx.CurPredicate,
				Object:    literal,
			},
			offsets: r.buildTextOffsets(
				encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
				encoding.PredicateStatementOffsets, ectx.CurPredicateLocation,
				encoding.ObjectStatementOffsets, token.Offsets,
			),
		})
	case r0.Rune == 't':
		r1, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_object.Err(r.newOffsetError(err, cursorio.NewDecodedRunes(r0), cursorio.DecodedRunes{}))
		} else if r1.Rune != 'r' {
			r.buf.BacktrackRunes(r0, r1)

			return readerStack{ectx, reader_scan_object_PrefixedName}, nil
		}

		r2, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_object.Err(r.newOffsetError(err, cursorio.NewDecodedRunes(r0, r1), cursorio.DecodedRunes{}))
		} else if r2.Rune != 'u' {
			r.buf.BacktrackRunes(r0, r1, r2)

			return readerStack{ectx, reader_scan_object_PrefixedName}, nil
		}

		r3, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_object.Err(r.newOffsetError(err, cursorio.NewDecodedRunes(r0, r1, r2), cursorio.DecodedRunes{}))
		} else if r3.Rune != 'e' {
			r.buf.BacktrackRunes(r0, r1, r2, r3)

			return readerStack{ectx, reader_scan_object_PrefixedName}, nil
		}

		// TODO verify next rune? avoid trueprefix:localname; need to figure out delimiters?

		return r.emit(&statement{
			triple: rdf.Triple{
				Subject:   ectx.CurSubject,
				Predicate: ectx.CurPredicate,
				Object: rdf.Literal{
					Datatype:    xsdiri.Boolean_Datatype,
					LexicalForm: "true",
				},
			},
			offsets: r.buildTextOffsets(
				encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
				encoding.PredicateStatementOffsets, ectx.CurPredicateLocation,
				encoding.ObjectStatementOffsets, r.commitForTextOffsetRange(cursorio.NewDecodedRunes(r0, r1, r2, r3)),
			),
		})
	case r0.Rune == 'f':
		r1, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_object.Err(r.newOffsetError(err, cursorio.NewDecodedRunes(r0), cursorio.DecodedRunes{}))
		} else if r1.Rune != 'a' {
			r.buf.BacktrackRunes(r0, r1)

			return readerStack{ectx, reader_scan_object_PrefixedName}, nil
		}

		r2, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_object.Err(r.newOffsetError(err, cursorio.NewDecodedRunes(r0, r1), cursorio.DecodedRunes{}))
		} else if r2.Rune != 'l' {
			r.buf.BacktrackRunes(r0, r1, r2)

			return readerStack{ectx, reader_scan_object_PrefixedName}, nil
		}

		r3, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_object.Err(r.newOffsetError(err, cursorio.NewDecodedRunes(r0, r1, r2), cursorio.DecodedRunes{}))
		} else if r3.Rune != 's' {
			r.buf.BacktrackRunes(r0, r1, r2, r3)

			return readerStack{ectx, reader_scan_object_PrefixedName}, nil
		}

		r4, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_object.Err(r.newOffsetError(err, cursorio.NewDecodedRunes(r0, r1, r2, r3), cursorio.DecodedRunes{}))
		} else if r4.Rune != 'e' {
			r.buf.BacktrackRunes(r0, r1, r2, r3, r4)

			return readerStack{ectx, reader_scan_object_PrefixedName}, nil
		}

		// TODO verify next rune? avoid trueprefix:localname; need to figure out delimiters?

		return r.emit(&statement{
			triple: rdf.Triple{
				Subject:   ectx.CurSubject,
				Predicate: ectx.CurPredicate,
				Object: rdf.Literal{
					Datatype:    xsdiri.Boolean_Datatype,
					LexicalForm: "false",
				},
			},
			offsets: r.buildTextOffsets(
				encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
				encoding.PredicateStatementOffsets, ectx.CurPredicateLocation,
				encoding.ObjectStatementOffsets, r.commitForTextOffsetRange(cursorio.NewDecodedRunes(r0, r1, r2, r3, r4)),
			),
		})
	case internal.IsRune_PN_CHARS_BASE(r0.Rune), r0.Rune == ':':
		r.buf.BacktrackRunes(r0)

		return readerStack{ectx, reader_scan_object_PrefixedName}, nil
	}

	return readerStack{}, grammar.R_object.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, cursorio.NewDecodedRunes(r0)))
}

func reader_scan_object_PrefixedName(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_object.Err(grammar.R_PrefixedName.Err(err))
	}

	token, err := r.producePrefixedName(r0)
	if err != nil {
		return readerStack{}, grammar.R_object.Err(err)
	}

	expanded, ok := ectx.Global.Prefixes.ExpandPrefix(token.NamespaceDecoded, token.LocalDecoded)
	if !ok {
		return readerStack{}, grammar.R_object.Err(grammar.R_PrefixedName.ErrWithTextOffsetRange(iriutil.NewUnknownPrefixError(token.NamespaceDecoded), token.Offsets))
	}

	return r.emit(&statement{
		triple: rdf.Triple{
			Subject:   ectx.CurSubject,
			Predicate: ectx.CurPredicate,
			Object:    expanded,
		},
		offsets: r.buildTextOffsets(
			encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
			encoding.PredicateStatementOffsets, ectx.CurPredicateLocation,
			encoding.ObjectStatementOffsets, token.Offsets,
		),
	})
}
