package rdf

// ObjectValue represents any value that can be used for an object property.
//
// This is a closed interface. See [BlankNode], [IRI], and [Literal].
type ObjectValue interface {
	Term

	isObjectValueBuiltin()
}

type ObjectValueList []ObjectValue
