package jsonldinternal

import (
	"errors"
	"fmt"
	"maps"
	"net/url"
	"slices"
	"strings"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type algorithmContextProcessing struct {
	ActiveContext *Context
	LocalContext  inspectjson.Value

	// [spec // 4.1.2] used when resolving relative context URLs
	BaseURL *iriutil.ParsedIRI

	// Optional

	// [spec // 4.1.2] defaulting to a new empty array, which is used to detect cyclical context inclusions
	RemoteContexts []string

	// [spec // 4.1.2] defaulting to false, which is used to allow changes to protected terms
	OverrideProtected bool

	// [spec // 4.1.2] defaulting to true to mark term definitions associated with non-propagated contexts
	Propagate bool

	// [spec // 4.1.2] defaulting to true, which is used to limit recursion when validating possibly recursive scoped contexts
	ValidateScopedContext bool
}

// func (c *Context) Parse(localContext inspectjson.Value) (*Context, error) {
// 	return algorithmContextProcessing{
// 		ActiveContext:         c,
// 		LocalContext:          localContext,
// 		BaseURL:               c.BaseURL,
// 		OverrideProtected:     false,
// 		Propagate:             true,
// 		ValidateScopedContext: true,
// 	}.Call()
// }

func (vars algorithmContextProcessing) Call() (*Context, error) {
	var localContext []inspectjson.Value

	{
		// [dpb] handled earlier from [4.1.2 // 4]

		if localContextArray, ok := vars.LocalContext.(inspectjson.ArrayValue); ok {
			localContext = localContextArray.Values
		} else {
			localContext = []inspectjson.Value{vars.LocalContext}
		}

		if len(localContext) == 0 {
			// unofficial optimization
			return vars.ActiveContext, nil
		}
	}

	// [spec // 4.1.2 // 1] Initialize *result* to the result of cloning *active context*, with inverse context set to `null`..
	// [dpb] typo: double period

	result := vars.ActiveContext.clone()
	result.InverseContext = nil

	// [spec // 4.1.2 // 2] If *local context* is an object containing the member @propagate, its value *MUST* be boolean `true` or `false`, set *propagate* to that value.
	// [spec // 4.1.2 // 2] NOTE Error handling is performed in step 5.11.

	if localContextObject, ok := vars.LocalContext.(inspectjson.ObjectValue); ok {
		if propagateObjectMember, ok := localContextObject.Members["@propagate"]; ok {
			if propagatePrimitive, ok := propagateObjectMember.Value.(inspectjson.BooleanValue); ok {
				vars.Propagate = propagatePrimitive.Value
			}
		}
	}

	// [spec // 4.1.2 // 3] If *propagate* is `false``, and *result* does not have a previous context, set previous context in *result* to *active context*.

	if !vars.Propagate && result.PreviousContext == nil {
		result.PreviousContext = vars.ActiveContext
	}

	// [spec // 4.1.2 // 4] If *local context* is not an array, set *local context* to an array containing only *local context*.

	// -- handled earlier

	// [spec // 4.1.2 // 5] For each item *context* in *local context*:

	for _, context := range localContext {

		// [spec // 4.1.2 // 5.1] If *context* is `null`:

		if _, ok := context.(inspectjson.NullValue); ok {

			// [spec // 4.1.2 // 5.1.1] If *override protected* is `false` and *active context* contains any protected term definitions, an `invalid context nullification` has been detected and processing is aborted.

			if !vars.OverrideProtected {
				for _, termDefinition := range result.TermDefinitions {
					if termDefinition.Protected {
						return nil, jsonldtype.Error{
							Code: jsonldtype.InvalidContextNullification,
						}
					}
				}
			}

			// [spec // 4.1.2 // 5.1.2] Initialize *result* as a newly-initialized active context, setting both base IRI and original base URL to the value of original base URL in *active context*, and, if *propagate* is `false`, previous context in *result* to the previous value of *result*.

			{
				nextResult := &Context{
					BaseURL:         vars.ActiveContext.OriginalBaseURL,
					OriginalBaseURL: vars.ActiveContext.OriginalBaseURL,
					//
					TermDefinitions: map[string]*TermDefinition{},
					_processor:      vars.ActiveContext._processor,
				}

				if !vars.Propagate {
					nextResult.PreviousContext = result
				}

				result = nextResult
			}

			// [spec // 4.1.2 // 5.1.3] Continue with the next *context*.

			continue
		}

		// [spec // 4.1.2 // 5.2] If *context* is a string,

		if contextPrimitive, ok := context.(inspectjson.StringValue); ok {

			// [spec // 4.1.2 // 5.2.1] Initialize *context* to the result of resolving *context* against *base URL*. If *base URL* is not a valid IRI, then *context* *MUST* be a valid IRI, otherwise a `loading document failed` error has been detected and processing is aborted.
			// [spec // 4.1.2 // 5.2.1] NOTE *base URL* is often not the same as base or the base IRI of the *active context*.

			contextURL, err := resolveURL(vars.BaseURL, contextPrimitive.Value)
			if err != nil {
				return nil, jsonldtype.Error{
					Code: jsonldtype.LoadingDocumentFailed,
					Err:  err,
				}
			}

			_contextURLString := contextURL.String()

			// [spec // 4.1.2 // 5.2.2] If *validate scoped context* is `false`, and `remote contexts` already includes `context` do not process `context` further and continue to any next `context` in `local context`.

			if !vars.ValidateScopedContext {
				var alreadyFound bool

				for _, already := range vars.RemoteContexts {
					if already == _contextURLString {
						alreadyFound = true

						break
					}
				}

				if alreadyFound {
					continue
				}
			}

			// [spec // 4.1.2 // 5.2.3] If the number of entries in the *remote contexts* array exceeds a processor defined limit, a `context overflow` error has been detected and processing is aborted; otherwise, add `context` to `remote contexts`.

			{
				if len(vars.RemoteContexts) > ProcessorRemoteContextsLimit {
					return nil, jsonldtype.Error{
						Code: jsonldtype.ContextOverflow,
						Err:  fmt.Errorf("limit exceeded: %d", ProcessorRemoteContextsLimit),
					}
				}

				vars.RemoteContexts = append(vars.RemoteContexts, _contextURLString)
			}

			// [spec // 4.1.2 // 5.2.4] If *context* was previously dereferenced, then the processor *MUST NOT* do a further dereference, and context is set to the previously established internal representation: set *context document* to the previously dereferenced document, and set *loaded context* to the value of the `@context` entry from the document in *context document*.

			var contextDocumentIRI *iriutil.ParsedIRI
			var loadedContext inspectjson.Value

			if dereferenced, ok := result._processor.dereferencedDocumentByIRI[_contextURLString]; ok {
				contextDocumentIRI = dereferenced.documentURL
				loadedContext = dereferenced.documentContextValue
			} else {

				// [spec // 4.1.2 // 5.2.5] Otherwise, set *context document* to the `RemoteDocument` obtained by dereferencing *context* using the `LoadDocumentCallback`, passing *context* for url, and `http://www.w3.org/ns/json-ld#context` for `profile` and for `requestProfile`.

				contextDocument, err := result._processor.documentLoader.LoadDocument(
					result._processor.ctx,
					contextURL.String(),
					jsonldtype.DocumentLoaderOptions{
						Profile:        &profileContextIRI,
						RequestProfile: []string{profileContextIRI},
					},
				)

				// [spec // 4.1.2 // 5.2.5.1] If *context* cannot be dereferenced, or the `document` from *context document* cannot be transformed into the internal representation , a `loading remote context failed` error has been detected and processing is aborted.

				if err != nil {
					return nil, jsonldtype.Error{
						Code: jsonldtype.LoadingRemoteContextFailed,
						Err:  err,
					}
				}

				// [spec // 4.1.2 // 5.2.5.2] If the `document` has no top-level map with an `@context` entry, an `invalid remote context` has been detected and processing is aborted.
				// [dpb] in theory, DocumentValue should support extractAllScripts flag; currently assumes single top-level JSON object

				contextDocumentObject, ok := contextDocument.Document.(inspectjson.ObjectValue)
				if !ok {
					return nil, jsonldtype.Error{
						Code: jsonldtype.InvalidRemoteContext,
						Err:  fmt.Errorf("invalid top-level type: %s", contextDocument.Document.GetGrammarName()),
					}
				}

				objectAtContext, ok := contextDocumentObject.Members["@context"]
				if !ok {
					return nil, jsonldtype.Error{
						Code: jsonldtype.InvalidRemoteContext,
						Err:  errors.New("object is missing @context entry"),
					}
				}

				contextDocumentIRI, err = iriutil.ParseIRI(contextDocument.DocumentURL.String())
				if err != nil {
					return nil, jsonldtype.Error{
						Code: jsonldtype.InvalidRemoteContext,
						Err:  fmt.Errorf("invalid document URL: %v", err),
					}
				}

				loadedContext = objectAtContext.Value

				// [dpb] a step to store dereferenced not mentioned in spec

				result._processor.dereferencedDocumentByIRI[_contextURLString] = dereferencedDocument{
					documentURL:          contextDocumentIRI,
					documentContextValue: loadedContext,
				}
			}

			// [spec // 4.1.2 // 5.2.6] Set *result* to the result of recursively calling this algorithm, passing *result* for *active context*, *loaded context* for *local context*, the `documentUrl` of *context document* for *base URL*, a copy of *remote contexts*, and *validate scoped context*.
			// [spec // 4.1.2 // 5.2.6] NOTE If *context* was previously dereferenced, processors MUST make provisions for retaining the *base URL* of that context for this step to enable the resolution of any relative context URLs that may be encountered during processing.

			nextResult, err := algorithmContextProcessing{
				ActiveContext:         result,
				LocalContext:          loadedContext,
				BaseURL:               contextDocumentIRI,
				RemoteContexts:        vars.RemoteContexts[:],
				ValidateScopedContext: vars.ValidateScopedContext,
				// defaults
				OverrideProtected: false,
				Propagate:         true,
			}.Call()
			if err != nil {
				return nil, err
			}

			result = nextResult

			// [spec // 4.1.2 // 5.2.7] Continue with the next *context*.

			continue
		}

		// [spec // 4.1.2 // 5.3] If *context* is not a map, an `invalid local context` error has been detected and processing is aborted.

		contextObject, ok := context.(inspectjson.ObjectValue)
		if !ok {
			return nil, jsonldtype.Error{
				Code: jsonldtype.InvalidLocalContext,
				Err:  fmt.Errorf("invalid type: %s", context.GetGrammarName()),
			}
		}

		// [spec // 4.1.2 // 5.4] Otherwise, *context* is a context definition.

		// [spec // 4.1.2 // 5.5] If *context* has an `@version` entry:

		if versionObjectMember, ok := contextObject.Members["@version"]; ok {

			// [spec // 4.1.2 // 5.5.1] If the associated value is not `1.1`, an `invalid @version value` has been detected, and processing is aborted.
			// [spec // 4.1.2 // 5.5.1] NOTE The use of `1.1` for the value of `@version` is intended to cause a JSON-LD 1.0 processor to stop processing. Although it is clearly meant to be related to JSON-LD 1.1, it does not otherwise adhere to the requirements for Semantic Versioning. Implementations may require special consideration when comparing the values of numbers with a non-zero fractional part.

			if versionPrimitive, ok := versionObjectMember.Value.(inspectjson.NumberValue); ok && versionPrimitive.Value != 1.1 {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidAtVersionValue,
					Err:  fmt.Errorf("invalid type: %s", versionObjectMember.Value.GetGrammarName()),
				}
			}

			// [spec // 4.1.2 // 5.5.2] If processing mode is set to `json-ld-1.0`, a `processing mode conflict` error has been detected and processing is aborted.

			if result._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
				return nil, jsonldtype.Error{
					Code: jsonldtype.ProcessingModeConflict,
					Err:  fmt.Errorf("invalid entry (processing mode %s): @version", result._processor.processingMode),
				}
			}
		}

		// [spec // 4.1.2 // 5.6] If *context* has an `@import` entry:

		if importObjectMember, ok := contextObject.Members["@import"]; ok {

			// [spec // 4.1.2 // 5.6.1] If processing mode is `json-ld-1.0``, an `invalid context entry` error has been detected and processing is aborted.

			if result._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidContextEntry,
					Err:  fmt.Errorf("invalid entry (processing mode %s): @import", result._processor.processingMode),
				}
			}

			// [spec // 4.1.2 // 5.6.2] Otherwise, if the value of `@import` is not a string, an `invalid @import value` error has been detected and processing is aborted.

			importString, ok := importObjectMember.Value.(inspectjson.StringValue)
			if !ok {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidAtImportValue,
					Err:  fmt.Errorf("invalid type: %s", importObjectMember.Value.GetGrammarName()),
				}
			}

			// [spec // 4.1.2 // 5.6.3] Initialize *import* to the result of resolving the value of `@import` against *base URL*.

			importURL, err := resolveURL(vars.BaseURL, importString.Value)
			if err != nil {
				return nil, jsonldtype.Error{
					Code: jsonldtype.LoadingDocumentFailed,
					Err:  err,
				}
			}

			// [spec // 4.1.2 // 5.6.4] Dereference *import* using the `LoadDocumentCallback`, passing *import* for url, and `http://www.w3.org/ns/json-ld#context` for `profile` and for `requestProfile`.

			importDocument, err := result._processor.documentLoader.LoadDocument(
				result._processor.ctx,
				importURL.String(),
				jsonldtype.DocumentLoaderOptions{
					Profile:        &profileContextIRI,
					RequestProfile: []string{profileContextIRI},
				},
			)

			// [spec // 4.1.2 // 5.6.5] If *import* cannot be dereferenced, or cannot be transformed into the internal representation, a `loading remote context failed` error has been detected and processing is aborted.

			if err != nil {
				return nil, jsonldtype.Error{
					Code: jsonldtype.LoadingRemoteContextFailed,
					Err:  err,
				}
			}

			// [spec // 4.1.2 // 5.6.6] If the dereferenced document has no top-level map with an `@context` entry, or if the value of `@context` is not a context definition (i.e., it is not an map), an `invalid remote context` has been detected and processing is aborted; otherwise, set *import context* to the value of that entry.

			importDocumentObject, ok := importDocument.Document.(inspectjson.ObjectValue)
			if !ok {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidRemoteContext,
					Err:  fmt.Errorf("invalid top-level type: %s", importDocument.Document.GetGrammarName()),
				}
			}

			objectAtContext, ok := importDocumentObject.Members["@context"]
			if !ok {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidRemoteContext,
					Err:  errors.New("object is missing @context entry"),
				}
			}

			atContextObject, ok := objectAtContext.Value.(inspectjson.ObjectValue)
			if !ok {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidRemoteContext,
					Err:  fmt.Errorf("invalid @context type: %s", objectAtContext.Value.GetGrammarName()),
				}
			}

			importContext := atContextObject

			// [spec // 4.1.2 // 5.6.7] If *import context* has a `@import` entry, an `invalid context entry` error has been detected and processing is aborted.

			if _, ok := importContext.Members["@import"]; ok {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidContextEntry,
					Err:  errors.New("invalid @import entry within imported context definition"),
				}
			}

			// [spec // 4.1.2 // 5.6.8] Set *context* to the result of merging *context* into *import context*, replacing common entries with those from *context*.

			for k, v := range contextObject.Members {
				atContextObject.Members[k] = v
			}

			contextObject = atContextObject
		}

		// [spec // 4.1.2 // 5.7] If *context* has an `@base` entry and `remote contexts` is empty, i.e., the currently being processed context is not a remote context:

		if baseObjectMember, ok := contextObject.Members["@base"]; ok && len(vars.RemoteContexts) == 0 {

			// [spec // 4.1.2 // 5.7.1] Initialize *value* to the value associated with the `@base` entry.

			value := baseObjectMember.Value

			// [spec // 4.1.2 // 5.7.2] If *value* is `null`, remove the base IRI of `result`.

			if _, ok := value.(inspectjson.NullValue); ok {
				result.BaseURL = nil
				result.BaseURLValue = nil

				// [spec // 4.1.2 // 5.7.3] Otherwise, if *value* is an IRI, the base IRI of *result* is set to *value*.
			} else if valueString, ok := value.(inspectjson.StringValue); ok {
				valueIRI, err := iriutil.ParseIRI(valueString.Value)
				if err != nil {
					return nil, jsonldtype.Error{
						Code: jsonldtype.InvalidBaseIRI, // per 5.7.5
						Err:  err,
					}
				} else if valueIRI.IsAbs() {
					result.BaseURL = valueIRI
					result.BaseURLValue = valueString
				} else if result.BaseURL != nil {

					// [spec // 4.1.2 // 5.7.4] Otherwise, if *value* is a relative IRI reference and the base IRI of *result* is not `null`, set the base IRI of *result* to the result of resolving *value* against the current base IRI of *result*.

					result.BaseURL = result.BaseURL.ResolveReference(valueIRI)
					result.BaseURLValue = valueString
				} else {

					// [spec // 4.1.2 // 5.7.5] Otherwise, an `invalid base IRI` error has been detected and processing is aborted.

					return nil, jsonldtype.Error{
						Code: jsonldtype.InvalidBaseIRI,
						Err:  errors.New("relative IRI reference without base IRI"),
					}
				}
			} else {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidBaseIRI, // per 5.7.5
					Err:  fmt.Errorf("invalid type: %s", value.GetGrammarName()),
				}
			}
		}

		// [spec // 4.1.2 // 5.8] If *context* has an `@vocab` entry:

		if vocabObjectMember, ok := contextObject.Members["@vocab"]; ok {

			// [spec // 4.1.2 // 5.8.1] Initialize *value* to the value associated with the `@vocab` entry.

			value := vocabObjectMember.Value

			// [spec // 4.1.2 // 5.8.2] If *value* is `null`, remove any vocabulary mapping from `result`.
			if _, ok := value.(inspectjson.NullValue); ok {
				result.VocabularyMapping = nil
				result.VocabularyMappingValue = nil

				// [spec // 4.1.2 // 5.8.3] Otherwise, if *value* is an IRI or blank node identifier, the vocabulary mapping of *result* is set to the result of IRI expanding *value* using `true` for *document relative* . If it is not an IRI, or a blank node identifier, an `invalid vocab mapping` error has been detected and processing is aborted.
				// [spec // 4.1.2 // 5.8.3] NOTE The use of blank node identifiers to value for @vocab is obsolete, and may be removed in a future version of JSON-LD.
			} else if valueString, ok := value.(inspectjson.StringValue); ok {
				// [dpb] additional validation fixes #t0115
				if result._processor.processingMode == ProcessingMode_JSON_LD_1_0 && !strings.HasPrefix(valueString.Value, "_:") {
					valueIRI, err := url.Parse(valueString.Value)
					if err != nil || !valueIRI.IsAbs() {
						return nil, jsonldtype.Error{
							Code: jsonldtype.InvalidVocabMapping,
							Err:  errors.New("@vocab must be an absolute IRI in JSON-LD 1.0"),
						}
					}
				}

				expandedVocab, err := algorithmIRIExpansion{
					value:            valueString,
					activeContext:    result,
					documentRelative: true,
					// assumed by t0125
					vocab: true,
				}.Call()
				if err != nil {
					return nil, err
				}

				switch t := expandedVocab.(type) {
				case ExpandedIRIasNil:
					result.VocabularyMapping = nil
					result.VocabularyMappingValue = nil
				case ExpandedIRIasBlankNode:
					// TODO warn obsolete

					result.VocabularyMapping = t
					result.VocabularyMappingValue = valueString
				case ExpandedIRIasIRI:
					result.VocabularyMapping = t
					result.VocabularyMappingValue = valueString
				default:
					return nil, jsonldtype.Error{
						Code: jsonldtype.InvalidVocabMapping,
						Err:  fmt.Errorf("invalid expanded type: %s", expandedVocab.ExpandedType()),
					}
				}
			} else {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidVocabMapping,
					Err:  fmt.Errorf("invalid type: %s", vocabObjectMember.Value.GetGrammarName()),
				}
			}
		}

		// [spec // 4.1.2 // 5.9] If *context* has an `@language` entry:

		if languageObjectMember, ok := contextObject.Members["@language"]; ok {

			// [spec // 4.1.2 // 5.9.1] Initialize *value* to the value associated with the `@language` entry.

			value := languageObjectMember.Value

			// [spec // 4.1.2 // 5.9.2] If *value* is `null`, remove any default language from *result*.

			if _, ok := value.(inspectjson.NullValue); ok {
				result.DefaultLanguageValue = nil

				// [spec // 4.1.2 // 5.9.3] Otherwise, if value is a string, the default language of result is set to value. If it is not a string, an invalid default language error has been detected and processing is aborted. If value is not well-formed according to section 2.2.9 of [BCP47], processors SHOULD issue a warning.
			} else if valueString, ok := value.(inspectjson.StringValue); ok {

				// TODO should validate and issue warning for not well-formed

				result.DefaultLanguageValue = &valueString
			} else {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidDefaultLanguage, // per 5.9.3
					Err:  fmt.Errorf("invalid type: %s", languageObjectMember.Value.GetGrammarName()),
				}
			}
		}

		// [spec // 4.1.2 // 5.10] If *context* has an `@direction` entry:

		if directionObjectMember, ok := contextObject.Members["@direction"]; ok {

			// [spec // 4.1.2 // 5.10.1] If processing mode is `json-ld-1.0`, an `invalid context entry` error has been detected and processing is aborted.

			if result._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidContextEntry,
					Err:  fmt.Errorf("invalid entry (processing mode %s): @direction", result._processor.processingMode),
				}
			}

			// [spec // 4.1.2 // 5.10.2] Initialize *value* to the value associated with the `@direction` entry.

			value := directionObjectMember.Value

			// [spec // 4.1.2 // 5.10.3] If *value* is `null`, remove any base direction from *result*.

			if _, ok := value.(inspectjson.NullValue); ok {
				result.DefaultDirectionValue = nil
			} else if valueString, ok := value.(inspectjson.StringValue); ok {

				// [spec // 4.1.2 // 5.10.4] Otherwise, if *value* is a string, the base direction of *result* is set to *value*. If it is not `null`, `"ltr"`, or `"rtl"`, an `invalid base direction` error has been detected and processing is aborted.

				if valueString.Value != "ltr" && valueString.Value != "rtl" {
					return nil, jsonldtype.Error{
						Code: jsonldtype.InvalidBaseDirection,
						Err:  fmt.Errorf("invalid value: %s", valueString.Value),
					}
				}

				result.DefaultDirectionValue = &valueString
			} else {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidBaseDirection, // per 5.10.4
					Err:  fmt.Errorf("invalid type: %s", directionObjectMember.Value.GetGrammarName()),
				}
			}
		}

		// [spec // 4.1.2 // 5.11] If *context* has an `@propagate` entry:

		if propagateObjectMember, ok := contextObject.Members["@propagate"]; ok {

			// [spec // 4.1.2 // 5.11.1] If processing mode is `json-ld-1.0`, an `invalid context entry` error has been detected and processing is aborted.

			if result._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidContextEntry,
					Err:  fmt.Errorf("invalid entry (processing mode %s): @propagate", result._processor.processingMode),
				}
			}

			// [spec // 4.1.2 // 5.11.2] Otherwise, if the value of `@propagate` is not boolean `true` or `false`, an `invalid @propagate value` error has been detected and processing is aborted.

			if _, ok := propagateObjectMember.Value.(inspectjson.BooleanValue); !ok {
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidAtPropagateValue,
					Err:  fmt.Errorf("invalid type: %s", propagateObjectMember.Value.GetGrammarName()),
				}
			}

			// [spec // 4.1.2 // 5.11.2] The previous context is actually set earlier in this algorithm; the previous two steps exist for error checking only.
		}

		// [spec // 4.1.2 // 5.12] Create a map *defined* to keep track of whether or not a term has already been defined or is currently being defined during recursion.

		defined := map[string]bool{}

		// [spec // 4.1.2 // 5.13] For each *key*-*value* pair in *context* where *key* is not `@base`, `@direction`, `@import`, `@language`, `@propagate`, `@protected`, `@version`, or `@vocab`, invoke the Create Term Definition algorithm, passing *result* for *active context*, *context* for *local context*, *key*, *defined*, *base URL*, the value of the `@protected` entry from *context*, if any, for *protected*, *override protected*, and a copy of *remote contexts*.

		var contextProtected bool

		if protectedObjectMember, ok := contextObject.Members["@protected"]; ok {
			if protectedBoolean, ok := protectedObjectMember.Value.(inspectjson.BooleanValue); ok && protectedBoolean.Value {
				contextProtected = true
			}
		}

		contextKeys := slices.Collect(maps.Keys(contextObject.Members))
		slices.SortFunc(contextKeys, strings.Compare)

		for _, key := range contextKeys {
			switch key {
			case "@base", "@direction", "@import", "@language", "@propagate", "@protected", "@version", "@vocab":
				continue
			}

			err := algorithmCreateTermDefinition{
				activeContext:     result,
				localContext:      &contextObject,
				term:              key,
				defined:           defined,
				baseURL:           vars.BaseURL,
				protected:         contextProtected,
				overrideProtected: vars.OverrideProtected,
				remoteContexts:    vars.RemoteContexts[:],
			}.Call()
			if err != nil {
				return nil, err
			}
		}
	}

	// [spec // 4.1.2 // 6] Return result.

	return result, nil
}
