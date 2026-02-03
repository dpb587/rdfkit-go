package jsonldinternal

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/iri"
)

const MagicKeywordPropertySourceOffsets = "@rdfkit.property.sourceOffsets"

var (
	tokenStringDirection = inspectjson.StringValue{
		Value: "@direction",
	}
	tokenStringGraph = inspectjson.StringValue{
		Value: "@graph",
	}
	tokenStringId = inspectjson.StringValue{
		Value: "@id",
	}
	tokenStringIndex = inspectjson.StringValue{
		Value: "@index",
	}
	tokenStringJson = inspectjson.StringValue{
		Value: "@json",
	}
	tokenStringLanguage = inspectjson.StringValue{
		Value: "@language",
	}
	tokenStringList = inspectjson.StringValue{
		Value: "@list",
	}
	tokenStringReverse = inspectjson.StringValue{
		Value: "@reverse",
	}
	tokenStringType = inspectjson.StringValue{
		Value: "@type",
	}
	tokenStringValue = inspectjson.StringValue{
		Value: "@value",
	}
)

var reKeywordABNF = regexp.MustCompile(`^@[a-zA-Z]+$`)

func Expand(input inspectjson.Value, opts jsonldtype.ProcessorOptions) (ExpandedValue, error) {
	// [spec // 9.1 // expand // 9] Set *expanded output* to the result of using the Expansion algorithm, passing the *active context*, `document` from *remote document* or input if there is no *remote document* as *element*, `null` as *active property*, `documentUrl` as *base URL*, if available, otherwise to the `base` option from options, and the `frameExpansion` and and `ordered` flags from options.

	if len(opts.ProcessingMode) == 0 {
		opts.ProcessingMode = ProcessingMode_JSON_LD_1_1
	}

	if opts.DocumentLoader == nil {
		opts.DocumentLoader = jsonldtype.DocumentLoaderFunc(func(ctx context.Context, url string, opts jsonldtype.DocumentLoaderOptions) (jsonldtype.RemoteDocument, error) {
			return jsonldtype.RemoteDocument{}, errors.New("no document loader configured")
		})
	}

	var baseIRI *iri.ParsedIRI

	if len(opts.BaseURL) > 0 {
		var err error

		baseIRI, err = iri.ParseIRI(opts.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("parse base url: %v", err)
		}
	}

	activeContext := &Context{
		BaseURL:         baseIRI,
		OriginalBaseURL: baseIRI,
		TermDefinitions: map[string]*TermDefinition{},
		_processor: &contextProcessor{
			ctx:                       context.Background(),
			processingMode:            opts.ProcessingMode,
			dereferencedDocumentByIRI: map[string]dereferencedDocument{},
			documentLoader:            opts.DocumentLoader,
		},
	}

	if opts.ExpandContext != nil {
		localContext := opts.ExpandContext

		if localContextMap, ok := localContext.(inspectjson.ObjectValue); ok {
			if contextMember, ok := localContextMap.Members["@context"]; ok {
				localContext = contextMember.Value
			}
		}

		expandedContext, err := algorithmContextProcessing{
			ActiveContext: activeContext,
			LocalContext:  localContext,
			BaseURL:       activeContext.BaseURL,
			// defaults
			RemoteContexts:        nil,
			OverrideProtected:     false,
			Propagate:             true,
			ValidateScopedContext: true,
		}.Call()
		if err != nil {
			return nil, err
		}

		activeContext = expandedContext
	}

	expandedOutput, err := algorithmExpansion{
		activeContext: activeContext,
		element:       input,
		baseURL:       baseIRI,
		ordered:       true,
	}.Call()
	if err != nil {
		return nil, err
	}

	// [spec // 9.1 // expand // 8.1] If *expanded output* is a map that contains only an `@graph` entry, set *expanded output* that value.

	if expandedOutput != nil {
		if expandedOutputMap, ok := expandedOutput.(*ExpandedObject); ok {
			if len(expandedOutputMap.Members) == 1 {
				if graphMember, ok := expandedOutputMap.Members["@graph"]; ok {
					expandedOutput = graphMember
				}
			}
		}
	} else {

		// [spec // 9.1 // expand // 8.2] If *expanded output* is `null`, set *expanded output* to an empty array.

		expandedOutput = &ExpandedArray{}
	}

	if _, ok := expandedOutput.(*ExpandedArray); !ok {
		expandedOutput = &ExpandedArray{
			Values: []ExpandedValue{expandedOutput},
		}
	}

	return expandedOutput, nil
}
