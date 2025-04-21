package htmlmicrodata

import (
	"bytes"
	"net/url"
	"regexp"
	"strings"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingutil"
	encodinghtml "github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdvalue"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdfio"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type DecoderOption interface {
	apply(s *DecoderConfig)
	newDecoder(doc *encodinghtml.Document) (*Decoder, error)
}

type Decoder struct {
	doc              *encodinghtml.Document
	docBaseURL       *iriutil.ParsedIRI
	buildTextOffsets encodingutil.TextOffsetsBuilderFunc

	captureOffsets     bool
	vocabularyResolver VocabularyResolver

	err          error
	statements   []*statement
	statementIdx int
}

var _ encoding.GraphDecoder = &Decoder{}

func NewDecoder(doc *encodinghtml.Document, opts ...DecoderOption) (*Decoder, error) {
	compiledOpts := DecoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newDecoder(doc)
}

func (r *Decoder) Close() error {
	return nil
}

func (w *Decoder) Err() error {
	return w.err
}

func (w *Decoder) Next() bool {
	if w.err != nil {
		return false
	} else if w.statementIdx == -1 {
		documentContainer := &DocumentResource{}
		ectx := evaluationContext{
			Global: &globalEvaluationContext{
				DocumentContainer:  documentContainer,
				ResolvedItemscopes: map[*html.Node]rdf.SubjectValue{},
				BlankNodeFactory:   blanknodeutil.NewFactory(),
			},
			BaseURL:          w.docBaseURL,
			CurrentContainer: documentContainer,
			RecursedItemrefs: map[string]struct{}{},
		}

		w.walk(ectx, w.doc.GetRoot())
	}

	w.statementIdx++

	return w.statementIdx < len(w.statements)
}

func (w *Decoder) GetTriple() rdf.Triple {
	return w.statements[w.statementIdx].triple
}

func (w *Decoder) GetStatement() rdfio.Statement {
	return w.statements[w.statementIdx]
}

