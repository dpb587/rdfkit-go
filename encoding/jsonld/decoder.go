package jsonld

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/internal/jsonldinternal"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
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
	r io.Reader

	defaultBase string

	captureTextOffsets bool
	initialTextOffset  cursorio.TextOffset

	parserOptions []inspectjson.ParserOption

	processingMode string
	documentLoader jsonldtype.DocumentLoader
	rdfDirection   string

	baseDirectiveListener   DecoderEvent_BaseDirective_ListenerFunc
	prefixDirectiveListener DecoderEvent_PrefixDirective_ListenerFunc
	buildTextOffsets        encodingutil.TextOffsetsBuilderFunc

	statements    []*statement
	statementsIdx int
	err           error
}

var _ encoding.DatasetDecoder = &Decoder{}

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
	return d.statements[d.statementsIdx].graphName
}

func (d *Decoder) GetTriple() rdf.Triple {
	return d.statements[d.statementsIdx].triple
}

func (d *Decoder) GetStatement() rdfio.Statement {
	return d.statements[d.statementsIdx]
}

func (r *Decoder) parseRoot() error {
	topt := inspectjson.TokenizerConfig{}

	if r.captureTextOffsets {
		topt = topt.SetSourceInitialOffset(r.initialTextOffset)
	}

	ts, err := inspectjson.Parse(r.r, append(r.parserOptions, topt)...)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	opts := jsonldtype.ProcessorOptions{
		ProcessingMode: r.processingMode,
		DocumentLoader: r.documentLoader,
	}

	if len(r.defaultBase) > 0 {
		opts.BaseURL = r.defaultBase
	}

	ets, err := jsonldinternal.Expand(ts, opts)
	if err != nil {
		return fmt.Errorf("expand: %w", err)
	}

	ectx := evaluationContext{
		global: &globalEvaluationContext{
			BlankNodeStringMapper: blanknodeutil.NewStringMapper(),
			BlankNodeFactory:      blanknodeutil.NewFactory(),
		},
		CurrentContainer: &DocumentResource{},
		ActiveGraph:      rdf.DefaultGraph,
	}

	return r.decodeElement(ectx, ets, false)
}

