package jsonldinternal

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
)

type algorithmExpansion struct {
	activeContext  *Context
	activeProperty *string

	// [spec] element to be expanded
	element inspectjson.Value

	// Optional
	// [spec] If not passed, the optional flags are set to false.

	// [spec] base URL associated with the documentUrl of the original document to expand
	baseURL *url.URL

	// [spec] frameExpansion flag allowing special forms of input used for frame expansion
	// [dpb] not implemented
	// frameExpansion bool

	// [spec] used to order map entry keys lexicographically, where noted
	ordered bool

	// [spec] used to control reverting previous term definitions in the active context associated with non-propagated contexts
	fromMap bool

	// [dpb] for tracking source offsets
	activePropertySourceOffsets *cursorio.TextOffsetRange
}

func (vars algorithmExpansion) Call() (ExpandedValue, error) {

	// [spec // 5.1.2 // 1] If element is `null`, return `null`.

	if _, ok := vars.element.(inspectjson.NullValue); ok {
		return nil, nil
	}

	// [spec // 5.1.2 // 2] If *active property* is `@default`, initialize the `frameExpansion` flag to `false`.

	if false {
		// [dpb] frameExpansion not implemented; always false

		if vars.activeProperty != nil && *vars.activeProperty == "@default" {
			// vars.frameExpansion = false
		}
	}

	// [spec // 5.1.2 // 3] If *active property* has a term definition in *active context* with a local context, initialize *property-scoped context* to that local context.

	var propertyScopedContext inspectjson.Value

	var termDefinition *TermDefinition

	if vars.activeProperty != nil {
		termDefinition = vars.activeContext.TermDefinitions[*vars.activeProperty]
	}

	if termDefinition != nil && termDefinition.Context != nil {
		propertyScopedContext = termDefinition.Context
	}

	// [spec // 5.1.2 // 4] If *element* is a scalar,

	// TODO dpb-revisit
	if elementGrammarName := vars.element.GetGrammarName(); elementGrammarName == "boolean" || elementGrammarName == "string" || elementGrammarName == "number" || elementGrammarName == "null" {

		// [spec // 5.1.2 // 4.1] If *active property* is `null` or `@graph`, drop the free-floating scalar by returning `null`.

		if vars.activeProperty == nil || *vars.activeProperty == "@graph" {
			return nil, nil
		}

		// [spec // 5.1.2 // 4.2] If *property-scoped context* is defined, set *active context* to the result of the Context Processing algorithm, passing *active context*, *property-scoped context* as *local context*, and *base URL* from the term definition for *active property* in *active context*.

		if propertyScopedContext != nil {
			activeContext, err := algorithmContextProcessing{
				ActiveContext: vars.activeContext,
				LocalContext:  propertyScopedContext,
				BaseURL: coalesceBaseURL(
					vars.activeContext.TermDefinitions[*vars.activeProperty].BaseURL,
					vars.activeContext.BaseURL,
				),
				// defaults
				RemoteContexts:        nil,
				OverrideProtected:     false,
				Propagate:             true,
				ValidateScopedContext: true,
			}.Call()
			if err != nil {
				return nil, err
			}

			vars.activeContext = activeContext
		}

		// [spec // 5.1.2 // 4.3] Return the result of the Value Expansion algorithm, passing the *active context*, *active property*, and *element* as *value*.

		return algorithmValueExpansion{
			activeContext:               vars.activeContext,
			activeProperty:              *vars.activeProperty,
			activePropertySourceOffsets: vars.activePropertySourceOffsets,
			value:                       vars.element,
			//
			processor: vars.activeContext._processor,
		}.Call(), nil
	}

	// [spec // 5.1.2 // 5] If *element* is an array,

	if elementArray, ok := vars.element.(inspectjson.ArrayValue); ok {

		// [spec // 5.1.2 // 5.1] Initialize an empty array, *result*.

		result := &ExpandedArray{}

		// [spec // 5.1.2 // 5.2] For each *item* in *element*:

		for _, elementItem := range elementArray.Values {

			// [spec // 5.1.2 // 5.2.1] Initialize *expanded item* to the result of using this algorithm recursively, passing *active context*, *active property*, *item* as *element*, *base URL*, the `frameExpansion` `ordered`, and *from map* flags.

			expandedItem, err := algorithmExpansion{
				activeContext:               vars.activeContext,
				activeProperty:              vars.activeProperty,
				activePropertySourceOffsets: vars.activePropertySourceOffsets,
				element:                     elementItem,
				baseURL:                     vars.baseURL,
				// frameExpansion: vars.frameExpansion,
				ordered: vars.ordered,
				fromMap: vars.fromMap,
			}.Call()
			if err != nil {
				return nil, err
			}

			// [spec // 5.1.2 // 5.2.2] If the container mapping of *active property* includes `@list`, and *expanded item* is an array, set *expanded item* to a new map containing the entry `@list` where the value is the original *expanded item*.

			if expandedItemArray, ok := expandedItem.(*ExpandedArray); ok && termDefinition != nil && len(termDefinition.ContainerMapping) > 0 {

				for _, containerMapping := range termDefinition.ContainerMapping {
					if containerMapping == "@list" {
						expandedItem = &ExpandedObject{
							Members: map[string]ExpandedValue{
								"@list": expandedItemArray,
							},
							PropertySourceOffsets: vars.activePropertySourceOffsets,
						}

						break
					}
				}
			}

			// [spec // 5.1.2 // 5.2.3] If *expanded item* is an array, append each of its items to *result*. Otherwise, if expanded item is not null, append it to result.

			if expandedItemArray, ok := expandedItem.(*ExpandedArray); ok {
				result.Values = append(result.Values, expandedItemArray.Values...)
			} else if expandedItem == nil {
				// skip
			} else {
				result.Values = append(result.Values, expandedItem)
			}
		}

		// [spec // 5.1.2 // 5.3] Return *result*.

		return result, nil
	}

	// [spec // 5.1.2 // 6] Otherwise *element* is a map.

	elementObject := vars.element.(inspectjson.ObjectValue)

	// [spec // 5.1.2 // 7] If *active context* has a previous context, the *active context* is not propagated. If *from map* is undefined or `false`, and *element* does not contain an entry expanding to `@value`, and *element* does not consist of a single entry expanding to `@id` (where entries are IRI expanded, set *active context* to previous context from *active context*, as the scope of a term-scoped context does not apply when processing new node objects.

	if vars.activeContext.PreviousContext != nil {

		if !vars.fromMap {

			var revert = true

			if len(elementObject.Members) == 1 {
				for _, member := range elementObject.Members {
					expandedKey := algorithmIRIExpansion{
						activeContext: vars.activeContext.PreviousContext,
						value:         member.Name,
						vocab:         true,
					}.Call()

					if keyKeyword, ok := expandedKey.(ExpandedIRIasKeyword); ok && keyKeyword == "@id" {
						revert = false
					}
				}
			}

			if revert {
				var expandedKeywords = map[string]struct{}{}

				for _, member := range elementObject.Members {
					expandedKey := algorithmIRIExpansion{
						activeContext: vars.activeContext,
						value:         member.Name,
						vocab:         true,
					}.Call()

					if keyKeyword, ok := expandedKey.(ExpandedIRIasKeyword); ok {
						expandedKeywords[string(keyKeyword)] = struct{}{}
					}
				}

				if _, ok := expandedKeywords["@value"]; ok {
					revert = false
				}
			}

			if revert {
				vars.activeContext = vars.activeContext.PreviousContext
			}
		}
	}

	// [spec // 5.1.2 // 8] If *property-scoped context* is defined, set *active context* to the result of the Context Processing algorithm, passing *active context*, *property-scoped context* as *local context*, *base URL* from the term definition for *active property*, in *active context* and `true` for *override protected*.

	if propertyScopedContext != nil {
		activeContext, err := algorithmContextProcessing{
			ActiveContext:     vars.activeContext,
			LocalContext:      propertyScopedContext,
			BaseURL:           coalesceBaseURL(termDefinition.BaseURL, vars.activeContext.BaseURL),
			OverrideProtected: true,
			// defaults
			RemoteContexts:        nil,
			Propagate:             true,
			ValidateScopedContext: true,
		}.Call()
		if err != nil {
			return nil, err
		}

		vars.activeContext = activeContext
	}

	// [spec // 5.1.2 // 9] If *element* contains the entry `@context`, set *active context* to the result of the Context Processing algorithm, passing *active context*, the value of the `@context` entry as *local context* and *base URL*.

	if contextEntry, ok := elementObject.Members["@context"]; ok {
		activeContext, err := algorithmContextProcessing{
			ActiveContext: vars.activeContext,
			LocalContext:  contextEntry.Value,
			BaseURL:       vars.baseURL,
			// defaults
			RemoteContexts:        nil,
			OverrideProtected:     false,
			Propagate:             true,
			ValidateScopedContext: true,
		}.Call()
		if err != nil {
			return nil, err
		}

		vars.activeContext = activeContext
	}

	// [spec // 5.1.2 // 10] Initialize *type-scoped context* to *active context*. This is used for expanding values that may be relevant to any previous type-scoped context.

	typeScopedContext := vars.activeContext

	// [spec // 5.1.2 // 11] For each *key* and *value* in *element* ordered lexicographically by *key* where *key* IRI expands to `@type`:

	var orderedElementKeys []string
	var expandedElementKeys = map[string]ExpandedIRI{}

	for key, member := range elementObject.Members {
		orderedElementKeys = append(orderedElementKeys, key)

		expandedElementKeys[key] = algorithmIRIExpansion{
			activeContext: vars.activeContext,
			value:         member.Name,
			// assumed: not explicitly documented
			vocab: true,
		}.Call()
	}

	slices.SortFunc(orderedElementKeys, strings.Compare)

	for _, key := range orderedElementKeys {
		if keyKeyword, ok := expandedElementKeys[key].(ExpandedIRIasKeyword); !ok || keyKeyword != "@type" {
			continue
		}

		// [spec // 5.1.2 // 11.1] Convert *value* into an array, if necessary.

		var valueArray []inspectjson.Value

		switch value := elementObject.Members[key].Value.(type) {
		case inspectjson.ArrayValue:
			valueArray = value.Values
		default:
			valueArray = []inspectjson.Value{value}
		}

		// [spec // 5.1.2 // 11.2] For each *term* which is a value of *value* ordered lexicographically, if *term* is a string, and *term*'s term definition in *type-scoped context* has a local context, set *active context* to the result Context Processing algorithm, passing *active context*, the value of the *term*'s local context as *local context*, *base URL* from the term definition for *value* in *active context*, and `false` for *propagate*.

		var orderedValueIdxs []int

		for termIdx, term := range valueArray {
			termString, ok := term.(inspectjson.StringValue)
			if !ok {
				continue
			}

			termDefinition, ok := typeScopedContext.TermDefinitions[termString.Value]
			if !ok || termDefinition.Context == nil {
				continue
			}

			orderedValueIdxs = append(orderedValueIdxs, termIdx)
		}

		slices.SortFunc(orderedValueIdxs, func(i, j int) int {
			return strings.Compare(
				valueArray[orderedValueIdxs[i]].(inspectjson.StringValue).Value,
				valueArray[orderedValueIdxs[j]].(inspectjson.StringValue).Value,
			)
		})

		for _, idx := range orderedValueIdxs {
			term := valueArray[idx].(inspectjson.StringValue)
			termDefinition := typeScopedContext.TermDefinitions[term.Value]

			nextActiveContext, err := algorithmContextProcessing{
				ActiveContext: vars.activeContext,
				LocalContext:  termDefinition.Context,
				BaseURL:       coalesceBaseURL(termDefinition.BaseURL, vars.activeContext.BaseURL),
				Propagate:     false,
				// defaults
				RemoteContexts:        nil,
				OverrideProtected:     false,
				ValidateScopedContext: true,
			}.Call()
			if err != nil {
				return nil, err
			}

			vars.activeContext = nextActiveContext
		}
	}

	// [spec // 5.1.2 // 12] Initialize two empty maps, *result* and *nests*. Initialize *input type* to expansion of the last value of the first entry in *element* expanding to `@type` (if any), ordering entries lexicographically by key. Both the key and value of the matched entry are IRI expanded.

	var resultObject = &ExpandedObject{
		Members:               map[string]ExpandedValue{},
		SourceOffsets:         vars.element.GetSourceOffsets(),
		PropertySourceOffsets: vars.activePropertySourceOffsets,
	}

	// [dpb] nests moved inside recursive function

	var inputType ExpandedIRI

	for _, key := range orderedElementKeys {
		if key == "@context" {
			continue
		} else if (algorithmIRIExpansion{
			activeContext: vars.activeContext,
			value:         elementObject.Members[key].Name,
			vocab:         true,
		}.Call()) != ExpandedIRIasKeyword("@type") {
			continue
		}

		value := elementObject.Members[key].Value

		if valueString, ok := value.(inspectjson.StringValue); ok {
			inputType = algorithmIRIExpansion{
				activeContext: vars.activeContext,
				value:         valueString,
				vocab:         true,
			}.Call()
		}

		break
	}

	// [dpb] scoped function to support recursion from [spec 14]
	// TODO this func-based recursion of just the two steps ended up seeming insufficient; remove?
	steps_13_14 := func(elementObject inspectjson.ObjectValue, orderedElementKeys []string) error {

		// [dpb] could probably be a simple map
		var nests = map[string]inspectjson.ObjectMember{}

		// [spec // 5.1.2 // 13] For each *key* and *value* in *element*, ordered lexicographically by *key* if `ordered` is `true`:
		// TODO always doing ordered? maps ordered anyway?

		for _, key := range orderedElementKeys {
			value := elementObject.Members[key].Value

			// [spec // 5.1.2 // 13.1] If *key* is `@context`, continue to the next *key*.
			if key == "@context" {
				continue
			}

			// [spec // 5.1.2 // 13.2] Initialize *expanded property* to the result of IRI expanding *key*.

			expandedProperty := algorithmIRIExpansion{
				activeContext: vars.activeContext,
				value:         elementObject.Members[key].Name,
				// assumed: not explicitly documented
				vocab: true,
			}.Call()

			// [spec // 5.1.2 // 13.3] If *expanded property* is `null` or it neither contains a colon (`:`) nor it is a keyword, drop *key* by continuing to the next *key*.

			switch t := expandedProperty.(type) {
			case ExpandedIRIasNil:
				continue
			case ExpandedIRIasBlankNode:
				// always contains a colon?
			case ExpandedIRIasKeyword:
				// useful
			case ExpandedIRIasIRI:
				if !strings.Contains(string(t), ":") {
					continue
				}
			default:
				return fmt.Errorf("unexpected expanded property type: %T", expandedProperty)
			}

			// [dpb] scope vars

			var expandedValue ExpandedValue

			// [spec // 5.1.2 // 13.4] If *expanded property* is a keyword:

			if expandedPropertyKeyword, ok := expandedProperty.(ExpandedIRIasKeyword); ok {

				// [spec // 5.1.2 // 13.4.1] If *active property* equals `@reverse`, an `invalid reverse property map` error has been detected and processing is aborted.

				if vars.activeProperty != nil && *vars.activeProperty == "@reverse" {
					return jsonldtype.Error{
						Code: jsonldtype.InvalidReversePropertyMap,
					}
				}

				// [spec // 5.1.2 // 13.4.2] If *result* already has an *expanded property* entry, other than `@included` or `@type` (unless processing mode is `json-ld-1.0`), a `colliding keywords` error has been detected and processing is aborted.

				if vars.activeContext._processor.processingMode != ProcessingMode_JSON_LD_1_0 {
					if expandedPropertyKeyword == "@included" || expandedPropertyKeyword == "@type" {
						// ignore
					} else if _, ok := resultObject.Members[string(expandedPropertyKeyword)]; ok {
						return jsonldtype.Error{
							Code: jsonldtype.CollidingKeywords,
							Err:  fmt.Errorf("invalid entry: %s", expandedPropertyKeyword),
						}
					}
				}

				// [spec // 5.1.2 // 13.4.3] If *expanded property* is `@id`:

				switch expandedPropertyKeyword {
				case "@id":

					// [spec // 5.1.2 // 13.4.3.1] If *value* is not a string, an `invalid @id value` error has been detected and processing is aborted. When the `frameExpansion` flag is set, *value* *MAY* be an empty map, or an array of one or more strings.

					valueString, ok := value.(inspectjson.StringValue)
					if !ok {
						return jsonldtype.Error{
							Code: jsonldtype.InvalidAtIDValue,
							Err:  fmt.Errorf("invalid type: %s", value.GetGrammarName()),
						}
					}

					// [spec // 5.1.2 // 13.4.3.2] Otherwise, set *expanded value* to the result of IRI expanding *value* using `true` for *document relative* and `false` for *vocab*. When the `frameExpansion` flag is set, *expanded value* will be an array of one or more of the values, with string values expanded using the IRI Expansion algorithm as above.

					expandedValue = algorithmIRIExpansion{
						value:            valueString,
						documentRelative: true,
						vocab:            false,
						// implicit
						activeContext: vars.activeContext,
					}.Call().NewPropertyValue(
						elementObject.Members[key].Name.SourceOffsets,
						valueString.SourceOffsets,
					)

				// [spec // 5.1.2 // 13.4.4] If *expanded property* is `@type`:

				case "@type":

					// [spec // 5.1.2 // 13.4.4.1] If *value* is neither a string nor an array of strings, an `invalid type value` error has been detected and processing is aborted. When the `frameExpansion` flag is set, *value* *MAY* be an empty map, or a default object where the value of `@default` is restricted to be an IRI. All other values mean that `invalid type value` error has been detected and processing is aborted.

					var expandTypeAsArray bool
					var valueArray []inspectjson.Value

					if valueArrayValue, ok := value.(inspectjson.ArrayValue); ok {
						expandTypeAsArray = true
						valueArray = valueArrayValue.Values
					} else {
						valueArray = []inspectjson.Value{value}
					}

					// type validation handled later with [spec 13.4.4.4]

					if false {
						// [dpb] unimplemented (only frameExpansion allows for a map)
						// [spec // 5.1.2 // 13.4.4.2] If *value* is an empty map, set *expanded value* to *value*.
						// [spec // 5.1.2 // 13.4.4.3] Otherwise, if *value* is a default object, set *expanded value* to a new default object with the value of `@default` set to the result of IRI expanding value using *type-scoped context* for *active context*, and `true` for *document relative*.
					}

					// [spec // 5.1.2 // 13.4.4.4] Otherwise, set *expanded value* to the result of IRI expanding each of its values using *type-scoped context* for *active context*, and `true` for *document relative*.

					var expandedValueArray []ExpandedValue

					for _, value := range valueArray {
						valueString, ok := value.(inspectjson.StringValue)
						if !ok {
							return jsonldtype.Error{
								Code: jsonldtype.InvalidTypeValue,
								Err:  fmt.Errorf("invalid type: %s", value.GetGrammarName()),
							}
						}

						expandedValueArray = append(expandedValueArray, algorithmIRIExpansion{
							value:            valueString,
							documentRelative: true,
							// assumed although not explicitly mentioned
							vocab: true,
							// implicit
							activeContext: typeScopedContext,
						}.Call().NewPropertyValue(
							elementObject.Members[key].Name.SourceOffsets,
							valueString.SourceOffsets,
						))
					}

					// [spec // 5.1.2 // 13.4.4.5] If *result* already has an entry for `@type`, prepend the value of `@type` in *result* to *expanded value*, transforming it into an array, if necessary.
					// [spec // 5.1.2 // 13.4.4.5] NOTE No transformation from a string *value* to an array *expanded value* is implied, and the form or *value* should be preserved in *expanded value*.

					if existingType, ok := resultObject.Members["@type"]; ok {
						expandTypeAsArray = true
						existingTypeArray, ok := existingType.(*ExpandedArray)
						if ok {
							expandedValueArray = append(existingTypeArray.Values, expandedValueArray...)
						} else {
							expandedValueArray = append([]ExpandedValue{existingType}, expandedValueArray...)
						}
					}

					if expandTypeAsArray {
						expandedValue = &ExpandedArray{
							Values: expandedValueArray,
						}
					} else {
						expandedValue = expandedValueArray[0]
					}

				// [spec // 5.1.2 // 13.4.5] If expanded property is @graph, set expanded value to the result of using this algorithm recursively passing active context, @graph for active property, value for element, base URL, and the frameExpansion and ordered flags, ensuring that expanded value is an array of one or more maps.

				case "@graph":

					var err error

					expandedValue, err = algorithmExpansion{
						activeContext:               vars.activeContext,
						activeProperty:              &key,
						activePropertySourceOffsets: elementObject.Members[key].Name.SourceOffsets,
						element:                     value,
						baseURL:                     vars.baseURL,
						// frameExpansion: vars.frameExpansion,
						ordered: vars.ordered,
					}.Call()
					if err != nil {
						return err
					}

					if _, ok := expandedValue.(*ExpandedArray); !ok {
						expandedValue = &ExpandedArray{
							Values: []ExpandedValue{expandedValue},
						}
					}

				// [spec // 5.1.2 // 13.4.6] If *expanded property* is `@included`:

				case "@included":

					// [spec // 5.1.2 // 13.4.6.1] If processing mode is `json-ld-1.0`, continue with the next *key* from *element*.

					if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
						continue
					}

					// [spec // 5.1.2 // 13.4.6.2] Set *expanded value* to the result of using this algorithm recursively passing *active context*, `null` for *active property*, *value* for *element*, *base URL*, and the `frameExpansion` and `ordered` flags, ensuring that the result is an array.

					var err error

					expandedValue, err = algorithmExpansion{
						activeContext:               vars.activeContext,
						activeProperty:              nil,
						activePropertySourceOffsets: nil,
						element:                     value,
						baseURL:                     vars.baseURL,
						// frameExpansion: vars.frameExpansion,
						ordered: vars.ordered,
					}.Call()
					if err != nil {
						return err
					}

					if _, ok := expandedValue.(*ExpandedArray); !ok {
						expandedValue = &ExpandedArray{
							Values: []ExpandedValue{expandedValue},
						}
					}

					// [spec // 5.1.2 // 13.4.6.3] If any element of *expanded value* is not a node object, an `invalid @included value` error has been detected and processing is aborted.

					for _, expandedValueItem := range expandedValue.(*ExpandedArray).Values {
						itemObject, ok := expandedValueItem.(*ExpandedObject)
						if !ok {
							return jsonldtype.Error{
								Code: jsonldtype.InvalidAtIncludedValue,
							}
						}

						var isInvalid string

						if _, ok := itemObject.Members["@value"]; ok {
							isInvalid = "@value"
						} else if _, ok := itemObject.Members["@list"]; ok {
							isInvalid = "@list"
						} else if _, ok := itemObject.Members["@set"]; ok {
							isInvalid = "@set"
						}

						if isInvalid != "" {
							return jsonldtype.Error{
								Code: jsonldtype.InvalidAtIncludedValue,
								Err:  fmt.Errorf("expected node object, but found keyword: %s", isInvalid),
							}
						}
					}

					// [spec // 5.1.2 // 13.4.6.4] If *result* already has an entry for `@included`, prepend the value of `@included` in *result* to *expanded value*.

					if existingIncluded, ok := resultObject.Members["@included"]; ok {
						expandedValueArray := expandedValue.(*ExpandedArray)
						expandedValueArray.Values = append(
							existingIncluded.(*ExpandedArray).Values,
							expandedValueArray.Values...,
						)

						expandedValue = expandedValueArray
					}

				// [spec // 5.1.2 // 13.4.7] If *expanded property* is `@value`:

				case "@value":

					// [spec // 5.1.2 // 13.4.7.1] If *input type* is `@json`, set *expanded value* to *value*. If processing mode is `json-ld-1.0`, an `invalid value object value` error has been detected and processing is aborted.

					if inputType == ExpandedIRIasKeyword("@json") {
						if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
							return jsonldtype.Error{
								Code: jsonldtype.InvalidValueObjectValue,
								Err:  fmt.Errorf("invalid keyword (processing mode %s): @json", vars.activeContext._processor.processingMode),
							}
						}

						expandedValue = &ExpandedScalarPrimitive{
							Value: value,
						}
					} else {

						// [spec // 5.1.2 // 13.4.7.2] Otherwise, if *value* is not a scalar or `null`, an `invalid value object value` error has been detected and processing is aborted. When the `frameExpansion` flag is set, *value* *MAY* be an empty map or an array of scalar values.

						if elementGrammarName := value.GetGrammarName(); elementGrammarName != "boolean" && elementGrammarName != "string" && elementGrammarName != "number" && elementGrammarName != "null" {
							return jsonldtype.Error{
								Code: jsonldtype.InvalidValueObjectValue,
								Err:  fmt.Errorf("invalid type: %s", elementGrammarName),
							}
						}

						// [spec // 5.1.2 // 13.4.7.3] Otherwise, set *expanded value* to *value*. When the `frameExpansion` flag is set, *expanded value* will be an array of one or more string values or an array containing an empty map.

						expandedValue = &ExpandedScalarPrimitive{
							Value: value,
						}

						// [spec // 5.1.2 // 13.4.7.4] If *expanded value* is `null`, set the `@value` entry of *result* to `null` and continue with the next *key* from *element*. Null values need to be preserved in this case as the meaning of an `@type` entry depends on the existence of an `@value` entry.

						if _, ok := expandedValue.(*ExpandedScalarPrimitive).Value.(inspectjson.NullValue); ok {
							resultObject.Members["@value"] = expandedValue

							continue
						}
					}

				// [spec // 5.1.2 // 13.4.8] If *expanded property* is `@language`:

				case "@language":

					// [spec // 5.1.2 // 13.4.8.1] If *value* is not a string, an `invalid language-tagged string` error has been detected and processing is aborted. When the frameExpansion flag is set, *value* *MAY* be an empty map or an array of zero or more strings.

					if _, ok := value.(inspectjson.StringValue); !ok {
						return jsonldtype.Error{
							Code: jsonldtype.InvalidLanguageTaggedString,
							Err:  fmt.Errorf("invalid type: %s", value.GetGrammarName()),
						}
					}

					// [spec // 5.1.2 // 13.4.8.2] Otherwise, set *expanded value* to *value*. If *value* is not well-formed according to section 2.2.9 of [BCP47], processors *SHOULD* issue a warning. When the `frameExpansion` flag is set, *expanded value* will be an array of one or more string values or an array containing an empty map.
					// [spec // 5.1.2 // 13.4.8.2] NOTE Processors MAY normalize language tags to lower case.

					// TODO warning

					expandedValue = &ExpandedScalarPrimitive{
						Value: value,
					}

				// [spec // 5.1.2 // 13.4.9] If *expanded property* is `@direction`:

				case "@direction":

					// [spec // 5.1.2 // 13.4.9.1] If processing mode is `json-ld-1.0`, continue with the next *key* from *element*.

					if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
						continue
					}

					// [spec // 5.1.2 // 13.4.9.2] If *value* is neither `"ltr"` nor `"rtl"`, an invalid base direction error has been detected and processing is aborted. When the `frameExpansion` flag is set, *value* MAY be an empty map or an array of zero or more strings.

					valueString, ok := value.(inspectjson.StringValue)
					if !ok {
						return jsonldtype.Error{
							Code: jsonldtype.InvalidBaseDirection,
							Err:  fmt.Errorf("invalid type: %s", value.GetGrammarName()),
						}
					} else if valueString.Value != "ltr" && valueString.Value != "rtl" {
						return jsonldtype.Error{
							Code: jsonldtype.InvalidBaseDirection,
							Err:  fmt.Errorf("invalid value: %s", valueString.Value),
						}
					}

					// [spec // 5.1.2 // 13.4.9.3] Otherwise, set *expanded value* to *value*. When the `frameExpansion` flag is set, *expanded value* will be an array of one or more string values or an array containing an empty map.

					expandedValue = &ExpandedScalarPrimitive{
						Value: value,
					}

				// [spec // 5.1.2 // 13.4.10] If *expanded property* is `@index`:

				case "@index":

					// [spec // 5.1.2 // 13.4.10.1] If *value* is not a string, an `invalid @index value` error has been detected and processing is aborted.

					if _, ok := value.(inspectjson.StringValue); !ok {
						return jsonldtype.Error{
							Code: jsonldtype.InvalidAtIndexValue,
							Err:  fmt.Errorf("invalid type: %s", value.GetGrammarName()),
						}
					}

					// [spec // 5.1.2 // 13.4.10.2] Otherwise, set *expanded value* to *value*.

					expandedValue = &ExpandedScalarPrimitive{
						Value: value,
					}

				// [spec // 5.1.2 // 13.4.11] If *expanded property* is `@list`:

				case "@list":

					// [spec // 5.1.2 // 13.4.11.1] If *active property* is `null` or `@graph`, continue with the next *key* from *element* to remove the free-floating list.

					if vars.activeProperty == nil || *vars.activeProperty == "@graph" {
						continue
					}

					// [spec // 5.1.2 // 13.4.11.2] Otherwise, initialize *expanded value* to the result of using this algorithm recursively passing *active context*, *active property*, *value* for *element*, *base URL*, and the `frameExpansion` and `ordered` flags, ensuring that the result is an array..

					var err error

					expandedValue, err = algorithmExpansion{
						activeContext:               vars.activeContext,
						activeProperty:              vars.activeProperty,
						activePropertySourceOffsets: vars.activePropertySourceOffsets,
						element:                     value,
						baseURL:                     vars.baseURL,
						// frameExpansion: vars.frameExpansion,
						ordered: vars.ordered,
					}.Call()
					if err != nil {
						return err
					}

					if _, ok := expandedValue.(*ExpandedArray); !ok {
						expandedValue = &ExpandedArray{
							Values: []ExpandedValue{expandedValue},
						}
					}

				// [spec // 5.1.2 // 13.4.12] If *expanded property* is `@set`, set *expanded value* to the result of using this algorithm recursively, passing *active context*, *active property*, *value* for *element*, *base URL*, and the `frameExpansion` and `ordered` flags.

				case "@set":

					var err error

					expandedValue, err = algorithmExpansion{
						activeContext:               vars.activeContext,
						activeProperty:              vars.activeProperty,
						activePropertySourceOffsets: vars.activePropertySourceOffsets,
						element:                     value,
						baseURL:                     vars.baseURL,
						// frameExpansion: vars.frameExpansion,
						ordered: vars.ordered,
					}.Call()
					if err != nil {
						return err
					}

				// [spec // 5.1.2 // 13.4.13] If *expanded property* is `@reverse`:

				case "@reverse":

					// [spec // 5.1.2 // 13.4.13.1] If *value* is not a map, an `invalid @reverse value` error has been detected and processing is aborted.

					valueObject, ok := value.(inspectjson.ObjectValue)
					if !ok {
						return jsonldtype.Error{
							Code: jsonldtype.InvalidAtReverseValue,
							Err:  fmt.Errorf("invalid type: %s", value.GetGrammarName()),
						}
					}

					// [spec // 5.1.2 // 13.4.13.2] Otherwise initialize *expanded value* to the result of using this algorithm recursively, passing *active context*, `@reverse` as *active property*, *value* as *element*, *base URL*, and the `frameExpansion` and `ordered flags`.

					var err error

					expandedValue, err = algorithmExpansion{
						activeContext:               vars.activeContext,
						activeProperty:              &key,
						activePropertySourceOffsets: elementObject.Members[key].Name.SourceOffsets,
						element:                     valueObject,
						baseURL:                     vars.baseURL,
						// frameExpansion: vars.frameExpansion,
						ordered: vars.ordered,
					}.Call()
					if err != nil {
						return err
					}

					// [spec // 5.1.2 // 13.4.13.3] If *expanded value* contains an `@reverse` entry, i.e., properties that are reversed twice, execute for each of its *property* and *item* the following steps:

					expandedValueObject := expandedValue.(*ExpandedObject)
					hasAtReverse := false

					if atReverseMember, ok := expandedValueObject.Members["@reverse"]; ok {

						hasAtReverse = true

						for property, item := range atReverseMember.(*ExpandedObject).Members {

							// [spec // 5.1.2 // 13.4.13.3.1] Use add value to add *item* to the *property* entry in *result* using `true` for *as array*.

							macroAddValue{
								Object:  resultObject,
								Key:     property,
								Value:   item,
								AsArray: true,
							}.Call()

						}

					}

					// [spec // 5.1.2 // 13.4.13.4] If expanded value contains an entry other than @reverse:

					if !hasAtReverse || (hasAtReverse && len(expandedValueObject.Members) > 1) {

						// [spec // 5.1.2 // 13.4.13.4.1] Set *reverse map* to the value of the `@reverse` entry in *result*, initializing it to an empty map, if necessary.

						reverseMapMember, ok := resultObject.Members["@reverse"]
						if !ok {
							reverseMapMember = &ExpandedObject{
								Members: map[string]ExpandedValue{},
							}

							resultObject.Members["@reverse"] = reverseMapMember
						}

						reverseMap := reverseMapMember.(*ExpandedObject)

						// [spec // 5.1.2 // 13.4.13.4.2] For each *property* and *items* in *expanded value* other than `@reverse`:

						for property, item := range expandedValueObject.Members {
							if property == "@reverse" {
								continue
							}

							// [spec // 5.1.2 // 13.4.13.4.2.1] For each *item* in *items*:

							for _, item := range item.(*ExpandedArray).Values {

								// [spec // 5.1.2 // 13.4.13.4.2.1.1] If *item* is a value object or list object, an `invalid reverse property value` has been detected and processing is aborted.

								if itemObject, ok := item.(*ExpandedObject); ok {
									if _, ok := itemObject.Members["@value"]; ok {
										return jsonldtype.Error{
											Code: jsonldtype.InvalidReversePropertyValue,
											Err:  fmt.Errorf("invalid entry key: @value"),
										}
									} else if _, ok := itemObject.Members["@list"]; ok {
										return jsonldtype.Error{
											Code: jsonldtype.InvalidReversePropertyValue,
											Err:  fmt.Errorf("invalid entry key: @list"),
										}
									}
								}

								// [spec // 5.1.2 // 13.4.13.4.2.1.2] Use add value to add *item* to the *property* entry in *reverse map* using `true` for *as array*.

								macroAddValue{
									Object:  reverseMap,
									Key:     property,
									Value:   item,
									AsArray: true,
								}.Call()

							}
						}
					}

					// [spec // 5.1.2 // 13.4.13.5] Continue with the next key from element.

					continue

				// [spec // 5.1.2 // 13.4.14] If *expanded property* is `@nest`, add *key* to *nests*, initializing it to an empty array, if necessary. Continue with the next *key* from *element*.

				case "@nest":

					if _, ok := nests[key]; !ok {
						nests[key] = inspectjson.ObjectMember{
							Name: inspectjson.StringValue{
								Value: key,
							},
							Value: inspectjson.ArrayValue{},
						}
					}

					// [dpb] spec does not explicitly state to add value
					// [dpb] spec later assumes values of nest's key are objects, so need to expand arrays

					valueArray, ok := value.(inspectjson.ArrayValue)
					if !ok {
						valueArray = inspectjson.ArrayValue{
							Values: []inspectjson.Value{value},
						}
					}

					nestsKeyValueArray := nests[key].Value.(inspectjson.ArrayValue)
					nestsKeyValueArray.Values = append(
						nestsKeyValueArray.Values,
						valueArray.Values...,
					)

					nestsKey := nests[key]
					nestsKey.Value = nestsKeyValueArray

					nests[key] = nestsKey

					continue

				// [spec // 5.1.2 // 13.4.15] When the `frameExpansion` flag is set, if *expanded property* is any other framing keyword (`@default`, `@embed`, `@explicit`, `@omitDefault`, or `@requireAll`), set *expanded value* to the result of performing the Expansion Algorithm recursively, passing *active context*, *active property*, *value* for *element*, *base URL*, and the `frameExpansion` and `ordered flags`.

				case "@default", "@embed", "@explicit", "@omitDefault", "@requireAll":

					// [dpb] unimplemented (frameExpansion)

				}

				// [spec // 5.1.2 // 13.4.16] Unless *expanded value* is `null`, *expanded property* is `@value`, and *input type* is not `@json`, set the *expanded property* entry of *result* to *expanded value*.

				if !(expandedValue == nil && expandedPropertyKeyword == "@value" && inputType != ExpandedIRIasKeyword("@json")) {
					resultObject.Members[string(expandedPropertyKeyword)] = expandedValue
				}

				// [spec // 5.1.2 // 13.4.17] Continue with the next *key* from *element*.

				continue
			}

			// [spec // 5.1.2 // 13.5] Initialize *container mapping* to *key*'s container mapping in *active context*.

			keyTermDefinition := vars.activeContext.TermDefinitions[key]

			var containerMapping []string

			if keyTermDefinition != nil {
				containerMapping = keyTermDefinition.ContainerMapping
			}

			// [spec // 5.1.2 // 13.6] If *key*'s term definition in *active context* has a type mapping of `@json`, set *expanded value* to a new map, set the entry `@value` to *value*, and set the entry `@type` to `@json`.

			if keyTermDefinition != nil && keyTermDefinition.TypeMapping == ExpandedIRIasKeyword("@json") {
				expandedValue = &ExpandedObject{
					Members: map[string]ExpandedValue{
						"@value": &ExpandedScalarPrimitive{
							Value: value,
						},
						"@type": &ExpandedScalarPrimitive{
							Value: tokenStringJson,
						},
					},
				}
			} else {

				// [spec // 5.1.2 // 13.7] Otherwise, if *container mapping* includes `@language` and *value* is a map then *value* is expanded from a language map as follows:

				if valueObject, ok := value.(inspectjson.ObjectValue); slices.Contains(containerMapping, "@language") && ok {

					// [spec // 5.1.2 // 13.7.1] Initialize *expanded value* to an empty array.

					expandedValue = &ExpandedArray{}

					// [spec // 5.1.2 // 13.7.2] Initialize *direction* to the default base direction from *active context*.

					var direction inspectjson.Value = vars.activeContext.DefaultDirectionValue

					// [spec // 5.1.2 // 13.7.3] If *key*'s term definition in active context has a direction mapping, update *direction* with that value.

					if keyTermDefinition != nil && keyTermDefinition.DirectionMappingValue != nil {
						direction = keyTermDefinition.DirectionMappingValue
					}

					// [spec // 5.1.2 // 13.7.4/ For each key-value pair *language*-*language value* in *value*, ordered lexicographically by *language* if `ordered` is `true`:

					var orderedLanguageKeys []string

					for language := range valueObject.Members {
						orderedLanguageKeys = append(orderedLanguageKeys, language)
					}

					if vars.ordered {
						slices.SortFunc(orderedLanguageKeys, strings.Compare)
					}

					for _, language := range orderedLanguageKeys {

						languageValue := valueObject.Members[language].Value

						// [spec // 5.1.2 // 13.7.4.1] If *language value* is not an array set *language value* to an array containing only *language value*.

						languageValueArray, ok := languageValue.(inspectjson.ArrayValue)
						if !ok {
							languageValueArray = inspectjson.ArrayValue{
								Values: []inspectjson.Value{languageValue},
							}
						}

						// [spec // 5.1.2 // 13.7.4.2] For each *item* in *language value*:

						for _, item := range languageValueArray.Values {

							// [spec // 5.1.2 // 13.7.4.2.1] If *item* is `null`, continue to the next entry in *language value*.

							if _, ok := item.(inspectjson.NullValue); ok {
								continue
							}

							// [spec // 5.1.2 // 13.7.4.2.2] *item* must be a string, otherwise an `invalid language map value` error has been detected and processing is aborted.

							if _, ok := item.(inspectjson.StringValue); !ok {
								return jsonldtype.Error{
									Code: jsonldtype.InvalidLanguageMapValue,
									Err:  fmt.Errorf("invalid type: %s", item.GetGrammarName()),
								}
							}

							// [spec // 5.1.2 // 13.7.4.2.3] Initialize a new map v consisting of two key-value pairs: (`@value`-*item*) and (`@language`-*language*). If *item* is neither `@none` nor well-formed according to section 2.2.9 of [BCP47], processors SHOULD issue a warning.
							// [spec // 5.1.2 // 13.7.4.2.3] NOTE Processors MAY normalize language tags to lower case.

							v := &ExpandedObject{
								Members: map[string]ExpandedValue{
									"@value": &ExpandedScalarPrimitive{
										Value: item,
									},
									"@language": &ExpandedScalarPrimitive{
										Value: inspectjson.StringValue{
											Value: language,
										},
									},
								},
							}

							// TODO warning

							// [spec // 5.1.2 // 13.7.4.2.4] If *language* is `@none`, or expands to `@none`, remove `@language` from v.

							if language == "@none" {
								delete(v.Members, "@language")
							} else if (algorithmIRIExpansion{
								value:         valueObject.Members[language].Name,
								activeContext: vars.activeContext,
							}).Call() == ExpandedIRIasKeyword("@none") {
								delete(v.Members, "@language")
							}

							// [spec // 5.1.2 // 13.7.4.2.5] If *direction* is not `null`, add an entry for `@direction` to *v* with *direction*.

							if _, ok := direction.(inspectjson.StringValue); ok {
								v.Members["@direction"] = &ExpandedScalarPrimitive{
									Value: direction,
								}
							}

							// [spec // 5.1.2 // 13.7.4.2.6] Append *v* to *expanded value*.

							expandedValueArray := expandedValue.(*ExpandedArray)
							expandedValueArray.Values = append(
								expandedValueArray.Values,
								v,
							)

							expandedValue = expandedValueArray
						}
					}
				} else {

					// [spec // 5.1.2 // 13.8] Otherwise, if *container mapping* includes `@index`, `@type`, or `@id` and *value* is a map then *value* is expanded from an map as follows:

					if valueObject, ok := value.(inspectjson.ObjectValue); (slices.Contains(containerMapping, "@index") || slices.Contains(containerMapping, "@type") || slices.Contains(containerMapping, "@id")) && ok {

						// [spec // 5.1.2 // 13.8.1] Initialize *expanded value* to an empty array.

						expandedValue = &ExpandedArray{}

						// [spec // 5.1.2 // 13.8.2] Initialize *index key* to the *key*'s index mapping in active context, or `@index`, if it does not exist.

						var indexKey ExpandedIRI
						var indexKeySourceOffsets *cursorio.TextOffsetRange

						if keyTermDefinition != nil && keyTermDefinition.IndexMapping != nil {
							indexKey = keyTermDefinition.IndexMapping
							indexKeySourceOffsets = keyTermDefinition.IndexMappingSourceOffsets
						} else {
							indexKey = ExpandedIRIasKeyword("@index")
							// indexKeySourceOffsets remains nil for default @index keyword
						}

						// [spec // 5.1.2 // 13.8.3] For each key-value pair *index*-*index value* in *value*, ordered lexicographically by *index* if `ordered` is `true`:

						var orderedIndexKeys []string

						for index := range valueObject.Members {
							orderedIndexKeys = append(orderedIndexKeys, index)
						}

						if vars.ordered {
							slices.SortFunc(orderedIndexKeys, strings.Compare)
						}

						for _, index := range orderedIndexKeys {

							indexValue := valueObject.Members[index].Value

							// [spec // 5.1.2 // 13.8.3.1] If *container mapping* includes `@id` or `@type`, initialize *map context* to the previous context from *active context* if it exists, otherwise, set *map context* to *active context*.

							var mapContext *Context

							if (slices.Contains(containerMapping, "@id") || slices.Contains(containerMapping, "@type")) && vars.activeContext.PreviousContext != nil {
								mapContext = vars.activeContext.PreviousContext
							} else {
								mapContext = vars.activeContext
							}

							// [spec // 5.1.2 // 13.8.3.2] If *container mapping* includes `@type` and *index*'s term definition in *map context* has a local context, update *map context* to the result of the Context Processing algorithm, passing *map context* as *active context* the value of the *index*'s local context as *local context* and *base URL* from the term definition for *index* in *map context*.

							indexTermDefinition := mapContext.TermDefinitions[index]

							if slices.Contains(containerMapping, "@type") && indexTermDefinition != nil && indexTermDefinition.Context != nil {
								var err error

								mapContext, err = algorithmContextProcessing{
									ActiveContext: mapContext,
									LocalContext:  indexTermDefinition.Context,
									BaseURL:       indexTermDefinition.BaseURL,
									// defaults
									RemoteContexts:        nil,
									OverrideProtected:     false,
									Propagate:             true,
									ValidateScopedContext: true,
								}.Call()
								if err != nil {
									return err
								}
							} else {

								// [spec // 5.1.2 // 13.8.3.3] Otherwise, set map context to active context.

								mapContext = vars.activeContext

							}

							// [spec // 5.1.2 // 13.8.3.4] Initialize *expanded index* to the result of IRI expanding *index*.

							expandedIndex := (algorithmIRIExpansion{
								value:         valueObject.Members[index].Name,
								activeContext: mapContext,
								// implied by #tm006
								vocab: true,
							}).Call()

							// [spec // 5.1.2 // 13.8.3.5] If *index value* is not an array set *index value* to an array containing only *index value*.

							if _, ok := indexValue.(inspectjson.ArrayValue); !ok {
								indexValue = inspectjson.ArrayValue{
									Values: []inspectjson.Value{indexValue},
								}

								{
									_v := valueObject.Members[index]
									_v.Value = indexValue
									valueObject.Members[index] = _v
								}
							}

							// [spec // 5.1.2 // 13.8.3.6] Initialize *index value* to the result of using this algorithm recursively, passing *map context* as *active context*, *key* as *active property*, *index value* as *element*, *base URL*, `true` for *from map*, and the `frameExpansion` and `ordered` flags.

							var err error

							expandedIndexValue, err := algorithmExpansion{
								activeContext:               mapContext,
								activeProperty:              &key,
								activePropertySourceOffsets: elementObject.Members[key].Name.SourceOffsets,
								element:                     indexValue,
								baseURL:                     vars.baseURL,
								fromMap:                     true,
								// frameExpansion: vars.frameExpansion,
								ordered: vars.ordered,
							}.Call()
							if err != nil {
								return err
							}

							indexValueArray := expandedIndexValue.(*ExpandedArray)

							// [spec // 5.1.2 // 13.8.3.7] For each *item* in *index value*:

							for itemIdx, item := range indexValueArray.Values {

								itemObject := item.(*ExpandedObject)

								// [spec // 5.1.2 // 13.8.3.7.1] If *container mapping* includes `@graph`, and *item* is not a graph object, set *item* to a new map containing the key-value pair `@graph`-*item*, ensuring that the value is represented using an array.

								if slices.Contains(containerMapping, "@graph") {
									var needsGraphObject bool

									if _, ok := itemObject.Members["@graph"]; !ok {
										needsGraphObject = true
									}

									if needsGraphObject {
										if _, ok := item.(*ExpandedArray); !ok {
											item = &ExpandedArray{
												Values: []ExpandedValue{item},
											}
										}

										itemObject = &ExpandedObject{
											Members: map[string]ExpandedValue{
												"@graph": item,
											},
										}

										item = itemObject
										indexValueArray.Values[itemIdx] = item
									}
								}

								// [spec // 5.1.2 // 13.8.3.7.2] If *container mapping* includes `@index`, *index key* is not `@index`, and *expanded index* is not `@none`:

								if slices.Contains(containerMapping, "@index") && indexKey != ExpandedIRIasKeyword("@index") && expandedIndex != ExpandedIRIasKeyword("@none") {

									// [spec // 5.1.2 // 13.8.3.7.2.1] Initialize *re-expanded index* to the result of calling the Value Expansion algorithm, passing the *active context*, *index key* as *active property*, and *index* as *value*.

									reExpandedIndex := (algorithmValueExpansion{
										activeContext:               mapContext,
										activeProperty:              indexKey.String(),
										activePropertySourceOffsets: indexKeySourceOffsets,
										value:                       valueObject.Members[index].Name,
									}).Call()

									// [spec // 5.1.2 // 13.8.3.7.2.2] Initialize *expanded index key* to the result of IRI expanding *index key*.

									expandedIndexKey := (algorithmIRIExpansion{
										value:         indexKey.NewPropertyValue(nil, nil).Value,
										activeContext: mapContext,
										// implied by #tpi06
										vocab: true,
									}).Call()

									// [spec // 5.1.2 // 13.8.3.7.2.3] Initialize *index property values* to an array consisting of *re-expanded index* followed by the existing values of the concatenation of *expanded index key* in *item*, if any.

									var indexPropertyValues = []ExpandedValue{
										reExpandedIndex,
									}

									if expandedIndexMember, ok := itemObject.Members[expandedIndexKey.String()]; ok {
										if expandedIndexArray, ok := expandedIndexMember.(*ExpandedArray); ok {
											indexPropertyValues = append(indexPropertyValues, expandedIndexArray.Values...)
										} else {
											indexPropertyValues = append(indexPropertyValues, expandedIndexMember)
										}
									}

									// [spec // 5.1.2 // 13.8.3.7.2.4] Add the key-value pair (*expanded index key*-*index property values*) to item.

									itemObject.Members[expandedIndexKey.String()] = &ExpandedArray{
										Values: indexPropertyValues,
									}

									// [spec // 5.1.2 // 13.8.3.7.2.5] If *item* is a value object, it *MUST NOT* contain any extra properties; an *invalid value object* error has been detected and processing is aborted.

									if _, ok := itemObject.Members["@value"]; ok && len(itemObject.Members) > 1 {
										for key := range itemObject.Members {
											if key != "@value" {
												return jsonldtype.Error{
													Code: jsonldtype.InvalidValueObject,
													Err:  fmt.Errorf("invalid entry: %s", key),
												}
											}
										}
									}
								} else {

									// [spec // 5.1.2 // 13.8.3.7.3] Otherwise, if *container mapping* includes `@index`, *item* does not have an entry `@index`, and *expanded index* is not `@none`, add the key-value pair (`@index`-*index*) to *item*.

									_, hasAtIndex := itemObject.Members["@index"]

									if slices.Contains(containerMapping, "@index") && !hasAtIndex && expandedIndex != ExpandedIRIasKeyword("@none") {
										itemObject.Members["@index"] = &ExpandedScalarPrimitive{
											Value: valueObject.Members[index].Name,
										}
									} else {

										// [spec // 5.1.2 // 13.8.3.7.4] Otherwise, if *container mapping* includes `@id` *item* does not have the entry `@id`, and *expanded index* is not `@none`, add the key-value pair (`@id`-*expanded index*) to *item*, where *expanded index* is set to the result of IRI expanding*index* using `true` for *document relative* and `false` for *vocab*.

										_, hasAtID := itemObject.Members["@id"]

										if slices.Contains(containerMapping, "@id") && !hasAtID && expandedIndex != ExpandedIRIasKeyword("@none") {
											itemObject.Members["@id"] = algorithmIRIExpansion{
												value:            valueObject.Members[index].Name,
												activeContext:    mapContext,
												documentRelative: true,
												vocab:            false,
											}.Call().NewPropertyValue(
												nil,
												valueObject.Members[index].Name.SourceOffsets,
											)
										} else {

											// [spec // 5.1.2 // 13.8.3.7.5] Otherwise, if *container mapping* includes `@type` and *expanded index* is not `@none`, initialize *types* to a new array consisting of *expanded index* followed by any existing values of `@type` in *item*. Add the key-value pair (`@type`-*types*) to *item*.

											if slices.Contains(containerMapping, "@type") && expandedIndex != ExpandedIRIasKeyword("@none") {

												var types = []ExpandedValue{
													expandedIndex.NewPropertyValue(nil, nil),
												}

												if expandedTypeMember, ok := itemObject.Members["@type"]; ok {
													if expandedTypeArray, ok := expandedTypeMember.(*ExpandedArray); ok {
														types = append(types, expandedTypeArray.Values...)
													} else {
														types = append(types, expandedTypeMember)
													}
												}

												itemObject.Members["@type"] = &ExpandedArray{
													Values: types,
												}
											}
										}
									}
								}

								// [spec // 5.1.2 // 13.8.3.7.6] Append *item* to *expanded value*.

								expandedValueArray := expandedValue.(*ExpandedArray)
								expandedValueArray.Values = append(
									expandedValueArray.Values,
									item,
								)

								expandedValue = expandedValueArray
							}
						}
					} else {

						// [spec // 5.1.2 // 13.9] Otherwise, initialize *expanded value* to the result of using this algorithm recursively, passing *active context*, *key* for *active property*, *value* for *element*, *base URL*, and the `frameExpansion` and `ordered flags`.

						var err error

						propertySourceOffsets := elementObject.Members[key].Name.SourceOffsets

						expandedValue, err = algorithmExpansion{
							activeContext:               vars.activeContext,
							activeProperty:              &key,
							activePropertySourceOffsets: propertySourceOffsets,
							element:                     value,
							baseURL:                     vars.baseURL,
							// frameExpansion: vars.frameExpansion,
							ordered: vars.ordered,
						}.Call()
						if err != nil {
							return err
						}

						// [dpb] hacky; ensure propertySourceOffsets get correctly attributed in this case
						if needsListWrap, ok := expandedValue.(*ExpandedArray); containerMapping != nil && slices.Contains(containerMapping, "@list") && ok {
							expandedValue = &ExpandedObject{
								Members: map[string]ExpandedValue{
									"@list": needsListWrap,
								},
								PropertySourceOffsets: propertySourceOffsets,
							}
						}
					}
				}
			}

			// [spec // 5.1.2 // 13.10] If *expanded value* is `null`, ignore *key* by continuing to the next *key* from *element*.

			if expandedValue == nil {
				continue
			} else if valuePrimitive, ok := expandedValue.(*ExpandedScalarPrimitive); ok {
				if _, ok := valuePrimitive.Value.(inspectjson.NullValue); ok {
					continue
				}
			}

			// [spec // 5.1.2 // 13.11] If *container mapping* includes `@list` and *expanded value* is not already a list object, convert *expanded value* to a list object by first setting it to an array containing only *expanded value* if it is not already an array, and then by setting it to a map containing the key-value pair `@list`-*expanded value*.

			if slices.Contains(containerMapping, "@list") {
				var needsListObject bool

				expandedValueObject, ok := expandedValue.(*ExpandedObject)
				if !ok {
					needsListObject = true
				} else if _, ok := expandedValueObject.Members["@list"]; !ok {
					needsListObject = true
				}

				if needsListObject {
					if _, ok := expandedValue.(*ExpandedArray); !ok {
						expandedValue = &ExpandedArray{
							Values: []ExpandedValue{expandedValue},
						}
					}

					expandedValue = &ExpandedObject{
						Members: map[string]ExpandedValue{
							"@list": expandedValue,
						},
						PropertySourceOffsets: vars.activePropertySourceOffsets,
					}
				}
			}

			// [spec // 5.1.2 // 13.12] If *container mapping* includes `@graph`, and includes neither `@id` nor `@index`, convert *expanded value* into an array, if necessary, then convert each value *ev* in *expanded value* into a graph object:

			if slices.Contains(containerMapping, "@graph") && !slices.Contains(containerMapping, "@id") && !slices.Contains(containerMapping, "@index") {

				expandedValueArray, ok := expandedValue.(*ExpandedArray)
				if !ok {
					expandedValueArray = &ExpandedArray{
						Values: []ExpandedValue{expandedValue},
					}

					expandedValue = expandedValueArray
				}

				for evIdx, ev := range expandedValueArray.Values {

					// [spec // 5.1.2 // 13.12.1] Convert *ev* into a graph object by creating a map containing the key-value pair `@graph`-*ev* where *ev* is represented as an array.
					// [spec // 5.1.2 // 13.12.1] This may lead to a graph object including another graph object, if *ev* was already in the form of a graph object.

					expandedValueArray.Values[evIdx] = &ExpandedObject{
						Members: map[string]ExpandedValue{
							"@graph": &ExpandedArray{
								Values: []ExpandedValue{ev},
							},
						},
					}
				}
			}

			// [spec // 5.1.2 // 13.13] If the term definition associated to *key* indicates that it is a reverse property

			if keyTermDefinition != nil && keyTermDefinition.ReverseProperty {

				// [spec // 5.1.2 // 13.13.1] If *result* has no `@reverse` entry, create one and initialize its value to an empty map.

				if _, ok := resultObject.Members["@reverse"]; !ok {
					resultObject.Members["@reverse"] = &ExpandedObject{
						Members: map[string]ExpandedValue{},
					}
				}

				// [spec // 5.1.2 // 13.13.2] Reference the value of the `@reverse` entry in *result* using the variable *reverse map*.

				reverseMap := resultObject.Members["@reverse"].(*ExpandedObject) // TODO unsafe assertion?

				// [spec // 5.1.2 // 13.13.3] If *expanded value* is not an array, set it to an array containing *expanded value*.

				expandedValueArray, ok := expandedValue.(*ExpandedArray)
				if !ok {
					expandedValueArray = &ExpandedArray{
						Values: []ExpandedValue{expandedValue},
					}

					expandedValue = expandedValueArray
				}

				// [spec // 5.1.2 // 13.13.4] For each item in expanded value

				expandedPropertyString := expandedProperty.String()

				for _, item := range expandedValueArray.Values {

					// [spec // 5.1.2 // 13.13.4.1] If *item* is a value object or list object, an `invalid reverse property value` has been detected and processing is aborted.

					if itemObject, ok := item.(*ExpandedObject); ok {
						if _, ok := itemObject.Members["@value"]; ok {
							return jsonldtype.Error{
								Code: jsonldtype.InvalidReversePropertyValue,
								Err:  fmt.Errorf("invalid entry key: @value"),
							}
						} else if _, ok := itemObject.Members["@list"]; ok {
							return jsonldtype.Error{
								Code: jsonldtype.InvalidReversePropertyValue,
								Err:  fmt.Errorf("invalid entry key: @list"),
							}
						}
					}

					// [spec // 5.1.2 // 13.13.4.2] If *reverse map* has no *expanded property* entry, create one and initialize its value to an empty array.

					if _, ok := reverseMap.Members[expandedPropertyString]; !ok {
						reverseMap.Members[expandedPropertyString] = &ExpandedArray{
							Values: []ExpandedValue{},
						}
					}

					// [spec // 5.1.2 // 13.13.4.3] Use add value to add *item* to the *expanded property* entry in *reverse map* using `true` for *as array*.

					macroAddValue{
						Value:   item,
						Key:     expandedPropertyString,
						Object:  reverseMap,
						AsArray: true,
					}.Call()
				}
			} else {

				// [spec // 5.1.2 // 13.14] Otherwise, *key* is not a reverse property use add value to add *expanded value* to the *expanded property* entry in *result* using `true` for *as array*.

				macroAddValue{
					Value:   expandedValue,
					Key:     expandedProperty.String(),
					Object:  resultObject,
					AsArray: true,
				}.Call()
			}
		}

		// [spec // 5.1.2 // 14] For each key *nesting-key* in *nests*, ordered lexicographically if `ordered` is `true`:

		{
			var orderedNestingKeys []string

			for nestingKey := range nests {
				orderedNestingKeys = append(orderedNestingKeys, nestingKey)
			}

			if vars.ordered {
				slices.SortFunc(orderedNestingKeys, strings.Compare)
			}

			for _, nestingKey := range orderedNestingKeys {

				// [spec // 5.1.2 // 14.1] Initialize *nested values* to the value of *nesting-key* in *element*, ensuring that it is an array.

				var nestedValues []inspectjson.Value

				if nestingValue, ok := nests[nestingKey].Value.(inspectjson.ArrayValue); ok {
					nestedValues = nestingValue.Values
				} else {
					nestedValues = []inspectjson.Value{nests[nestingKey].Value}
				}

				// [spec // 5.1.2 // 14.2] For each *nested value* in *nested values*:

				for _, nestedValue := range nestedValues {

					// [spec // 5.1.2 // 14.2.1] If *nested value* is not a map, or any key within *nested value* expands to `@value`, an `invalid @nest value` error has been detected and processing is aborted.

					nestedValueObject, ok := nestedValue.(inspectjson.ObjectValue)
					if !ok {
						return jsonldtype.Error{
							Code: jsonldtype.InvalidAtNestValue,
							Err:  fmt.Errorf("invalid type: %s", nestedValue.GetGrammarName()),
						}
					}

					var orderedNestedValueKeys []string

					for key := range nestedValueObject.Members {
						if (algorithmIRIExpansion{
							value:         nestedValueObject.Members[key].Name,
							activeContext: vars.activeContext,
						}.Call()) == ExpandedIRIasKeyword("@value") {
							return jsonldtype.Error{
								Code: jsonldtype.InvalidAtNestValue,
								Err:  fmt.Errorf("invalid entry key: %s", key),
							}
						}

						orderedNestedValueKeys = append(orderedNestedValueKeys, key)
					}

					if vars.ordered {
						slices.SortFunc(orderedNestedValueKeys, strings.Compare)
					}

					// [spec // 5.1.2 // 14.2.2] Recursively repeat steps 13 and 14 using *nested value* for *element*.
					// [spec // 5.1.2 // 14.2.2] NOTE By invoking steps 13 and 14 on *nested value* we are able to unfold arbitrary levels of nesting, with results being merged into *result*. Step 13 iterates through each entry in *nested value* and expands it, while collecting new *nested values* found at each level, until all nesting has been extracted.

					// [dpb] doing only steps 13 and 14 seems insufficient when nesting key provides its own context which is otherwise not maintained; e.g. #tc037

					expandedNestValue, err := algorithmExpansion{
						activeContext:               vars.activeContext,
						activeProperty:              &nestingKey,
						activePropertySourceOffsets: elementObject.Members[nestingKey].Name.SourceOffsets,
						element:                     nestedValue,
						baseURL:                     vars.baseURL,
						// frameExpansion: vars.frameExpansion,
						ordered: vars.ordered,
					}.Call()
					if err != nil {
						return err
					}

					for key, expandedNestValueItem := range expandedNestValue.(*ExpandedObject).Members {
						macroAddValue{
							Value:  expandedNestValueItem,
							Key:    key,
							Object: resultObject,
							// [dpb] this was previously false, but added @id check for #tin06
							AsArray: key != "@id",
						}.Call()
					}
				}
			}
		}

		return nil
	}
	err := steps_13_14(elementObject, orderedElementKeys)
	if err != nil {
		return nil, err
	}

	// [dpb] done with treating it as an object

	var result ExpandedValue = resultObject

	// [spec // 5.1.2 // 15] If *result* contains the entry `@value`:

	if resultValueMember, ok := resultObject.Members["@value"]; ok {

		// [spec // 5.1.2 // 15.1] The *result* must not contain any entries other than `@direction`, `@index`, `@language`, `@type`, and `@value`. It must not contain an `@type` entry if it contains either `@language` or `@direction` entries. Otherwise, an `invalid value object` error has been detected and processing is aborted.

		var typeValue inspectjson.Value
		var hasTypeJSON bool
		var hasLanguageOrDirection bool

		for key, value := range resultObject.Members {
			switch key {
			case "@direction", "@language":
				hasLanguageOrDirection = true
			case "@index", "@value":
				// nop
			case "@type":
				var badGrammarName inspectjson.GrammarName

				if valuePrimitive, ok := value.(*ExpandedScalarPrimitive); ok {
					switch vT := valuePrimitive.Value.(type) {
					case inspectjson.StringValue:
						if vT.Value == "@json" {
							hasTypeJSON = true
						}

						typeValue = vT
					case inspectjson.NullValue:
						typeValue = vT
					default:
						badGrammarName = valuePrimitive.Value.GetGrammarName()
					}
				} else if _, ok := value.(*ExpandedArray); ok {
					badGrammarName = "array"
				} else if _, ok := value.(*ExpandedObject); ok {
					badGrammarName = "object"
				} else {
					badGrammarName = "unknown"
				}

				if len(badGrammarName) > 0 {
					return nil, jsonldtype.Error{
						Code: jsonldtype.InvalidValueObject,
						Err:  fmt.Errorf("invalid @type type: %s", badGrammarName),
					}
				}

				break
			default:
				return nil, jsonldtype.Error{
					Code: jsonldtype.InvalidValueObject,
					Err:  fmt.Errorf("invalid entry key: %s", key),
				}
			}
		}

		if hasLanguageOrDirection && typeValue != nil {
			return nil, jsonldtype.Error{
				Code: jsonldtype.InvalidValueObject,
				Err:  fmt.Errorf("invalid entry key: %s", "@type"),
			}
		}

		// [spec // 5.1.2 // 15.2] If the *result*'s `@type` entry is `@json`, then the `@value` entry may contain any value, and is treated as a JSON literal.

		if hasTypeJSON {
			// nop
		} else {

			// [spec // 5.1.2 // 15.3] Otherwise, if the value of *result*'s `@value` entry is `null`, or an empty array, return `null`.

			if valuePrimitive, ok := resultValueMember.(*ExpandedScalarPrimitive); ok {
				if _, ok := valuePrimitive.Value.(inspectjson.NullValue); ok {
					return nil, nil
				}
			}

			if memberArray, ok := resultValueMember.(*ExpandedArray); ok && len(memberArray.Values) == 0 {
				return nil, nil
			}

			// [spec // 5.1.2 // 15.4] Otherwise, if the value of *result*'s `@value` entry is not a string and *result* contains the entry `@language`, an `invalid language-tagged value` error has been detected (only strings can be language-tagged) and processing is aborted.

			if valuePrimitive, ok := resultValueMember.(*ExpandedScalarPrimitive); ok {
				if _, ok := valuePrimitive.Value.(inspectjson.StringValue); !ok {
					if _, ok := resultObject.Members["@language"]; ok {
						return nil, jsonldtype.Error{
							Code: jsonldtype.InvalidLanguageTaggedValue,
							Err:  fmt.Errorf("invalid @value type: %s", valuePrimitive.Value.GetGrammarName()),
						}
					}
				}
			}

			// [spec // 5.1.2 // 15.5] Otherwise, if the *result* has an `@type` entry and its value is not an IRI, an `invalid typed value` error has been detected and processing is aborted.

			if typeValue != nil {

				typeValueString, ok := typeValue.(inspectjson.StringValue)
				if !ok {
					return nil, jsonldtype.Error{
						Code: jsonldtype.InvalidTypedValue,
						Err:  fmt.Errorf("invalid @type type: %s", typeValue.GetGrammarName()),
					}
				} else if !isIRI(vars.activeContext._processor.processingMode, typeValueString.Value) {
					return nil, jsonldtype.Error{
						Code: jsonldtype.InvalidTypedValue,
						Err:  fmt.Errorf("invalid @type value: expected iri: found %q", typeValueString.Value),
					}
				}
			}
		}
	} else {

		// [spec // 5.1.2 // 16] Otherwise, if *result* contains the entry `@type` and its associated value is not an array, set it to an array containing only the associated value.

		var hasScalarType bool

		resultType, ok := resultObject.Members["@type"]
		if ok {
			if _, ok := resultType.(*ExpandedArray); !ok {
				hasScalarType = true
			}
		}

		if hasScalarType {
			resultObject.Members["@type"] = &ExpandedArray{
				Values: []ExpandedValue{resultType},
			}
		} else {

			// [spec // 5.1.2 // 17] Otherwise, if *result* contains the entry `@set` or `@list`:

			_, hasSet := resultObject.Members["@set"]
			_, hasList := resultObject.Members["@list"]

			if hasSet || hasList {

				// [spec // 5.1.2 // 17.1] The *result* must contain at most one other entry which must be `@index`. Otherwise, an `invalid set or list object` error has been detected and processing is aborted.

				for key := range resultObject.Members {
					if key != "@set" && key != "@list" && key != "@index" {
						return nil, jsonldtype.Error{
							Code: jsonldtype.InvalidSetOrListObject,
							Err:  fmt.Errorf("invalid key: %s", key),
						}
					}
				}

				// [spec // 5.1.2 // 17.2] If *result* contains the entry `@set`, then set *result* to the entry's associated value.

				if hasSet {
					// TODO unsafe assertion?
					result = resultObject.Members["@set"]
				}
			}
		}
	}

	// [spec // 5.1.2 // 18] If *result* is a map that contains only the entry `@language`, return `null`.

	if resultObject, ok := result.(*ExpandedObject); ok && len(resultObject.Members) == 1 {
		if _, ok := resultObject.Members["@language"]; ok {
			return nil, nil
		}
	}

	// [spec // 5.1.2 // 19] If *active property* is `null` or `@graph`, drop free-floating values as follows:

	if vars.activeProperty == nil || *vars.activeProperty == "@graph" {

		// [spec // 5.1.2 // 19.1] If *result* is a map which is empty, or contains only the entries `@value` or `@list`, set result to null.
		// [dpb] testsuites define that `@value` + related language/type should be dropped, too; not strictly "only the entry `@value`" but closer to exclusively a value object?

		var isEmptyishMap bool

		if resultObject, ok := result.(*ExpandedObject); ok {
			if len(resultObject.Members) == 0 {
				isEmptyishMap = true
			} else {
				isEmptyishMap = true

				var hasValue, hasType, hasLanguage, hasDirection bool

				for key := range resultObject.Members {
					switch key {
					case "@value":
						hasValue = true
					case "@list":
						// nop
					case "@type":
						hasType = true
					case "@language":
						hasLanguage = true
					case "@direction":
						hasDirection = true
					default:
						isEmptyishMap = false
					}

					if !isEmptyishMap {
						break
					}
				}

				if isEmptyishMap {
					if hasType || hasLanguage || hasDirection {
						if !hasValue {
							// not sure this is actually a valid state, but for the sake of minimal divergence
							isEmptyishMap = false
						}
					}
				}
			}
		}

		if isEmptyishMap {
			result = nil
		} else {

			// [spec // 5.1.2 // 19.2] Otherwise, if *result* is a map whose only entry is @id, set result to null. When the frameExpansion flag is set, a map containing only the @id entry is retained.

			if resultObject, ok := result.(*ExpandedObject); ok && len(resultObject.Members) == 1 {
				if _, ok := resultObject.Members["@id"]; ok {
					result = nil
				}
			}
		}
	}

	return result, nil
}