func (w *Decoder) walk(ectx evaluationContext, n *html.Node) {
	if n.Namespace == "" { // http://www.w3.org/1999/xhtml
		var attrItemid, attrItemprop, attrItemref, attrItemtype string
		var attrItemidIdx, attrItempropIdx, attrItemtypeIdx int
		var attrItemscope bool

		for attrIdx, attr := range n.Attr {
			if attr.Namespace != "" {
				continue
			}

			switch attr.Key {
			case "itemid":
				attrItemid = attr.Val
				attrItemidIdx = attrIdx
			case "itemprop":
				attrItemprop = attr.Val
				attrItempropIdx = attrIdx
			case "itemref":
				attrItemref = attr.Val
			case "itemscope":
				attrItemscope = true
			case "itemtype":
				attrItemtype = attr.Val
				attrItemtypeIdx = attrIdx
			}
		}

		if attrItemscope {
			if ectx.Global.ResolvedItemscopes[n] == nil {
				nodeProfile, _ := w.doc.GetNodeMetadata(n)

				var nextContainer encoding.ContainerResource
				var nextSubject rdf.SubjectValue
				var nextSubjectRange *cursorio.TextOffsetRange

				if ectx.CurrentSubject == nil {
					nextContainer = ectx.Global.DocumentContainer
				}

				if len(attrItemid) > 0 {
					var sValue = strings.TrimSpace(attrItemid)

					if resolvedValue, err := ectx.ResolveURL(sValue); err == nil {
						sValue = resolvedValue
					}

					nextSubject = rdf.IRI(sValue)

					if w.captureOffsets {
						if attrProfile := nodeProfile.TagAttr[attrItemidIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
							nextSubjectRange = attrProfile.ValueOffsets
						}
					}
				} else {
					nextSubject = ectx.Global.BlankNodeFactory.NewBlankNode()

					if w.captureOffsets {
						if nodeProfile.EndTagTokenOffsets != nil {
							nextSubjectRange = &cursorio.TextOffsetRange{
								From:  nodeProfile.TokenOffsets.From,
								Until: nodeProfile.EndTagTokenOffsets.Until,
							}
						} else {
							nextSubjectRange = &nodeProfile.TokenOffsets
						}
					}
				}

				if len(attrItemprop) > 0 {
					if ectx.CurrentSubject == nil {
						// TODO warn
						// TODO warn iff not referenced by itemref parent id?
					} else {
						var attrCursorRange *cursorio.TextOffsetRange

						if w.captureOffsets {
							if attrProfile := nodeProfile.TagAttr[attrItempropIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
								attrCursorRange = attrProfile.ValueOffsets
							}
						}

						w.iterateItemprops(ectx, attrItemprop, attrCursorRange, func(p string, pCursorRange *cursorio.TextOffsetRange) {
							w.statements = append(w.statements, &statement{
								triple: rdf.Triple{
									Subject:   ectx.CurrentSubject,
									Predicate: rdf.IRI(p),
									Object:    nextSubject,
								},
								offsets: w.buildTextOffsets(
									encoding.SubjectStatementOffsets, ectx.CurrentSubjectRange,
									encoding.PredicateStatementOffsets, pCursorRange,
									encoding.ObjectStatementOffsets, nextSubjectRange,
								),
								containerResource: nextContainer,
							})
						})
					}
				}

				var nextItemtypes []string

				if len(attrItemtype) > 0 {
					var attrValOffset int
					var attrVal = attrItemtype

					var attrItemtypeKeyCursorRange *cursorio.TextOffsetRange

					if w.captureOffsets {
						if attrProfile := nodeProfile.TagAttr[attrItemtypeIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
							attrItemtypeKeyCursorRange = &attrProfile.KeyOffsets
						}
					}

					for len(attrVal) > 0 {
						if mm := regexp.MustCompile(`^\s+`).FindString(attrVal); len(mm) > 0 {
							attrValOffset += len(mm)
							attrVal = attrVal[len(mm):]

							continue
						}

						mm := regexp.MustCompile(`^[^\s]+`).FindString(attrVal)
						if len(mm) == 0 {
							panic("should not have found an empty match")
						}

						var oValue = mm

						oValueURL, err := url.Parse(oValue)
						if err != nil {
							// TODO warning
						} else {
							oValue = oValueURL.String()
						}

						nextItemtypes = append(nextItemtypes, oValue)

						// TODO recursive offset
						var attrCursorRange *cursorio.TextOffsetRange

						if w.captureOffsets {
							if attrProfile := nodeProfile.TagAttr[attrItemtypeIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
								attrCursorRange = attrProfile.ValueOffsets
							}
						}

						w.statements = append(w.statements, &statement{
							triple: rdf.Triple{
								Subject:   nextSubject,
								Predicate: rdfiri.Type_Property,
								Object:    rdf.IRI(oValue),
							},
							offsets: w.buildTextOffsets(
								encoding.SubjectStatementOffsets, nextSubjectRange,
								encoding.PredicateStatementOffsets, attrItemtypeKeyCursorRange,
								encoding.ObjectStatementOffsets, attrCursorRange,
							),
							containerResource: nextContainer,
						})

						attrValOffset += len(mm)
						attrVal = attrVal[len(mm):]
					}
				}

				ectx.CurrentSubject = nextSubject
				ectx.CurrentSubjectRange = nextSubjectRange
				ectx.CurrentItemtypes = nextItemtypes

				ectx.Global.ResolvedItemscopes[n] = nextSubject

				if len(attrItemref) > 0 {
					for _, itemref := range strings.Fields(strings.TrimSpace(attrItemref)) {
						if len(itemref) == 0 {
							continue
						}

						itemrefNodes := w.doc.GetNodesByID(itemref)
						if len(itemrefNodes) == 0 {
							// TODO warning

							continue
						} else if n == itemrefNodes[0] {
							// TODO warning

							continue
						} else if _, known := ectx.RecursedItemrefs[itemref]; known {
							// TODO warning

							continue
						} else if len(itemrefNodes) > 1 {
							// TODO warning
						}

						nectx := ectx
						nectx.RecursedItemrefs = map[string]struct{}{
							itemref: {},
						}

						for k, v := range ectx.RecursedItemrefs {
							nectx.RecursedItemrefs[k] = v
						}

						w.walk(nectx, itemrefNodes[0])
					}
				}
			}
		} else {
			if len(attrItemid) > 0 {
				// WARN
			}

			if len(attrItemref) > 0 {
				// WARN
			}

			if len(attrItemtype) > 0 {
				// WARN
			}

			if len(attrItemprop) > 0 && ectx.CurrentSubject != nil {
				nodeProfile, _ := w.doc.GetNodeMetadata(n)

				objectValue, objectValueCursorRange := w.parseMicrodataItemvalue(ectx, n)
				var attrCursorRange *cursorio.TextOffsetRange

				if w.captureOffsets {
					if attrProfile := nodeProfile.TagAttr[attrItempropIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
						attrCursorRange = attrProfile.ValueOffsets
					}
				}

				w.iterateItemprops(ectx, attrItemprop, attrCursorRange, func(p string, pCursorRange *cursorio.TextOffsetRange) {
					w.statements = append(w.statements, &statement{
						triple: rdf.Triple{
							Subject:   ectx.CurrentSubject,
							Predicate: rdf.IRI(p),
							Object:    objectValue,
						},
						offsets: w.buildTextOffsets(
							encoding.SubjectStatementOffsets, ectx.CurrentSubjectRange,
							encoding.PredicateStatementOffsets, pCursorRange,
							encoding.ObjectStatementOffsets, objectValueCursorRange,
						),
					})
				})
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		w.walk(ectx, c)
	}
}

func (w *Decoder) iterateItemprops(ectx evaluationContext, attrValue string, attrCursorRange *cursorio.TextOffsetRange, f func(p string, pCursorRange *cursorio.TextOffsetRange)) {
	knownItemprops := map[string]struct{}{}

	for _, itemprop := range strings.Fields(strings.TrimSpace(attrValue)) {
		if len(itemprop) == 0 {
			continue
		} else if _, known := knownItemprops[itemprop]; known {
			// TODO warning
			continue
		}

		knownItemprops[itemprop] = struct{}{}

		pValue, err := w.vocabularyResolver.ResolveMicrodataProperty(ectx.CurrentItemtypes, itemprop)
		if err != nil {
			// TODO warning
			continue
		}

		f(pValue, attrCursorRange)
	}
}

func (w *Decoder) collectTextContent(buf *bytes.Buffer, n *html.Node) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		w.collectTextContent(buf, c)
	}
}

