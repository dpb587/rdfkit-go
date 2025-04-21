package rdfjson

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type DecoderOption interface {
	apply(s *DecoderConfig)
	newDecoder(r io.Reader) (*Decoder, error)
}

type Decoder struct {
	r                io.Reader
	topts            []inspectjson.TokenizerOptionsApplier
	buildTextOffsets encodingutil.TextOffsetsBuilderFunc

	statements    []*statement
	statementsIdx int
	err           error
}

var _ encoding.GraphDecoder = &Decoder{}

func NewDecoder(r io.Reader, opts ...DecoderOption) (*Decoder, error) {
	compiledOpts := DecoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newDecoder(r)
}

func (d *Decoder) Close() error {
	return nil
}

func (d *Decoder) Err() error {
	return d.err
}

func (d *Decoder) Next() bool {
	if d.err != nil {
		return false
	} else if d.statementsIdx == -1 {
		d.err = d.parseRoot()
	}

	d.statementsIdx++

	return d.statementsIdx < len(d.statements)
}

func (d *Decoder) GetGraphName() rdf.GraphNameValue {
	return rdf.DefaultGraph
}

func (d *Decoder) GetTriple() rdf.Triple {
	return d.statements[d.statementsIdx].triple
}

func (d *Decoder) GetStatement() rdfio.Statement {
	return d.statements[d.statementsIdx]
}

