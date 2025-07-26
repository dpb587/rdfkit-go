package schemaliteral

import (
	"net/url"
	"regexp"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/schema/schemairi"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdutil"
	"github.com/dpb587/rdfkit-go/rdf"
)

func NewURL(v string) rdf.ObjectValue {
	return rdf.Literal{
		Datatype:    schemairi.URL_Class,
		LexicalForm: v,
	}
}

var textIsLikelyURL = regexp.MustCompile(`(^/|^\.\.?/|^#|^\?|^[a-z]{2,12}:)`)

func CastURL(t rdf.Term, opts CastOptions) (rdf.ObjectValue, bool) {
	var valueString string

	switch t := t.(type) {
	case rdf.Literal:
		switch t.Datatype {
		case schemairi.URL_Class:
			valueString = t.LexicalForm // ensure resolved
		case rdfiri.LangString_Datatype,
			xsdiri.String_Datatype:
			valueString = xsdutil.WhiteSpaceCollapse(t.LexicalForm)

			if len(valueString) == 0 {
				// empty string more often used for empty value, not relative URL
				return nil, false
			} else if !textIsLikelyURL.MatchString(valueString) {
				// URL data type will be preferred before Text data type, so trying to avoid non-URLs
				return nil, false
			}
		default:
			return nil, false
		}
	case rdf.IRI:
		valueString = string(t)
	default:
		return nil, false
	}

	parsed, err := url.Parse(valueString)
	if err != nil {
		return nil, false
	}

	if !parsed.IsAbs() && opts.BaseURL != nil {
		parsed = opts.BaseURL.ResolveReference(parsed)
	}

	return rdf.Literal{
		Datatype:    schemairi.URL_Class,
		LexicalForm: parsed.String(),
	}, true
}
