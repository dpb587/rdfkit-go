package rdfa

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"regexp"
	"strings"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	encodinghtml "github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/encoding/rdfa/rdfacontent"
	"github.com/dpb587/rdfkit-go/internal/ptr"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/rdfa/rdfairi"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdobject"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil/rdfacontext"
	"github.com/dpb587/rdfkit-go/x/storage/inmemory"
	"github.com/dpb587/rdfkit-go/x/storage/inmemory/simplequery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type DecoderOption interface {
	apply(s *DecoderConfig)
	newDecoder(doc *encodinghtml.Document) (*Decoder, error)
}

type statement struct {
	triple            rdf.Triple
	textOffsets       encoding.StatementTextOffsets
	containerResource encoding.ContainerResource
}

type Decoder struct {
	doc        *encodinghtml.Document
	docBaseURL *iriutil.ParsedIRI

	captureOffsets        bool
	htmlProcessingProfile HtmlProcessingProfile
	defaultVocabulary     *string
	defaultPrefixes       iriutil.PrefixMap
	bnStringFactory       blanknodes.StringFactory
	buildTextOffsets      encodingutil.TextOffsetsBuilderFunc

	err error

	statements    []statement
	statementsIdx int
}

var _ encoding.TriplesDecoder = &Decoder{}
var _ encoding.StatementTextOffsetsProvider = &Decoder{}

func NewDecoder(doc *encodinghtml.Document, opts ...DecoderOption) (*Decoder, error) {
	compiledOpts := DecoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newDecoder(doc)
}

func (v *Decoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return rdfacontent.TypeIdentifier
}

func (v *Decoder) Close() error {
	return nil
}

func (v *Decoder) Err() error {
	return v.err
}

