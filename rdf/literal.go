package rdf

type LiteralTag int

const (
	LanguageLiteralTag LiteralTag = iota

	// BaseDirectionLiteralTag behavior is still undefined for this Go project and pending RDF 1.2.
	BaseDirectionLiteralTag
)

// A Literal is used for values such as strings, numbers, and dates.
type Literal struct {
	// Datatype determines how the lexical form maps to a literal value.
	Datatype IRI

	// LexicalForm is the datatype-based encoding of the literal value.
	LexicalForm string

	// Tags, for specific datatypes, will contain values that further describe the literal.
	//
	// For a datatype of `http://www.w3.org/1999/02/22-rdf-syntax-ns#langString`, the [LanguageLiteralTag] should be
	// present and contain a well-formed [BCP47] tag value.
	//
	// For all other datatypes, the value should be nil.
	//
	// [BCP47]: https://www.rfc-editor.org/rfc/rfc5646
	Tags map[LiteralTag]string
}

var _ Term = Literal{}
var _ ObjectValue = Literal{}

func (Literal) isTermBuiltin()        {}
func (Literal) isObjectValueBuiltin() {}

func (Literal) TermKind() TermKind {
	return TermKindLiteral
}

func (t Literal) TermEquals(a Term) bool {
	dLiteral, ok := a.(Literal)
	if !ok {
		return false
	} else if t.Datatype != dLiteral.Datatype {
		return false
	}

	if len(t.Tags) > 0 {
		if len(t.Tags) != len(dLiteral.Tags) {
			return false
		}

		for k, v := range t.Tags {
			if dv, ok := dLiteral.Tags[k]; !ok || dv != v {
				return false
			}
		}
	}

	return t.LexicalForm == dLiteral.LexicalForm
}
