package rdf

// BlankNodeIdentifier is a constraint to ensure expected types being used as identifiers. Values are considered the
// internal representation of a blank node used to determine whether two [BlankNode] values are the same.
//
// Unlike encoding formats, such as Turtle, a BlankNodeIdentifier does not have an inherent string representation.
//
// Custom implementations are not generally necessary, but, if needed, they must embed
// [UnimplementedBlankNodeIdentifier] to fulfill this interface.
type BlankNodeIdentifier interface {
	isBlankNodeIdentifier()
}

type UnimplementedBlankNodeIdentifier struct{}

func (UnimplementedBlankNodeIdentifier) isBlankNodeIdentifier() {}
