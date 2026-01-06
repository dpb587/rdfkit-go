package jsonld

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"slices"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldcontent"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type EncoderOption interface {
	apply(s *EncoderConfig)
	newEncoder(w io.Writer) (*Encoder, error)
}

type Encoder struct {
	w                 *json.Encoder
	base              *iriutil.BaseIRI
	prefixes          *iriutil.PrefixTracker
	buffered          bool
	blankNodeStringer blanknodeutil.Stringer

	err     error
	builder *rdfdescription.DatasetResourceListBuilder
}

var _ encoding.QuadsEncoder = &Encoder{}

func NewEncoder(w io.Writer, opts ...EncoderOption) (*Encoder, error) {
	compiledOpts := EncoderConfig{}

	for _, opt := range opts {
		opt.apply(&compiledOpts)
	}

	return compiledOpts.newEncoder(w)
}

func (e *Encoder) GetContentMetadata() encoding.ContentMetadata {
	return jsonldcontent.DefaultMetadata
}

func (e *Encoder) GetContentTypeIdentifier() encoding.ContentTypeIdentifier {
	return jsonldcontent.TypeIdentifier
}

func (e *Encoder) Close() error {
	if e.err != nil {
		return e.err
	}

	var graphItems = []any{}

	for _, graphName := range e.builder.GetGraphNames() {
		if graphName != nil {
			continue // TODO multi-graph support
		}

		builder := e.builder.GetResourceListBuilder(graphName)

		for _, resource := range builder.GetResources() {
			graphItems = append(graphItems, e.buildResource(builder, resource, true))
		}
	}

	if e.buffered && len(graphItems) > 1 {
		for i, item := range graphItems {
			marshaled, _ := json.Marshal(item)
			graphItems[i] = json.RawMessage(marshaled)
		}

		slices.SortFunc(graphItems, func(a, b any) int {
			return bytes.Compare(a.(json.RawMessage), b.(json.RawMessage))
		})
	}

	var wrapped map[string]any

	if len(graphItems) == 1 {
		wrapped = graphItems[0].(map[string]any)
	} else {
		wrapped = map[string]any{
			"@graph": graphItems,
		}
	}

	var wrappedContext = map[string]any{}

	if e.base != nil {
		wrappedContext["@base"] = e.base.String()
	}

	if prefixMappings := e.prefixes.GetUsedPrefixMappings(); len(prefixMappings) > 0 {
		for _, pm := range prefixMappings {
			wrappedContext[pm.Prefix] = pm.Expanded
		}
	}

	if len(wrappedContext) > 0 {
		wrapped["@context"] = wrappedContext
	}

	err := e.w.Encode(wrapped)
	if err != nil {
		return fmt.Errorf("encode: %v", err)
	}

	e.err = io.ErrClosedPipe

	return nil
}

func (w *Encoder) AddQuad(ctx context.Context, t rdf.Quad) error {
	if w.err != nil {
		return w.err
	}

	w.builder.Add(t)

	return nil
}

