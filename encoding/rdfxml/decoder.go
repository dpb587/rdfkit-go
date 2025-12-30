package rdfxml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectxml-go/inspectxml"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	"github.com/dpb587/rdfkit-go/encoding/rdfxml/internal"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type DecoderOption interface {
	apply(s *DecoderConfig)
	newDecoder(r io.Reader) (*Decoder, error)
}

type statement struct {
	triple            rdf.Triple
	textOffsets       encoding.StatementTextOffsets
	containerResource encoding.ContainerResource
}

type Decoder struct {
	r io.Reader

	baseURL *iriutil.ParsedIRI

	blankNodeStringMapper blanknodeutil.StringMapper

	captureTextOffsets bool
	initialTextOffset  cursorio.TextOffset

	baseDirectiveListener   DecoderEvent_BaseDirective_ListenerFunc
	prefixDirectiveListener DecoderEvent_PrefixDirective_ListenerFunc
	warningListener         func(err error)

	buildTextOffsets encodingutil.TextOffsetsBuilderFunc

	statements    []statement
	statementsIdx int
	err           error

	tokenNext     func() (xml.Token, error)
	tokenMetadata func() (*inspectxml.TokenMetadata, bool)
}

var _ rdf.TripleIterator = &Decoder{}
var _ encoding.StatementTextOffsetsProvider = &Decoder{}

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
	if d.statementsIdx == -1 {
		d.parseAll()
	}

	if d.err != nil {
		return false
	}

	d.statementsIdx++

	return d.statementsIdx < len(d.statements)
}

func (r *Decoder) Triple() rdf.Triple {
	return r.statements[r.statementsIdx].triple
}

func (r *Decoder) Statement() rdf.Statement {
	return r.Triple()
}

func (r *Decoder) StatementTextOffsets() encoding.StatementTextOffsets {
	return r.statements[r.statementsIdx].textOffsets
}

func (d *Decoder) parseAll() {
	if d.captureTextOffsets {
		xmlDecoder := inspectxml.NewDecoder(d.r, inspectxml.DecoderOptions{
			InitialCursor: d.initialTextOffset,
		})

		d.tokenNext = xmlDecoder.Token
		d.tokenMetadata = xmlDecoder.GetTokenMetadata
	} else {
		xmlDecoder := xml.NewDecoder(d.r)

		d.tokenNext = xmlDecoder.Token
	}

	defer func() {
		d.tokenNext = nil
		d.tokenMetadata = nil
	}()

	err := d.decodeRoot(
		evaluationContext{
			Base: d.baseURL,
			Global: &globalEvaluationContext{
				BlankNodeStringMapper: d.blankNodeStringMapper,
				BlankNodeFactory:      blanknodeutil.NewFactory(),
				uniqRefID:             map[uniqRefID]struct{}{},
			},
			CurrentContainer: &DocumentResource{},
			UsedIDs:          map[string]struct{}{},
		},
	)
	if err != nil {
		d.err = err
	}
}

func (d *Decoder) processCommonAttr(ectx evaluationContext, startElement xml.StartElement, tokenMetadata *inspectxml.TokenMetadata) (evaluationContext, []unifiedAttr, []unifiedAttr, error) {
	var rdfAttrList []unifiedAttr
	var otherAttrList []unifiedAttr

	for _, attr := range d.getUnifiedAttributes(startElement, tokenMetadata) {
		switch attr.Name.Space {
		case internal.Space:
			rdfAttrList = append(rdfAttrList, attr)
		case "http://www.w3.org/XML/1998/namespace":
			switch attr.Name.Local {
			case "lang":
				ectx.Language = &attr.Value
			case "base":
				baseIRI := ectx.ResolveIRI(attr.Value)

				// TODO inefficient
				valueIRI, err := iriutil.ParseIRI(string(baseIRI))
				if err != nil {
					return evaluationContext{}, nil, nil, fmt.Errorf("parse base: %w", err)
				}

				ectx.Base = valueIRI
				ectx.UsedIDs = map[string]struct{}{}
			}
		case "xmlns":
			// handled by [xml.Decoder]
			continue
		default:
			if attr.Name == (xml.Name{Space: "", Local: "xmlns"}) {
				// handled by [xml.Decoder]
				continue
			}

			otherAttrList = append(otherAttrList, attr)
		}
	}

	return ectx, rdfAttrList, otherAttrList, nil
}

