package turtle

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
	"github.com/dpb587/rdfkit-go/rdf/terms"
)

type TermFormatterOptions struct {
	ASCII                   bool
	Base                    *iri.BaseIRI
	Prefixes                iri.PrefixMapper
	BlankNodeStringProvider blanknodes.StringProvider
}

type termFormatter struct {
	ascii             bool
	base              *iri.BaseIRI
	prefixes          iri.PrefixMapper
	blankNodeStringer blanknodes.StringProvider
}

func NewTermFormatter(options TermFormatterOptions) terms.Formatter {
	return &termFormatter{
		ascii:             options.ASCII,
		base:              options.Base,
		prefixes:          options.Prefixes,
		blankNodeStringer: options.BlankNodeStringProvider,
	}
}

func (tf *termFormatter) FormatTerm(t rdf.Term) string {
	switch t := t.(type) {
	case rdf.IRI:
		if pr, ok := tf.prefixes.CompactPrefix(string(t)); ok {
			return pr.Prefix + ":" + format_PN_LOCAL(pr.Reference)
		} else if tf.base != nil {
			if reference, ok := tf.base.RelativizeIRI(string(t)); ok {
				return "<" + formatIRI(reference, tf.ascii) + ">"
			}
		}

		return "<" + formatIRI(string(t), tf.ascii) + ">"
	case rdf.BlankNode:
		if tf.blankNodeStringer != nil {
			return "_:" + tf.blankNodeStringer.GetBlankNodeString(t)
		}

		return "_:b0x" + strconv.FormatUint(uint64(reflect.ValueOf(t).Pointer()), 16)
	case rdf.Literal:
		switch t.Datatype {
		case xsdiri.String_Datatype:
			return formatLiteralLexicalForm(string(t.LexicalForm), tf.ascii)
		case xsdiri.Boolean_Datatype:
			if t.LexicalForm == "true" {
				return "true"
			} else if t.LexicalForm == "false" {
				return "false"
			}
		}

		sb := &strings.Builder{}
		sb.WriteString(formatLiteralLexicalForm(string(t.LexicalForm), tf.ascii))
		if t.Datatype == rdfiri.LangString_Datatype {
			if langTag, ok := t.Tag.(rdf.LanguageLiteralTag); ok {
				sb.WriteString("@")
				sb.WriteString(langTag.Language)
			}
		} else {
			sb.WriteString("^^")
			sb.WriteString(tf.FormatTerm(t.Datatype))
		}

		return sb.String()
	}

	panic(fmt.Errorf("unknown term type %T", t))
}

//

var defaultTermFormatter = &termFormatter{}

// FormatTerm formats the term according to Turtle syntax.
//
// Blank Nodes are formatted as "_:b0x" followed by their memory address.
func FormatTerm(t rdf.Term) string {
	return defaultTermFormatter.FormatTerm(t)
}

//

var defaultTermFormatterASCII = &termFormatter{
	ascii: true,
}

// FormatTerm formats the term according to Turtle syntax with only ASCII characters.
//
// Blank Nodes are formatted as "_:b0x" followed by their memory address.
func FormatTermASCII(t rdf.Term) string {
	return defaultTermFormatterASCII.FormatTerm(t)
}
