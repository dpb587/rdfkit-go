package rdf

// GraphNameValue represents any value that can be used for a graph name property.
//
// Normative values are a [BlankNode] or [IRI] type; or the [DefaultGraph] constant.
type GraphNameValue interface {
	isGraphNameValueBuiltin()
}

//

type GraphNameValueList []GraphNameValue
