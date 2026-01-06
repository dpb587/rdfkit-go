package jsonld

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/internal/jsonldinternal"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldcontent"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

type statement struct {
	quad              rdf.Quad
	textOffsets       encoding.StatementTextOffsets
	containerResource encoding.ContainerResource
}

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
	expandContext  inspectjson.Value
	rdfDirection   string

	baseDirectiveListener   DecoderEvent_BaseDirective_ListenerFunc
	prefixDirectiveListener DecoderEvent_PrefixDirective_ListenerFunc
	buildTextOffsets        encodingutil.TextOffsetsBuilderFunc

	err error

	statements    []statement
	statementsIdx int
}

var _ encoding.QuadsDecoder = &Decoder{}
var _ encoding.StatementTextOffsetsProvider = &Decoder{}

func NewDecoder(r io.Reader, opts ...DecoderOption) (*Decoder, error) {
	compiledOpts := DecoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newDecoder(r)
}

func (d *Decoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return jsonldcontent.TypeIdentifier
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

func (r *Decoder) Quad() rdf.Quad {
	return r.statements[r.statementsIdx].quad
}

func (r *Decoder) Statement() rdf.Statement {
	return r.Quad()
}

func (r *Decoder) StatementTextOffsets() encoding.StatementTextOffsets {
	return r.statements[r.statementsIdx].textOffsets
}

// non-comprehensive validation functions
// TODO improve thoroughness? extract to shared package?

func isWellFormedIRI(s string) bool {
	var foundSchemeSeparator bool
	var foundFragment bool
	var foundQuery bool

	for _, r := range s {
		if r < 0x20 {
			return false
		}

		switch r {
		case ' ', '<', '>', '"', '{', '}', '|', '\\', '^', '`':
			return false
		case ':':
			foundSchemeSeparator = true
		case '?':
			if foundQuery {
				return false
			}

			foundQuery = true
		case '#':
			if foundFragment {
				return false
			}

			foundFragment = true
		}
	}

	return foundSchemeSeparator
}

func isWellFormedLiteralLanguageTag(s string) bool {
	return !strings.Contains(s, " ")
}

func isWellFormedLiteralBaseDirectionTag(s string) bool {
	return s == "ltr" || s == "rtl"
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
		ExpandContext:  r.expandContext,
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
	}

	return r.decodeElement(ectx, ets, false)
}