func (r *Decoder) parseRoot() error {
	t := inspectjson.NewTokenizer(r.r, r.topts...)

	bnodeMap := blanknodeutil.NewStringMapper()

	var subjectValue rdf.SubjectValue
	var subjectOffset *cursorio.TextOffsetRange
	var predicateValue rdf.PredicateValue
	var predicateOffset *cursorio.TextOffsetRange

	token, err := t.Next()
	if err != nil {
		return err
	} else if _, ok := token.(inspectjson.BeginObjectToken); !ok {
		return fmt.Errorf("unexpected token: %v", token.GetGrammarName())
	}

	for {
		token, err = t.Next()
		if err != nil {
			return err
		} else if _, ok := token.(inspectjson.EndObjectToken); ok {
			break
		} else if _, ok := token.(inspectjson.ValueSeparatorToken); ok {
			continue
		}

		tokenScalarString, ok := token.(inspectjson.StringToken)
		if !ok {
			return fmt.Errorf("unexpected token: %v", token.GetGrammarName())
		}

		if strings.HasPrefix(tokenScalarString.Content, "_:") {
			subjectValue = bnodeMap.MapBlankNodeIdentifier(tokenScalarString.Content[2:])
		} else {
			subjectValue = rdf.IRI(tokenScalarString.Content)
		}

		subjectOffset = tokenScalarString.SourceOffsets

		token, err = t.Next()
		if err != nil {
			return err
		} else if _, ok := token.(inspectjson.NameSeparatorToken); !ok {
			return fmt.Errorf("unexpected token: %v", token.GetGrammarName())
		}

		token, err = t.Next()
		if err != nil {
			return err
		} else if _, ok := token.(inspectjson.BeginObjectToken); !ok {
			return fmt.Errorf("unexpected token: %v", token.GetGrammarName())
		}

		for {
			token, err = t.Next()
			if err != nil {
				return err
			} else if _, ok := token.(inspectjson.EndObjectToken); ok {
				break
			} else if _, ok := token.(inspectjson.ValueSeparatorToken); ok {
				continue
			}

			tokenString := token.(inspectjson.StringToken)
			if !ok {
				return fmt.Errorf("unexpected token: %v", token.GetGrammarName())
			}

			predicateValue = rdf.IRI(tokenString.Content)
			predicateOffset = tokenString.SourceOffsets

			token, err = t.Next()
			if err != nil {
				return err
			} else if _, ok := token.(inspectjson.NameSeparatorToken); !ok {
				return fmt.Errorf("unexpected token: %v", token.GetGrammarName())
			}

			token, err = t.Next()
			if err != nil {
				return err
			} else if _, ok := token.(inspectjson.BeginArrayToken); !ok {
				return fmt.Errorf("unexpected token: %v", token.GetGrammarName())
			}

			for {
				token, err = t.Next()
				if err != nil {
					return err
				} else if _, ok := token.(inspectjson.EndArrayToken); ok {
					break
				} else if _, ok := token.(inspectjson.ValueSeparatorToken); ok {
					continue
				}

				tokenObjectBegin, ok := token.(inspectjson.BeginObjectToken)
				if !ok {
					return fmt.Errorf("unexpected token: %v", token.GetGrammarName())
				}

				var objectOffsetRange *cursorio.TextOffsetRange = tokenObjectBegin.SourceOffsets

				objectMembers := struct {
					Datatype *inspectjson.StringToken
					Lang     *inspectjson.StringToken
					Type     *inspectjson.StringToken
					Value    *inspectjson.StringToken
				}{}

				for {
					nameToken, err := t.Next()
					if err != nil {
						return err
					} else if nameTokenObjectEnd, ok := nameToken.(inspectjson.EndObjectToken); ok {
						if objectOffsetRange != nil {
							objectOffsetRange = &cursorio.TextOffsetRange{
								From:  objectOffsetRange.From,
								Until: nameTokenObjectEnd.SourceOffsets.Until,
							}
						}

						break
					} else if _, ok := nameToken.(inspectjson.ValueSeparatorToken); ok {
						continue
					}

					nameTokenString := nameToken.(inspectjson.StringToken)
					if !ok {
						return fmt.Errorf("unexpected token: %v", nameToken.GetGrammarName())
					}

					switch nameTokenString.Content {
					case "datatype", "lang", "type", "value":
						// ok
					default:
						return fmt.Errorf("unexpected key: %s", nameTokenString.Content)
					}

					token, err = t.Next()
					if err != nil {
						return err
					} else if _, ok := token.(inspectjson.NameSeparatorToken); !ok {
						return fmt.Errorf("unexpected token: %v", token.GetGrammarName())
					}

					valueToken, err := t.Next()
					if err != nil {
						return err
					}

					valueTokenString, ok := valueToken.(inspectjson.StringToken)
					if !ok {
						return fmt.Errorf("unexpected token: %v", valueToken.GetGrammarName())
					}

					switch nameTokenString.Content {
					case "datatype":
						objectMembers.Datatype = &valueTokenString
					case "lang":
						objectMembers.Lang = &valueTokenString
					case "type":
						objectMembers.Type = &valueTokenString
					case "value":
						objectMembers.Value = &valueTokenString
					}
				}

				if objectMembers.Type == nil {
					return fmt.Errorf("missing key: type")
				}

				if objectMembers.Value == nil {
					return fmt.Errorf("missing key: value")
				}

				switch objectMembers.Type.Content {
				case "literal":
					if objectMembers.Datatype != nil {
						r.statements = append(r.statements, &statement{
							triple: rdf.Triple{
								Subject:   subjectValue,
								Predicate: predicateValue,
								Object: rdf.Literal{
									LexicalForm: objectMembers.Value.Content,
									Datatype:    rdf.IRI(objectMembers.Datatype.Content),
								},
							},
							offsets: r.buildTextOffsets(
								encoding.SubjectStatementOffsets, subjectOffset,
								encoding.PredicateStatementOffsets, predicateOffset,
								encoding.ObjectStatementOffsets, objectOffsetRange,
							),
						})
					} else if objectMembers.Lang != nil {
						r.statements = append(r.statements, &statement{
							triple: rdf.Triple{
								Subject:   subjectValue,
								Predicate: predicateValue,
								Object: rdf.Literal{
									LexicalForm: objectMembers.Value.Content,
									Datatype:    rdfiri.LangString_Datatype,
									Tags: map[rdf.LiteralTag]string{
										rdf.LanguageLiteralTag: objectMembers.Lang.Content,
									},
								},
							},
							offsets: r.buildTextOffsets(
								encoding.SubjectStatementOffsets, subjectOffset,
								encoding.PredicateStatementOffsets, predicateOffset,
								encoding.ObjectStatementOffsets, objectOffsetRange,
							),
						})
					} else {
						r.statements = append(r.statements, &statement{
							triple: rdf.Triple{
								Subject:   subjectValue,
								Predicate: predicateValue,
								Object: rdf.Literal{
									LexicalForm: objectMembers.Value.Content,
									Datatype:    xsdiri.String_Datatype,
								},
							},
							offsets: r.buildTextOffsets(
								encoding.SubjectStatementOffsets, subjectOffset,
								encoding.PredicateStatementOffsets, predicateOffset,
								encoding.ObjectStatementOffsets, objectOffsetRange,
							),
						})
					}
				case "uri":
					r.statements = append(r.statements, &statement{
						triple: rdf.Triple{
							Subject:   subjectValue,
							Predicate: predicateValue,
							Object:    rdf.IRI(objectMembers.Value.Content),
						},
						offsets: r.buildTextOffsets(
							encoding.SubjectStatementOffsets, subjectOffset,
							encoding.PredicateStatementOffsets, predicateOffset,
							encoding.ObjectStatementOffsets, objectOffsetRange,
						),
					})
				case "bnode":
					v := objectMembers.Value.Content
					if !strings.HasPrefix(v, "_:") {
						return fmt.Errorf("invalid bnode value: %s", v)
					}

					r.statements = append(r.statements, &statement{
						triple: rdf.Triple{
							Subject:   subjectValue,
							Predicate: predicateValue,
							Object:    bnodeMap.MapBlankNodeIdentifier(v[2:]),
						},
						offsets: r.buildTextOffsets(
							encoding.SubjectStatementOffsets, subjectOffset,
							encoding.PredicateStatementOffsets, predicateOffset,
							encoding.ObjectStatementOffsets, objectOffsetRange,
						),
					})
				default:
					return fmt.Errorf("unexpected type: %s", objectMembers.Type.Content)
				}
			}
		}
	}

	token, err = t.Next()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}

		return err
	}

	return fmt.Errorf("unexpected token: %v", token.GetGrammarName())
}
