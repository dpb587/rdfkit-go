package rdf

// GraphNameValue represents any value that can be used for a graph name property.
//
// Normative values are a [BlankNode] or [IRI] type; or nil for the default graph.
type GraphNameValue interface {
	Term // not inherently a term in and of itself, but simplifies usage of Term-related convenience functions

	isGraphNameValueBuiltin()
}

//

type GraphNameValueList []GraphNameValue
