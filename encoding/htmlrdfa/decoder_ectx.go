package htmlrdfa

import (
	"strings"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

var InitialContext = iri.PrefixMappingList{}

type incompleteTripleDirection int

const (
	incompleteTripleDirectionNone incompleteTripleDirection = iota
	incompleteTripleDirectionForward
	incompleteTripleDirectionReverse
)

type incompleteTriple struct {
	List           *listMappingBuilder
	Predicate      rdf.IRI
	PredicateRange *cursorio.TextOffsetRange
	Direction      incompleteTripleDirection
}

type HtmlProcessingProfile int

const (
	UnspecifiedHtmlProcessingProfile HtmlProcessingProfile = 0b0

	DisabledHtmlProcessingProfile HtmlProcessingProfile = 0b1

	ActiveHtmlProcessingProfile HtmlProcessingProfile = 0b10

	// XHTML1_RDFa10_HtmlProcessProfile HtmlProcessingProfile = 0b10
	XHTML1_RDFa11_HtmlProcessProfile HtmlProcessingProfile = 0b110
	XHTML5_RDFa11_HtmlProcessProfile HtmlProcessingProfile = 0b1110

	HTML4_RDFa11_HtmlProcessProfile HtmlProcessingProfile = 0b11110
	HTML5_RDFa11_HtmlProcessProfile HtmlProcessingProfile = 0b111110
)

func (p HtmlProcessingProfile) IsXHTML() bool {
	return p.IsXHTML1() || p.IsXHTML5()
}

func (p HtmlProcessingProfile) IsXHTML1() bool {
	// return p == XHTML1_RDFa10_HtmlProcessProfile || p == XHTML1_RDFa11_HtmlProcessProfile
	return p == XHTML1_RDFa11_HtmlProcessProfile
}

func (p HtmlProcessingProfile) IsXHTML5() bool {
	return p == XHTML5_RDFa11_HtmlProcessProfile
}

func (p HtmlProcessingProfile) IsHTML() bool {
	return p.IsHTML4() || p.IsHTML5()
}

func (p HtmlProcessingProfile) IsHTML4() bool {
	// return p == XHTML1_RDFa10_HtmlProcessProfile || p == XHTML1_RDFa11_HtmlProcessProfile
	return p == HTML4_RDFa11_HtmlProcessProfile
}

func (p HtmlProcessingProfile) IsHTML5() bool {
	return p == HTML5_RDFa11_HtmlProcessProfile
}

type globalEvaluationContext struct {
	DisableBackcompatXmlnsPrefixes bool

	HostDefaultVocabulary *string

	// not enumerated as a host option by specifications? seems useful though
	HostDefaultPrefixes *iri.PrefixManager

	// https://www.w3.org/TR/rdfa-in-html/
	HtmlProcessing HtmlProcessingProfile
	HtmlFoundBase  bool

	BlankNodeStringFactory blanknodes.StringFactory

	// As a special case, `_:` is also a valid reference for *one* specific bnode.
	BlankNodeSpecialSpecific rdf.BlankNode
}

type listMappingBuilder struct {
	Objects rdf.ObjectValueList
}

type evaluationContext struct {
	BaseURL           *iri.ParsedIRI
	ParentSubject     rdf.SubjectValue
	ParentObject      rdf.ObjectValue
	IncompleteTriples []incompleteTriple
	ListMapping       map[rdf.IRI]*listMappingBuilder
	Language          *string
	PrefixMapping     *iri.PrefixManager
	TermMappings      map[string]rdf.IRI
	DefaultVocabulary *string

	ParentSubjectAnno *cursorio.TextOffsetRange
	ParentObjectAnno  *cursorio.TextOffsetRange

	CurrentContainer encoding.ContainerResource

	Global *globalEvaluationContext
}

func resolveIRI(ectx evaluationContext, prefixes *iri.PrefixManager, value string, baseURL *iri.ParsedIRI, defaultVocabulary *string, allowSafeCurie bool, allowTerms bool) rdf.SubjectValue {
	if len(value) == 0 {
		if baseURL != nil {
			return rdf.IRI(baseURL.String())
		}

		return nil
	}

	isSafeCurie := value[0] == '[' && value[len(value)-1] == ']'
	if isSafeCurie {
		value = value[1 : len(value)-1]

		// [rdfa-core // 7.4] A related consequence of this is that when the value of an attribute of this datatype is an empty SafeCURIE (e.g., @about="[]"), that value does not result in an IRI and therefore the value is ignored.
		// [rdfa-info-test-suite] rdfa1.1/xhtml5/manifest#0121
		if allowSafeCurie && len(value) == 0 {
			return nil
		}
	}

	// apparently safe values can still be blank nodes (i.e. test-suites/rdfa1.0/html4/0017.html)
	if strings.HasPrefix(value, "_:") {
		if value == "_:" {
			// [rdfa-core] As a special case, `_:` is also a valid reference for *one* specific bnode.
			// [rdfa-info-test-suite] rdfa1.1/xhtml5/manifest#0088
			if ectx.Global.BlankNodeSpecialSpecific.Identifier == nil {
				ectx.Global.BlankNodeSpecialSpecific = ectx.Global.BlankNodeStringFactory.NewBlankNode()
			}

			return ectx.Global.BlankNodeSpecialSpecific
		}

		return ectx.Global.BlankNodeStringFactory.NewStringBlankNode(value[2:])
	}

	valueSplit := strings.SplitN(value, ":", 2)

	if len(valueSplit) == 1 {
		// spec does not seem to be explicit about resolution order between term vs relative IRI

		if baseURL != nil {
			resolved, err := baseURL.Parse(value)
			if err == nil {
				return rdf.IRI(resolved.String())
			}
		}

		if !strings.ContainsRune(valueSplit[0], ':') {
			// [rdfa-core // 7.4.3] If there is a local default vocabulary the IRI is obtained by concatenating that value and the term.

			if defaultVocabulary != nil {
				if ectx.Global.HostDefaultVocabulary != defaultVocabulary {
					return rdf.IRI(*defaultVocabulary + value)
				}
			}

			if allowTerms {
				// [rdfa-core // 7.4.3] Otherwise, check if the term matches an item in the list of local term mappings. First compare against the list case-sensitively, and if there is no match then compare case-insensitively. If there is a match, use the associated IRI.

				if vv, ok := ectx.TermMappings[value]; ok {
					return vv
				}

				valueLower := strings.ToLower(value)

				for k, v := range ectx.TermMappings {
					if strings.ToLower(k) == valueLower {
						return v
					}
				}
			}

			// [rdfa-core // 7.4.3] Otherwise, the term has no associated IRI and must be ignored.

			return nil
		}
	}

	expanded, ok := prefixes.ExpandPrefix(iri.PrefixReference{
		Prefix:    valueSplit[0],
		Reference: valueSplit[1],
	})
	if ok {
		return rdf.IRI(expanded)
	}

	// If the prefix is empty and not found in prefix mappings, use default vocabulary
	if valueSplit[0] == "" && defaultVocabulary != nil {
		if len(valueSplit[1]) == 0 {
			return rdf.IRI(*defaultVocabulary)
		}

		return rdf.IRI(*defaultVocabulary + valueSplit[1])
	}

	return rdf.IRI(value)
}