func (w *Decoder) parseMicrodataItemvalue(ectx evaluationContext, n *html.Node) (rdf.ObjectValue, *cursorio.TextOffsetRange) {
	switch n.DataAtom {
	case atom.Meta:
		// [spec] The value is the value of the element's content attribute, if any, or the empty string if there is no such attribute.
		if v, offset := w.parseMicrodataItempropAttr(ectx, n, "content", itempropAttrString); v != nil {
			return v, offset
		}

		return rdf.Literal{
			Datatype: xsdiri.String_Datatype,
		}, nil
	case atom.Audio, atom.Embed, atom.Iframe, atom.Img, atom.Source, atom.Track, atom.Video:
		// [spec] The value is the result of encoding-parsing-and-serializing a URL given the element's src attribute's value, relative to the element's node document, at the time the attribute is set, or the empty string if there is no such attribute or the result is failure.

		if v, offset := w.parseMicrodataItempropAttr(ectx, n, "src", itempropAttrIRI); v != nil {
			return v, offset
		}

		return rdf.Literal{
			Datatype: xsdiri.String_Datatype,
		}, nil
	case atom.A, atom.Area, atom.Link:
		// [spec] The value is the result of encoding-parsing-and-serializing a URL given the element's href attribute's value, relative to the element's node document, at the time the attribute is set, or the empty string if there is no such attribute or the result is failure.

		if v, offset := w.parseMicrodataItempropAttr(ectx, n, "href", itempropAttrIRI); v != nil {
			return v, offset
		}

		return rdf.Literal{
			Datatype: xsdiri.String_Datatype,
		}, nil
	case atom.Object:
		// [spec] The value is the result of encoding-parsing-and-serializing a URL given the element's data attribute's value, relative to the element's node document, at the time the attribute is set, or the empty string if there is no such attribute or the result is failure.

		if v, offset := w.parseMicrodataItempropAttr(ectx, n, "data", itempropAttrIRI); v != nil {
			return v, offset
		}

		return rdf.Literal{
			Datatype: xsdiri.String_Datatype,
		}, nil
	case atom.Data:
		// [spec] The value is the value of the element's value attribute, if it has one, or the empty string otherwise.

		if v, offset := w.parseMicrodataItempropAttr(ectx, n, "value", itempropAttrString); v != nil {
			return v, offset
		}

		return rdf.Literal{
			Datatype: xsdiri.String_Datatype,
		}, nil
	case atom.Meter:
		// [spec] The value is the value of the element's value attribute, if it has one, or the empty string otherwise.
		// [dpb] unofficial cast meter to number?

		if v, offset := w.parseMicrodataItempropAttr(ectx, n, "value", itempropAttrMeter); v != nil {
			return v, offset
		}

		return rdf.Literal{
			Datatype: xsdiri.String_Datatype,
		}, nil
	case atom.Time:
		// [spec] The datetime value of a time element is the value of the element's datetime content attribute, if it has one, otherwise the child text content of the time element.
		if v, offset := w.parseMicrodataItempropAttr(ectx, n, "datetime", itempropAttrTime); v != nil {
			return v, offset
		}

		// fallthrough text content
	}

	buf := bytes.Buffer{}

	w.collectTextContent(&buf, n)

	var termCursorRange *cursorio.TextOffsetRange

	if buf.Len() == 0 {
		// avoid checking for offsets
		// also avoids GetInner which will panic if node had no children (aka missing end tag)
	} else if w.captureOffsets {
		if nodeProfile, ok := w.doc.GetNodeMetadata(n); ok {
			innerOffsets := nodeProfile.GetInnerOffsets()
			termCursorRange = &innerOffsets
		}
	}

	return rdf.Literal{
		Datatype:    xsdiri.String_Datatype,
		LexicalForm: buf.String(),
	}, termCursorRange
}

