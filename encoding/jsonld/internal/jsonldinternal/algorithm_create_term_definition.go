package jsonldinternal

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
)

type algorithmCreateTermDefinition struct {
	activeContext *Context
	localContext  *inspectjson.ObjectValue
	term          string
	defined       map[string]bool

	// [spec // 4.2.2] *base URL* defaulting to `null``
	baseURL *url.URL

	// [spec // 4.2.2] protected which defaults to `false`
	protected bool

	// [spec // 4.2.2] *override protected*, defaulting to `false`, which is used to allow changes to protected terms
	overrideProtected bool

	// [spec // 4.2.2] defaulting to a new empty array, which is used to detect cyclical context inclusions
	remoteContexts []string

	// [spec // 4.2.2] defaulting to `true`, which is used to limit recursion when validating possibly recursive scoped contexts..
	// [dpb] spec typo double period
	validateScopedContext bool
}

func (vars algorithmCreateTermDefinition) Call() error {

	// [spec // 4.2.2 // 1] If *defined* contains the entry *term* and the associated value is `true` (indicating that the term definition has already been created), return. Otherwise, if the value is `false`, a `cyclic IRI mapping` error has been detected and processing is aborted.

	if defined, ok := vars.defined[vars.term]; ok {
		if defined {
			return nil
		}

		return jsonldtype.Error{
			Code: jsonldtype.CyclicIRIMapping,
		}
	}

	// [spec // 4.2.2 // 2] If *term* is the empty string (`""`), an `invalid term definition` error has been detected and processing is aborted. Otherwise, set the value associated with *defined*'s *term* entry to `false`. This indicates that the term definition is now being created but is not yet complete.

	vars.defined[vars.term] = false

	// [spec // 4.2.2 // 3] Initialize *value* to a copy of the value associated with the entry *term* in *local context*.

	valueValue := vars.localContext.Members[vars.term].Value

	// [spec // 4.2.2 // 4] If *term* is `@type`, and processing mode is `json-ld-1.0`, a `keyword redefinition` error has been detected and processing is aborted. At this point, *value* *MUST* be a map with only either or both of the following entries:
	// [spec // 4.2.2 // 4] * An entry for `@container` with value `@set`.
	// [spec // 4.2.2 // 4] * An entry for `@protected`.
	// [spec // 4.2.2 // 4] Any other value means that a `keyword redefinition` error has been detected and processing is aborted.

	if vars.term == "@type" {
		if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
			return jsonldtype.Error{
				Code: jsonldtype.KeywordRedefinition,
				Err:  fmt.Errorf("invalid @type entry (processing mode %s)", vars.activeContext._processor.processingMode),
			}
		}

		valueObject, ok := valueValue.(inspectjson.ObjectValue)
		if !ok {
			return jsonldtype.Error{
				Code: jsonldtype.KeywordRedefinition,
				Err:  fmt.Errorf("invalid type: %s", valueObject.GetGrammarName()),
			}
		}

		var validEntries int

		if containerValue, ok := valueObject.Members["@container"]; ok {
			if containerString, ok := containerValue.Value.(inspectjson.StringValue); !ok || containerString.Value != "@set" {
				return jsonldtype.Error{
					Code: jsonldtype.KeywordRedefinition,
					Err:  fmt.Errorf("invalid @container value: %s", containerValue.Value),
				}
			}

			validEntries++
		}

		if _, ok := valueObject.Members["@protected"]; ok {
			validEntries++
		}

		if validEntries == 0 || len(valueObject.Members) != validEntries {
			return jsonldtype.Error{
				Code: jsonldtype.KeywordRedefinition,
				Err:  errors.New("unexpected entries within @type"),
			}
		}
	} else {

		// [spec // 4.2.2 // 5] Otherwise, since keywords cannot be overridden, *term* *MUST NOT* be a keyword and a `keyword redefinition` error has been detected and processing is aborted. If *term* has the form of a keyword (i.e., it matches the ABNF rule `"@"1*ALPHA` from [RFC5234]), return; processors *SHOULD* generate a warning.

		if _, known := definedKeywords[vars.term]; known {
			return jsonldtype.Error{
				Code: jsonldtype.KeywordRedefinition,
				Err:  errors.New(vars.term),
			}
		} else if reKeywordABNF.MatchString(vars.term) {
			// TODO warning

			return nil
		}
	}

	// [spec // 4.2.2 // 6] Initialize *previous definition* to any existing term definition for *term* in *active context*, removing that term definition from *active context*.

	previousDefinition := vars.activeContext.TermDefinitions[vars.term]
	delete(vars.activeContext.TermDefinitions, vars.term)

	//

	var simpleTerm bool
	var valueObject inspectjson.ObjectValue

	// [spec // 4.2.2 // 7] If *value* is `null`, convert it to a map consisting of a single entry whose key is `@id` and whose value is `null`.

	if _, ok := valueValue.(inspectjson.NullValue); ok {
		valueObject = inspectjson.ObjectValue{
			Members: map[string]inspectjson.ObjectMember{
				"@id": {
					Name:  tokenStringId,
					Value: valueValue,
				},
			},
		}

		// [spec // 4.2.2 // 8] Otherwise, if *value* is a string, convert it to a map consisting of a single entry whose key is `@id` and whose value is *value*. Set *simple term* to `true`.
	} else if _, ok := valueValue.(inspectjson.StringValue); ok {
		valueObject = inspectjson.ObjectValue{
			Members: map[string]inspectjson.ObjectMember{
				"@id": {
					Name:  tokenStringId,
					Value: valueValue,
				},
			},
		}

		simpleTerm = true
	} else {

		// [spec // 4.2.2 // 9] Otherwise, *value* *MUST* be a map, if not, an `invalid term definition` error has been detected and processing is aborted. Set *simple term* to `false`.

		valueObjectAssertion, ok := valueValue.(inspectjson.ObjectValue)
		if !ok {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidTermDefinition,
				Err:  fmt.Errorf("invalid type: %s", valueValue.GetGrammarName()),
			}
		}

		valueObject = valueObjectAssertion
	}

	// [spec // 4.2.2 // 10] Create a new term definition, *definition*, initializing prefix flag to `false`, protected to *protected*, and reverse property to `false`.

	definition := &TermDefinition{
		Prefix:          false,
		Protected:       vars.protected,
		ReverseProperty: false,
	}

	// [spec // 4.2.2 // 11] If *value* has an `@protected` entry, set the protected flag in *definition* to the value of this entry. If the value of `@protected` is not a boolean, an `invalid @protected value` error has been detected and processing is aborted. If processing mode is `json-ld-1.0`, an `invalid term definition` has been detected and processing is aborted.

	if protectedValue, ok := valueObject.Members["@protected"]; ok {
		if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidTermDefinition,
				Err:  fmt.Errorf("invalid @protected entry (processing mode %s)", vars.activeContext._processor.processingMode),
			}
		}

		protectedBoolean, ok := protectedValue.Value.(inspectjson.BooleanValue)
		if ok {
			definition.Protected = protectedBoolean.Value
		} else {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidAtProtectedValue,
				Err:  fmt.Errorf("invalid type: %s", protectedValue.Value.GetGrammarName()),
			}
		}
	}

	// [spec // 4.2.2 // 12] If *value* contains the entry `@type`:

	if typeMember, ok := valueObject.Members["@type"]; ok {

		// [spec // 4.2.2 // 12.1] Initialize *type* to the value associated with the `@type` entry, which *MUST* be a string. Otherwise, an `invalid type mapping` error has been detected and processing is aborted.

		typeString, ok := typeMember.Value.(inspectjson.StringValue)
		if !ok {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidTypeMapping,
				Err:  fmt.Errorf("invalid type: %s", typeMember.Value.GetGrammarName()),
			}
		}

		// [spec // 4.2.2 // 12.2] Set *type* to the result of IRI expanding *type*, using *local context*, and *defined*.

		expandedType := algorithmIRIExpansion{
			value:        typeString,
			localContext: vars.localContext,
			defined:      vars.defined,
			// assumed
			vocab:         true,
			activeContext: vars.activeContext,
		}.Call()

		// [spec // 4.2.2 // 12.3] If the expanded *type* is `@json` or `@none`, and processing mode is `json-ld-1.0`, an `invalid type mapping` error has been detected and processing is aborted.

		if expandedTypeKeyword, ok := expandedType.(ExpandedIRIasKeyword); ok {
			if expandedTypeKeyword == "@json" || expandedTypeKeyword == "@none" {
				if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
					return jsonldtype.Error{
						Code: jsonldtype.InvalidTypeMapping,
						Err:  fmt.Errorf("invalid value (processing mode %s): %s", vars.activeContext._processor.processingMode, expandedTypeKeyword),
					}
				}
			}
		}

		// [spec // 4.2.2 // 12.4] Otherwise, if the expanded *type* is neither `@id`, nor `@json`, nor `@none`, nor `@vocab`, nor an IRI, an `invalid type mapping` error has been detected and processing is aborted.

		switch t := expandedType.(type) {
		case ExpandedIRIasKeyword:
			if t != "@id" && t != "@json" && t != "@none" && t != "@vocab" {
				return jsonldtype.Error{
					Code: jsonldtype.InvalidTypeMapping,
					Err:  fmt.Errorf("invalid value: %s", t),
				}
			}
		case ExpandedIRIasIRI:
			// valid
		default:
			return jsonldtype.Error{
				Code: jsonldtype.InvalidTypeMapping,
				Err:  fmt.Errorf("invalid expanded type: %s", expandedType.ExpandedType()),
			}
		}

		// [spec // 4.2.2 // 12.5] Set the type mapping for *definition* to *type*.

		definition.TypeMapping = expandedType
		definition.TypeMappingValue = typeString
	}

	// [spec // 4.2.2 // 13] If *value* contains the entry `@reverse`:

	if reverseMember, ok := valueObject.Members["@reverse"]; ok {

		// [spec // 4.2.2 // 13.1] If *value* contains `@id` or `@nest`, entries, an `invalid reverse property` error has been detected and processing is aborted.

		if _, ok := valueObject.Members["@id"]; ok {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidReverseProperty,
				Err:  errors.New("found @id entry"),
			}
		} else if _, ok := valueObject.Members["@nest"]; ok {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidReverseProperty,
				Err:  errors.New("found @nest entry"),
			}
		}

		// [spec // 4.2.2 // 13.2] If the value associated with the `@reverse` entry is not a string, an `invalid IRI mapping` error has been detected and processing is aborted.

		reverseString, ok := reverseMember.Value.(inspectjson.StringValue)
		if !ok {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidIRIMapping,
				Err:  fmt.Errorf("invalid type: %s", reverseMember.Value.GetGrammarName()),
			}
		}

		// [spec // 4.2.2 // 13.3] If the value associated with the `@reverse` entry is a string having the form of a keyword (i.e., it matches the ABNF rule `"@"1*ALPHA` from [RFC5234]), return; processors *SHOULD* generate a warning.

		if reKeywordABNF.MatchString(reverseString.Value) {
			// TODO warning

			return nil
		}

		// [spec // 4.2.2 // 13.4] Otherwise, set the IRI mapping of *definition* to the result of IRI expanding the value associated with the `@reverse` entry, using *local context*, and *defined*. If the result does not have the form of an IRI or a blank node identifier, an `invalid IRI mapping` error has been detected and processing is aborted.

		expandedReverse := algorithmIRIExpansion{
			value:        reverseString,
			localContext: vars.localContext,
			defined:      vars.defined,
			// assumed
			activeContext: vars.activeContext,
			// vocab
			vocab: true,
		}.Call()

		switch expandedReverse.(type) {
		case ExpandedIRIasIRI:
			// valid
		case ExpandedIRIasBlankNode:
			// valid
		default:
			return jsonldtype.Error{
				Code: jsonldtype.InvalidIRIMapping,
				Err:  fmt.Errorf("invalid expanded type: %s", expandedReverse.ExpandedType()),
			}
		}

		definition.IRI = expandedReverse
		definition.IRIValue = reverseString

		// [spec // 4.2.2 // 13.5] If *value* contains an `@container` entry, set the container mapping of *definition* to an array containing its value; if its value is neither `@set`, nor `@index`, nor `null`, an `invalid reverse property` error has been detected (reverse properties only support set- and index-containers) and processing is aborted.

		if containerMember, ok := valueObject.Members["@container"]; ok {
			if _, ok := containerMember.Value.(inspectjson.NullValue); ok {
				definition.ContainerMapping = nil // TODO test suite; spec implies array with value of null?
			} else if containerString, ok := containerMember.Value.(inspectjson.StringValue); ok {
				if containerString.Value != "@set" && containerString.Value != "@index" {
					return jsonldtype.Error{
						Code: jsonldtype.InvalidReverseProperty,
						Err:  fmt.Errorf("invalid @container value: %s", containerString.Value),
					}
				}

				definition.ContainerMapping = []string{containerString.Value}
			} else {
				return jsonldtype.Error{
					Code: jsonldtype.InvalidReverseProperty,
					Err:  fmt.Errorf("invalid type: %s", containerMember.Value.GetGrammarName()),
				}
			}
		}

		// [spec // 4.2.2 // 13.6] Set the reverse property flag of *definition* to `true`.

		definition.ReverseProperty = true

		// [spec // 4.2.2 // 13.7] Set the term definition of *term* in *active context* to *definition* and the value associated with *defined*'s entry *term* to `true` and return.

		vars.activeContext.TermDefinitions[vars.term] = definition
		vars.defined[vars.term] = true

		// [dpb] testcase #t0131 suggests a sibling @index must be expanded
		// [dpb] this block copied from later where standard @index expansion is handled
		// [dpb] this is not described in the spec, but otherwise the return from 13.7 skips it?

		if indexMember, ok := valueObject.Members["@index"]; ok {

			// [spec // 4.2.2 // 20.2] Initialize *index* to the value associated with the `@index` entry. If the result of IRI expanding that value is not an IRI, an `invalid term definition` has been detected and processing is aborted.

			expandedIndex := algorithmIRIExpansion{
				value: indexMember.Value,
				// spec doesn't describe propagation like other descriptions
				activeContext: vars.activeContext,
				localContext:  vars.localContext,
				defined:       vars.defined,
			}.Call()

			expandedIndexIRI, ok := expandedIndex.(ExpandedIRIasIRI)
			if !ok {
				return jsonldtype.Error{
					Code: jsonldtype.InvalidTermDefinition,
					Err:  fmt.Errorf("invalid expanded type: %s", expandedIndex.ExpandedType()),
				}
			}

			// [spec // 4.2.2 // 20.3] Set the index mapping of *definition* to *index*

			definition.IndexMapping = &expandedIndexIRI
		}

		return nil
	}

	// [spec // 4.2.2 // 14] If *value* contains the entry `@id` and its value does not equal `term`:

	if idMember, ok := valueObject.Members["@id"]; ok {

		// [spec // 4.2.2 // 14.1] If the `@id` entry of *value* is `null`, the term is not used for IRI expansion, but is retained to be able to detect future redefinitions of this term.

		if _, ok := idMember.Value.(inspectjson.NullValue); ok {
			// e.g. #t0028

			// [dpb] this is used in iri expansion to detect non-usage
			definition.IRI = ExpandedIRIasNil{}
			definition.IRIValue = idMember.Value
		} else {
			// [spec // 4.2.2 // 14.2] Otherwise:

			// [spec // 4.2.2 // 14.2.1] If the value associated with the `@id` entry is not a string, an `invalid IRI mapping` error has been detected and processing is aborted.

			idString, ok := idMember.Value.(inspectjson.StringValue)
			if !ok {
				return jsonldtype.Error{
					Code: jsonldtype.InvalidIRIMapping,
					Err:  fmt.Errorf("invalid type: %s", idMember.Value.GetGrammarName()),
				}
			}

			// [spec // 4.2.2 // 14.2.2] If the value associated with the `@id` entry is not a keyword, but has the form of a keyword (i.e., it matches the ABNF rule `"@"1*ALPHA` from [RFC5234]), return; processors *SHOULD* generate a warning.

			if _, ok := definedKeywords[idString.Value]; !ok && reKeywordABNF.MatchString(idString.Value) {
				// TODO warning

				return nil
			}

			// [spec // 4.2.2 // 14.2.3] Otherwise, set the IRI mapping of *definition* to the result of IRI expanding the value associated with the `@id` entry, using *local context*, and *defined*. If the resulting IRI mapping is neither a keyword, nor an IRI, nor a blank node identifier, an `invalid IRI mapping` error has been detected and processing is aborted; if it equals `@context`, an `invalid keyword alias` error has been detected and processing is aborted.

			expandedID := algorithmIRIExpansion{
				value:        idMember.Value,
				localContext: vars.localContext,
				defined:      vars.defined,
				// assumed
				activeContext: vars.activeContext,
				// assumed via testsuites
				vocab: true,
			}.Call()

			switch t := expandedID.(type) {
			case ExpandedIRIasKeyword:
				if t == "@context" {
					return jsonldtype.Error{
						Code: jsonldtype.InvalidKeywordAlias,
						Err:  fmt.Errorf("invalid value: %s", t),
					}
				} else if t == "@type" {
					if !simpleTerm && definition.TypeMapping == ExpandedIRIasKeyword("@id") && vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_1 {
						return jsonldtype.Error{
							Code: jsonldtype.InvalidIRIMapping,
							Err:  fmt.Errorf("invalid value (processing mode %s): %s", vars.activeContext._processor.processingMode, t),
						}
					}
				}
			case ExpandedIRIasIRI:
				// valid
			case ExpandedIRIasBlankNode:
				// valid
			default:
				return jsonldtype.Error{
					Code: jsonldtype.InvalidIRIMapping,
					Err:  fmt.Errorf("invalid expanded type: %s", expandedID.ExpandedType()),
				}
			}

			definition.IRI = expandedID
			definition.IRIValue = idMember.Value

			// [spec // 4.2.2 // 14.2.4] If the *term* contains a colon (`:`) anywhere but as the first or last character of *term*, or if it contains a slash (`/`) anywhere:

			termSubstr := vars.term
			if len(termSubstr) > 1 {
				termSubstr = termSubstr[1 : len(termSubstr)-1]
			}

			hasColonOrSlash := strings.Contains(termSubstr, ":") || strings.Contains(vars.term, "/")

			if hasColonOrSlash {
				// [spec // 4.2.2 // 14.2.4.1] Set the value associated with *defined*'s *term* entry to `true`.

				vars.defined[vars.term] = true

				// [spec // 4.2.2 // 14.2.4.2] If the result of IRI expanding *term* using *local context*, and *defined*, is not the same as the IRI mapping of *definition*, an `invalid IRI mapping` error has been detected and processing is aborted.
				// [dpb] spec refers to expanding *term*, but that is expected to be different; seems the original id value was intended to validate rewrites?

				expandedTerm := algorithmIRIExpansion{
					value:        idMember.Value,
					localContext: vars.localContext,
					defined:      vars.defined,
					// assumed
					activeContext: vars.activeContext,
					// assumed as vocab?
				}.Call()

				if !expandedTerm.Equals(definition.IRI) {
					return jsonldtype.Error{
						Code: jsonldtype.InvalidIRIMapping,
						Err:  fmt.Errorf("expanded term does not match expanded iri (%s != %s)", expandedTerm, definition.IRI),
					}
				}
			}

			// [spec // 4.2.2 // 14.2.5] If *term* contains neither a colon (:) nor a slash (/), simple term is true, and if the IRI mapping of definition is either an IRI ending with a gen-delim character, or a blank node identifier, set the prefix flag in definition to true.

			if hasColonOrSlash || simpleTerm {
				definition.Prefix = true
			} else {
				switch t := definition.IRI.(type) {
				case ExpandedIRIasIRI:
					if len(t) > 0 {
						switch t[len(t)-1] {
						case ':', '/', '?', '#', '[', ']', '@':
							definition.Prefix = true
						}
					}

					definition.Prefix = true
				case ExpandedIRIasBlankNode:
					definition.Prefix = true
				}
			}
		}
	} else {

		// [spec // 4.2.2 // 15] Otherwise if the *term* contains a colon (`:`) anywhere after the first character:

		termPrefixSuffix := strings.SplitN(vars.term, ":", 2)

		if len(termPrefixSuffix) == 2 && len(termPrefixSuffix[0]) > 0 {

			// [spec // 4.2.2 // 15.1] If *term* is a compact IRI with a prefix that is an entry in *local context* a dependency has been found. Use this algorithm recursively passing *active context*, *local context*, the prefix as *term*, and *defined*.

			if _, ok := vars.localContext.Members[termPrefixSuffix[0]]; ok {
				err := algorithmCreateTermDefinition{
					activeContext: vars.activeContext,
					localContext:  vars.localContext,
					term:          termPrefixSuffix[0],
					defined:       vars.defined,
					// defaults
					baseURL:               nil,
					protected:             false,
					overrideProtected:     false,
					remoteContexts:        nil,
					validateScopedContext: true,
				}.Call()
				if err != nil {
					return err
				}
			}

			// [spec // 4.2.2 // 15.2] If *term*'s prefix has a term definition in *active context*, set the IRI mapping of *definition* to the result of concatenating the value associated with the prefix's IRI mapping and the term's *suffix*.

			if prefixDefinition, ok := vars.activeContext.TermDefinitions[termPrefixSuffix[0]]; ok {
				switch t := prefixDefinition.IRI.(type) {
				case ExpandedIRIasIRI:
					definition.IRI = ExpandedIRIasIRI(fmt.Sprintf("%s%s", t, termPrefixSuffix[1]))
				default:
					return fmt.Errorf("unexpected prefix definition iri type: %T", t) // not described in spec
				}
			} else {

				// [spec // 4.2.2 // 15.3] Otherwise, *term* is an IRI or blank node identifier. Set the IRI mapping of *definition* to *term*.

				if strings.HasPrefix(vars.term, "_:") {
					definition.IRI = ExpandedIRIasBlankNode(vars.term)
				} else {
					definition.IRI = ExpandedIRIasIRI(vars.term)
				}

			}
		} else {

			// [spec // 4.2.2 // 16] Otherwise if the *term* contains a slash (`/`):

			if strings.Contains(vars.term, "/") {

				// [spec // 4.2.2 // 16.1] *Term* is a relative IRI reference.

				// [spec // 4.2.2 // 16.2] Set the IRI mapping of *definition* to the result of IRI expanding *term*. If the resulting IRI mapping is not an IRI, an `invalid IRI mapping` error has been detected and processing is aborted.

				expandedTerm := algorithmIRIExpansion{
					value: inspectjson.StringValue{
						Value: vars.term,
					},
					// spec doesn't describe propagation like other descriptions
					activeContext: vars.activeContext,
					localContext:  vars.localContext,
					defined:       vars.defined,
				}.Call()

				if t, ok := expandedTerm.(ExpandedIRIasIRI); ok {
					definition.IRI = t
				} else {
					return jsonldtype.Error{
						Code: jsonldtype.InvalidIRIMapping,
						Err:  fmt.Errorf("invalid expanded type: %s", expandedTerm.ExpandedType()),
					}
				}
			} else {

				// [spec // 4.2.2 // 17] Otherwise, if term is `@type`, set the IRI mapping of *definition* to `@type`.

				if vars.term == "@type" {
					definition.IRI = ExpandedIRIasKeyword("@type")
				} else {

					// [spec // 4.2.2 // 18] Otherwise, if *active context* has a vocabulary mapping, the IRI mapping of *definition* is set to the result of concatenating the value associated with the vocabulary mapping and *term*. If it does not have a vocabulary mapping, an `invalid IRI mapping` error been detected and processing is aborted.

					if vars.activeContext.VocabularyMapping != nil {
						switch t := vars.activeContext.VocabularyMapping.(type) {
						case ExpandedIRIasIRI:
							definition.IRI = ExpandedIRIasIRI(fmt.Sprintf("%s%s", t, vars.term))
						default:
							return fmt.Errorf("unexpected vocab mapping iri type: %T", t) // not described in spec
						}
					} else {
						return jsonldtype.Error{
							Code: jsonldtype.InvalidIRIMapping,
							Err:  errors.New("missing vocabulary mapping"),
						}
					}
				}
			}
		}
	}

	// [spec // 4.2.2 // 19] If *value* contains the entry `@container`:

	if containerMember, ok := valueObject.Members["@container"]; ok {

		// [spec // 4.2.2 // 19.1] Initialize *container* to the value associated with the `@container` entry, which *MUST* be either `@graph`, `@id`, `@index`, `@language`, `@list`, `@set`, `@type`, or an array containing exactly any one of those keywords, an array containing `@graph` and either `@id` or `@index` optionally including `@set`, or an array containing a combination of `@set` and any of `@index`, `@graph`, `@id`, `@type`, `@language` in any order . Otherwise, an `invalid container mapping` has been detected and processing is aborted.
		// [spec // 4.2.2 // 19.2] If the container value is `@graph`, `@id`, or `@type`, or is otherwise not a string, generate an `invalid container mapping` error and abort processing if processing mode is `json-ld-1.0`.

		var containerValues []string

		switch t := containerMember.Value.(type) {
		case inspectjson.StringValue:
			valueString := t.Value

			switch valueString {
			case "@graph", "@id", "@type":
				if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
					return jsonldtype.Error{
						Code: jsonldtype.InvalidContainerMapping,
						Err:  fmt.Errorf("invalid value (processing mode %s): %s", vars.activeContext._processor.processingMode, valueString),
					}
				}

				fallthrough
			case "@index", "@language", "@list", "@set":
				containerValues = append(containerValues, valueString)
			default:
				return jsonldtype.Error{
					Code: jsonldtype.InvalidContainerMapping,
					Err:  fmt.Errorf("invalid value: %s", valueString),
				}
			}
		case inspectjson.ArrayValue:
			if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
				return jsonldtype.Error{
					Code: jsonldtype.InvalidContainerMapping,
					Err:  fmt.Errorf("invalid type (processing mode %s): %s", vars.activeContext._processor.processingMode, t.GetGrammarName()),
				}
			}

			var uniqueContainerValues = map[string]struct{}{}

			for _, containerValue := range t.Values {
				containerString, ok := containerValue.(inspectjson.StringValue)
				if !ok {
					return jsonldtype.Error{
						Code: jsonldtype.InvalidContainerMapping,
						Err:  fmt.Errorf("invalid type: %s", containerValue.GetGrammarName()),
					}
				}

				valueString := containerString.Value

				switch valueString {
				case "@graph", "@id", "@index", "@language", "@list", "@set", "@type":
					uniqueContainerValues[valueString] = struct{}{}
					containerValues = append(containerValues, valueString)
				default:
					return jsonldtype.Error{
						Code: jsonldtype.InvalidContainerMapping,
						Err:  fmt.Errorf("invalid value: %s", valueString),
					}
				}
			}

			switch len(uniqueContainerValues) {
			case 0:
				return jsonldtype.Error{
					Code: jsonldtype.InvalidContainerMapping,
					Err:  errors.New("empty array"),
				}
			case 1:
				// valid
			default:
				if _, ok := uniqueContainerValues["@set"]; ok {
					for container := range uniqueContainerValues {
						switch container {
						case "@set", "@index", "@graph", "@id", "@type", "@language":
							// valid
						default:
							return jsonldtype.Error{
								Code: jsonldtype.InvalidContainerMapping,
								Err:  fmt.Errorf("invalid entry (with @set): %s", container),
							}
						}
					}
				} else if _, ok := uniqueContainerValues["@graph"]; ok {
					_, hasAtID := uniqueContainerValues["@id"]
					_, hasAtIndex := uniqueContainerValues["@index"]

					if hasAtID && hasAtIndex {
						return jsonldtype.Error{
							Code: jsonldtype.InvalidContainerMapping,
							Err:  errors.New("invalid combination (with @graph): @id, @index"),
						}
					} else if !hasAtID && !hasAtIndex {
						return jsonldtype.Error{
							Code: jsonldtype.InvalidContainerMapping,
							Err:  errors.New("missing entry (with @graph): @id or @index"),
						}
					}

					if len(uniqueContainerValues) > 2 {
						for container := range uniqueContainerValues {
							switch container {
							case "@graph", "@id", "@index": // @set is handled in prior statement
								// valid
							default:
								return jsonldtype.Error{
									Code: jsonldtype.InvalidContainerMapping,
									Err:  fmt.Errorf("invalid entry: %s", container),
								}
							}
						}
					}
				}
			}
		default:
			return jsonldtype.Error{
				Code: jsonldtype.InvalidContainerMapping,
				Err:  fmt.Errorf("invalid type: %s", containerMember.Value.GetGrammarName()),
			}
		}

		// [spec // 4.2.2 // 19.3] Set the container mapping of *definition* to *container* coercing to an array, if necessary.

		definition.ContainerMapping = containerValues

		// [spec // 4.2.2 // 19.4] If the container mapping of *definition* includes `@type`:

		if slices.Contains(definition.ContainerMapping, "@type") {

			// [spec // 4.2.2 // 19.4.1] If *type mapping* in *definition* is undefined, set it to `@id`.

			if definition.TypeMapping == nil {
				definition.TypeMapping = ExpandedIRIasKeyword("@id")
			}

			// [spec // 4.2.2 // 19.4.2] If *type mapping* in *definition* is neither `@id` nor `@vocab`, an `invalid type mapping` error has been detected and processing is aborted.

			if definition.TypeMapping != ExpandedIRIasKeyword("@id") && definition.TypeMapping != ExpandedIRIasKeyword("@vocab") {
				return jsonldtype.Error{
					Code: jsonldtype.InvalidTypeMapping,
					Err:  fmt.Errorf("invalid value: %s", definition.TypeMapping),
				}
			}
		}
	}

	// [spec // 4.2.2 // 20] If *value* contains the entry `@index`:

	if indexMember, ok := valueObject.Members["@index"]; ok {

		// [spec // 4.2.2 // 20.1] If processing mode is `json-ld-1.0` or container mapping does not include `@index`, an `invalid term definition` has been detected and processing is aborted.

		if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidTermDefinition,
				Err:  fmt.Errorf("invalid entry (processing mode %s): @index", vars.activeContext._processor.processingMode),
			}
		}

		{
			var foundContainerIndex bool

			for _, container := range definition.ContainerMapping {
				if container == "@index" {
					foundContainerIndex = true

					break
				}
			}

			if !foundContainerIndex {
				return jsonldtype.Error{
					Code: jsonldtype.InvalidTermDefinition,
					Err:  errors.New("missing container mapping value: @index"),
				}
			}
		}

		// [spec // 4.2.2 // 20.2] Initialize *index* to the value associated with the `@index` entry. If the result of IRI expanding that value is not an IRI, an `invalid term definition` has been detected and processing is aborted.

		expandedIndex := algorithmIRIExpansion{
			value: indexMember.Value,
			// spec doesn't describe propagation like other descriptions
			activeContext: vars.activeContext,
			localContext:  vars.localContext,
			defined:       vars.defined,
		}.Call()

		expandedIndexIRI, ok := expandedIndex.(ExpandedIRIasIRI)
		if !ok {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidTermDefinition,
				Err:  fmt.Errorf("invalid expanded type: %s", expandedIndex.ExpandedType()),
			}
		}

		// [spec // 4.2.2 // 20.3] Set the index mapping of *definition* to *index*

		definition.IndexMapping = &expandedIndexIRI
	}

	// [spec // 4.2.2 // 21] If *value* contains the entry `@context`:

	if contextMember, ok := valueObject.Members["@context"]; ok {

		// [spec // 4.2.2 // 21.1] If processing mode is `json-ld-1.0`, an `invalid term definition` has been detected and processing is aborted.

		if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidTermDefinition,
				Err:  fmt.Errorf("invalid entry (processing mode %s): @context", vars.activeContext._processor.processingMode),
			}
		}

		// [spec // 4.2.2 // 21.2] Initialize *context* to the value associated with the `@context` entry, which is treated as a local context.

		contextValue := contextMember.Value

		// [spec // 4.2.2 // 21.3] Invoke the Context Processing algorithm using the *active context*, *context* as *local context*, *base URL*, `true` for *override protected*, a copy of *remote contexts*, and `false` for *validate scoped context*. If any error is detected, an `invalid scoped context` error has been detected and processing is aborted.
		// [spec // 4.2.2 // 21.3] NOTE The result of the Context Processing algorithm is discarded; it is called to detect errors at definition time. If used, the context will be re-processed and applied to the active context as part of expansion or compaction.

		_, err := algorithmContextProcessing{
			ActiveContext:         vars.activeContext,
			LocalContext:          contextValue,
			BaseURL:               vars.baseURL,
			OverrideProtected:     true,
			RemoteContexts:        vars.remoteContexts[:],
			ValidateScopedContext: false,
			// defaults
			Propagate: true,
		}.Call()
		if err != nil {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidScopedContext,
				Err:  err,
			}
		}

		// [spec // 4.2.2 // 21.4] Set the local context of *definition* to *context*, and base URL to *base URL*.

		definition.Context = contextValue
		definition.BaseURL = vars.baseURL
	}

	// [spec // 4.2.2 // 22] If *value* contains the entry `@language` and does not contain the entry `@type`:

	if languageMember, ok := valueObject.Members["@language"]; ok {
		if _, ok := valueObject.Members["@type"]; !ok {

			// [spec // 4.2.2 // 22.1] Initialize *language* to the value associated with the `@language` entry, which MUST be either `null` or a string. If *language* is not well-formed according to section 2.2.9 of [BCP47], processors *SHOULD* issue a warning. Otherwise, an `invalid language mapping` error has been detected and processing is aborted.
			// [spec // 4.2.2 // 22.1] Set the language mapping of *definition* to *language*.

			switch languageType := languageMember.Value.(type) {
			case inspectjson.NullValue:
				definition.LanguageMappingValue = languageType
			case inspectjson.StringValue:
				// TODO well-formed validation warning

				definition.LanguageMappingValue = languageType
			default:
				return jsonldtype.Error{
					Code: jsonldtype.InvalidLanguageMapping,
					Err:  fmt.Errorf("invalid type: %s", languageMember.Value.GetGrammarName()),
				}
			}
		}
	}

	// [spec // 4.2.2 // 23] If *value* contains the entry `@direction` and does not contain the entry `@type`:

	if directionMember, ok := valueObject.Members["@direction"]; ok {

		// [spec // 4.2.2 // 23.1] Initialize *direction* to the value associated with the `@direction` entry, which *MUST* be either `null`, `"ltr"`, or `"rtl"`. Otherwise, an `invalid base direction` error has been detected and processing is aborted.
		// [spec // 4.2.2 // 23.2] Set the direction mapping of *definition* to *direction*.

		switch directionT := directionMember.Value.(type) {
		case inspectjson.NullValue:
			definition.DirectionMappingValue = directionT
		case inspectjson.StringValue:
			if directionT.Value != "ltr" && directionT.Value != "rtl" {
				return jsonldtype.Error{
					Code: jsonldtype.InvalidBaseDirection,
					Err:  fmt.Errorf("invalid value: %s", directionT.Value),
				}
			}

			definition.DirectionMappingValue = directionT
		default:
			return jsonldtype.Error{
				Code: jsonldtype.InvalidBaseDirection,
				Err:  fmt.Errorf("invalid type: %s", directionMember.Value.GetGrammarName()),
			}
		}
	}

	// [spec // 4.2.2 // 24] If *value* contains the entry `@nest`:

	if nestMember, ok := valueObject.Members["@nest"]; ok {

		// [spec // 4.2.2 // 24.1] If processing mode is `json-ld-1.0`, an `invalid term definition` has been detected and processing is aborted.

		if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidTermDefinition,
				Err:  fmt.Errorf("invalid entry (processing mode %s): @nest", vars.activeContext._processor.processingMode),
			}
		}

		// [spec // 4.2.2 // 24.2] Initialize nest value in *definition* to the value associated with the `@nest` entry, which *MUST* be a string and *MUST NOT* be a keyword other than `@nest`. Otherwise, an `invalid @nest value` error has been detected and processing is aborted.

		if nestString, ok := nestMember.Value.(inspectjson.StringValue); ok {
			if reKeywordABNF.MatchString(nestString.Value) && nestString.Value != "@nest" {
				return jsonldtype.Error{
					Code: jsonldtype.InvalidAtNestValue,
					Err:  fmt.Errorf("invalid value: %s", nestString.Value),
				}
			}

			definition.NestValue = &nestString.Value
		} else {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidAtNestValue,
				Err:  fmt.Errorf("invalid type: %s", nestMember.Value.GetGrammarName()),
			}
		}
	}

	// [spec // 4.2.2 // 25] If *value* contains the entry `@prefix`:

	if prefixMember, ok := valueObject.Members["@prefix"]; ok {

		// [spec // 4.2.2 // 25.1] If processing mode is `json-ld-1.0`, or if *term* contains a colon (`:`) or slash (`/`), an `invalid term definition` has been detected and processing is aborted.

		if vars.activeContext._processor.processingMode == ProcessingMode_JSON_LD_1_0 {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidTermDefinition,
				Err:  fmt.Errorf("invalid entry (processing mode %s): @prefix", vars.activeContext._processor.processingMode),
			}
		} else if strings.Contains(vars.term, ":") || strings.Contains(vars.term, "/") {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidTermDefinition,
				Err:  errors.New("invalid @prefix entry: term contains colon or slash"),
			}
		}

		// [spec // 4.2.2 // 25.2] Set the prefix flag to the value associated with the @prefix` entry, which *MUST* be a boolean. Otherwise, an `invalid @prefix value` error has been detected and processing is aborted.

		if prefixBoolean, ok := prefixMember.Value.(inspectjson.BooleanValue); ok {
			if prefixBoolean.Value {
				// [spec // 4.2.2 // 25.3] If the prefix flag of *definition* is set to `true`, and its IRI mapping is a keyword, an `invalid term definition` has been detected and processing is aborted.

				if definition.IRI == nil {
					// TODO undefined behavior?
				} else if _, ok := definition.IRI.(ExpandedIRIasKeyword); ok {
					return jsonldtype.Error{
						Code: jsonldtype.InvalidTermDefinition,
						Err:  fmt.Errorf("invalid @prefix entry: invalid iri mapping type: keyword"),
					}
				}

				definition.Prefix = true
			} else {
				definition.Prefix = false
			}
		} else {
			return jsonldtype.Error{
				Code: jsonldtype.InvalidAtPrefixValue,
				Err:  fmt.Errorf("invalid type: %s", prefixMember.Value.GetGrammarName()),
			}
		}
	}

	// [spec // 4.2.2 // 26] If *value* contains any entry other than `@id`, `@reverse`, `@container`, `@context`, `@direction`, `@index`, `@language`, `@nest`, `@prefix`, `@protected`, or `@type`, an `invalid term definition` error has been detected and processing is aborted.

	for key := range valueObject.Members {
		switch key {
		case "@id", "@reverse", "@container", "@context", "@direction", "@index", "@language", "@nest", "@prefix", "@protected", "@type":
			// valid
		default:
			return jsonldtype.Error{
				Code: jsonldtype.InvalidTermDefinition,
				Err:  fmt.Errorf("invalid entry: %s", key),
			}
		}
	}

	// [spec // 4.2.2 // 27] If *override protected* is `false` and *previous definition* exists and is protected;

	if vars.overrideProtected == false && previousDefinition != nil && previousDefinition.Protected {

		// [spec // 4.2.2 // 27.1] If *definition* is not the same as *previous definition* (other than the value of protected), a `protected term redefinition` error has been detected, and processing is aborted.

		definition.Protected = previousDefinition.Protected

		if !definition.Equals(previousDefinition) {
			return jsonldtype.Error{
				Code: jsonldtype.ProtectedTermRedefinition,
			}
		}

		// [spec // 4.2.2 // 27.2] Set *definition* to *previous definition* to retain the value of protected.

		definition = previousDefinition
	}

	// [spec // 4.2.2 // 28] Set the *term* definition of term in *active context* to *definition* and set the value associated with *defined*'s entry *term* to `true`.

	vars.activeContext.TermDefinitions[vars.term] = definition
	vars.defined[vars.term] = true

	return nil
}