func (d *Decoder) decodeRoot(ectx evaluationContext) error {
	for {
		token, err := d.tokenNext()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		switch tokenT := token.(type) {
		case xml.ProcInst:
			// should only be first; https://github.com/golang/go/issues/65691
			continue // TODO warn?
		case xml.Directive:
			return d.newTokenError(ErrDirectivesNotSupported, tokenT)
		case xml.CharData:
			continue // expected (whitespace), ignored
		case xml.Comment:
			continue // expected, ignored
		case xml.StartElement:
			var tokenMetadata *inspectxml.TokenMetadata

			if d.tokenMetadata != nil {
				tokenMetadata, _ = d.tokenMetadata()
			}

			for _, attr := range d.getUnifiedAttributes(tokenT, tokenMetadata) {
				if attr.Name.Space == "http://www.w3.org/XML/1998/namespace" && attr.Name.Local == "base" {
					if d.baseDirectiveListener != nil {
						d.baseDirectiveListener(DecoderEvent_BaseDirective_Data{
							Value:        string(ectx.ResolveIRI(attr.Value)),
							ValueOffsets: attr.Metadata.Value,
						})
					}
				} else if attr.Name.Space == "xmlns" {
					if d.prefixDirectiveListener != nil {
						d.prefixDirectiveListener(DecoderEvent_PrefixDirective_Data{
							Prefix:          attr.Name.Local,
							PrefixOffsets:   &attr.Metadata.Name,
							Expanded:        attr.Value,
							ExpandedOffsets: attr.Metadata.Value,
						})
					}
				}
			}

			if tokenT.Name == (xml.Name{Space: internal.Space, Local: internal.Local_RDF_Syntax}) {
				err := d.decodeRDF(ectx, tokenT, tokenMetadata)
				if err != nil {
					return fmt.Errorf("rdf: %w", err)
				}
			} else {
				_, err := d.processNodeElt(ectx, tokenT, tokenMetadata)
				if err != nil {
					return fmt.Errorf("nodeElement: %w", err)
				}
			}
		}
	}
}

func (d *Decoder) decodeRDF(ectx evaluationContext, startElement xml.StartElement, startElementMetadata *inspectxml.TokenMetadata) error {
	ectx, rdfAttrList, _, err := d.processCommonAttr(ectx, startElement, startElementMetadata)
	if err != nil {
		return err
	} else if len(rdfAttrList) != 0 {
		return fmt.Errorf("unexpected attr: %v", rdfAttrList[0].Name)
	}

	for {
		token, err := d.tokenNext()
		if err != nil {
			return err
		}

		switch tokenT := token.(type) {
		case xml.StartElement:
			switch tokenT.Name {
			// nodeElementURIs = anyURI - ( coreSyntaxTerms | rdf:li | oldTerms )
			case
				// coreSyntaxTerms
				xml.Name{Space: internal.Space, Local: internal.Local_RDF_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_ID_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_About_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_ParseType_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_Resource_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_NodeID_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_Datatype_Syntax},
				// explicit
				xml.Name{Space: internal.Space, Local: internal.Local_Li_Syntax},
				// oldTerms
				xml.Name{Space: internal.Space, Local: internal.Local_AboutEach_Old},
				xml.Name{Space: internal.Space, Local: internal.Local_AboutEachPrefix_Old},
				xml.Name{Space: internal.Space, Local: internal.Local_BagID_Old}:
				// not allowed
				return d.newTokenNameError(
					ElementNotAllowedError{
						Name: tokenT.Name,
					},
					tokenT,
				)
			}

			var tokenMetadata *inspectxml.TokenMetadata

			if d.tokenMetadata != nil {
				tokenMetadata, _ = d.tokenMetadata()
			}

			_, err := d.processNodeElt(ectx, tokenT, tokenMetadata)
			if err != nil {
				return fmt.Errorf("description: %w", err)
			}
		case xml.EndElement:
			return nil
		}
	}
}

// ./rdfms-rdf-id/error002.rdf
func (d *Decoder) validateID(v string) error {
	if !reXmlNamespaceName.MatchString(v) || strings.Contains(v, ":") {
		return InvalidNameError{
			Name: v,
		}
	}

	return nil
}

