package rdf

// SubjectValue represents any value that can be used for a subject property.
//
// This is a closed interface. See [BlankNode] and [IRI].
type SubjectValue interface {
	Term

	isSubjectValueBuiltin()

	isObjectValueBuiltin() // presumptuous; not explicitly stated in the spec?
}

type SubjectValueList []SubjectValue