func (r *Decoder) decodeElement(ectx evaluationContext, element inspectjson.Value, dropValuePropertyRange bool) error {
	if elementArray, ok := element.(inspectjson.ArrayValue); ok {
		for _, item := range elementArray.Values {
			err := r.decodeElement(ectx, item, dropValuePropertyRange)
			if err != nil {
				return err
			}
		}

		return nil
	}

	elementObject := element.(inspectjson.ObjectValue)

	if _, ok := elementObject.Members["@value"]; ok {
		return r.decodeValueNode(ectx, elementObject, dropValuePropertyRange)
	}

	if ectx.ActiveProperty != nil {
		// hacky to drop outer document container
		ectx.CurrentContainer = nil
	}

	if atList, ok := elementObject.Members["@list"]; ok {
		listArray := atList.Value.(inspectjson.ArrayValue)

		if len(listArray.Values) == 0 {
			r.statements = append(r.statements, &statement{
				graphName: ectx.ActiveGraph,
				triple: rdf.Triple{
					Subject:   ectx.ActiveSubject,
					Predicate: ectx.ActiveProperty,
					Object:    rdfiri.Nil_List,
				},
				offsets: r.buildTextOffsets(
					encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
					encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
					encoding.PredicateStatementOffsets, ectx.ActivePropertyRange,
					// TODO range iff BeginToken/EndToken known
					// Object,    atList.Value.BeginToken.OffsetRange.NewUntilOffset(atListArray.EndToken.OffsetRange.UntilOffset()),
				),
				containerResource: ectx.CurrentContainer,
			})
		} else {
			listSubject := ectx.global.BlankNodeFactory.NewBlankNode()

			var propagatePropertyRange *cursorio.TextOffsetRange

			if vObject, ok := listArray.Values[0].(inspectjson.ObjectValue); ok {
				propagatePropertyRange = r.tryReplacedMember(vObject)
			}

			r.statements = append(r.statements, &statement{
				graphName: ectx.ActiveGraph,
				triple: rdf.Triple{
					Subject:   ectx.ActiveSubject,
					Predicate: ectx.ActiveProperty,
					Object:    listSubject,
				},
				offsets: r.buildTextOffsets(
					encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
					encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
					encoding.PredicateStatementOffsets, propagatePropertyRange,
				),
				containerResource: ectx.CurrentContainer,
			})

			for listIdx, listValue := range listArray.Values {
				if listIdx > 0 {
					nextListSubject := ectx.global.BlankNodeFactory.NewBlankNode()

					r.statements = append(r.statements, &statement{
						graphName: ectx.ActiveGraph,
						triple: rdf.Triple{
							Subject:   listSubject,
							Predicate: rdfiri.Rest_Property,
							Object:    nextListSubject,
						},
						offsets: r.buildTextOffsets(
							encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
						),
						containerResource: ectx.CurrentContainer,
					})

					listSubject = nextListSubject
				}

				lctx := ectx
				lctx.ActiveSubject = listSubject
				lctx.ActiveSubjectRange = nil
				lctx.ActiveProperty = rdfiri.First_Property
				lctx.ActivePropertyRange = propagatePropertyRange

				err := r.decodeElement(lctx, listValue, true)
				if err != nil {
					return fmt.Errorf("decode list item: %v", err)
				}
			}

			r.statements = append(r.statements, &statement{
				graphName: ectx.ActiveGraph,
				triple: rdf.Triple{
					Subject:   listSubject,
					Predicate: rdfiri.Rest_Property,
					Object:    rdfiri.Nil_List,
				},
				offsets: r.buildTextOffsets(
					encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
				),
				containerResource: ectx.CurrentContainer,
			})
		}

		return nil
	}

	var selfSubject rdf.SubjectValue
	var selfSubjectRange *cursorio.TextOffsetRange

	if atID, ok := elementObject.Members["@id"]; ok {
		if _, ok := atID.Value.(inspectjson.NullValue); ok {
			return nil
		}

		idString := atID.Value.(inspectjson.StringValue).Value

		if strings.HasPrefix(idString, "_:") {
			selfSubject = ectx.global.BlankNodeStringMapper.MapBlankNodeIdentifier(idString[2:])
		} else {
			selfSubject = rdf.IRI(idString)
		}

		selfSubjectRange = atID.Value.GetSourceOffsets()
	} else {
		selfSubject = ectx.global.BlankNodeFactory.NewBlankNode()
	}

	if ectx.ActiveProperty != nil {
		if ectx.Reverse {
			r.statements = append(r.statements, &statement{
				graphName: ectx.ActiveGraph,
				triple: rdf.Triple{
					Subject:   selfSubject,
					Predicate: ectx.ActiveProperty,
					Object:    ectx.ActiveSubject,
				},
				offsets: r.buildTextOffsets(
					encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
					encoding.SubjectStatementOffsets, selfSubjectRange,
					encoding.PredicateStatementOffsets, ectx.ActivePropertyRange,
					encoding.ObjectStatementOffsets, ectx.ActiveSubjectRange,
				),
				containerResource: ectx.CurrentContainer,
			})

			ectx.Reverse = false
		} else {
			r.statements = append(r.statements, &statement{
				graphName: ectx.ActiveGraph,
				triple: rdf.Triple{
					Subject:   ectx.ActiveSubject,
					Predicate: ectx.ActiveProperty,
					Object:    selfSubject,
				},
				offsets: r.buildTextOffsets(
					encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
					encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
					encoding.PredicateStatementOffsets, r.tryReplacedMember(elementObject),
					encoding.ObjectStatementOffsets, selfSubjectRange,
				),
				containerResource: ectx.CurrentContainer,
			})
		}
	}

	ectx.ActiveSubject = selfSubject
	ectx.ActiveSubjectRange = selfSubjectRange

	if atReverse, ok := elementObject.Members["@reverse"]; ok {
		// TODO double reverse

		reverseObject := atReverse.Value.(inspectjson.ObjectValue)

		for key, member := range reverseObject.Members {
			// if len(key) > 1 && key[0] == '@' {
			// 	continue
			// }

			nectx := ectx
			nectx.ActiveProperty = rdf.IRI(key)
			nectx.ActivePropertyRange = member.Name.SourceOffsets
			nectx.Reverse = true

			for _, item := range member.Value.(inspectjson.ArrayValue).Values {
				err := r.decodeElement(nectx, item, false)
				if err != nil {
					return err
				}
			}
		}
	}

	if atType, ok := elementObject.Members["@type"]; ok {
		for _, typeValue := range atType.Value.(inspectjson.ArrayValue).Values {
			typePrimitive := typeValue.(inspectjson.StringValue)

			r.statements = append(r.statements, &statement{
				graphName: ectx.ActiveGraph,
				triple: rdf.Triple{
					Subject:   ectx.ActiveSubject,
					Predicate: rdfiri.Type_Property,
					Object:    rdf.IRI(typePrimitive.Value),
				},
				offsets: r.buildTextOffsets(
					encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
					encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
					// TODO need an alternative approach for propagating object types
					// i.e. value is always a string, so no equivalent to ReplacedMembers hack
					// encoding.PredicateStatementOffsets, atType.Name.SourceOffsets,
					encoding.ObjectStatementOffsets, typePrimitive.SourceOffsets,
				),
				containerResource: ectx.CurrentContainer,
			})
		}
	}

	if atGraph, ok := elementObject.Members["@graph"]; ok {
		nectx := ectx
		nectx.ActiveGraph = ectx.ActiveSubject.(rdf.GraphNameValue)
		nectx.ActiveGraphRange = ectx.ActiveSubjectRange
		nectx.ActiveSubject = nil
		nectx.ActiveSubjectRange = nil
		nectx.ActiveProperty = nil
		nectx.ActivePropertyRange = nil
		nectx.CurrentContainer = nil

		for _, item := range atGraph.Value.(inspectjson.ArrayValue).Values {
			err := r.decodeElement(nectx, item, false)
			if err != nil {
				return err
			}
		}
	}

	for key, member := range elementObject.Members {
		if len(key) > 1 && key[0] == '@' {
			continue
		}

		nectx := ectx
		nectx.ActiveProperty = rdf.IRI(key)
		nectx.ActivePropertyRange = member.Name.SourceOffsets

		for _, item := range member.Value.(inspectjson.ArrayValue).Values {
			err := r.decodeElement(nectx, item, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Decoder) decodeValueNode(ectx evaluationContext, v inspectjson.ObjectValue, dropPropertyRange bool) error {
	var lit rdf.Literal

	if atType, ok := v.Members["@type"]; ok {
		lit.Datatype = rdf.IRI(atType.Value.(inspectjson.StringValue).Value)
	}

	if lit.Datatype == "@json" {
		buf, err := json.Marshal(v.Members["@value"].Value.AsBuiltin())
		if err != nil {
			return fmt.Errorf("marshal for @json: %v", err)
		}

		lit.Datatype = rdfiri.JSON_Datatype
		lit.LexicalForm = string(buf)
	} else {
		switch valuePrimitive := v.Members["@value"].Value.(type) {
		case inspectjson.StringValue:
			lit.LexicalForm = valuePrimitive.Value

			atLanguage, atLangageKnown := v.Members["@language"]
			atDirection, atDirectionKnown := v.Members["@direction"]

			if atLangageKnown || atDirectionKnown {
				lit.Tags = map[rdf.LiteralTag]string{}

				if atLangageKnown {
					lit.Tags[rdf.LanguageLiteralTag] = atLanguage.Value.(inspectjson.StringValue).Value
				}

				if atDirectionKnown {
					lit.Tags[rdf.BaseDirectionLiteralTag] = atDirection.Value.(inspectjson.StringValue).Value
				}

				if len(lit.Datatype) == 0 {
					if atDirectionKnown && len(r.rdfDirection) > 0 {
						if r.rdfDirection == "i18n-datatype" {
							lit.Datatype = rdf.IRI(fmt.Sprintf(
								"https://www.w3.org/ns/i18n#%s_%s",
								strings.ToLower(lit.Tags[rdf.LanguageLiteralTag]),
								lit.Tags[rdf.BaseDirectionLiteralTag],
							))
						} else if r.rdfDirection == "compound-literal" {
							compoundNode := ectx.global.BlankNodeFactory.NewBlankNode()

							lit.Datatype = xsdiri.String_Datatype

							r.statements = append(r.statements,
								&statement{
									graphName: ectx.ActiveGraph,
									triple: rdf.Triple{
										Subject:   ectx.ActiveSubject,
										Predicate: ectx.ActiveProperty,
										Object:    compoundNode,
									},
									offsets: r.buildTextOffsets(
										encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
										encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
										encoding.PredicateStatementOffsets, r.tryReplacedMember(v),
									),
									containerResource: ectx.CurrentContainer,
								},
								&statement{
									graphName: ectx.ActiveGraph,
									triple: rdf.Triple{
										Subject:   compoundNode,
										Predicate: rdfiri.Value_Property,
										Object: rdf.Literal{
											Datatype:    xsdiri.String_Datatype,
											LexicalForm: lit.LexicalForm,
										},
									},
									offsets: r.buildTextOffsets(
										encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
										encoding.ObjectStatementOffsets, v.Members["@value"].Value.GetSourceOffsets(),
									),
									containerResource: ectx.CurrentContainer,
								},
								&statement{
									graphName: ectx.ActiveGraph,
									triple: rdf.Triple{
										Subject:   compoundNode,
										Predicate: rdfiri.Direction_Property,
										Object: rdf.Literal{
											Datatype:    xsdiri.String_Datatype,
											LexicalForm: lit.Tags[rdf.BaseDirectionLiteralTag],
										},
									},
									offsets: r.buildTextOffsets(
										encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
										encoding.ObjectStatementOffsets, atDirection.Value.GetSourceOffsets(),
									),
									containerResource: ectx.CurrentContainer,
								},
							)

							if atLangageKnown {
								r.statements = append(r.statements,
									&statement{
										graphName: ectx.ActiveGraph,
										triple: rdf.Triple{
											Subject:   compoundNode,
											Predicate: rdfiri.Language_Property,
											Object: rdf.Literal{
												Datatype:    xsdiri.String_Datatype,
												LexicalForm: strings.ToLower(lit.Tags[rdf.LanguageLiteralTag]),
											},
										},
										offsets: r.buildTextOffsets(
											encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
											encoding.ObjectStatementOffsets, atLanguage.Value.GetSourceOffsets(),
										),
										containerResource: ectx.CurrentContainer,
									},
								)
							}

							return nil
						}
					} else if atLangageKnown {
						lit.Datatype = rdfiri.LangString_Datatype
					}
				}
			}
		case inspectjson.NumberValue:
			lit.LexicalForm = strconv.FormatFloat(valuePrimitive.Value, 'f', -1, 64)

			if len(lit.Datatype) == 0 {
				if strings.Contains(lit.LexicalForm, ".") {
					lit.Datatype = xsdiri.Double_Datatype
				} else {
					lit.Datatype = xsdiri.Integer_Datatype
				}
			}

			// based on testsuites (#t0035, #te061)
			// probably related to canonicalization recommendations in spec and XML datatypes, though didn't find exact
			// reference for this behavior

			// if !strings.Contains(lit.LexicalForm, ".") {
			// 	lit.LexicalForm += ".0"
			// }

			if strings.Contains(lit.LexicalForm, ".") && !strings.Contains(lit.LexicalForm, "e") && !strings.Contains(lit.LexicalForm, "E") {
				lit.LexicalForm += "E0"
			}
		case inspectjson.BooleanValue:
			if len(lit.Datatype) == 0 {
				lit.Datatype = xsdiri.Boolean_Datatype
			}

			if valuePrimitive.Value {
				lit.LexicalForm = "true"
			} else {
				lit.LexicalForm = "false"
			}
		default:
			return fmt.Errorf("unexpected value type: %v", valuePrimitive.GetGrammarName())
		}
	}

	if len(lit.Datatype) == 0 {
		lit.Datatype = xsdiri.String_Datatype
	}

	var predicateOffsets *cursorio.TextOffsetRange

	if !dropPropertyRange {
		predicateOffsets = r.tryReplacedMember(v)
	}

	r.statements = append(r.statements, &statement{
		graphName: ectx.ActiveGraph,
		triple: rdf.Triple{
			Subject:   ectx.ActiveSubject,
			Predicate: ectx.ActiveProperty,
			Object:    lit,
		},
		offsets: r.buildTextOffsets(
			encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
			encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
			encoding.PredicateStatementOffsets, predicateOffsets,
			encoding.ObjectStatementOffsets, v.Members["@value"].Value.GetSourceOffsets(),
		),
		containerResource: ectx.CurrentContainer,
	})

	return nil
}

func (r *Decoder) tryReplacedMember(v inspectjson.ObjectValue) *cursorio.TextOffsetRange {
	if len(v.ReplacedMembers) == 0 {
		return nil
	}

	rm := v.ReplacedMembers[len(v.ReplacedMembers)-1].Name

	if rm.Value == jsonldinternal.MagicKeywordPropertySourceOffsets {
		return rm.SourceOffsets
	}

	return nil
}