type itempropAttrObjectValue func(ectx evaluationContext, v string) rdf.ObjectValue

func itempropAttrString(ectx evaluationContext, v string) rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    xsdiri.String_Datatype,
		LexicalForm: v,
	}
}

func itempropAttrTime(ectx evaluationContext, v string) rdf.ObjectValue {
	if mapped, err := xsdvalue.MapDate(v); err == nil {
		return mapped.AsLiteralTerm()
	} else if mapped, err := xsdvalue.MapTime(v); err == nil {
		return mapped.AsLiteralTerm()
	} else if mapped, err := xsdvalue.MapDateTime(v); err == nil {
		return mapped.AsLiteralTerm()
	} else if mapped, err := xsdvalue.MapGYearMonth(v); err == nil {
		return mapped.AsLiteralTerm()
	} else if mapped, err := xsdvalue.MapGYear(v); err == nil {
		return mapped.AsLiteralTerm()
	} else if mapped, err := xsdvalue.MapDuration(v); err == nil {
		return mapped.AsLiteralTerm()
	}

	return rdf.Literal{
		Datatype:    xsdiri.String_Datatype,
		LexicalForm: v,
	}
}

func itempropAttrMeter(ectx evaluationContext, v string) rdf.ObjectValue {
	{
		vv, err := xsdvalue.MapInteger(v)
		if err == nil {
			return vv.AsLiteralTerm()
		}
	}

	{
		vv, err := xsdvalue.MapDecimal(v)
		if err == nil {
			return vv.AsLiteralTerm()
		}
	}

	return rdf.Literal{
		Datatype:    xsdiri.String_Datatype,
		LexicalForm: v,
	}
}

func itempropAttrIRI(ectx evaluationContext, v string) rdf.ObjectValue {
	if resolvedValue, err := ectx.ResolveURL(v); err == nil {
		v = resolvedValue
	} else {
		// TODO warn
	}

	return rdf.IRI(v)
}

func (w *Decoder) parseMicrodataItempropAttr(ectx evaluationContext, n *html.Node, attrKey string, attrValuer itempropAttrObjectValue) (rdf.ObjectValue, *cursorio.TextOffsetRange) {
	for attrIdx, attr := range n.Attr {
		if attr.Namespace != "" {
			continue
		} else if attr.Key == attrKey {
			var oValue = attr.Val

			var termCursorRange *cursorio.TextOffsetRange

			if w.captureOffsets {
				if nodeProfile, ok := w.doc.GetNodeMetadata(n); ok {
					if attrProfile := nodeProfile.TagAttr[attrIdx]; attrProfile != nil && attrProfile.ValueOffsets != nil {
						termCursorRange = attrProfile.ValueOffsets
					}
				}
			}

			return attrValuer(ectx, oValue), termCursorRange
		}
	}

	// TODO warning missing

	return nil, nil
}
