package rdf

// A Literal is used for values such as strings, numbers, and dates.
type Literal struct {
	// Datatype determines how the lexical form maps to a literal value.
	Datatype IRI

	// LexicalForm is the datatype-based encoding of the literal value.
	LexicalForm string

	// Tag, for specific datatypes, may contain properties that further describe the literal.
	//
	// For `http://www.w3.org/1999/02/22-rdf-syntax-ns#langString`, [LanguageLiteralTag] is expected.
	//
	// For all other datatypes, the value should be nil.
	Tag LiteralTag
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

	if t.Tag == nil && dLiteral.Tag == nil {
		// ok
	} else if t.Tag == nil || dLiteral.Tag == nil {
		return false
	} else if !t.Tag.Equals(dLiteral.Tag) {
		return false
	}

	return t.LexicalForm == dLiteral.LexicalForm
}
