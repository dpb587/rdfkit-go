package rdf

type BlankNodeIdentifier interface {
	EqualsBlankNodeIdentifier(bni BlankNodeIdentifier) bool
}

// BlankNode does not identify a specific resource and is disjoint from an [IRI] and a [Literal]. Two BlankNode values
// are equal if both of their Identifier fields are non-nil and equivalent.
//
// You should use [NewBlankNode] or a [BlankNodeFactory] to construct a new BlankNode.
type BlankNode struct {
	Identifier BlankNodeIdentifier
}

var _ Term = BlankNode{}
var _ SubjectValue = BlankNode{}
var _ ObjectValue = BlankNode{}
var _ GraphNameValue = BlankNode{}

func (BlankNode) isTermBuiltin()           {}
func (BlankNode) isSubjectValueBuiltin()   {}
func (BlankNode) isObjectValueBuiltin()    {}
func (BlankNode) isGraphNameValueBuiltin() {}

func (t BlankNode) AsObjectValue() ObjectValue {
	return t
}

func (BlankNode) TermKind() TermKind {
	return TermKindBlankNode
}

func (t BlankNode) TermEquals(a Term) bool {
	if t.Identifier == nil {
		return false
	}

	aBlankNode, ok := a.(BlankNode)
	if !ok {
		return false
	} else if aBlankNode.Identifier == nil {
		return false
	}

	return t.Identifier.EqualsBlankNodeIdentifier(aBlankNode.Identifier)
}