// 7.2.11 Production nodeElement
func (d *Decoder) processNodeElt(ectx evaluationContext, startElement xml.StartElement, startElementMetadata *inspectxml.TokenMetadata) (rdf.SubjectValue, error) {
	ectx, rdfAttrList, otherAttrList, err := d.processCommonAttr(ectx, startElement, startElementMetadata)
	if err != nil {
		return nil, err
	}

	var eSubject rdf.SubjectValue
	var eSubjectLocation *cursorio.TextOffsetRange

	var usedNameAttributes []string

	for _, attr := range rdfAttrList {
		switch attr.Name.Local {
		case internal.Local_ID_Syntax:
			if err := d.validateID(attr.Value); err != nil {
				return nil, d.newTokenAttrError(err, attr)
			}

			eSubject = ectx.ResolveIRI("#" + attr.Value)
			eSubjectLocation = attr.Metadata.Value

			// ./rdfms-difference-between-ID-and-about/error1.rdf
			if _, known := ectx.UsedIDs[attr.Value]; known {
				return nil, d.newTokenAttrError(
					DuplicateScopedNameError{
						Name: string(eSubject.(rdf.IRI)),
					},
					attr,
				)
			}

			ectx.UsedIDs[attr.Value] = struct{}{}

			usedNameAttributes = append(usedNameAttributes, "rdf:ID")
		case internal.Local_NodeID_Syntax:
			// ./rdfms-syntax-incomplete/error001.rdf
			if err := d.validateID(attr.Value); err != nil {
				return nil, d.newTokenAttrError(err, attr)
			}

			eSubject = ectx.Global.BlankNodeStringMapper.MapBlankNodeIdentifier(attr.Value)
			eSubjectLocation = attr.Metadata.Value

			usedNameAttributes = append(usedNameAttributes, "rdf:nodeID")
		case internal.Local_About_Syntax:
			eSubject = ectx.ResolveIRI(attr.Value)
			eSubjectLocation = attr.Metadata.Value

			usedNameAttributes = append(usedNameAttributes, "rdf:about")
		}
	}

	// ./rdfms-syntax-incomplete/error004.rdf
	if len(usedNameAttributes) > 1 {
		return nil, fmt.Errorf("multiple name attributes found: %s", strings.Join(usedNameAttributes, ", "))
	}

	if eSubject == nil {
		eSubject = ectx.Global.BlankNodeFactory.NewBlankNode()
	}

	if startElement.Name != (xml.Name{Space: internal.Space, Local: internal.Local_Description_Syntax}) {
		t := statement{
			triple: rdf.Triple{
				Subject:   eSubject,
				Predicate: rdfiri.Type_Property,
				Object:    rdf.IRI(startElement.Name.Space + startElement.Name.Local),
			},
			containerResource: ectx.CurrentContainer,
		}

		if d.captureTextOffsets {
			t.textOffsets = encoding.StatementTextOffsets{}

			if eSubjectLocation != nil {
				t.textOffsets[encoding.SubjectStatementOffsets] = *eSubjectLocation
			}

			if tokenMetadata, ok := d.tokenMetadata(); ok && tokenMetadata.TagName != nil {
				t.textOffsets[encoding.ObjectStatementOffsets] = *tokenMetadata.TagName
			}
		}

		d.statements = append(d.statements, t)
	}

	for _, attr := range rdfAttrList {
		switch attr.Name.Local {
		case internal.Local_ID_Syntax, internal.Local_NodeID_Syntax, internal.Local_About_Syntax:
			continue // already handled

		case internal.Local_Type_Property:
			t := statement{
				triple: rdf.Triple{
					Subject:   eSubject,
					Predicate: rdfiri.Type_Property,
					Object:    ectx.ResolveIRI(attr.Value),
				},
				containerResource: ectx.CurrentContainer,
			}

			if d.captureTextOffsets {
				t.textOffsets = encoding.StatementTextOffsets{}

				if eSubjectLocation != nil {
					t.textOffsets[encoding.SubjectStatementOffsets] = *eSubjectLocation
				}

				if attr.Metadata != nil {
					t.textOffsets[encoding.PredicateStatementOffsets] = attr.Metadata.Name

					if attr.Metadata.Value != nil {
						t.textOffsets[encoding.ObjectStatementOffsets] = *attr.Metadata.Value
					}
				}
			}

			d.statements = append(d.statements, t)

		// https://www.w3.org/2000/03/rdf-tracking/#rdfms-rdf-names-use
		// reserved vs allowed property names (via default case)
		case internal.Local_RDF_Syntax, internal.Local_Resource_Syntax, internal.Local_BagID_Old, internal.Local_ParseType_Syntax, internal.Local_AboutEach_Old, internal.Local_AboutEachPrefix_Old, internal.Local_Li_Syntax:
			return nil, d.newTokenAttrError(
				AttributeNotAllowedError{
					Name: attr.Name,
				},
				attr,
			)
		default:
			otherAttrList = append(otherAttrList, attr)
		}
	}

	for _, attr := range otherAttrList {
		lit := rdf.Literal{
			Datatype:    xsdiri.String_Datatype,
			LexicalForm: attr.Value,
		}

		if ectx.Language != nil {
			lit.Datatype = rdfiri.LangString_Datatype
			lit.Tag = rdf.LanguageLiteralTag{
				Language: *ectx.Language,
			}
		}

		t := statement{
			triple: rdf.Triple{
				Subject:   eSubject,
				Predicate: rdf.IRI(attr.Name.Space + attr.Name.Local),
				Object:    lit,
			},
			containerResource: ectx.CurrentContainer,
		}

		if d.captureTextOffsets {
			t.textOffsets = encoding.StatementTextOffsets{}

			if eSubjectLocation != nil {
				t.textOffsets[encoding.SubjectStatementOffsets] = *eSubjectLocation
			}

			if attr.Metadata != nil {
				t.textOffsets[encoding.PredicateStatementOffsets] = attr.Metadata.Name

				if attr.Metadata.Value != nil {
					t.textOffsets[encoding.ObjectStatementOffsets] = *attr.Metadata.Value
				}
			}
		}

		d.statements = append(d.statements, t)
	}

	ectx.ParentSubject = eSubject
	ectx.ParentSubjectLocation = eSubjectLocation

	parentContainerIndex := 0
	ectx.ParentContainerIndex = &parentContainerIndex

	err = d.processChildren_PropertyEltList(ectx)
	if err != nil {
		return nil, fmt.Errorf("propertyEltList: %w", err)
	}

	return eSubject, nil
}

