package rdf

// BlankNode does not identify a specific resource and is disjoint from an [IRI] and a [Literal]. Two BlankNode values
// are equal if and only if they have the same internal [BlankNodeIdentifier].
//
// [NewBlankNode] must be used to create a blank node.
type BlankNode interface {
	Term
	SubjectValue
	ObjectValue
	GraphNameValue

	GetBlankNodeIdentifier() BlankNodeIdentifier
}

type blankNode struct {
	identifier BlankNodeIdentifier
}

// NewBlankNode creates a new blank node using [DefaultBlankNodeFactory].
func NewBlankNode() BlankNode {
	return DefaultBlankNodeFactory.NewBlankNode()
}

// NewBlankNodeWithIdentifier creates a new blank node using the provided identifier. This should only be called by
// custom [BlankNodeFactory] implementations.
func NewBlankNodeWithIdentifier(identifier BlankNodeIdentifier) BlankNode {
	if identifier == nil {
		panic("identifier cannot be nil")
	}

	return blankNode{
		identifier: identifier,
	}
}

var _ BlankNode = blankNode{}

func (blankNode) isTermBuiltin()           {}
func (blankNode) isSubjectValueBuiltin()   {}
func (blankNode) isObjectValueBuiltin()    {}
func (blankNode) isGraphNameValueBuiltin() {}

func (t blankNode) AsObjectValue() ObjectValue {
	return t
}

func (blankNode) TermKind() TermKind {
	return TermKindBlankNode
}

func (t blankNode) TermEquals(a Term) bool {
	if t.identifier == nil {
		return false
	}

	aBlankNode, ok := a.(blankNode)
	if !ok {
		return false
	} else if aBlankNode.identifier == nil {
		return false
	}

	return t.identifier == aBlankNode.identifier
}

func (t blankNode) GetBlankNodeIdentifier() BlankNodeIdentifier {
	return t.identifier
}
