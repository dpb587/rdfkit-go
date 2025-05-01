package turtle

import (
	"errors"
	"io"
	"unicode"

	"github.com/dpb587/cursorio-go/x/cursorioutil"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
	"github.com/dpb587/rdfkit-go/rdf"
)

func reader_scanStatement_Subject_AnonOrBlankNode(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_subject.Err(err)
	} else if r0 == ']' {
		r.commit([]rune{r0})

		r.pushState(ectx, reader_scan_Triples_End)
		r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

		return readerStack{ectx, reader_scan_PredicateObjectList}, nil
	}

	r.buf.BacktrackRunes(r0)

	r.pushState(ectx, reader_scan_Triples_End)
	r.pushState(ectx, reader_scan_PredicateObjectList)
	r.pushState(ectx, reader_scan_blankNodePropertyList_End)
	r.pushState(ectx, reader_scan_PredicateObjectList_Continue)

	return readerStack{ectx, reader_scan_PredicateObjectList}, nil
}

func reader_scanStatement(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
	if err != nil {
		if errors.Is(err, io.EOF) {
			return r.terminate()
		}

		return readerStack{}, grammar.R_statement.Err(err)
	}

	r.pushState(ectx, reader_scanStatement)

	switch r0 {
	case '@':
		r1, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, err
		}

		switch r1 {

		// @base directive
		case 'b':
			r2, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, []rune{r0}, nil))
			} else if r2 != 'a' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2}, []rune{r0, r1}, []rune{r2}))
			}

			r3, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, []rune{r0, r1}, nil))
			} else if r3 != 's' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r3}, []rune{r0, r1, r2}, []rune{r3}))
			}

			r4, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, []rune{r0, r1, r2}, nil))
			} else if r4 != 'e' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r4}, []rune{r0, r1, r2, r3}, []rune{r4}))
			}

			r.commit([]rune{r0, r1, r2, r3, r4})

			return readerStack{
				ectx,
				scanFunc(func(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
					if err != nil {
						return readerStack{}, grammar.R_directive.Err(grammar.R_base.Err(err))
					}

					baseToken, err := r.produceIRIREF(r0)
					if err != nil {
						return readerStack{}, grammar.R_directive.Err(grammar.R_base.Err(err))
					}

					resolvedBase, err := ectx.ResolveURL(baseToken.Decoded)
					if err != nil {
						return readerStack{}, grammar.R_directive.Err(grammar.R_base.Err(grammar.R_IRIREF.ErrCursorRange(err, baseToken.Offsets)))
					}

					return readerStack{
						ectx,
						scanFunc(func(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
							if err != nil {
								return readerStack{}, grammar.R_directive.Err(grammar.R_base.Err(r.newOffsetError(err, nil, nil)))
							} else if r0 != '.' {
								return readerStack{}, grammar.R_directive.Err(grammar.R_base.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0})))
							}

							r.commit([]rune{r0})

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
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, []rune{r0}, nil))
			} else if r2 != 'r' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r2}, []rune{r0, r1}, []rune{r2}))
			}

			r3, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, []rune{r0, r1}, nil))
			} else if r3 != 'e' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r3}, []rune{r0, r1, r2}, []rune{r3}))
			}

			r4, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, []rune{r0, r1, r2}, nil))
			} else if r4 != 'f' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r4}, []rune{r0, r1, r2, r3}, []rune{r4}))
			}

			r5, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, []rune{r0, r1, r2, r3}, nil))
			} else if r5 != 'i' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r5}, []rune{r0, r1, r2, r3, r4}, []rune{r5}))
			}

			r6, err := r.buf.NextRune()
			if err != nil {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(err, []rune{r0, r1, r2, r3, r4}, nil))
			} else if r6 != 'x' {
				return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r6}, []rune{r0, r1, r2, r3, r4, r5}, []rune{r6}))
			}

			r.commit([]rune{r0, r1, r2, r3, r4, r5, r6})

			return readerStack{
				ectx,
				scanFunc(func(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
					if err != nil {
						return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(err))
					}

					prefixToken, err := r.producePNAME_NS(r0)
					if err != nil {
						return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(err))
					}

					return readerStack{
						ectx,
						scanFunc(func(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
							if err != nil {
								return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(err))
							}

							expandedToken, err := r.produceIRIREF(r0)
							if err != nil {
								return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(err))
							}

							resolvedExpanded, err := ectx.ResolveURL(expandedToken.Decoded)
							if err != nil {
								return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(grammar.R_IRIREF.ErrCursorRange(err, expandedToken.Offsets)))
							}

							return readerStack{
								ectx,
								scanFunc(func(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
									if err != nil {
										return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(r.newOffsetError(err, nil, nil)))
									} else if r0 != '.' {
										return readerStack{}, grammar.R_directive.Err(grammar.R_prefixID.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0})))
									}

									r.commit([]rune{r0})

									ectx.Global.Prefixes[prefixToken.Decoded] = rdf.IRI(resolvedExpanded.String())

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
			return readerStack{}, grammar.R_directive.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r1}, []rune{r0}, []rune{r1}))
		}

	// BASE directive
	case 'B', 'b':
		r1, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, []rune{r0}, nil))
		} else if r1 != 'A' && r1 != 'a' {
			r.buf.BacktrackRunes(r0, r1)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r2, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, []rune{r0, r1}, nil))
		} else if r2 != 'S' && r2 != 's' {
			r.buf.BacktrackRunes(r0, r1, r2)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r3, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, []rune{r0, r1, r2}, nil))
		} else if r3 != 'E' && r3 != 'e' {
			r.buf.BacktrackRunes(r0, r1, r2, r3)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r4, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, []rune{r0, r1, r2, r3}, nil))
		} else if r4 == '<' {
			r.commit([]rune{r0, r1, r2, r3})
			r.buf.BacktrackRunes(r4)
		} else if !unicode.IsSpace(r4) { // TODO IsRune_WS
			r.buf.BacktrackRunes(r0, r1, r2, r3, r4)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		} else {
			r.commit([]rune{r0, r1, r2, r3, r4})
		}

		return readerStack{
			ectx,
			scanFunc(func(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
				if err != nil {
					return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlBase.Err(err))
				}

				token, err := r.produceIRIREF(r0)
				if err != nil {
					return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlBase.Err(err))
				}

				resolvedBase, err := ectx.ResolveURL(token.Decoded)
				if err != nil {
					return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlBase.Err(grammar.R_IRIREF.ErrCursorRange(err, token.Offsets)))
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
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, []rune{r0}, nil))
		} else if r1 != 'R' && r1 != 'r' {
			r.buf.BacktrackRunes(r0, r1)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r2, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, []rune{r0, r1}, nil))
		} else if r2 != 'E' && r2 != 'e' {
			r.buf.BacktrackRunes(r0, r1, r2)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r3, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, []rune{r0, r1, r2}, nil))
		} else if r3 != 'F' && r3 != 'f' {
			r.buf.BacktrackRunes(r0, r1, r2, r3)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r4, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, []rune{r0, r1, r2, r3}, nil))
		} else if r4 != 'I' && r4 != 'i' {
			r.buf.BacktrackRunes(r0, r1, r2, r3, r4)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r5, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, []rune{r0, r1, r2, r3, r4}, nil))
		} else if r5 != 'X' && r5 != 'x' {
			r.buf.BacktrackRunes(r0, r1, r2, r3, r4, r5)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r6, err := r.buf.NextRune()
		if err != nil {
			return readerStack{}, grammar.R_statement.Err(r.newOffsetError(err, []rune{r0, r1, r2, r3, r4, r5}, nil))
		} else if !unicode.IsSpace(r6) { // TODO IsRune_WS
			r.buf.BacktrackRunes(r0, r1, r2, r3, r4, r5, r6)

			r.pushState(ectx, reader_scan_Triples_End)

			return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
		}

		r.commit([]rune{r0, r1, r2, r3, r4, r5, r6})

		return readerStack{
			ectx,
			scanFunc(func(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
				if err != nil {
					return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlPrefix.Err(err))
				}

				prefixToken, err := r.producePNAME_NS(r0)
				if err != nil {
					return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlPrefix.Err(err))
				}

				return readerStack{
					ectx,
					scanFunc(func(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
						if err != nil {
							return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlPrefix.Err(err))
						}

						expandedToken, err := r.produceIRIREF(r0)
						if err != nil {
							return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlPrefix.Err(err))
						}

						resolvedExpanded, err := ectx.ResolveURL(expandedToken.Decoded)
						if err != nil {
							return readerStack{}, grammar.R_statement.Err(grammar.R_sparqlPrefix.Err(grammar.R_IRIREF.ErrCursorRange(err, expandedToken.Offsets)))
						}

						ectx.Global.Prefixes[prefixToken.Decoded] = rdf.IRI(resolvedExpanded.String())

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
		blankNode := ectx.Global.BlankNodeFactory.NewBlankNode()
		blankNodeRange := r.commitForTextOffsetRange([]rune{r0})

		ectx.CurSubject = blankNode
		ectx.CurSubjectLocation = blankNodeRange

		return readerStack{ectx, reader_scanStatement_Subject_AnonOrBlankNode}, nil
	case '(':
		r.pushState(ectx, reader_scan_Triples_End)

		nectx := ectx
		nectx.CurSubject = ectx.Global.BlankNodeFactory.NewBlankNode()
		nectx.CurSubjectLocation = r.commitForTextOffsetRange([]rune{r0})

		r.pushState(nectx, reader_scan_PredicateObjectList)

		fn := scanFunc(func(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
			return reader_scan_collection(r, ectx, r0, nectx.CurSubject, nectx.CurSubjectLocation)
		})

		return readerStack{ectx, fn}, nil
	}

	if r0 == ':' || internal.IsRune_PN_CHARS_BASE(r0) {
		r.buf.BacktrackRunes(r0)

		r.pushState(ectx, reader_scan_Triples_End)

		return readerStack{ectx, reader_scan_Triples_Subject_PrefixedName}, nil
	}

	return readerStack{}, grammar.R_statement.Err(r.newOffsetError(cursorioutil.UnexpectedRuneError{Rune: r0}, nil, []rune{r0}))
}