// 7.2.13 Production propertyEltList
func (d *Decoder) processChildren_PropertyEltList(ectx evaluationContext) error {
	for {
		token, err := d.tokenNext()
		if err != nil {
			return err
		}

		switch tokenT := token.(type) {
		case xml.StartElement:
			// propertyElementURIs = anyURI - ( coreSyntaxTerms | rdf:Description | oldTerms )
			switch tokenT.Name {
			case
				// coreSyntaxTerms
				xml.Name{Space: internal.Space, Local: internal.Local_RDF_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_ID_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_About_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_ParseType_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_Resource_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_NodeID_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_Datatype_Syntax},
				// explicit
				xml.Name{Space: internal.Space, Local: internal.Local_Description_Syntax},
				// oldTerms
				xml.Name{Space: internal.Space, Local: internal.Local_AboutEach_Old},
				xml.Name{Space: internal.Space, Local: internal.Local_AboutEachPrefix_Old},
				xml.Name{Space: internal.Space, Local: internal.Local_BagID_Old}:
				// not allowed
				return d.newTokenNameError(
					ElementNotAllowedError{
						Name: tokenT.Name,
					},
					tokenT,
				)
			}

			var tokenMetadata *inspectxml.TokenMetadata

			if d.tokenMetadata != nil {
				tokenMetadata, _ = d.tokenMetadata()
			}

			err := d.processPropertyElt(ectx, tokenT, tokenMetadata)
			if err != nil {
				return fmt.Errorf("property: %w", err)
			}
		case xml.EndElement:
			return nil
		}
	}
}