func (r *Decoder) decodeElement(ectx evaluationContext, element jsonldinternal.ExpandedValue, dropValuePropertyRange bool) error {
	if elementArray, ok := element.(*jsonldinternal.ExpandedArray); ok {
		for _, item := range elementArray.Values {
			err := r.decodeElement(ectx, item, dropValuePropertyRange)
			if err != nil {
				return err
			}
		}

		return nil
	}

	elementObject := element.(*jsonldinternal.ExpandedObject)

	if ectx.ActiveProperty != nil {
		// hacky to drop outer document container
		ectx.CurrentContainer = nil

		if _, ok := elementObject.Members["@value"]; ok {
			return r.decodeValueNode(ectx, elementObject, dropValuePropertyRange)
		}
	}

	if atList, ok := elementObject.Members["@list"]; ok {
		listArray := atList.(*jsonldinternal.ExpandedArray)

		if len(listArray.Values) == 0 {
			if ectx.ActiveProperty != nil {
				r.statements = append(r.statements, statement{
					quad: rdf.Quad{
						Triple: rdf.Triple{
							Subject:   ectx.ActiveSubject,
							Predicate: ectx.ActiveProperty,
							Object:    rdfiri.Nil_List,
						},
						GraphName: ectx.ActiveGraph,
					},
					textOffsets: r.buildTextOffsets(
						encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
						encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
						encoding.PredicateStatementOffsets, elementObject.PropertySourceOffsets,
						// TODO range iff BeginToken/EndToken known
						// Object,    atList.Value.BeginToken.OffsetRange.NewUntilOffset(atListArray.EndToken.OffsetRange.UntilOffset()),
					),
					containerResource: ectx.CurrentContainer,
				})
			}
		} else {
			listSubject := ectx.global.BlankNodeFactory.NewBlankNode()

			propagatePropertyRange := elementObject.PropertySourceOffsets

			if ectx.ActiveProperty != nil {
				r.statements = append(r.statements, statement{
					quad: rdf.Quad{
						Triple: rdf.Triple{
							Subject:   ectx.ActiveSubject,
							Predicate: ectx.ActiveProperty,
							Object:    listSubject,
						},
						GraphName: ectx.ActiveGraph,
					},
					textOffsets: r.buildTextOffsets(
						encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
						encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
						encoding.PredicateStatementOffsets, propagatePropertyRange,
					),
					containerResource: ectx.CurrentContainer,
				})
			}

			for listIdx, listValue := range listArray.Values {
				if listIdx > 0 && ectx.ActiveProperty != nil {
					nextListSubject := ectx.global.BlankNodeFactory.NewBlankNode()

					r.statements = append(r.statements, statement{
						quad: rdf.Quad{
							Triple: rdf.Triple{
								Subject:   listSubject,
								Predicate: rdfiri.Rest_Property,
								Object:    nextListSubject,
							},
							GraphName: ectx.ActiveGraph,
						},
						textOffsets: r.buildTextOffsets(
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

			if ectx.ActiveProperty != nil {
				r.statements = append(r.statements, statement{
					quad: rdf.Quad{
						Triple: rdf.Triple{
							Subject:   listSubject,
							Predicate: rdfiri.Rest_Property,
							Object:    rdfiri.Nil_List,
						},
						GraphName: ectx.ActiveGraph,
					},
					textOffsets: r.buildTextOffsets(
						encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
					),
					containerResource: ectx.CurrentContainer,
				})
			}
		}

		return nil
	}

	var selfSubject rdf.SubjectValue
	var selfSubjectRange *cursorio.TextOffsetRange

	if atID, ok := elementObject.Members["@id"]; ok {
		valuePrimitive := atID.(*jsonldinternal.ExpandedScalarPrimitive)

		if _, ok := valuePrimitive.Value.(inspectjson.NullValue); ok {
			return nil
		}

		idString := valuePrimitive.Value.(inspectjson.StringValue).Value

		if strings.HasPrefix(idString, "_:") {
			selfSubject = ectx.global.BlankNodeStringMapper.MapBlankNodeIdentifier(idString[2:])
		} else {
			if !isWellFormedIRI(idString) {
				// TODO warn
				return nil
			}

			selfSubject = rdf.IRI(idString)
		}

		selfSubjectRange = valuePrimitive.Value.GetSourceOffsets()
	} else {
		selfSubject = ectx.global.BlankNodeFactory.NewBlankNode()
		selfSubjectRange = elementObject.SourceOffsets
	}

	if ectx.ActiveProperty != nil {
		if ectx.Reverse {
			r.statements = append(r.statements, statement{
				quad: rdf.Quad{
					Triple: rdf.Triple{
						Subject:   selfSubject,
						Predicate: ectx.ActiveProperty,
						Object:    ectx.ActiveSubject,
					},
					GraphName: ectx.ActiveGraph,
				},
				textOffsets: r.buildTextOffsets(
					encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
					encoding.SubjectStatementOffsets, selfSubjectRange,
					encoding.PredicateStatementOffsets, ectx.ActivePropertyRange,
					encoding.ObjectStatementOffsets, ectx.ActiveSubjectRange,
				),
				containerResource: ectx.CurrentContainer,
			})

			ectx.Reverse = false
		} else {
			r.statements = append(r.statements, statement{
				quad: rdf.Quad{
					Triple: rdf.Triple{
						Subject:   ectx.ActiveSubject,
						Predicate: ectx.ActiveProperty,
						Object:    selfSubject,
					},
					GraphName: ectx.ActiveGraph,
				},
				textOffsets: r.buildTextOffsets(
					encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
					encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
					encoding.PredicateStatementOffsets, elementObject.PropertySourceOffsets,
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

		reverseObject := atReverse.(*jsonldinternal.ExpandedObject)

		// [dpb] sort keys for deterministic iteration; not found in spec?
		reverseKeys := slices.Collect(maps.Keys(reverseObject.Members))
		slices.SortFunc(reverseKeys, strings.Compare)

		for _, key := range reverseKeys {
			if len(key) > 1 && key[0] == '@' {
				continue
			}

			nectx := ectx

			if len(key) > 2 && key[:2] == "_:" {
				nectx.ActiveProperty = nil // not supported
			} else {
				if !isWellFormedIRI(key) {
					// TODO warn
					continue
				}

				nectx.ActiveProperty = rdf.IRI(key)
				// nectx.ActivePropertyRange = member.Name.SourceOffsets
			}

			nectx.Reverse = true

			for _, item := range reverseObject.Members[key].(*jsonldinternal.ExpandedArray).Values {
				err := r.decodeElement(nectx, item, false)
				if err != nil {
					return err
				}
			}
		}
	}

	if atType, ok := elementObject.Members["@type"]; ok {
		for _, typeValue := range atType.(*jsonldinternal.ExpandedArray).Values {
			typePrimitive := typeValue.(*jsonldinternal.ExpandedScalarPrimitive)
			typeString := typePrimitive.Value.(inspectjson.StringValue)

			var effectiveObject rdf.ObjectValue

			if len(typeString.Value) > 2 && typeString.Value[:2] == "_:" {
				effectiveObject = ectx.global.BlankNodeStringMapper.MapBlankNodeIdentifier(typeString.Value[2:])
			} else {
				if !isWellFormedIRI(typeString.Value) {
					// TODO warn
					continue
				}

				effectiveObject = rdf.IRI(typeString.Value)
			}

			r.statements = append(r.statements, statement{
				quad: rdf.Quad{
					Triple: rdf.Triple{
						Subject:   ectx.ActiveSubject,
						Predicate: rdfiri.Type_Property,
						Object:    effectiveObject,
					},
					GraphName: ectx.ActiveGraph,
				},
				textOffsets: r.buildTextOffsets(
					encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
					encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
					encoding.PredicateStatementOffsets, typePrimitive.PropertySourceOffsets,
					encoding.ObjectStatementOffsets, typePrimitive.Value.GetSourceOffsets(),
				),
				containerResource: ectx.CurrentContainer,
			})
		}
	}

	if atGraph, ok := elementObject.Members["@graph"]; ok {
		var validGraphName = true

		if graphIRI, ok := ectx.ActiveSubject.(rdf.IRI); ok {
			validGraphName = isWellFormedIRI(string(graphIRI))
		}

		if !validGraphName {
			// TODO warn
		} else {
			nectx := ectx
			nectx.ActiveGraph = ectx.ActiveSubject.(rdf.GraphNameValue)
			nectx.ActiveGraphRange = ectx.ActiveSubjectRange
			nectx.ActiveSubject = nil
			nectx.ActiveSubjectRange = nil
			nectx.ActiveProperty = nil
			nectx.ActivePropertyRange = nil
			nectx.CurrentContainer = nil

			for _, item := range atGraph.(*jsonldinternal.ExpandedArray).Values {
				err := r.decodeElement(nectx, item, false)
				if err != nil {
					return err
				}
			}
		}
	}

	if atIncluded, ok := elementObject.Members["@included"]; ok {
		nectx := ectx
		nectx.ActiveSubject = nil
		nectx.ActiveSubjectRange = nil
		nectx.ActiveProperty = nil
		nectx.ActivePropertyRange = nil

		for _, item := range atIncluded.(*jsonldinternal.ExpandedArray).Values {
			err := r.decodeElement(nectx, item, false)
			if err != nil {
				return err
			}
		}
	}

	// [dpb] Sort keys for deterministic iteration; not found in spec?
	memberKeys := slices.Collect(maps.Keys(elementObject.Members))
	slices.Sort(memberKeys)

	for _, key := range memberKeys {
		if len(key) > 1 && key[0] == '@' {
			continue
		}

		nectx := ectx

		if len(key) > 2 && key[:2] == "_:" {
			nectx.ActiveProperty = nil // not supported
		} else {
			if !isWellFormedIRI(key) {
				// TODO warn
				continue
			}

			nectx.ActiveProperty = rdf.IRI(key)
			// nectx.ActivePropertyRange = member.Name.SourceOffsets
		}

		for _, item := range elementObject.Members[key].(*jsonldinternal.ExpandedArray).Values {
			err := r.decodeElement(nectx, item, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Decoder) decodeValueNode(ectx evaluationContext, v *jsonldinternal.ExpandedObject, dropPropertyRange bool) error {
	var lit rdf.Literal

	if atType, ok := v.Members["@type"]; ok {
		lit.Datatype = rdf.IRI(atType.(*jsonldinternal.ExpandedScalarPrimitive).Value.(inspectjson.StringValue).Value)
	}

	atValuePrimitive := v.Members["@value"].(*jsonldinternal.ExpandedScalarPrimitive)

	if lit.Datatype == "@json" {
		buf := &bytes.Buffer{}

		jsonEncoder := json.NewEncoder(buf)
		jsonEncoder.SetEscapeHTML(false)

		err := jsonEncoder.Encode(atValuePrimitive.Value.AsBuiltin())
		if err != nil {
			return fmt.Errorf("marshal for @json: %v", err)
		}

		lit.Datatype = rdfiri.JSON_Datatype

		// drop the trailing \n
		lit.LexicalForm = string(buf.Bytes()[0 : len(buf.Bytes())-1])
	} else {
		switch valuePrimitive := atValuePrimitive.Value.(type) {
		case inspectjson.StringValue:
			lit.LexicalForm = valuePrimitive.Value

			atLanguage, atLangageKnown := v.Members["@language"]
			atDirection, atDirectionKnown := v.Members["@direction"]

			if atLangageKnown || atDirectionKnown {
				var litTagLanguage, litTagBaseDirection string

				if atLangageKnown {
					litTagLanguage = atLanguage.(*jsonldinternal.ExpandedScalarPrimitive).Value.(inspectjson.StringValue).Value
					if !isWellFormedLiteralLanguageTag(litTagLanguage) {
						// TODO warn
						return nil
					}
				}

				if atDirectionKnown {
					litTagBaseDirection = atDirection.(*jsonldinternal.ExpandedScalarPrimitive).Value.(inspectjson.StringValue).Value

					// spec does not explicitly call for a well-formed base direction?
					if !isWellFormedLiteralBaseDirectionTag(litTagBaseDirection) {
						// TODO warn
						return nil
					}
				}

				if len(lit.Datatype) == 0 {
					if atDirectionKnown && len(r.rdfDirection) > 0 {
						if r.rdfDirection == "i18n-datatype" {
							lit.Datatype = rdf.IRI(fmt.Sprintf(
								"https://www.w3.org/ns/i18n#%s_%s",
								strings.ToLower(litTagLanguage),
								litTagBaseDirection,
							))
						} else if r.rdfDirection == "compound-literal" {
							compoundNode := ectx.global.BlankNodeFactory.NewBlankNode()

							lit.Datatype = xsdiri.String_Datatype

							r.statements = append(r.statements,
								statement{
									quad: rdf.Quad{
										Triple: rdf.Triple{
											Subject:   ectx.ActiveSubject,
											Predicate: ectx.ActiveProperty,
											Object:    compoundNode,
										},
										GraphName: ectx.ActiveGraph,
									},
									textOffsets: r.buildTextOffsets(
										encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
										encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
										encoding.PredicateStatementOffsets, v.SourceOffsets,
									),
									containerResource: ectx.CurrentContainer,
								},
								statement{
									quad: rdf.Quad{
										Triple: rdf.Triple{
											Subject:   compoundNode,
											Predicate: rdfiri.Value_Property,
											Object: rdf.Literal{
												Datatype:    xsdiri.String_Datatype,
												LexicalForm: lit.LexicalForm,
											},
										},
										GraphName: ectx.ActiveGraph,
									},
									textOffsets: r.buildTextOffsets(
										encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
										encoding.ObjectStatementOffsets, valuePrimitive.SourceOffsets,
									),
									containerResource: ectx.CurrentContainer,
								},
								statement{
									quad: rdf.Quad{
										Triple: rdf.Triple{
											Subject:   compoundNode,
											Predicate: rdfiri.Direction_Property,
											Object: rdf.Literal{
												Datatype:    xsdiri.String_Datatype,
												LexicalForm: litTagBaseDirection,
											},
										},
										GraphName: ectx.ActiveGraph,
									},
									textOffsets: r.buildTextOffsets(
										encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
										encoding.ObjectStatementOffsets, atDirection.(*jsonldinternal.ExpandedScalarPrimitive).Value.GetSourceOffsets(),
									),
									containerResource: ectx.CurrentContainer,
								},
							)

							if atLangageKnown {
								r.statements = append(r.statements,
									statement{
										quad: rdf.Quad{
											Triple: rdf.Triple{
												Subject:   compoundNode,
												Predicate: rdfiri.Language_Property,
												Object: rdf.Literal{
													Datatype:    xsdiri.String_Datatype,
													LexicalForm: strings.ToLower(litTagLanguage),
												},
											},
											GraphName: ectx.ActiveGraph,
										},
										textOffsets: r.buildTextOffsets(
											encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
											encoding.ObjectStatementOffsets, atLanguage.(*jsonldinternal.ExpandedScalarPrimitive).Value.GetSourceOffsets(),
										),
										containerResource: ectx.CurrentContainer,
									},
								)
							}

							return nil
						} else {
							lit.Datatype = "http://www.w3.org/1999/02/22-rdf-syntax-ns#dirLangString" // RDF 1.2
							lit.Tag = rdf.DirectionalLanguageLiteralTag{
								Language:      litTagLanguage,
								BaseDirection: litTagBaseDirection,
							}
						}
					} else if atLangageKnown {
						lit.Datatype = rdfiri.LangString_Datatype
						lit.Tag = rdf.LanguageLiteralTag{
							Language: litTagLanguage,
						}
					}
				}
			}
		case inspectjson.NumberValue:
			fixedForm := strconv.FormatFloat(valuePrimitive.Value, 'f', -1, 64)
			hasDecimal := strings.Contains(fixedForm, ".")

			// zero not expected to be negative (#trt01)
			if valuePrimitive.Value == 0 && strings.HasPrefix(fixedForm, "-") {
				fixedForm = "0"
			}

			datatypeExplicitlySet := len(lit.Datatype) > 0

			if !datatypeExplicitlySet {
				if hasDecimal || (valuePrimitive.Value < math.MinInt32 || valuePrimitive.Value > math.MaxInt32) {
					lit.Datatype = xsdiri.Double_Datatype
				} else {
					lit.Datatype = xsdiri.Integer_Datatype
				}
			}

			// apparently explicit datatypes should follow double conventions, too (#te061)
			if lit.Datatype == xsdiri.Double_Datatype || (datatypeExplicitlySet && hasDecimal) {
				sciForm := strconv.FormatFloat(valuePrimitive.Value, 'E', -1, 64)
				parts := strings.Split(sciForm, "E")

				if len(parts) == 2 {
					mantissa := parts[0]
					expStr := parts[1]

					if !strings.Contains(mantissa, ".") {
						mantissa += ".0"
					}

					// remove leading plus sign and zeros
					expVal := 0
					fmt.Sscanf(expStr, "%d", &expVal)

					lit.LexicalForm = fmt.Sprintf("%sE%d", mantissa, expVal)
				} else {
					lit.LexicalForm = sciForm
				}
			} else {
				lit.LexicalForm = fixedForm
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
		predicateOffsets = v.PropertySourceOffsets
	}

	r.statements = append(r.statements, statement{
		quad: rdf.Quad{
			Triple: rdf.Triple{
				Subject:   ectx.ActiveSubject,
				Predicate: ectx.ActiveProperty,
				Object:    lit,
			},
			GraphName: ectx.ActiveGraph,
		},
		textOffsets: r.buildTextOffsets(
			encoding.GraphNameStatementOffsets, ectx.ActiveGraphRange,
			encoding.SubjectStatementOffsets, ectx.ActiveSubjectRange,
			encoding.PredicateStatementOffsets, predicateOffsets,
			encoding.ObjectStatementOffsets, atValuePrimitive.Value.GetSourceOffsets(),
		),
		containerResource: ectx.CurrentContainer,
	})

	return nil
}
