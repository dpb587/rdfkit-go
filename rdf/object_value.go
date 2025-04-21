package rdf

// ObjectValue represents any value that can be used for an object property.
//
// Normative values are a [BlankNode], [IRI], or [Literal] type.
type ObjectValue interface {
	Term

	isObjectValueBuiltin()
}

type ObjectValueList []ObjectValue