// 7.2.14 Production propertyElt
func (d *Decoder) processPropertyElt(ectx evaluationContext, startElement xml.StartElement, startElementMetadata *inspectxml.TokenMetadata) error {
	var attrParseType string
	var rdfID *string
	var rdfResource *string

	for attrIdx, attr := range startElement.Attr {
		switch attr.Name {
		case xml.Name{Space: internal.Space, Local: internal.Local_ID_Syntax}:
			if err := d.validateID(attr.Value); err != nil {
				if startElementMetadata != nil {
					return d.newTokenAttrError(err, unifiedAttr{
						Attr:     attr,
						Metadata: startElementMetadata.TagAttr[attrIdx],
					})
				}

				return err
			}

			rdfID = &attr.Value
		case xml.Name{Space: internal.Space, Local: internal.Local_Resource_Syntax}:
			rdfResource = &attr.Value
		case xml.Name{Space: internal.Space, Local: internal.Local_ParseType_Syntax}:
			switch attr.Value {
			case "Literal", "Resource", "Collection":
				attrParseType = attr.Value
			default:
				if d.warningListener != nil {
					d.warningListener(AttributeNotAllowedError{
						Name: attr.Name,
					})
				}

				attrParseType = "Literal"
			}
		}
	}

	// ./rdfms-empty-property-elements/error001.rdf
	if attrParseType == "Literal" && rdfResource != nil {
		return errors.New("rdf:resource cannot be used for a literal type")
	}

	if startElement.Name == (xml.Name{Space: internal.Space, Local: internal.Local_Li_Syntax}) {
		*ectx.ParentContainerIndex = *ectx.ParentContainerIndex + 1

		ectx.ParentPredicate = rdf.IRI(fmt.Sprintf("%s_%d", internal.Space, *ectx.ParentContainerIndex))

		if d.tokenMetadata != nil {
			if tokenMetadata, ok := d.tokenMetadata(); ok {
				ectx.ParentPredicateLocation = tokenMetadata.TagName
			} else {
				ectx.ParentPredicateLocation = nil
			}
		} else {
			ectx.ParentPredicateLocation = nil
		}
	} else {
		if startElement.Name.Space == internal.Space && len(startElement.Name.Local) > 2 && startElement.Name.Local[0] == '_' {
			// TODO validate? avoid mixing rdf:li & rdf:_n? avoid duplicates?
		}

		ectx.ParentPredicate = rdf.IRI(startElement.Name.Space + startElement.Name.Local)

		if d.tokenMetadata != nil {
			if tokenMetadata, ok := d.tokenMetadata(); ok {
				ectx.ParentPredicateLocation = tokenMetadata.TagName
			} else {
				ectx.ParentPredicateLocation = nil
			}
		} else {
			ectx.ParentPredicateLocation = nil
		}
	}

	switch attrParseType {
	case "Literal":
		return d.processParseTypeLiteralPropertyElt(ectx, startElement, startElementMetadata)
	case "Resource":
		ectx.CurrentContainer = nil
		return d.processParseTypeResourcePropertyElt(ectx, startElement, startElementMetadata)
	case "Collection":
		ectx.CurrentContainer = nil
		return d.processParseTypeCollectionPropertyElt(ectx, startElement, startElementMetadata)
	}

	var found string
	var foundCharData []byte

	for {
		token, err := d.tokenNext()
		if err != nil {
			return err
		}

		switch tokenT := token.(type) {
		case xml.EndElement:
			if len(found) > 0 {
				return nil
			} else if len(foundCharData) > 0 {
				ectx, rdfAttrList, _, err := d.processCommonAttr(ectx, startElement, startElementMetadata)
				if err != nil {
					return err
				}

				var explicitDatatype bool

				lit := rdf.Literal{
					Datatype:    xsdiri.String_Datatype,
					LexicalForm: string(foundCharData),
				}

				for _, attr := range rdfAttrList {
					switch attr.Name.Local {
					case internal.Local_Datatype_Syntax:
						explicitDatatype = true
						lit.Datatype = ectx.ResolveIRI(attr.Value)
					}
				}

				if !explicitDatatype && ectx.Language != nil {
					lit.Datatype = rdfiri.LangString_Datatype
					lit.Tag = rdf.LanguageLiteralTag{
						Language: *ectx.Language,
					}
				}

				t := statement{
					triple: rdf.Triple{
						Subject:   ectx.ParentSubject,
						Predicate: ectx.ParentPredicate,
						Object:    lit,
					},
					containerResource: ectx.CurrentContainer,
				}

				if d.captureTextOffsets {
					t.textOffsets = encoding.StatementTextOffsets{}

					if ectx.ParentSubjectLocation != nil {
						t.textOffsets[encoding.SubjectStatementOffsets] = *ectx.ParentSubjectLocation
					}

					if ectx.ParentPredicateLocation != nil {
						t.textOffsets[encoding.PredicateStatementOffsets] = *ectx.ParentPredicateLocation
					}

					if endTokenMetadata, ok := d.tokenMetadata(); ok && startElementMetadata != nil {
						t.textOffsets[encoding.ObjectStatementOffsets] = cursorio.TextOffsetRange{
							From:  startElementMetadata.Token.Until,
							Until: endTokenMetadata.Token.From,
						}
					}
				}

				d.statements = append(d.statements, t)

				if rdfID != nil {
					d.addReify(ectx, *rdfID, d.statements[len(d.statements)-1])
				}

				return nil
			}

			ectx, rdfAttrList, otherAttrList, err := d.processCommonAttr(ectx, startElement, startElementMetadata)
			if err != nil {
				return err
			}

			var resourceAttr, nodeIdAttr, datatypeAttr *unifiedAttr

			var usedNameAttributes []string

			for _, attr := range rdfAttrList {
				switch attr.Name.Local {
				case internal.Local_ID_Syntax:
					// already handled
				case internal.Local_Resource_Syntax:
					resourceAttr = &attr

					usedNameAttributes = append(usedNameAttributes, "rdf:resource")
				case internal.Local_NodeID_Syntax:
					// ./rdfms-syntax-incomplete/error003.rdf
					if err := d.validateID(attr.Value); err != nil {
						return d.newTokenAttrError(err, attr)
					}

					nodeIdAttr = &attr

					usedNameAttributes = append(usedNameAttributes, "rdf:nodeID")
				case internal.Local_Datatype_Syntax:
					datatypeAttr = &attr
				default:
					return d.newTokenAttrError(
						AttributeNotAllowedError{
							Name: attr.Name,
						},
						attr,
					)
				}
			}

			if len(usedNameAttributes) > 1 {
				return fmt.Errorf("multiple name attributes found: %s", strings.Join(usedNameAttributes, ", "))
			}

			ot := statement{
				triple: rdf.Triple{
					Subject:   ectx.ParentSubject,
					Predicate: ectx.ParentPredicate,
				},
				containerResource: ectx.CurrentContainer,
			}

			if d.captureTextOffsets {
				ot.textOffsets = encoding.StatementTextOffsets{}

				if ectx.ParentSubjectLocation != nil {
					ot.textOffsets[encoding.SubjectStatementOffsets] = *ectx.ParentSubjectLocation
				}

				if ectx.ParentPredicateLocation != nil {
					ot.textOffsets[encoding.PredicateStatementOffsets] = *ectx.ParentPredicateLocation
				}
			}

			if len(otherAttrList) == 0 && resourceAttr == nil && nodeIdAttr == nil && datatypeAttr == nil {
				ot.triple.Object = rdf.Literal{
					Datatype:    xsdiri.String_Datatype,
					LexicalForm: "",
				}
			} else {
				if resourceAttr != nil {
					ot.triple.Object = ectx.ResolveIRI((*resourceAttr).Value)

					if d.captureTextOffsets && (*resourceAttr).Metadata.Value != nil {
						ot.textOffsets[encoding.ObjectStatementOffsets] = *(*resourceAttr).Metadata.Value
					}
				} else if nodeIdAttr != nil {
					ot.triple.Object = ectx.Global.BlankNodeStringMapper.MapBlankNodeIdentifier((*nodeIdAttr).Value)

					if d.captureTextOffsets && (*nodeIdAttr).Metadata.Value != nil {
						ot.textOffsets[encoding.ObjectStatementOffsets] = *(*nodeIdAttr).Metadata.Value
					}
				} else {
					ot.triple.Object = ectx.Global.BlankNodeFactory.NewBlankNode()
				}

				for _, attr := range append(rdfAttrList, otherAttrList...) {
					switch attr.Name.Space {
					case internal.Space:
						switch attr.Name.Local {
						case internal.Local_ID_Syntax, internal.Local_Resource_Syntax, internal.Local_NodeID_Syntax, internal.Local_Datatype_Syntax:
							// already handled
						case internal.Local_Type_Property:
							t := statement{
								triple: rdf.Triple{
									Subject:   ot.triple.Object.(rdf.SubjectValue), // definitely iri or blank node
									Predicate: rdfiri.Type_Property,
									Object:    ectx.ResolveIRI(attr.Value),
								},
								containerResource: ectx.CurrentContainer,
							}

							if d.captureTextOffsets {
								t.textOffsets = encoding.StatementTextOffsets{}

								if otv, ok := ot.textOffsets[encoding.ObjectStatementOffsets]; ok {
									t.textOffsets[encoding.SubjectStatementOffsets] = otv
								}

								if attr.Metadata != nil {
									t.textOffsets[encoding.PredicateStatementOffsets] = attr.Metadata.Name

									if attr.Metadata.Value != nil {
										t.textOffsets[encoding.ObjectStatementOffsets] = *attr.Metadata.Value
									}
								}
							}

							d.statements = append(d.statements, t)
						// propertyAttributeURIs = anyURI - ( coreSyntaxTerms | rdf:Description | rdf:li | oldTerms )
						case
							// coreSyntaxTerms
							internal.Local_RDF_Syntax,
							// internal.ID_Syntax,
							internal.Local_About_Syntax,
							internal.Local_ParseType_Syntax,
							// internal.Resource_Syntax,
							// internal.NodeID_Syntax,
							// internal.Datatype_Syntax,
							// explicit
							internal.Local_Description_Syntax,
							internal.Local_Li_Syntax,
							// oldTerms
							internal.Local_AboutEach_Old,
							internal.Local_AboutEachPrefix_Old,
							internal.Local_BagID_Old:
							// not allowed
							return d.newTokenAttrError(
								AttributeNotAllowedError{
									Name: attr.Name,
								},
								attr,
							)
						}
					default:
						lit := rdf.Literal{
							Datatype:    xsdiri.String_Datatype,
							LexicalForm: attr.Value,
						}

						if ectx.Language != nil {
							lit.Datatype = rdfiri.LangString_Datatype
							lit.Tag = rdf.LanguageLiteralTag{
								Language: *ectx.Language,
							}
						}

						t := statement{
							triple: rdf.Triple{
								Subject:   ot.triple.Object.(rdf.SubjectValue), // definitely iri or blank node
								Predicate: ectx.ResolveIRI(attr.Name.Space + attr.Name.Local),
								Object:    lit,
							},
							containerResource: ectx.CurrentContainer,
						}

						if d.captureTextOffsets {
							t.textOffsets = encoding.StatementTextOffsets{}

							if otv, ok := ot.textOffsets[encoding.ObjectStatementOffsets]; ok {
								t.textOffsets[encoding.SubjectStatementOffsets] = otv
							}

							if attr.Metadata != nil {
								t.textOffsets[encoding.PredicateStatementOffsets] = attr.Metadata.Name

								if attr.Metadata.Value != nil {
									t.textOffsets[encoding.ObjectStatementOffsets] = *attr.Metadata.Value
								}
							}
						}

						d.statements = append(d.statements, t)
					}
				}
			}

			d.statements = append(d.statements, ot)

			if rdfID != nil {
				d.addReify(ectx, *rdfID, d.statements[len(d.statements)-1])
			}

			return nil
		case xml.StartElement: // resourcePropertyElt
			if len(found) != 0 {
				return fmt.Errorf("already found property value (%s) but found: %v", found, tokenT.Name)
			}

			switch tokenT.Name {
			// nodeElementURIs = anyURI - ( coreSyntaxTerms | rdf:li | oldTerms )
			case
				// coreSyntaxTerms
				xml.Name{Space: internal.Space, Local: internal.Local_RDF_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_ID_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_About_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_ParseType_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_Resource_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_NodeID_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_Datatype_Syntax},
				// explicit
				xml.Name{Space: internal.Space, Local: internal.Local_Li_Syntax},
				// oldTerms
				xml.Name{Space: internal.Space, Local: internal.Local_AboutEach_Old},
				xml.Name{Space: internal.Space, Local: internal.Local_AboutEachPrefix_Old},
				xml.Name{Space: internal.Space, Local: internal.Local_BagID_Old}:
				// not allowed
				return d.newTokenNameError(
					ElementNotAllowedError{
						Name: tokenT.Name,
					},
					tokenT,
				)
			}

			var tokenMetadata *inspectxml.TokenMetadata

			if d.tokenMetadata != nil {
				tokenMetadata, _ = d.tokenMetadata()
			}

			s, err := d.processNodeElt(ectx, tokenT, tokenMetadata)
			if err != nil {
				return fmt.Errorf("nodeElement: %w", err)
			}

			t := statement{
				triple: rdf.Triple{
					Subject:   ectx.ParentSubject,
					Predicate: ectx.ParentPredicate,
					Object:    s,
				},
				containerResource: ectx.CurrentContainer,
			}

			if d.captureTextOffsets {
				t.textOffsets = d.buildTextOffsets(
					encoding.SubjectStatementOffsets, ectx.ParentSubjectLocation,
					encoding.PredicateStatementOffsets, ectx.ParentPredicateLocation,
				)

				// TODO processNodeElt
			}

			d.statements = append(d.statements, t)

			if rdfID != nil {
				d.addReify(ectx, *rdfID, d.statements[len(d.statements)-1])
			}

			found = tokenT.Name.Space + tokenT.Name.Local
		case xml.CharData:
			foundCharData = append(foundCharData, tokenT...)
		}
	}
}

