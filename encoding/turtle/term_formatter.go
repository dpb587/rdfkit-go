package turtle

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
)

type TermFormatterOptions struct {
	ASCII             bool
	Base              *iriutil.BaseIRI
	Prefixes          iriutil.PrefixMap
	BlankNodeStringer blanknodeutil.Stringer
}

type termFormatter struct {
	ascii             bool
	base              *iriutil.BaseIRI
	prefixes          iriutil.PrefixMap
	blankNodeStringer blanknodeutil.Stringer
}

func NewTermFormatter(options TermFormatterOptions) termutil.Formatter {
	return &termFormatter{
		ascii:             options.ASCII,
		base:              options.Base,
		prefixes:          options.Prefixes,
		blankNodeStringer: options.BlankNodeStringer,
	}
}

func (tf *termFormatter) FormatTerm(t rdf.Term) string {
	switch t := t.(type) {
	case rdf.IRI:
		if prefix, localName, ok := tf.prefixes.CompactPrefix(t); ok {
			return prefix + ":" + formatIRI(localName, tf.ascii)
		} else if tf.base != nil {
			if reference, ok := tf.base.RelativizeIRI(t); ok {
				return "<" + formatIRI(reference, tf.ascii) + ">"
			}
		}

		return "<" + formatIRI(string(t), tf.ascii) + ">"
	case rdf.BlankNode:
		if tf.blankNodeStringer != nil {
			return tf.blankNodeStringer.GetBlankNodeIdentifier(t)
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
			sb.WriteString("@")
			sb.WriteString(t.Tags[rdf.LanguageLiteralTag])
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
