package turtle

import (
	"errors"
	"io"
	"unicode"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
	"github.com/dpb587/rdfkit-go/iri"
)

func reader_scanStatement_Subject_AnonOrBlankNode(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_subject.Err(err)
	} else if r0.Rune == ']' {
		r.commit(r0.AsDecodedRunes())

		r.pushState(ectx, reader_scan_Triples_End)
		r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

		return readerStack{ectx, reader_scan_PredicateObjectList_Required}, nil
	}

	r.buf.BacktrackRunes(r0)

	r.pushState(ectx, reader_scan_Triples_End)
	r.pushState(ectx, reader_scan_PredicateObjectList)
	r.pushState(ectx, reader_scan_blankNodePropertyList_End)
	r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

	return readerStack{ectx, reader_scan_PredicateObjectList_Required}, nil
}

func reader_scanStatement(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
	if err != nil {
		if errors.Is(err, io.EOF) {
			return r.terminate()
		}

		return readerStack{}, grammar.R_statement.Err(err)
	}

	r.pushState(ectx, reader_scanStatement)

	switch r0.Rune {
	case '@':
		r1, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, err
		}

		switch r1.Rune {

		// @base directive
		case 'b':
			r2, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, r0.AsDecodedRunes(), cursorio.DecodedRunes{}))
			} else if r2.Rune != 'a' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2.Rune}, cursorio.DecodedRuneList{r0, r1}.AsDecodedRunes(), r2.AsDecodedRunes()))
			}

			r3, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1}.AsDecodedRunes(), cursorio.DecodedRunes{}))
			} else if r3.Rune != 's' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r3.Rune}, cursorio.DecodedRuneList{r0, r1, r2}.AsDecodedRunes(), r3.AsDecodedRunes()))
			}

			r4, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1, r2}.AsDecodedRunes(), cursorio.DecodedRunes{}))
			} else if r4.Rune != 'e' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r4.Rune}, cursorio.DecodedRuneList{r0, r1, r2, r3}.AsDecodedRunes(), r4.AsDecodedRunes()))
			}

			r.commit(cursorio.DecodedRuneList{r0, r1, r2, r3, r4}.AsDecodedRunes())

			return readerStack{
				ectx,
				scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
					if err != nil {
						return readerStack{}, grammar.R_directive.Err(grammar.R_base.Err(err))
					}

					baseToken, err := r.produceIRIREF(r0)
					if err != nil {
						return readerStack{}, grammar.R_directive.Err(grammar.R_base.Err(err))
					}

					resolvedBase, err := ectx.ResolveURL(baseToken.Decoded)
					if err != nil {
						return readerStack{}, grammar.R_directive.Err(grammar.R_base.Err(grammar.R_IRIREF.ErrWithTextOffsetRange(err, baseToken.Offsets)))
					}

					return readerStack{
						ectx,
						scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
							if err != nil {
								return readerStack{}, grammar.R_directive.Err(grammar.R_base.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{})))
							} else if r0.Rune != '.' {
								return readerStack{}, grammar.R_directive.Err(grammar.R_base.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes())))
							}

							r.commit(r0.AsDecodedRunes())

							ectx.Global.Base = resolvedBase

							if r.baseDirectiveListener != nil {
								r.baseDirectiveListener(DecoderEvent_BaseDirective_Data{
									Value:        resolvedBase.String(),
									ValueOffsets: baseToken.Offsets,
								})
							}

							return readerStack{ectx, reader_scanStatement}, nil
						}),
					}, nil
				}),
			}, nil

		// @prefix directive
		case 'p':
			r2, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, r0.AsDecodedRunes(), cursorio.DecodedRunes{}))
			} else if r2.Rune != 'r' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2.Rune}, cursorio.DecodedRuneList{r0, r1}.AsDecodedRunes(), r2.AsDecodedRunes()))
			}

			r3, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1}.AsDecodedRunes(), cursorio.DecodedRunes{}))
			} else if r3.Rune != 'e' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r3.Rune}, cursorio.DecodedRuneList{r0, r1, r2}.AsDecodedRunes(), r3.AsDecodedRunes()))
			}

			r4, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1, r2}.AsDecodedRunes(), cursorio.DecodedRunes{}))
			} else if r4.Rune != 'f' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r4.Rune}, cursorio.DecodedRuneList{r0, r1, r2, r3}.AsDecodedRunes(), r4.AsDecodedRunes()))
			}

			r5, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1, r2, r3}.AsDecodedRunes(), cursorio.DecodedRunes{}))
			} else if r5.Rune != 'i' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r5.Rune}, cursorio.DecodedRuneList{r0, r1, r2, r3, r4}.AsDecodedRunes(), r5.AsDecodedRunes()))
			}

			r6, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1, r2, r3, r4}.AsDecodedRunes(), cursorio.DecodedRunes{}))
			} else if r6.Rune != 'x' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r6.Rune}, cursorio.DecodedRuneList{r0, r1, r2, r3, r4, r5}.AsDecodedRunes(), r6.AsDecodedRunes()))
			}

			r.commit(cursorio.DecodedRuneList{r0, r1, r2, r3, r4, r5, r6}.AsDecodedRunes())

			return readerStack{
				ectx,
				scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
					if err != nil {
						return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(err))
					}

					prefixToken, err := r.producePNAME_NS(r0)
					if err != nil {
						return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(err))
					}

					return readerStack{
						ectx,
						scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
							if err != nil {
								return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(err))
							}

							expandedToken, err := r.produceIRIREF(r0)
							if err != nil {
								return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(err))
							}

							resolvedExpanded, err := ectx.ResolveURL(expandedToken.Decoded)
							if err != nil {
								return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(grammar.R_IRIREF.ErrWithTextOffsetRange(err, expandedToken.Offsets)))
							}

							return readerStack{
								ectx,
								scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
									if err != nil {
										return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(r.newOffsetError(err, cursorio.DecodedRunes{}, cursorio.DecodedRunes{})))
									} else if r0.Rune != '.' {
										return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes())))
									}

									r.commit(r0.AsDecodedRunes())

									ectx.Global.Prefixes.AddPrefixMappings(iri.PrefixMapping{
										Prefix:   prefixToken.Decoded,
										Expanded: resolvedExpanded.String(),
									})

									if r.prefixDirectiveListener != nil {
										r.prefixDirectiveListener(DecoderEvent_PrefixDirective_Data{
											Prefix:          prefixToken.Decoded,
											PrefixOffsets:   prefixToken.Offsets,
											Expanded:        resolvedExpanded.String(),
											ExpandedOffsets: expandedToken.Offsets,
										})
									}

									return readerStack{ectx, reader_scanStatement}, nil
								}),
							}, nil
						}),
					}, nil
				}),
			}, nil
		default:
			return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1.Rune}, r0.AsDecodedRunes(), r1.AsDecodedRunes()))
		}

	// BASE directive
	case 'B', 'b':
		r1, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, r0.AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if r1.Rune != 'A' && r1.Rune != 'a' {
			r.buf.BacktrackRunes(r0, r1)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r2, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1}.AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if r2.Rune != 'S' && r2.Rune != 's' {
			r.buf.BacktrackRunes(r0, r1, r2)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r3, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1, r2}.AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if r3.Rune != 'E' && r3.Rune != 'e' {
			r.buf.BacktrackRunes(r0, r1, r2, r3)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r4, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1, r2, r3}.AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if r4.Rune == '<' {
			r.commit(cursorio.DecodedRuneList{r0, r1, r2, r3}.AsDecodedRunes())
			r.buf.BacktrackRunes(r4)
		} else if !unicode.IsSpace(r4.Rune) { // TODO IsRune_WS
			r.buf.BacktrackRunes(r0, r1, r2, r3, r4)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		} else {
			r.commit(cursorio.DecodedRuneList{r0, r1, r2, r3, r4}.AsDecodedRunes())
		}

		return readerStack{
			ectx,
			scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
				if err != nil {
					return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlBase.Err(err))
				}

				token, err := r.produceIRIREF(r0)
				if err != nil {
					return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlBase.Err(err))
				}

				resolvedBase, err := ectx.ResolveURL(token.Decoded)
				if err != nil {
					return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlBase.Err(grammar.R_IRIREF.ErrWithTextOffsetRange(err, token.Offsets)))
				}

				ectx.Global.Base = resolvedBase

				if r.baseDirectiveListener != nil {
					r.baseDirectiveListener(DecoderEvent_BaseDirective_Data{
						Value:        resolvedBase.String(),
						ValueOffsets: token.Offsets,
					})
				}

				return readerStack{ectx, reader_scanStatement}, nil
			}),
		}, nil

	// PREFIX directive
	case 'P', 'p':
		r1, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, r0.AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if r1.Rune != 'R' && r1.Rune != 'r' {
			r.buf.BacktrackRunes(r0, r1)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r2, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1}.AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if r2.Rune != 'E' && r2.Rune != 'e' {
			r.buf.BacktrackRunes(r0, r1, r2)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r3, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1, r2}.AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if r3.Rune != 'F' && r3.Rune != 'f' {
			r.buf.BacktrackRunes(r0, r1, r2, r3)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r4, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1, r2, r3}.AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if r4.Rune != 'I' && r4.Rune != 'i' {
			r.buf.BacktrackRunes(r0, r1, r2, r3, r4)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r5, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1, r2, r3, r4}.AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if r5.Rune != 'X' && r5.Rune != 'x' {
			r.buf.BacktrackRunes(r0, r1, r2, r3, r4, r5)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r6, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, cursorio.DecodedRuneList{r0, r1, r2, r3, r4, r5}.AsDecodedRunes(), cursorio.DecodedRunes{}))
		} else if !unicode.IsSpace(r6.Rune) { // TODO IsRune_WS
			r.buf.BacktrackRunes(r0, r1, r2, r3, r4, r5, r6)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r.commit(cursorio.DecodedRuneList{r0, r1, r2, r3, r4, r5, r6}.AsDecodedRunes())

		return readerStack{
			ectx,
			scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
				if err != nil {
					return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlPrefix.Err(err))
				}

				prefixToken, err := r.producePNAME_NS(r0)
				if err != nil {
					return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlPrefix.Err(err))
				}

				return readerStack{
					ectx,
					scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
						if err != nil {
							return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlPrefix.Err(err))
						}

						expandedToken, err := r.produceIRIREF(r0)
						if err != nil {
							return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlPrefix.Err(err))
						}

						resolvedExpanded, err := ectx.ResolveURL(expandedToken.Decoded)
						if err != nil {
							return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlPrefix.Err(grammar.R_IRIREF.ErrWithTextOffsetRange(err, expandedToken.Offsets)))
						}

						ectx.Global.Prefixes.AddPrefixMappings(iri.PrefixMapping{
							Prefix:   prefixToken.Decoded,
							Expanded: resolvedExpanded.String(),
						})

						if r.prefixDirectiveListener != nil {
							r.prefixDirectiveListener(DecoderEvent_PrefixDirective_Data{
								Prefix:          prefixToken.Decoded,
								PrefixOffsets:   prefixToken.Offsets,
								Expanded:        resolvedExpanded.String(),
								ExpandedOffsets: expandedToken.Offsets,
							})
						}

						return readerStack{ectx, reader_scanStatement}, nil
					}),
				}, nil
			}),
		}, nil
	case '<':
		r.buf.BacktrackRunes(r0)

		r.pushState(ectx, reader_scan_Triples_End)

		return readerStack{ectx, reader_scan_Triples_Subject_IRIREF}, nil
	case '_':
		// TODO require : since static?
		r.buf.BacktrackRunes(r0)

		r.pushState(ectx, reader_scan_Triples_End)

		return readerStack{ectx, reader_scan_Triples_Subject_BlankNode}, nil
	case '[':
		blankNode := ectx.Global.BlankNodeStringFactory.NewBlankNode()
		blankNodeRange := r.commitForTextOffsetRange(r0.AsDecodedRunes())

		ectx.CurSubject = blankNode
		ectx.CurSubjectLocation = blankNodeRange

		return readerStack{ectx, reader_scanStatement_Subject_AnonOrBlankNode}, nil
	case '(':
		r.pushState(ectx, reader_scan_Triples_End)

		nectx := ectx
		nectx.CurSubject = ectx.Global.BlankNodeStringFactory.NewBlankNode()
		nectx.CurSubjectLocation = r.commitForTextOffsetRange(r0.AsDecodedRunes())

		r.pushState(nectx, reader_scan_PredicateObjectList)

		fn := scanFunc(func(r *Decoder, ectx evaluationContext, r0 cursorio.DecodedRune, err error) (readerStack, error) {
			return reader_scan_collection(r, ectx, r0, nectx.CurSubject, nectx.CurSubjectLocation)
		})

		return readerStack{ectx, fn}, nil
	}

	if r0.Rune == ':' || internal.IsRune_PN_CHARS_BASE(r0.Rune) {
		r.buf.BacktrackRunes(r0)

		r.pushState(ectx, reader_scan_Triples_End)

		return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
	}

	return readerStack{}, grammar.R_statement.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0.Rune}, cursorio.DecodedRunes{}, r0.AsDecodedRunes()))
}