func (d *Decoder) processParseTypeLiteralPropertyElt(ectx evaluationContext, startElement xml.StartElement, startElementMetadata *inspectxml.TokenMetadata) error {
	ectx, rdfAttrList, _, err := d.processCommonAttr(ectx, startElement, startElementMetadata)
	if err != nil {
		return err
	}

	lexicalForm, _, err := d.xmlRender()
	if err != nil {
		return fmt.Errorf("render xml: %w", err)
	}

	t := statement{
		triple: rdf.Triple{
			Subject:   ectx.ParentSubject,
			Predicate: ectx.ParentPredicate,
			Object: rdf.Literal{
				Datatype:    rdfiri.XMLLiteral_Datatype,
				LexicalForm: string(lexicalForm),
			},
		},
		containerResource: ectx.CurrentContainer,
	}

	if d.captureTextOffsets {
		t.textOffsets = d.buildTextOffsets(
			encoding.SubjectStatementOffsets, ectx.ParentSubjectLocation,
			encoding.PredicateStatementOffsets, ectx.ParentPredicateLocation,
		)

		// TODO o
	}

	d.statements = append(d.statements, t)

	for _, attr := range rdfAttrList {
		switch attr.Name.Local {
		case internal.Local_ID_Syntax:
			d.addReify(ectx, attr.Value, d.statements[len(d.statements)-1])
		case internal.Local_ParseType_Syntax:
			// already handled
		default:
			if d.warningListener != nil {
				d.warningListener(AttributeNotAllowedError{
					Name: attr.Name,
				})
			}

			continue
		}
	}

	return nil
}

