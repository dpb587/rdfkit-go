package rdfa

import (
	"strings"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

var InitialContext = iriutil.PrefixMap{}

type incompleteTripleDirection int

const (
	incompleteTripleDirectionNone incompleteTripleDirection = iota
	incompleteTripleDirectionForward
	incompleteTripleDirectionReverse
)

type incompleteTriple struct {
	List           rdf.SubjectValue
	Predicate      rdf.IRI
	PredicateRange *cursorio.TextOffsetRange
	Direction      incompleteTripleDirection
}

type HtmlProcessingProfile int

const (
	UnspecifiedHtmlProcessingProfile HtmlProcessingProfile = 0b0

	DisabledHtmlProcessingProfile HtmlProcessingProfile = 0b1

	ActiveHtmlProcessingProfile HtmlProcessingProfile = 0b10

	XHTML1_RDFa10_HtmlProcessProfile HtmlProcessingProfile = 0b10
	XHTML1_RDFa11_HtmlProcessProfile HtmlProcessingProfile = 0b110
	XHTML5_RDFa11_HtmlProcessProfile HtmlProcessingProfile = 0b1110
)

type globalEvaluationContext struct {
	DisableBackcompatXmlnsPrefixes bool

	HostDefaultVocabulary *string

	// not enumerated as a host option by specifications? seems useful though
	HostDefaultPrefixes iriutil.PrefixMap

	// https://www.w3.org/TR/rdfa-in-html/
	HtmlProcessing HtmlProcessingProfile
	HtmlFoundBase  bool

	BlankNodeStringMapper blanknodeutil.StringMapper
	BlankNodeFactory      rdf.BlankNodeFactory
}

type listMappingBuilder struct {
	Objects rdf.ObjectValueList
}

type evaluationContext struct {
	BaseURL           *iriutil.ParsedIRI
	ParentSubject     rdf.SubjectValue
	ParentObject      rdf.ObjectValue
	IncompleteTriples []incompleteTriple
	ListMapping       map[rdf.IRI]*listMappingBuilder
	Language          *string
	PrefixMapping     iriutil.PrefixMap
	TermMappings      map[string]rdf.IRI
	DefaultVocabulary *string

	ParentSubjectAnno *cursorio.TextOffsetRange
	ParentObjectAnno  *cursorio.TextOffsetRange

	CurrentContainer encoding.ContainerResource

	Global *globalEvaluationContext
	XMLNS  map[string]string // for propagating in XML/HTML node serializations
}

func resolveSubjectIRI(g *globalEvaluationContext, prefixes iriutil.PrefixMap, termMappings map[string]rdf.IRI, value string, baseURL *iriutil.ParsedIRI) rdf.SubjectValue {
	return resolveIRI(g, prefixes, termMappings, value, baseURL, nil)
}

// func (ectx EvaluationContext) SafeCURIEorCURIEorIRIorBlankNode(value string) term.TermPrimitive {
// 	return ectx.resolveIRI(value, nil, ectx.DefaultVocabulary)
// }

func resolveSafeCURIEorCURIEorIRI(g *globalEvaluationContext, prefixes iriutil.PrefixMap, termMappings map[string]rdf.IRI, value string, defaultVocabulary *string) (rdf.IRI, bool) {
	v, ok := resolveIRI(g, prefixes, termMappings, value, nil, defaultVocabulary).(rdf.IRI)
	if !ok {
		return "", false
	}

	return v, true
}

func resolveIRI(g *globalEvaluationContext, prefixes iriutil.PrefixMap, termMappings map[string]rdf.IRI, value string, baseURL *iriutil.ParsedIRI, defaultVocabulary *string) rdf.SubjectValue {
	if len(value) == 0 {
		if baseURL != nil {
			return rdf.IRI(baseURL.String())
		}

		return nil
	}

	var isSafe bool

	if value[0] == '[' && value[len(value)-1] == ']' {
		value = value[1 : len(value)-1]
		isSafe = true
	}

	// apparently safe values can still be blank nodes (i.e. test-suites/rdfa1.0/html4/0017.html)
	if strings.HasPrefix(value, "_:") {
		if value == "_:" {
			// [rdfa-core] As a special case, _: is also a valid reference for one specific bnode.
			return g.BlankNodeFactory.NewBlankNode()
		}

		return g.BlankNodeStringMapper.MapBlankNodeIdentifier(value[2:])
	}

	if isSafe {
		// TODO this is not following spec; should be returning empty + bool=false?
		return rdf.IRI(value)
	}

	valueSplit := strings.SplitN(value, ":", 2)

	if len(valueSplit) == 1 {
		if vv, ok := termMappings[string(value)]; ok {
			return vv
		}

		if baseURL != nil {
			resolved, err := baseURL.Parse(value)
			if err == nil {
				return rdf.IRI(resolved.String())
			}
		}

		if defaultVocabulary != nil {
			parsedDefaultVocabulary, err := iriutil.ParseIRI(string(*defaultVocabulary))
			if err == nil {
				resolved, err := parsedDefaultVocabulary.Parse(value)
				if err == nil {
					return rdf.IRI(resolved.String())
				}
			}

			// fallback?
			return rdf.IRI(*defaultVocabulary + value)
		}

		return nil
	}

	expanded, ok := prefixes.ExpandPrefix(valueSplit[0], valueSplit[1])
	if ok {
		return rdf.IRI(expanded)
	}

	return rdf.IRI(value)
}