func (v *Decoder) Next() bool {
	if v.err != nil {
		return false
	} else if v.statementsIdx == -1 {
		rootNode := v.doc.GetRoot()

		gectx := &globalEvaluationContext{
			HostDefaultVocabulary:  ptr.Value("http://www.w3.org/1999/xhtml/vocab#"),
			HostDefaultPrefixes:    v.defaultPrefixes,
			HtmlProcessing:         v.htmlProcessingProfile,
			BlankNodeStringFactory: v.bnStringFactory,
		}

		if gectx.HostDefaultPrefixes == nil {
			gectx.HostDefaultPrefixes = iriutil.NewPrefixMap(rdfacontext.WidelyUsedInitialContext()...)
		}

		ectx := evaluationContext{
			Global:            gectx,
			BaseURL:           v.docBaseURL,
			ListMapping:       map[rdf.IRI]*listMappingBuilder{},
			TermMappings:      map[string]rdf.IRI{},
			DefaultVocabulary: gectx.HostDefaultVocabulary,
			PrefixMapping:     gectx.HostDefaultPrefixes,
			CurrentContainer:  &DocumentResource{},
		}

		if v.defaultVocabulary != nil {
			ectx.DefaultVocabulary = v.defaultVocabulary
		}

		err := v.walkNode(ectx, rootNode)
		if err != nil {
			v.err = err

			return false
		}
	}

	v.statementsIdx++

	return v.statementsIdx < len(v.statements)
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

func (v *Decoder) walkNode(ectx evaluationContext, n *html.Node) error {
	// TODO choose only one root condition
	isRootElement := n.Parent == nil || (n.Parent.Type == html.DocumentNode && n.Type == html.ElementNode)

	if n.Type == html.DoctypeNode && ectx.Global.HtmlProcessing == UnspecifiedHtmlProcessingProfile {
		// rdfa-in-html // 3.1 // Additional Processing Rule 5

		// TODO log auto-detect behavior

		if strings.Contains(n.Data, "//DTD XHTML+RDFa 1.0//") && strings.Contains(n.Data, `"http://www.w3.org/MarkUp/DTD/xhtml-rdfa-1.dtd"`) {
			// [dpb] not currently testing against 1.0
			// ectx.Global.HtmlProcessing = XHTML1_RDFa10_HtmlProcessProfile
			ectx.Global.HtmlProcessing = XHTML1_RDFa11_HtmlProcessProfile
		} else if strings.Contains(n.Data, "//DTD XHTML+RDFa 1.1//") && strings.Contains(n.Data, `"http://www.w3.org/MarkUp/DTD/xhtml-rdfa-2.dtd"`) {
			ectx.Global.HtmlProcessing = XHTML1_RDFa11_HtmlProcessProfile
		} else {
			ectx.Global.HtmlProcessing = XHTML5_RDFa11_HtmlProcessProfile
		}
	} else if n.DataAtom == atom.Html {
		if ectx.Global.HtmlProcessing == UnspecifiedHtmlProcessingProfile {
			ectx.Global.HtmlProcessing = XHTML5_RDFa11_HtmlProcessProfile

			for _, attr := range n.Attr {
				if attr.Namespace == "" && attr.Key == "version" {
					// TODO not sure where these are defined in spec
					switch attr.Val {
					case "HTML+RDFa 1.0", "XHTML+RDFa 1.0":
						// [dpb] not currently testing against 1.0
						// ectx.Global.HtmlProcessing = XHTML1_RDFa10_HtmlProcessProfile
						ectx.Global.HtmlProcessing = XHTML1_RDFa11_HtmlProcessProfile
					case "HTML+RDFa 1.1", "XHTML+RDFa 1.1":
						ectx.Global.HtmlProcessing = XHTML1_RDFa11_HtmlProcessProfile
					default:
						// TODO warning
					}

					break
				}
			}
		}

		// only once we encounter the html tag and presumably all declarations of the processing profile have occurred
		// do we configure the initial contexts; this helps avoid profile configuration in multiple locations
		if ectx.Global.HtmlProcessing&ActiveHtmlProcessingProfile > 0 {
			// rdfa-in-html // 3.1 // Additional Processing Rule 1
			ectx.DefaultVocabulary = ectx.Global.HostDefaultVocabulary

			// rdfa-in-html // 3.1 // Additional Processing Rule 2
			ectx.PrefixMapping = ectx.PrefixMapping.NewPrefixMap(InitialContext.AsPrefixMappingList()...)
			maps.Copy(ectx.TermMappings, w3_2011_rdfacontext_rdfa11_TermMappings)

			// xhtml-rdfa // Additional RDFa Processing Rules
			if ectx.Global.HtmlProcessing == XHTML1_RDFa11_HtmlProcessProfile {
				// [xhtml-rdfa] The default vocabulary IRI is undefined.
				// [dpb] not actually true? rdfa-info-test-suite still expects rdfa-in-html default
				// ectx.DefaultVocabulary = nil

				// [xhtml-rdfa] XHTML+RDFa uses an additional initial context by default, http://www.w3.org/2011/rdfa-context/xhtml-rdfa-1.1, which must be applied after the initial context for [RDFA-CORE] (http://www.w3.org/2011/rdfa-context/rdfa-1.1).
				maps.Copy(ectx.TermMappings, w3_2011_rdfacontext_xhtmlrdfa11_TermMappings)
			}

			// custom behavior not described in spec
			if v.defaultVocabulary != nil {
				ectx.DefaultVocabulary = v.defaultVocabulary
			}
		}
	}

	if n.DataAtom == atom.Base {
		// rdfa-in-html // 3.1 // Additional Processing Rule 3

		if ectx.Global.HtmlFoundBase {
			// TODO warn; per-spec, only the first base tag is respected
		} else {
			var baseHref *string

			for _, attr := range n.Attr {
				if attr.Namespace == "" && attr.Key == "href" {
					baseHref = &attr.Val

					break
				}
			}

			if baseHref == nil {
				// TODO warn; per-spec, href is a required attribute
			} else {
				// duplicate behavior from newDecoder where base is pulled from parsed html docInfo?
				baseURL, err := iriutil.ParseIRI(*baseHref)
				if err != nil {
					// TODO warn; invalid data
				} else {
					// test[0117.html] fragment must be dropped; not documented in specs?
					baseURL.DropFragment()

					ectx.BaseURL = baseURL
					ectx.Global.HtmlFoundBase = true
				}
			}
		}
	}

	nodeProfile, _ := v.doc.GetNodeMetadata(n)

	// rdfa-core // 7.5 // Processing Rule 1

	var skipElement bool
	var newSubject rdf.SubjectValue
	var currentObjectResource rdf.ObjectValue
	var typedResource rdf.ObjectValue
	var localPrefixMappings = ectx.PrefixMapping
	var localIncompleteTriples = []incompleteTriple{}
	var listMapping = ectx.ListMapping
	var currentLanguage = ectx.Language
	// var localTermMappings = ectx.TermMappings // TODO use this variable instead
	var localDefaultVocabulary = ectx.DefaultVocabulary

	// rdfa-in-html (not described for processing)
	var localBaseURL = ectx.BaseURL

	// anno features
	var newSubjectAnno *cursorio.TextOffsetRange
	var currentObjectResourceAnno *cursorio.TextOffsetRange
	var typedResourceAnno *cursorio.TextOffsetRange

	//

	var attrPrefix string
	var attrAbout, attrContent, attrDatatype, attrDatetime, attrHref, attrInlist, attrLang, attrLangXml, attrProperty, attrRel, attrResource, attrRev, attrSrc, attrTypeof, attrVocab *string
	var attrAboutIdx, attrContentIdx, attrDatetimeIdx, attrHrefIdx, attrPropertyIdx, attrRelIdx, attrResourceIdx, attrRevIdx, attrSrcIdx, attrTypeofIdx int
	var attrPrefixEntries []iriutil.PrefixMapping

	for attrIdx, attr := range n.Attr {
		switch attr.Namespace {
		case "":
			switch attr.Key {
			case "about":
				attrAbout = &attr.Val
				attrAboutIdx = attrIdx
			case "content":
				attrContent = &attr.Val
				attrContentIdx = attrIdx
			case "datetime":
				attrDatetime = &attr.Val // rdfa-in-html // 3.1 // Additional Processing Rule 9
				attrDatetimeIdx = attrIdx
			case "datatype":
				attrDatatype = &attr.Val
			case "href":
				attrHref = &attr.Val
				attrHrefIdx = attrIdx
			case "inlist":
				attrInlist = &attr.Val
			case "lang":
				attrLang = &attr.Val
			case "prefix":
				attrPrefix = attr.Val
			case "property":
				attrProperty = &attr.Val
				attrPropertyIdx = attrIdx
			case "rel":
				attrRel = &attr.Val
				attrRelIdx = attrIdx
			case "resource":
				attrResource = &attr.Val
				attrResourceIdx = attrIdx
			case "rev":
				attrRev = &attr.Val
				attrRevIdx = attrIdx
			case "src":
				attrSrc = &attr.Val
				attrSrcIdx = attrIdx
			case "typeof":
				attrTypeof = &attr.Val
				attrTypeofIdx = attrIdx
			case "vocab":
				attrVocab = &attr.Val
			case "xml:base":
				if ectx.Global.HtmlProcessing.IsXHTML() {
					// rdfa-in-html // 3.1 // Additional Processing Rule 3

					baseURL, err := iriutil.ParseIRI(attr.Val)
					if err != nil {
						// TODO warning
					} else if localBaseURL != nil {
						// technically, only one html base tag is allowed, but xml allows nested resolves
						// we don't differentiate html/xml base url, so an xml tag might end up resolving against html
						// ambiguous or unexpected behavior?
						localBaseURL = localBaseURL.ResolveReference(baseURL)
					} else {
						localBaseURL = baseURL
					}
				}
			case "xml:lang":
				if ectx.Global.HtmlProcessing.IsXHTML() {
					// rdfa-in-html // 3.1 // Additional Processing Rule 4
					attrLangXml = &attr.Val
				}
			default:
				if strings.HasPrefix(attr.Key, "xmlns:") {
					if !ectx.Global.DisableBackcompatXmlnsPrefixes {
						prefixName := strings.ToLower(attr.Key[6:])

						if len(prefixName) == 0 {
							// [dpb] would be weird; ignore for simpler code; TODO verify if actually invalid
							// [rdfa-info-test-suite rdfa1.1/xhtml1/manifest#0180] dislikes empty @prefix values anyway
							continue
						} else if prefixName[0] == '_' {
							// [dpb] not explicitly mentioned in spec, but would introduce ambiguity as blank nodes
							// [rdfa-info-test-suite rdfa1.1/xhtml1/manifest#0258]
							// TODO warn
							continue
						}

						attrPrefixEntries = append(attrPrefixEntries, iriutil.PrefixMapping{
							Prefix:   prefixName,
							Expanded: rdf.IRI(attr.Val),
						})
					}
				}
			}
		}
	}

	{
		// rdfa-core // 7.5 // Processing Rule 2

		if attrVocab != nil {
			if len(*attrVocab) == 0 {
				localDefaultVocabulary = ectx.Global.HostDefaultVocabulary
			} else {
				vocabIRI, ok := resolveIRI(ectx, localPrefixMappings, *attrVocab, nil, localDefaultVocabulary, false, true).(rdf.IRI)
				if !ok {
					// TODO warning
				} else {
					localDefaultVocabulary = ptr.Value(string(vocabIRI))

					var objectOffsets *cursorio.TextOffsetRange

					if v.captureOffsets {
						if nodeProfile.EndTagTokenOffsets != nil {
							typedResourceAnno = &cursorio.TextOffsetRange{
								From:  nodeProfile.TokenOffsets.From,
								Until: nodeProfile.EndTagTokenOffsets.Until,
							}
						} else {
							typedResourceAnno = &nodeProfile.TokenOffsets
						}
					}

					v.statements = append(v.statements, statement{
						triple: rdf.Triple{
							Subject:   rdf.IRI(ectx.BaseURL.String()),
							Predicate: rdfairi.UsesVocabulary_Property,
							Object:    vocabIRI,
						},
						textOffsets: v.buildTextOffsets(
							encoding.ObjectStatementOffsets, objectOffsets,
						),
					})
				}
			}
		}
	}

	{
		// rdfa-core // 7.5 // Processing Rule 3

		if len(attrPrefix) > 0 {
			fields := strings.Fields(strings.TrimSpace(attrPrefix))

			if len(fields)%2 != 0 {
				// TODO warn
			} else {
				for fieldIdx := 0; fieldIdx < len(fields); fieldIdx += 2 {
					prefixTerm := strings.ToLower(fields[fieldIdx])

					if !strings.HasSuffix(prefixTerm, ":") {
						// TODO warn
						continue
					} else if prefixTerm[0] == '_' {
						// [dpb] not explicitly mentioned in spec, but would introduce ambiguity as blank nodes
						// [rdfa-info-test-suite rdfa1.1/xhtml1/manifest#0258]
						// TODO warn
						continue
					} else if len(prefixTerm) == 1 {
						// [rdfa-info-test-suite rdfa1.1/xhtml1/manifest#0180] suggests empty prefixes are not allowed
						// [dpb] #0180 seems incorrect anyway since the @property reference should resolve against default vocabulary?
						// [dpb] not based in spec? but still seems good to ignore as causing surprising behavior vs @vocab
						continue
					}

					// per spec, @prefix will be processed after xmlns (which added to attrPrefixEntries earlier)
					attrPrefixEntries = append(attrPrefixEntries, iriutil.PrefixMapping{
						Prefix:   prefixTerm[:len(prefixTerm)-1],
						Expanded: rdf.IRI(fields[fieldIdx+1]),
					})
				}
			}
		}

		if len(attrPrefixEntries) > 0 {
			if ectx.Global.HtmlProcessing&ActiveHtmlProcessingProfile > 0 {
				// rdfa-in-html // 3.1 // Additional Processing Rule 6

				for _, prefixEntry := range attrPrefixEntries {
					if expanded, known := localPrefixMappings.ExpandPrefix(prefixEntry.Prefix, ""); !known || expanded != rdf.IRI(prefixEntry.Expanded) {
						// TODO emit rdfa:PrefixRedefinition
					}
				}
			}

			localPrefixMappings = localPrefixMappings.NewPrefixMap(attrPrefixEntries...)
		}
	}

	{
		// rdfa-core // 7.5 // Processing Rule 4

		if attrLangXml != nil && ectx.Global.HtmlProcessing&ActiveHtmlProcessingProfile > 0 {
			// rdfa-in-html // 3.1 // Additional Processing Rule 4

			if len(*attrLangXml) == 0 {
				currentLanguage = nil
			} else {
				currentLanguage = attrLangXml

				if attrLang != nil && *attrLangXml != *attrLang {
					// TODO warning; spec says must be equal
				}
			}
		} else if attrLang != nil {
			if len(*attrLang) == 0 {
				currentLanguage = nil
			} else {
				currentLanguage = attrLang
			}
		}
	}

	{
		// rdfa-in-html // 3.1 // Additional Processing Rule 7

		if attrProperty != nil {
			filterNonCURIENonURI := func(value string) string {
				var resolvedList []string

				for _, r := range strings.Fields(strings.TrimSpace(value)) {
					resolved, ok := resolveIRI(ectx, localPrefixMappings, r, nil, nil, false, false).(rdf.IRI)
					if ok {
						resolvedList = append(resolvedList, "["+string(resolved)+"]")
					}
				}

				return strings.Join(resolvedList, " ")
			}

			if attrRel != nil {
				if resolved := filterNonCURIENonURI(*attrRel); len(resolved) > 0 {
					attrRel = &resolved
				} else {
					attrRel = nil
				}
			}

			if attrRev != nil {
				if resolved := filterNonCURIENonURI(*attrRev); len(resolved) > 0 {
					attrRev = &resolved
				} else {
					attrRev = nil
				}
			}
		}
	}

	{
		// rdfa-core // 7.5 // Processing Rule 5

		if attrRel == nil && attrRev == nil {
			if attrProperty != nil && attrContent == nil && attrDatatype == nil {
				// rdfa-core // 7.5 // Processing Rule 5, Option 1

				var attrAboutResolved rdf.SubjectValue

				if attrAbout != nil {
					attrAboutResolved = resolveIRI(ectx, localPrefixMappings, *attrAbout, localBaseURL, localDefaultVocabulary, true, true)
				}

				if attrAboutResolved != nil {
					newSubject = attrAboutResolved
					newSubjectAnno = nil

					if v.captureOffsets {
						if attrProfile := nodeProfile.TagAttr[attrAboutIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
							newSubjectAnno = attrProfile.ValueOffsets
						}
					}
				} else if ectx.Global.HtmlProcessing&ActiveHtmlProcessingProfile > 0 && (n.DataAtom == atom.Head || n.DataAtom == atom.Body) {
					// rdfa-in-html // 3.1 // Additional Processing Rule 8
					newSubject = ectx.ParentObject.(rdf.SubjectValue)
					newSubjectAnno = ectx.ParentObjectAnno
				} else if isRootElement {
					if s := resolveIRI(ectx, localPrefixMappings, "", localBaseURL, localDefaultVocabulary, false, true); s != nil {
						newSubject = s
						newSubjectAnno = nil
					}
				} else if ectx.ParentObject != nil {
					newSubject = ectx.ParentObject.(rdf.SubjectValue)
					newSubjectAnno = ectx.ParentObjectAnno
				}

				if attrTypeof != nil {
					if attrAbout != nil && attrAboutResolved != nil {
						typedResource = newSubject
						typedResourceAnno = newSubjectAnno
					} else if isRootElement {
						typedResource = resolveIRI(ectx, localPrefixMappings, "", localBaseURL, localDefaultVocabulary, false, true)
						typedResourceAnno = nil
					} else {
						if attrResource != nil {
							typedResource = resolveIRI(ectx, localPrefixMappings, *attrResource, localBaseURL, localDefaultVocabulary, true, true)
							typedResourceAnno = nil

							if v.captureOffsets {
								if attrProfile := nodeProfile.TagAttr[attrResourceIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
									typedResourceAnno = attrProfile.ValueOffsets
								}
							}
						} else if attrHref != nil {
							typedResource = resolveIRI(ectx, localPrefixMappings, *attrHref, localBaseURL, localDefaultVocabulary, false, true)
							typedResourceAnno = nil

							if v.captureOffsets {
								if attrProfile := nodeProfile.TagAttr[attrHrefIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
									typedResourceAnno = attrProfile.ValueOffsets
								}
							}
						} else if attrSrc != nil {
							typedResource = resolveIRI(ectx, localPrefixMappings, *attrSrc, localBaseURL, localDefaultVocabulary, false, true)
							typedResourceAnno = nil

							if v.captureOffsets {
								if attrProfile := nodeProfile.TagAttr[attrSrcIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
									typedResourceAnno = attrProfile.ValueOffsets
								}
							}
						} else {
							typedResource = ectx.Global.BlankNodeStringFactory.NewBlankNode()
							typedResourceAnno = nil

							if v.captureOffsets {
								if nodeProfile.EndTagTokenOffsets != nil {
									typedResourceAnno = &cursorio.TextOffsetRange{
										From:  nodeProfile.TokenOffsets.From,
										Until: nodeProfile.EndTagTokenOffsets.Until,
									}
								} else {
									typedResourceAnno = &nodeProfile.TokenOffsets
								}
							}
						}

						currentObjectResource = typedResource
						currentObjectResourceAnno = typedResourceAnno
					}
				}
			} else {
				// rdfa-core // 7.5 // Processing Rule 5, Option 2

				if attrAbout != nil {
					if s := resolveIRI(ectx, localPrefixMappings, *attrAbout, localBaseURL, localDefaultVocabulary, true, true); s != nil {
						newSubject = s
						newSubjectAnno = nil

						if v.captureOffsets {
							if attrProfile := nodeProfile.TagAttr[attrAboutIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
								newSubjectAnno = attrProfile.ValueOffsets
							}
						}
					}
				}

				if newSubject == nil && attrResource != nil {
					if s := resolveIRI(ectx, localPrefixMappings, *attrResource, localBaseURL, localDefaultVocabulary, true, true); s != nil {
						newSubject = s
						newSubjectAnno = nil

						if v.captureOffsets {
							if attrProfile := nodeProfile.TagAttr[attrResourceIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
								newSubjectAnno = attrProfile.ValueOffsets
							}
						}
					}
				}

				if newSubject == nil && attrHref != nil {
					if s := resolveIRI(ectx, localPrefixMappings, *attrHref, localBaseURL, localDefaultVocabulary, false, true); s != nil {
						newSubject = s
						newSubjectAnno = nil

						if v.captureOffsets {
							if attrProfile := nodeProfile.TagAttr[attrHrefIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
								newSubjectAnno = attrProfile.ValueOffsets
							}
						}
					}
				}

				if newSubject == nil && attrSrc != nil {
					if s := resolveIRI(ectx, localPrefixMappings, *attrSrc, localBaseURL, localDefaultVocabulary, false, true); s != nil {
						newSubject = s
						newSubjectAnno = nil

						if v.captureOffsets {
							if attrProfile := nodeProfile.TagAttr[attrSrcIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
								newSubjectAnno = attrProfile.ValueOffsets
							}
						}
					}
				}

				if newSubject != nil {
					// done
				} else if ectx.Global.HtmlProcessing&ActiveHtmlProcessingProfile > 0 && (n.DataAtom == atom.Head || n.DataAtom == atom.Body) {
					// rdfa-in-html // 3.1 // Additional Processing Rule 8
					newSubject = ectx.ParentObject.(rdf.SubjectValue)
					newSubjectAnno = ectx.ParentObjectAnno
				} else {
					if isRootElement {
						if s := resolveIRI(ectx, localPrefixMappings, "", localBaseURL, localDefaultVocabulary, false, true); s != nil {
							newSubject = s
							newSubjectAnno = nil
						}
					} else if attrTypeof != nil {
						newSubject = ectx.Global.BlankNodeStringFactory.NewBlankNode()
						newSubjectAnno = nil

						if v.captureOffsets {
							if nodeProfile.EndTagTokenOffsets != nil {
								newSubjectAnno = &cursorio.TextOffsetRange{
									From:  nodeProfile.TokenOffsets.From,
									Until: nodeProfile.EndTagTokenOffsets.Until,
								}
							} else {
								newSubjectAnno = &nodeProfile.TokenOffsets
							}
						}
					} else if ectx.ParentObject != nil {
						newSubject = ectx.ParentObject.(rdf.SubjectValue)
						newSubjectAnno = ectx.ParentObjectAnno

						skipElement = attrProperty == nil
					}
				}

				if attrTypeof != nil {
					typedResource = newSubject
					typedResourceAnno = newSubjectAnno
				}
			}
		}
	}

	{
		// rdfa-core // 7.5 // Processing Rule 6

		if attrRel != nil || attrRev != nil {
			if attrAbout != nil {
				if s := resolveIRI(ectx, localPrefixMappings, *attrAbout, localBaseURL, localDefaultVocabulary, true, true); s != nil {
					newSubject = s
					newSubjectAnno = nil

					if v.captureOffsets {
						if attrProfile := nodeProfile.TagAttr[attrAboutIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
							newSubjectAnno = attrProfile.ValueOffsets
						}
					}

					if attrTypeof != nil {
						typedResource = newSubject
						typedResourceAnno = newSubjectAnno
					}
				}
			} else {
				if isRootElement {
					if s := resolveIRI(ectx, localPrefixMappings, "", localBaseURL, localDefaultVocabulary, false, true); s != nil {
						newSubject = s
						newSubjectAnno = nil
					}
				} else if ectx.ParentObject != nil {
					newSubject = ectx.ParentObject.(rdf.SubjectValue)
					newSubjectAnno = ectx.ParentObjectAnno
				}
			}

			if attrResource != nil {
				if s := resolveIRI(ectx, localPrefixMappings, *attrResource, localBaseURL, localDefaultVocabulary, true, true); s != nil {
					currentObjectResource = s
					currentObjectResourceAnno = nil

					if v.captureOffsets {
						if attrProfile := nodeProfile.TagAttr[attrResourceIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
							currentObjectResourceAnno = attrProfile.ValueOffsets
						}
					}
				}
			}

			if currentObjectResource == nil && attrHref != nil {
				if s := resolveIRI(ectx, localPrefixMappings, *attrHref, localBaseURL, localDefaultVocabulary, false, true); s != nil {
					currentObjectResource = s
					currentObjectResourceAnno = nil

					if v.captureOffsets {
						if attrProfile := nodeProfile.TagAttr[attrHrefIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
							currentObjectResourceAnno = attrProfile.ValueOffsets
						}
					}
				}
			}

			if currentObjectResource == nil && attrSrc != nil {
				if s := resolveIRI(ectx, localPrefixMappings, *attrSrc, localBaseURL, localDefaultVocabulary, false, true); s != nil {
					currentObjectResource = s
					currentObjectResourceAnno = nil

					if v.captureOffsets {
						if attrProfile := nodeProfile.TagAttr[attrSrcIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
							currentObjectResourceAnno = attrProfile.ValueOffsets
						}
					}
				}
			}

			if currentObjectResource == nil && attrTypeof != nil && attrAbout == nil {
				currentObjectResource = ectx.Global.BlankNodeStringFactory.NewBlankNode()
				currentObjectResourceAnno = nil
			}

			if attrTypeof != nil && attrAbout == nil && typedResource == nil {
				typedResource = currentObjectResource
				typedResourceAnno = currentObjectResourceAnno
			}
		}
	}

	{
		// rdfa-core // 7.5 // Processing Rule 7

		if typedResource != nil && attrTypeof != nil {
			var attrValOffset int
			var attrVal = *attrTypeof

			var predicateRange *cursorio.TextOffsetRange

			if v.captureOffsets {
				if attrProfile := nodeProfile.TagAttr[attrTypeofIdx]; attrProfile != nil {
					predicateRange = &attrProfile.KeyOffsets
				}
			}

			for len(attrVal) > 0 {
				if mm := fieldsNextSpace.FindString(attrVal); len(mm) > 0 {
					attrValOffset += len(mm)
					attrVal = attrVal[len(mm):]

					continue
				}

				fieldLexical := fieldsNextNonSpace.FindString(attrVal)

				fieldIRI, ok := resolveIRI(ectx, localPrefixMappings, fieldLexical, nil, localDefaultVocabulary, false, true).(rdf.IRI)
				if !ok {
					// TODO warning
				} else {
					var anno *cursorio.TextOffsetRange

					if v.captureOffsets {
						if attrProfile := nodeProfile.TagAttr[attrTypeofIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
							// TODO resurrect offset writer

							anno = attrProfile.ValueOffsets
						}
					}

					v.statements = append(v.statements, statement{
						triple: rdf.Triple{
							Subject:   typedResource.(rdf.SubjectValue),
							Predicate: rdfiri.Type_Property,
							Object:    fieldIRI,
						},
						textOffsets: v.buildTextOffsets(
							encoding.SubjectStatementOffsets, typedResourceAnno,
							encoding.PredicateStatementOffsets, predicateRange,
							encoding.ObjectStatementOffsets, anno,
						),
						containerResource: ectx.CurrentContainer,
					})
				}

				attrValOffset += len(fieldLexical)
				attrVal = attrVal[len(fieldLexical):]
			}
		}
	}

	{
		// rdfa-core // 7.5 // Processing Rule 8
		// [dpb] spec says compare newSubject with ParentObject, however this seems incorrect

		if newSubject != nil && (ectx.ParentSubject == nil || !ectx.ParentSubject.TermEquals(newSubject)) {
			listMapping = map[rdf.IRI]*listMappingBuilder{}
		}
	}

	// some @rel attributes apparently are ignored
	var attrRelValid bool

	if currentObjectResource != nil {
		// rdfa-core // 7.5 // Processing Rule 9

		if attrInlist != nil && attrRel != nil {
			for _, relField := range strings.Fields(strings.TrimSpace(*attrRel)) {
				if ectx.Global.HtmlProcessing&ActiveHtmlProcessingProfile > 0 {
					if ectx.Global.HtmlProcessing.IsXHTML1() {
						// not suppressing (per rdfa-info-test-suite cases)
					} else {
						// custom behavior, not found in spec yet?

						switch n.DataAtom {
						case atom.Form, atom.A, atom.Area, atom.Link:
							if _, known := htmlIgnoredLinkRels[fmt.Sprintf("%s/%s", strings.ToLower(relField), n.DataAtom.String())]; known {
								continue
							}
						}
					}
				}

				relValue, ok := resolveIRI(ectx, localPrefixMappings, relField, nil, localDefaultVocabulary, false, true).(rdf.IRI)
				if !ok {
					// TODO warning
					continue
				}

				if _, known := listMapping[relValue]; !known {
					listMapping[relValue] = &listMappingBuilder{}
				}

				listMapping[relValue].Objects = append(listMapping[relValue].Objects, currentObjectResource)
			}
		}

		{
			if attrRel != nil && attrInlist == nil {
				var predicateRange *cursorio.TextOffsetRange

				if v.captureOffsets {
					if attrProfile := nodeProfile.TagAttr[attrRelIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
						predicateRange = attrProfile.ValueOffsets
					}
				}

				for _, relField := range strings.Fields(strings.TrimSpace(*attrRel)) {
					if ectx.Global.HtmlProcessing&ActiveHtmlProcessingProfile > 0 {
						if ectx.Global.HtmlProcessing.IsXHTML1() {
							// not suppressing (per rdfa-info-test-suite cases)
						} else {
							// custom behavior, not found in spec yet?

							switch n.DataAtom {
							case atom.Form, atom.A, atom.Area, atom.Link:
								if _, known := htmlIgnoredLinkRels[fmt.Sprintf("%s/%s", strings.ToLower(relField), n.DataAtom.String())]; known {
									continue
								}
							}
						}
					}

					relValue, ok := resolveIRI(ectx, localPrefixMappings, relField, nil, localDefaultVocabulary, false, true).(rdf.IRI)
					if !ok {
						// TODO warning
						continue
					}

					v.statements = append(v.statements, statement{
						triple: rdf.Triple{
							Subject:   newSubject,
							Predicate: relValue,
							Object:    currentObjectResource,
						},
						textOffsets: v.buildTextOffsets(
							encoding.SubjectStatementOffsets, newSubjectAnno,
							encoding.PredicateStatementOffsets, predicateRange,
							encoding.ObjectStatementOffsets, currentObjectResourceAnno,
						),
						containerResource: ectx.CurrentContainer,
					})

					attrRelValid = true
				}
			}

			if attrRev != nil {
				var predicateRange *cursorio.TextOffsetRange

				if v.captureOffsets {
					if attrProfile := nodeProfile.TagAttr[attrRevIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
						predicateRange = attrProfile.ValueOffsets
					}
				}

				for _, revField := range strings.Fields(strings.TrimSpace(*attrRev)) {
					revValue, ok := resolveIRI(ectx, localPrefixMappings, revField, nil, localDefaultVocabulary, false, true).(rdf.IRI)
					if !ok {
						// TODO warning
						continue
					}

					v.statements = append(v.statements, statement{
						triple: rdf.Triple{
							Subject:   currentObjectResource.(rdf.SubjectValue),
							Predicate: revValue,
							Object:    newSubject,
						},
						textOffsets: v.buildTextOffsets(
							encoding.SubjectStatementOffsets, currentObjectResourceAnno,
							encoding.PredicateStatementOffsets, predicateRange,
							encoding.ObjectStatementOffsets, newSubjectAnno,
						),
						containerResource: ectx.CurrentContainer,
					})
				}
			}
		}
	} else if attrRel != nil || attrRev != nil {
		// Processing, Step 10

		currentObjectResource = ectx.Global.BlankNodeStringFactory.NewBlankNode()
		currentObjectResourceAnno = nil

		if v.captureOffsets {
			if attrRel != nil {
				if attrProfile := nodeProfile.TagAttr[attrRelIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
					currentObjectResourceAnno = attrProfile.ValueOffsets
				}
			} else {
				if attrProfile := nodeProfile.TagAttr[attrRevIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
					currentObjectResourceAnno = attrProfile.ValueOffsets
				}
			}
		}

		if attrRel != nil {
			var predicateRange *cursorio.TextOffsetRange

			if v.captureOffsets {
				if attrProfile := nodeProfile.TagAttr[attrRelIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
					predicateRange = attrProfile.ValueOffsets
				}
			}

			for _, relField := range strings.Fields(strings.TrimSpace(*attrRel)) {
				if ectx.Global.HtmlProcessing&ActiveHtmlProcessingProfile > 0 {
					if ectx.Global.HtmlProcessing.IsXHTML1() {
						// not suppressing (per rdfa-info-test-suite cases)
					} else {
						// custom behavior, not found in spec yet?

						switch n.DataAtom {
						case atom.Form, atom.A, atom.Area, atom.Link:
							if _, known := htmlIgnoredLinkRels[fmt.Sprintf("%s/%s", strings.ToLower(relField), n.DataAtom.String())]; known {
								continue
							}
						}
					}
				}

				relValue, ok := resolveIRI(ectx, localPrefixMappings, relField, nil, localDefaultVocabulary, false, true).(rdf.IRI)
				if !ok {
					// TODO warning
					continue
				}

				if attrInlist != nil {
					if _, known := listMapping[relValue]; !known {
						listMapping[relValue] = &listMappingBuilder{}
					}

					localIncompleteTriples = append(localIncompleteTriples, incompleteTriple{
						List:      listMapping[relValue],
						Direction: incompleteTripleDirectionNone,
					})
				} else {
					localIncompleteTriples = append(localIncompleteTriples, incompleteTriple{
						Predicate:      relValue,
						PredicateRange: predicateRange,
						Direction:      incompleteTripleDirectionForward,
					})
				}

				attrRelValid = true
			}
		}

		if attrRev != nil {
			var predicateRange *cursorio.TextOffsetRange

			if v.captureOffsets {
				if attrProfile := nodeProfile.TagAttr[attrRevIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
					predicateRange = attrProfile.ValueOffsets
				}
			}

			for _, revField := range strings.Fields(strings.TrimSpace(*attrRev)) {
				revValue, ok := resolveIRI(ectx, localPrefixMappings, revField, nil, localDefaultVocabulary, false, true).(rdf.IRI)
				if !ok {
					// TODO warning
					continue
				}

				localIncompleteTriples = append(localIncompleteTriples, incompleteTriple{
					Predicate:      revValue,
					PredicateRange: predicateRange,
					Direction:      incompleteTripleDirectionReverse,
				})
			}
		}
	}

	if attrProperty != nil {
		// Processing, Step 11

		var currentPropertyValue rdf.ObjectValue
		var currentPropertyValueAnno *cursorio.TextOffsetRange

		var datatypeIRI rdf.IRI

		if attrDatatype != nil {
			datatypeValue, ok := resolveIRI(ectx, localPrefixMappings, *attrDatatype, nil, localDefaultVocabulary, false, true).(rdf.IRI)
			if !ok {
				// TODO warning
			} else {
				datatypeIRI = datatypeValue
			}
		}

		if ectx.Global.HtmlProcessing&ActiveHtmlProcessingProfile > 0 && n.DataAtom == atom.Time && attrContent == nil {
			if attrDatetime == nil {
				// rdfa-in-html // 3.1 // Additional Processing Rule 10

				buf := bytes.Buffer{}

				v.collectTextContent(&buf, n)

				attrDatetime = ptr.Value(buf.String())

				if v.captureOffsets {
					if innerOffsets := nodeProfile.GetInnerOffsets(); innerOffsets != nil {
						currentPropertyValueAnno = innerOffsets
					}
				}
			} else {
				if v.captureOffsets {
					if attrProfile := nodeProfile.TagAttr[attrDatetimeIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
						currentPropertyValueAnno = attrProfile.ValueOffsets
					}
				}
			}

			// rdfa-in-html // 3.1 // Additional Processing Rule 9

			if len(datatypeIRI) > 0 {
				currentPropertyValue = rdf.Literal{
					LexicalForm: *attrDatetime,
					Datatype:    datatypeIRI,
				}
			} else if mapped, err := xsdobject.MapDuration(*attrDatetime); err == nil {
				currentPropertyValue = mapped
			} else if mapped, err := xsdobject.MapDateTime(*attrDatetime); err == nil {
				currentPropertyValue = mapped
			} else if mapped, err := xsdobject.MapDate(*attrDatetime); err == nil {
				currentPropertyValue = mapped
			} else if mapped, err := xsdobject.MapTime(*attrDatetime); err == nil {
				currentPropertyValue = mapped
			} else if mapped, err := xsdobject.MapGYearMonth(*attrDatetime); err == nil {
				currentPropertyValue = mapped
			} else if mapped, err := xsdobject.MapGYear(*attrDatetime); err == nil {
				currentPropertyValue = mapped
			} else {
				// TODO warning

				currentPropertyValue = rdf.Literal{
					Datatype:    xsdiri.String_Datatype,
					LexicalForm: *attrDatetime,
				}
			}

		} else if attrDatatype != nil && len(datatypeIRI) > 0 && datatypeIRI != rdfiri.XMLLiteral_Datatype {
			if attrContent != nil {
				currentPropertyValue = rdf.Literal{
					Datatype:    datatypeIRI,
					LexicalForm: *attrContent,
				}

				if v.captureOffsets {
					if attrProfile := nodeProfile.TagAttr[attrContentIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
						currentPropertyValueAnno = attrProfile.ValueOffsets
					}
				}
			} else {
				buf := bytes.Buffer{}

				v.collectTextContent(&buf, n)

				currentPropertyValue = rdf.Literal{
					Datatype:    datatypeIRI,
					LexicalForm: buf.String(),
				}

				if v.captureOffsets {
					if innerOffsets := nodeProfile.GetInnerOffsets(); innerOffsets != nil {
						currentPropertyValueAnno = innerOffsets
					}
				}
			}
		} else if attrDatatype != nil && len(datatypeIRI) == 0 {
			if attrContent != nil {
				currentPropertyValue = rdf.Literal{
					Datatype:    xsdiri.String_Datatype,
					LexicalForm: *attrContent,
				}

				if v.captureOffsets {
					if attrProfile := nodeProfile.TagAttr[attrContentIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
						currentPropertyValueAnno = attrProfile.ValueOffsets
					}
				}
			} else {
				buf := bytes.Buffer{}

				v.collectTextContent(&buf, n)

				currentPropertyValue = rdf.Literal{
					Datatype:    xsdiri.String_Datatype,
					LexicalForm: buf.String(),
				}

				if v.captureOffsets {
					if innerOffsets := nodeProfile.GetInnerOffsets(); innerOffsets != nil {
						currentPropertyValueAnno = innerOffsets
					}
				}
			}
		} else if attrDatatype != nil && datatypeIRI == rdfiri.XMLLiteral_Datatype {
			lexicalForm, err := v.xmlRender(n)
			if err != nil {
				return fmt.Errorf("xml render: %v", err)
			}

			currentPropertyValue = rdf.Literal{
				Datatype:    rdfiri.XMLLiteral_Datatype,
				LexicalForm: lexicalForm,
			}

			if v.captureOffsets {
				// TODO use first/last child instead?
				if innerOffsets := nodeProfile.GetInnerOffsets(); innerOffsets != nil {
					currentPropertyValueAnno = innerOffsets
				}
			}
		} else if attrDatatype != nil && datatypeIRI == rdfiri.HTML_Datatype {
			// rdfa-in-html // 3.1 // Additional Processing Rule 1

			lexicalForm, err := v.htmlRender(n)
			if err != nil {
				return fmt.Errorf("html render: %v", err)
			}

			currentPropertyValue = rdf.Literal{
				Datatype:    rdfiri.HTML_Datatype,
				LexicalForm: lexicalForm,
			}

			if v.captureOffsets {
				if innerOffsets := nodeProfile.GetInnerOffsets(); innerOffsets != nil {
					currentPropertyValueAnno = innerOffsets
				}
			}
		} else if attrContent != nil {
			currentPropertyValue = rdf.Literal{
				Datatype:    xsdiri.String_Datatype,
				LexicalForm: *attrContent,
			}

			if v.captureOffsets {
				if attrProfile := nodeProfile.TagAttr[attrContentIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
					currentPropertyValueAnno = attrProfile.ValueOffsets
				}
			}
		} else if !attrRelValid && attrRev == nil && attrContent == nil {
			if attrResource != nil {
				if s := resolveIRI(ectx, localPrefixMappings, *attrResource, localBaseURL, localDefaultVocabulary, true, true); s != nil {
					currentPropertyValue = s
					currentPropertyValueAnno = nil

					if v.captureOffsets {
						if attrProfile := nodeProfile.TagAttr[attrResourceIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
							currentPropertyValueAnno = attrProfile.ValueOffsets
						}
					}
				}
			}

			if currentPropertyValue == nil && attrHref != nil {
				if s := resolveIRI(ectx, localPrefixMappings, *attrHref, localBaseURL, localDefaultVocabulary, false, true); s != nil {
					currentPropertyValue = s
					currentPropertyValueAnno = nil

					if v.captureOffsets {
						if attrProfile := nodeProfile.TagAttr[attrHrefIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
							currentPropertyValueAnno = attrProfile.ValueOffsets
						}
					}
				}
			}

			if currentPropertyValue == nil && attrSrc != nil {
				if s := resolveIRI(ectx, localPrefixMappings, *attrSrc, localBaseURL, localDefaultVocabulary, false, true); s != nil {
					currentPropertyValue = s
					currentPropertyValueAnno = nil

					if v.captureOffsets {
						if attrProfile := nodeProfile.TagAttr[attrSrcIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
							currentPropertyValueAnno = attrProfile.ValueOffsets
						}
					}
				}
			}
		}

		if currentPropertyValue == nil {
			if attrTypeof != nil && attrAbout == nil {
				currentPropertyValue = typedResource
				currentPropertyValueAnno = typedResourceAnno
			} else {
				buf := bytes.Buffer{}

				v.collectTextContent(&buf, n)

				currentPropertyValue = rdf.Literal{
					Datatype:    xsdiri.String_Datatype,
					LexicalForm: buf.String(),
				}

				if v.captureOffsets {
					if innerOffsets := nodeProfile.GetInnerOffsets(); innerOffsets != nil {
						currentPropertyValueAnno = innerOffsets
					}
				}
			}
		}

		if currentLanguage != nil {
			if cpvLiteral, ok := currentPropertyValue.(rdf.Literal); ok && cpvLiteral.Datatype == xsdiri.String_Datatype {
				cpvLiteral.Datatype = rdfiri.LangString_Datatype
				cpvLiteral.Tag = rdf.LanguageLiteralTag{
					Language: *currentLanguage,
				}

				currentPropertyValue = cpvLiteral
			}
		}

		var predicateRange *cursorio.TextOffsetRange

		if v.captureOffsets {
			if attrProfile := nodeProfile.TagAttr[attrPropertyIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
				predicateRange = attrProfile.ValueOffsets
			}
		}

		for _, propertyField := range strings.Fields(strings.TrimSpace(*attrProperty)) {
			propertyValue, ok := resolveIRI(ectx, localPrefixMappings, propertyField, nil, localDefaultVocabulary, false, true).(rdf.IRI)
			if !ok {
				// TODO warning
				continue
			}

			if attrInlist != nil {
				if _, known := listMapping[propertyValue]; !known {
					listMapping[propertyValue] = &listMappingBuilder{}
				}

				listMapping[propertyValue].Objects = append(listMapping[propertyValue].Objects, currentPropertyValue)
			} else {
				v.statements = append(v.statements, statement{
					triple: rdf.Triple{
						Subject:   newSubject,
						Predicate: propertyValue,
						Object:    currentPropertyValue,
					},
					textOffsets: v.buildTextOffsets(
						encoding.SubjectStatementOffsets, newSubjectAnno,
						encoding.PredicateStatementOffsets, predicateRange,
						encoding.ObjectStatementOffsets, currentPropertyValueAnno,
					),
					containerResource: ectx.CurrentContainer,
				})
			}
		}
	}

	{ // Processing, Step 12
		if !skipElement && newSubject != nil {
			for _, incompleteTriple := range ectx.IncompleteTriples {
				switch incompleteTriple.Direction {
				case incompleteTripleDirectionNone:
					incompleteTriple.List.Objects = append(incompleteTriple.List.Objects, newSubject)
				case incompleteTripleDirectionForward:
					v.statements = append(v.statements, statement{
						triple: rdf.Triple{
							Subject:   ectx.ParentSubject,
							Predicate: incompleteTriple.Predicate,
							Object:    newSubject,
						},
						textOffsets: v.buildTextOffsets(
							encoding.SubjectStatementOffsets, ectx.ParentSubjectAnno,
							encoding.PredicateStatementOffsets, incompleteTriple.PredicateRange,
							encoding.ObjectStatementOffsets, newSubjectAnno,
						),
						containerResource: ectx.CurrentContainer,
					})
				case incompleteTripleDirectionReverse:
					v.statements = append(v.statements, statement{
						triple: rdf.Triple{
							Subject:   newSubject,
							Predicate: incompleteTriple.Predicate,
							Object:    ectx.ParentSubject,
						},
						textOffsets: v.buildTextOffsets(
							encoding.SubjectStatementOffsets, newSubjectAnno,
							encoding.PredicateStatementOffsets, incompleteTriple.PredicateRange,
							encoding.ObjectStatementOffsets, ectx.ParentSubjectAnno,
						),
						containerResource: ectx.CurrentContainer,
					})
				}
			}
		}
	}

	{ // Processing, Step 13
		var childectx = ectx

		if skipElement {
			childectx.Language = currentLanguage
			childectx.PrefixMapping = localPrefixMappings
			childectx.DefaultVocabulary = localDefaultVocabulary
		} else {
			childectx.BaseURL = localBaseURL

			if newSubject != nil {
				childectx.ParentSubject = newSubject
				childectx.ParentSubjectAnno = newSubjectAnno
			}

			if currentObjectResource != nil {
				childectx.ParentObject = currentObjectResource
				childectx.ParentObjectAnno = currentObjectResourceAnno
			} else if newSubject != nil {
				childectx.ParentObject = newSubject
				childectx.ParentObjectAnno = newSubjectAnno
			} else {
				childectx.ParentObject = childectx.ParentSubject
				childectx.ParentObjectAnno = childectx.ParentSubjectAnno
			}

			if ectx.ParentObject != nil && childectx.ParentObject != nil && ectx.ParentObject != childectx.ParentObject {
				// hacky conditions; not sure if this is fully correct
				childectx.CurrentContainer = nil
			}

			childectx.PrefixMapping = localPrefixMappings
			childectx.IncompleteTriples = localIncompleteTriples
			childectx.ListMapping = listMapping
			childectx.Language = currentLanguage
			childectx.DefaultVocabulary = localDefaultVocabulary
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			err := v.walkNode(childectx, c)
			if err != nil {
				return err
			}
		}
	}

	{ // Processing, Step 14
		for listPredicate, listItems := range listMapping {
			if ectx.ListMapping[listPredicate] == listItems {
				continue
			}

			if len(listItems.Objects) == 0 {
				v.statements = append(v.statements, statement{
					triple: rdf.Triple{
						Subject:   newSubject, // "current subject" in spec
						Predicate: listPredicate,
						Object:    rdfiri.Nil_List,
					},
					textOffsets: v.buildTextOffsets(
						encoding.SubjectStatementOffsets, newSubjectAnno,
					),
					containerResource: ectx.CurrentContainer,
				})

				continue
			}

			var listBnodes = make([]rdf.BlankNode, len(listItems.Objects))

			for listItemIdx := range listItems.Objects {
				listBnodes[listItemIdx] = ectx.Global.BlankNodeStringFactory.NewBlankNode()
			}

			for listItemIdx, listItem := range listItems.Objects {
				v.statements = append(v.statements, statement{
					triple: rdf.Triple{
						Subject:   listBnodes[listItemIdx],
						Predicate: rdfiri.First_Property,
						Object:    listItem,
					},
					containerResource: ectx.CurrentContainer,
				})

				if listItemIdx == len(listItems.Objects)-1 {
					v.statements = append(v.statements, statement{
						triple: rdf.Triple{
							Subject:   listBnodes[listItemIdx],
							Predicate: rdfiri.Rest_Property,
							Object:    rdfiri.Nil_List,
						},
						containerResource: ectx.CurrentContainer,
					})
				} else {
					v.statements = append(v.statements, statement{
						triple: rdf.Triple{
							Subject:   listBnodes[listItemIdx],
							Predicate: rdfiri.Rest_Property,
							Object:    listBnodes[listItemIdx+1],
						},
						containerResource: ectx.CurrentContainer,
					})
				}
			}

			v.statements = append(v.statements, statement{
				triple: rdf.Triple{
					Subject:   newSubject, // "current subject" in spec
					Predicate: listPredicate,
					Object:    listBnodes[0],
				},
				textOffsets: v.buildTextOffsets(
					encoding.SubjectStatementOffsets, newSubjectAnno,
				),
				containerResource: ectx.CurrentContainer,
			})
		}
	}

	if isRootElement && ectx.Global.HtmlProcessing&ActiveHtmlProcessingProfile > 0 {
		// rdfa-in-html // 3.5 // Property Copying
		// TODO this feels like a naive implementation; should decide how to handle:
		// - nested objects (probably need to rewrite blank nodes)
		// - maybe should outer-iterate over target subjects?
		// - might need multiple iterations for some case, suggested by spec comment?

		propertyCopyGraph := inmemory.NewDataset()

		{
			ctx := context.Background()

			for _, stmt := range v.statements {
				err := propertyCopyGraph.AddQuad(ctx, stmt.triple.AsQuad(nil))
				if err != nil {
					return fmt.Errorf("add: %v", err)
				}
			}
		}

		iter, err := propertyCopyGraph.QuerySimple(
			context.Background(),
			simplequery.Query{
				Select: []simplequery.Var{
					"subject",
					"target",
					"predicate",
					"object",
				},
				Where: simplequery.WhereTripleList{
					{
						Subject:   simplequery.Var("subject"),
						Predicate: simplequery.Term{Term: rdfairi.Copy_Property},
						Object:    simplequery.Var("target"),
					},
					{
						Subject:   simplequery.Var("target"),
						Predicate: simplequery.Term{Term: rdfiri.Type_Property},
						Object:    simplequery.Term{Term: rdfairi.Pattern_Class},
					},
					{
						Subject:   simplequery.Var("target"),
						Predicate: simplequery.Var("predicate"),
						Object:    simplequery.Var("object"),
					},
				},
			},
			simplequery.QueryOptions{},
		)
		if err != nil {
			return fmt.Errorf("query: %v", err)
		}

		deleteSubjectCopyTarget := map[[2]rdf.Term]struct{}{}
		deleteTargetTypePattern := map[rdf.Term]struct{}{}
		deleteTargetPredicateObject := map[[2]rdf.Term][]rdf.Term{}

		for iter.Next() {
			binding := iter.GetBinding()
			subject, predicate, object := binding.Get("subject"), binding.Get("predicate"), binding.Get("object")

			if predicate == rdfiri.Type_Property && object == rdfairi.Pattern_Class {
				// ambiguous: this exception is not described by the spec
				// exclusion matches spec output graph example, though
			} else {
				v.statements = append(v.statements, statement{
					triple: rdf.Triple{
						Subject:   subject.(rdf.SubjectValue),
						Predicate: predicate.(rdf.PredicateValue),
						Object:    object.(rdf.ObjectValue),
					},
					containerResource: ectx.CurrentContainer,
				})
			}

			target := binding.Get("target")

			deleteSubjectCopyTarget[[2]rdf.Term{subject, target}] = struct{}{}
			deleteTargetTypePattern[target] = struct{}{}
			deleteTargetPredicateObject[[2]rdf.Term{target, predicate.(rdf.IRI)}] = append(deleteTargetPredicateObject[[2]rdf.Term{target, predicate.(rdf.IRI)}], object)
		}

		if iter.Err() != nil {
			return fmt.Errorf("iter: %v", iter.Err())
		}

		var nextTuples []statement

		for _, tuple := range v.statements {
			if tuple.triple.Predicate == rdfairi.Copy_Property {
				if _, known := deleteSubjectCopyTarget[[2]rdf.Term{tuple.triple.Subject, tuple.triple.Object.(rdf.Term)}]; known {
					continue
				}
			} else if tuple.triple.Predicate == rdfiri.Type_Property && tuple.triple.Object == rdfairi.Pattern_Class {
				if _, known := deleteTargetTypePattern[tuple.triple.Subject]; known {
					continue
				}
			} else if matches, known := deleteTargetPredicateObject[[2]rdf.Term{tuple.triple.Subject, tuple.triple.Predicate}]; known {
				var matched bool

				for _, match := range matches {
					if tuple.triple.Object.TermEquals(match) {
						matched = true

						break
					}
				}

				if matched {
					continue
				}
			}

			nextTuples = append(nextTuples, tuple)
		}

		v.statements = nextTuples
	}

	return nil
}

func (v *Decoder) collectTextContent(buf *bytes.Buffer, n *html.Node) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		v.collectTextContent(buf, c)
	}
}

var (
	fieldsNextSpace    = regexp.MustCompile(`^\s+`)
	fieldsNextNonSpace = regexp.MustCompile(`^[^\s]+`)
)