func (d *Decoder) processParseTypeResourcePropertyElt(ectx evaluationContext, startElement xml.StartElement, startElementMetadata *inspectxml.TokenMetadata) error {
	ectx, rdfAttrList, _, err := d.processCommonAttr(ectx, startElement, startElementMetadata)
	if err != nil {
		return err
	}

	n := ectx.Global.BlankNodeFactory.NewBlankNode()

	t := statement{
		triple: rdf.Triple{
			Subject:   ectx.ParentSubject,
			Predicate: ectx.ParentPredicate,
			Object:    n,
		},
		containerResource: ectx.CurrentContainer,
	}

	if d.captureTextOffsets {
		t.textOffsets = d.buildTextOffsets(
			encoding.SubjectStatementOffsets, ectx.ParentSubjectLocation,
			encoding.PredicateStatementOffsets, ectx.ParentPredicateLocation,
		)
	}

	d.statements = append(d.statements, t)

	for _, attr := range rdfAttrList {
		switch attr.Name.Local {
		case internal.Local_ID_Syntax:
			d.addReify(ectx, attr.Value, d.statements[len(d.statements)-1])
		case internal.Local_ParseType_Syntax:
			// already handled
		default:
			if d.warningListener != nil {
				d.warningListener(AttributeNotAllowedError{
					Name: attr.Name,
				})
			}

			continue
		}
	}

	ectx.ParentSubject = n
	// TODO location
	ectx.ParentPredicate = nil

	parentContainerIndex := 0
	ectx.ParentContainerIndex = &parentContainerIndex

	err = d.processChildren_PropertyEltList(ectx)
	if err != nil {
		return fmt.Errorf("propertyEltList: %w", err)
	}

	return nil
}

