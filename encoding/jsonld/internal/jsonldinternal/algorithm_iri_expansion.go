package jsonldinternal

import (
	"fmt"
	"strings"

	"github.com/dpb587/inspectjson-go/inspectjson"
)

type algorithmIRIExpansion struct {
	activeContext *Context

	// TODO switch to inspectjson.Token? some calls started to wrap with Primitive
	value inspectjson.Value

	// [spec // 5.2.2] whether value can be interpreted as a relative IRI reference against the document's base IRI
	documentRelative bool

	// [spec // 5.2.2] whether value can be interpreted as a relative IRI reference against the active context's vocabulary mapping
	vocab bool

	localContext *inspectjson.ObjectValue
	defined      map[string]bool
}

func (opts algorithmIRIExpansion) Call() ExpandedIRI {

	// [spec // 5.2.2 // 1] If *value* is a keyword or `null`, return `value` as is.

	if _, ok := opts.value.(inspectjson.NullValue); ok {
		return ExpandedIRIasNil{}
	}

	valueString, isValueString := opts.value.(inspectjson.StringValue)
	if isValueString {
		if _, known := definedKeywords[valueString.Value]; known {
			return ExpandedIRIasKeyword(valueString.Value)
		}
	}

	// [dpb] steps after seem to assume string? could simplify

	// [spec // 5.2.2 // 2] If *value* has the form of a keyword (i.e., it matches the ABNF rule `"@"1*ALPHA` from [RFC5234]), a processor SHOULD generate a warning and return `null`.

	if isValueString && reKeywordABNF.MatchString(valueString.Value) {
		// TODO warning

		return ExpandedIRIasNil{}
	}

	// [spec // 5.2.2 // 3] If *local context* is not `null`, it contains an entry with a key that equals *value*, and the *value* of the entry for *value* in *defined* is not `true`, invoke the Create Term Definition algorithm, passing *active context*, *local context*, *value* as *term*, and *defined*. This will ensure that a term definition is created for *value* in *active context* during Context Processing.

	if isValueString && opts.localContext != nil {
		valueString := valueString.Value

		if _, ok := opts.localContext.Members[valueString]; ok && opts.defined[valueString] != true {
			// TODO error handling?
			_ = algorithmCreateTermDefinition{
				activeContext: opts.activeContext,
				localContext:  opts.localContext,
				term:          valueString,
				defined:       opts.defined,
				// defaults
				baseURL:               nil,
				protected:             false,
				overrideProtected:     false,
				remoteContexts:        nil,
				validateScopedContext: true,
			}.Call()
		}
	}

	// [spec // 5.2.2 // 4] If *active context* has a term definition for *value*, and the associated IRI mapping is a keyword, return that keyword.

	var termDefinition *TermDefinition

	if isValueString {
		termDefinition = opts.activeContext.TermDefinitions[valueString.Value]
	}

	if termDefinition != nil {
		if termDefinition.IRI != nil {
			switch termDefinition.IRI.(type) {
			case ExpandedIRIasKeyword:
				return termDefinition.IRI
			case ExpandedIRIasNil:
				// [dpb] unofficial method explicitly null'd values are ignored
				return termDefinition.IRI
			}
		}
	}

	// [spec // 5.2.2 // 5] If *vocab* is `true` and the *active context* has a term definition for *value*, return the associated IRI mapping.

	if opts.vocab && termDefinition != nil && termDefinition.IRI != nil {
		return termDefinition.IRI
	}

	// [spec // 5.2.2 // 6] If *value* contains a colon (`:`) anywhere after the first character, it is either an IRI, a compact IRI, or a blank node identifier:

	if isValueString {
		valueString := valueString.Value

		if len(valueString) > 1 && strings.Contains(valueString[1:], ":") {

			// [spec // 5.2.2 // 6.1] Split *value* into a *prefix* and *suffix* at the first occurrence of a colon (`:`).

			valuePrefixSuffix := strings.SplitN(valueString, ":", 2)

			// [spec // 5.2.2 // 6.2] If *prefix* is underscore (`_`) or *suffix* begins with double-forward-slash (`//`), return *value* as it is already an IRI or a blank node identifier.

			if valuePrefixSuffix[0] == "_" {
				return ExpandedIRIasBlankNode(valueString)
			} else if len(valuePrefixSuffix[1]) > 2 && valuePrefixSuffix[1][:2] == "//" {
				return ExpandedIRIasIRI(valueString)
			}

			// [spec // 5.2.2 // 6.3] If *local context* is not `null`, it contains a *prefix* entry, and the value of the *prefix* entry in *defined* is not true, invoke the Create Term Definition algorithm, passing *active context*, *local context*, *prefix* as *term*, and *defined*. This will ensure that a term definition is created for prefix in active context during Context Processing.

			if opts.localContext != nil {
				if _, ok := opts.localContext.Members[valuePrefixSuffix[0]]; ok && opts.defined[valuePrefixSuffix[0]] != true {
					// TODO error handling?
					_ = algorithmCreateTermDefinition{
						activeContext: opts.activeContext,
						localContext:  opts.localContext,
						term:          valuePrefixSuffix[0],
						defined:       opts.defined,
						// defaults
						baseURL:               nil,
						protected:             false,
						overrideProtected:     false,
						remoteContexts:        nil,
						validateScopedContext: true,
					}.Call()
				}
			}

			// [spec // 5.2.2 // 6.4] If *active context* contains a term definition for *prefix* having a non-`null` IRI mapping and the prefix flag of the term definition is `true`, return the result of concatenating the IRI mapping associated with *prefix* and *suffix*.

			if termDefinition := opts.activeContext.TermDefinitions[valuePrefixSuffix[0]]; termDefinition != nil && termDefinition.IRI != nil && termDefinition.Prefix {
				switch t := termDefinition.IRI.(type) {
				case ExpandedIRIasIRI:
					return t + ExpandedIRIasIRI(valuePrefixSuffix[1])
				case ExpandedIRIasBlankNode:
					return t + ExpandedIRIasBlankNode(valuePrefixSuffix[1])
				}

				panic(fmt.Errorf("unexpected term definition IRI type: %T", termDefinition.IRI))
			}

			// [spec // 5.2.2 // 6.5] If value has the form of an IRI, return value.
			// [dpb] spec seems amgiguous given [6.2] checked it as an IRI, and relative IRI validation seems difficult to get right (vs unspecified compact IRI)
			// [dpb] following seems hacky; #t0118, #tc022, t0109

			if !strings.Contains(valueString, "/") && !strings.Contains(valueString, "#") && !strings.Contains(valueString, "?") {
				return ExpandedIRIasIRI(valueString)
			}
		}

		// [spec // 5.2.2 // 7] If *vocab* is `true`, and *active context* has a vocabulary mapping, return the result of concatenating the vocabulary mapping with *value*.

		if opts.vocab && opts.activeContext.VocabularyMapping != nil {
			// TODO unsafe assertion; possible BlankNode?
			switch t := opts.activeContext.VocabularyMapping.(type) {
			case ExpandedIRIasIRI:
				return t + ExpandedIRIasIRI(valueString)
			case ExpandedIRIasBlankNode:
				return t + ExpandedIRIasBlankNode(valueString)
			}
		}

		// [spec // 5.2.2 // 8] Otherwise, if *document relative* is `true` set *value* to the result of resolving *value* against the base IRI from *active context*. Only the basic algorithm in section 5.2 of [RFC3986] is used; neither Syntax-Based Normalization nor Scheme-Based Normalization are performed. Characters additionally allowed in IRI references are treated in the same way that unreserved characters are treated in URI references, per section 6.5 of [RFC3987].

		if opts.documentRelative && opts.activeContext.BaseURL != nil {
			resolvedURL, err := opts.activeContext.BaseURL.Parse(valueString)
			if err == nil {
				// [dpb] early return, but should be assigning to value
				return ExpandedIRIasIRI(resolvedURL.String())
			}
		}
	}

	// [spec // 5.2.2 // 9] Return *value* as is.
	// [dpb] TODO refactor to avoid RawValue?

	if isValueString {
		// TODO only cast to iri if value qualifies as IRI? always false by this point?
		return ExpandedIRIasIRI(valueString.Value)
	}

	return ExpandedIRIasRawValue{opts.value}
}