func (e *Encoder) buildResource(builder *rdfdescription.ResourceListBuilder, resource rdfdescription.Resource, root bool) map[string]any {
	graphProperties := make(map[string][]any)

	for _, statement := range resource.GetResourceStatements() {
		var statementObject any
		var predicate rdf.IRI

		switch statementT := statement.(type) {
		case rdfdescription.AnonResourceStatement:
			predicate = statementT.Predicate.(rdf.IRI)
			statementObject = e.buildResource(builder, statementT.AnonResource, false)
		case rdfdescription.ObjectStatement:
			predicate = statementT.Predicate.(rdf.IRI)

			switch obj := statementT.Object.(type) {
			case rdf.IRI:
				var wrapID string

				if iriPrefix, iriSuffix, ok := e.prefixes.CompactPrefix(obj); ok {
					wrapID = iriPrefix + ":" + iriSuffix
				} else if e.base != nil {
					if rel, ok := e.base.RelativizeIRI(obj); ok {
						wrapID = rel
					} else {
						wrapID = string(obj)
					}
				} else {
					wrapID = string(obj)
				}

				if predicate == rdfiri.Type_Property {
					statementObject = wrapID
				} else {
					statementObject = map[string]any{
						"@id": wrapID,
					}
				}
			case rdf.BlankNode:
				statementObject = map[string]any{
					"@id": "_:" + e.blankNodeStringer.GetBlankNodeIdentifier(obj),
				}
			case rdf.Literal:
				switch obj.Datatype {
				case xsdiri.String_Datatype:
					statementObject = obj.LexicalForm
				case xsdiri.Integer_Datatype, xsdiri.Double_Datatype:
					// TODO avoid number overflow
					statementObject = json.Number(obj.LexicalForm)
				case xsdiri.Boolean_Datatype:
					switch obj.LexicalForm {
					case "true":
						statementObject = true
					case "false":
						statementObject = false
					default:
						typePrefix, typeSuffix, ok := e.prefixes.CompactPrefix(obj.Datatype)
						if ok {
							statementObject = map[string]any{
								"@value": obj.LexicalForm,
								"@type":  typePrefix + ":" + typeSuffix,
							}
						} else {
							statementObject = map[string]any{
								"@value": obj.LexicalForm,
								"@type":  string(obj.Datatype),
							}
						}
					}
				default:
					typePrefix, typeSuffix, ok := e.prefixes.CompactPrefix(obj.Datatype)
					if ok {
						statementObject = map[string]any{
							"@value": obj.LexicalForm,
							"@type":  typePrefix + ":" + typeSuffix,
						}
					} else {
						statementObject = map[string]any{
							"@value": obj.LexicalForm,
							"@type":  string(obj.Datatype),
						}
					}

					if obj.Datatype == rdfiri.LangString_Datatype {
						if tag, ok := obj.Tag.(rdf.LanguageLiteralTag); ok {
							statementObject.(map[string]any)["@language"] = tag.Language
						}
					}
				}
			}
		default:
			panic(fmt.Errorf("unsupported statement type: %T", statementT))
		}

		var key string = string(predicate)

		if predicate == rdfiri.Type_Property {
			key = "@type"
		} else if predicatePrefix, predicateSuffix, ok := e.prefixes.CompactPrefix(predicate); ok {
			key = predicatePrefix + ":" + predicateSuffix
		}

		graphProperties[key] = append(graphProperties[key], statementObject)
	}

	var graphItem = map[string]any{}

	switch v := resource.GetResourceSubject().(type) {
	case rdf.IRI:
		if iriPrefix, iriSuffix, ok := e.prefixes.CompactPrefix(v); ok {
			graphItem["@id"] = iriPrefix + ":" + iriSuffix
		} else if e.base != nil {
			if rel, ok := e.base.RelativizeIRI(v); ok {
				graphItem["@id"] = rel
			} else {
				graphItem["@id"] = string(v)
			}
		} else {
			graphItem["@id"] = string(v)
		}
	case rdf.BlankNode:
		if root {
			if builder.GetBlankNodeReferences(v) > 0 {
				graphItem["@id"] = "_:" + e.blankNodeStringer.GetBlankNodeIdentifier(v)
			}
		} else if builder.GetBlankNodeReferences(v) > 1 {
			graphItem["@id"] = "_:" + e.blankNodeStringer.GetBlankNodeIdentifier(v)
		}
	case nil:
		// AnonResource
	default:
		panic(fmt.Errorf("unsupported resource subject type: %T", v))
	}

	for key, values := range graphProperties {
		if len(values) == 1 {
			graphItem[key] = values[0]
		} else if len(values) > 1 {
			if e.buffered {
				for i, value := range values {
					marshaled, _ := json.Marshal(value)
					values[i] = json.RawMessage(marshaled)
				}

				slices.SortFunc(values, func(a, b any) int {
					return bytes.Compare(a.(json.RawMessage), b.(json.RawMessage))
				})
			}

			graphItem[key] = values
		}
	}

	return graphItem
}