func (d *Decoder) processParseTypeCollectionPropertyElt(ectx evaluationContext, startElement xml.StartElement, startElementMetadata *inspectxml.TokenMetadata) error {
	ectx, rdfAttrList, _, err := d.processCommonAttr(ectx, startElement, startElementMetadata)
	if err != nil {
		return err
	}

	var rdfID *string

	for _, attr := range rdfAttrList {
		switch attr.Name.Local {
		case internal.Local_ID_Syntax:
			rdfID = &attr.Value
		default:
			if d.warningListener != nil {
				d.warningListener(AttributeNotAllowedError{
					Name: attr.Name,
				})
			}

			continue
		}
	}

	var lastContainer rdf.SubjectValue

	for {
		token, err := d.tokenNext()
		if err != nil {
			return err
		}

		switch tokenT := token.(type) {
		case xml.StartElement:
			switch tokenT.Name {
			// nodeElementURIs = anyURI - ( coreSyntaxTerms | rdf:li | oldTerms )
			case
				// coreSyntaxTerms
				xml.Name{Space: internal.Space, Local: internal.Local_RDF_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_ID_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_About_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_ParseType_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_Resource_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_NodeID_Syntax},
				xml.Name{Space: internal.Space, Local: internal.Local_Datatype_Syntax},
				// explicit
				xml.Name{Space: internal.Space, Local: internal.Local_Li_Syntax},
				// oldTerms
				xml.Name{Space: internal.Space, Local: internal.Local_AboutEach_Old},
				xml.Name{Space: internal.Space, Local: internal.Local_AboutEachPrefix_Old},
				xml.Name{Space: internal.Space, Local: internal.Local_BagID_Old}:
				// not allowed
				return d.newTokenNameError(
					ElementNotAllowedError{
						Name: tokenT.Name,
					},
					tokenT,
				)
			}

			var tokenMetadata *inspectxml.TokenMetadata

			if d.tokenMetadata != nil {
				tokenMetadata, _ = d.tokenMetadata()
			}

			eSubject, err := d.processNodeElt(ectx, tokenT, tokenMetadata)
			if err != nil {
				return fmt.Errorf("description: %w", err)
			}

			nextContainer := ectx.Global.BlankNodeFactory.NewBlankNode()

			if lastContainer == nil {
				d.statements = append(d.statements,
					statement{
						triple: rdf.Triple{
							Subject:   ectx.ParentSubject,
							Predicate: ectx.ParentPredicate,
							Object:    nextContainer,
						},
						containerResource: ectx.CurrentContainer,
					},
					statement{
						triple: rdf.Triple{
							Subject:   nextContainer,
							Predicate: rdfiri.First_Property,
							Object:    eSubject,
						},
						containerResource: ectx.CurrentContainer,
					},
				)

				if rdfID != nil {
					d.addReify(ectx, *rdfID, d.statements[len(d.statements)-2])
				}
			} else {
				d.statements = append(d.statements,
					statement{
						triple: rdf.Triple{
							Subject:   lastContainer,
							Predicate: rdfiri.Rest_Property,
							Object:    nextContainer,
						},
						containerResource: ectx.CurrentContainer,
					},
					statement{
						triple: rdf.Triple{
							Subject:   nextContainer,
							Predicate: rdfiri.First_Property,
							Object:    eSubject,
						},
						containerResource: ectx.CurrentContainer,
					},
				)
			}

			lastContainer = nextContainer
		case xml.EndElement:
			if lastContainer == nil {
				d.statements = append(d.statements, statement{
					triple: rdf.Triple{
						Subject:   ectx.ParentSubject,
						Predicate: ectx.ParentPredicate,
						Object:    rdfiri.Nil_List,
					},
					containerResource: ectx.CurrentContainer,
				})

				if rdfID != nil {
					d.addReify(ectx, *rdfID, d.statements[len(d.statements)-1])
				}
			} else {
				d.statements = append(d.statements, statement{
					triple: rdf.Triple{
						Subject:   lastContainer,
						Predicate: rdfiri.Rest_Property,
						Object:    rdfiri.Nil_List,
					},
					containerResource: ectx.CurrentContainer,
				})
			}

			return nil
		}
	}
}

func (d *Decoder) addReify(ectx evaluationContext, id string, t statement) {
	idr := ectx.ResolveIRI("#" + id)

	d.statements = append(d.statements,
		statement{
			triple: rdf.Triple{
				Subject:   idr,
				Predicate: rdfiri.Type_Property,
				Object:    rdfiri.Statement_Class,
			},
			containerResource: ectx.CurrentContainer,
		},
		statement{
			triple: rdf.Triple{
				Subject:   idr,
				Predicate: rdfiri.Subject_Property,
				Object:    t.triple.Subject,
			},
			textOffsets: encoding.StatementTextOffsets{
				encoding.ObjectStatementOffsets: t.textOffsets[encoding.SubjectStatementOffsets],
			},
			containerResource: ectx.CurrentContainer,
		},
		statement{
			triple: rdf.Triple{
				Subject:   idr,
				Predicate: rdfiri.Predicate_Property,
				Object:    t.triple.Predicate,
			},
			textOffsets: encoding.StatementTextOffsets{
				encoding.ObjectStatementOffsets: t.textOffsets[encoding.PredicateStatementOffsets],
			},
			containerResource: ectx.CurrentContainer,
		},
		statement{
			triple: rdf.Triple{
				Subject:   idr,
				Predicate: rdfiri.Object_Property,
				Object:    t.triple.Object,
			},
			textOffsets: encoding.StatementTextOffsets{
				encoding.ObjectStatementOffsets: t.textOffsets[encoding.ObjectStatementOffsets],
			},
			containerResource: ectx.CurrentContainer,
		},
	)
}
