package rdf

// SubjectValue represents any value that can be used for a subject property.
//
// Normative values are a [BlankNode] or [IRI] type.
type SubjectValue interface {
	Term

	isSubjectValueBuiltin()

	isObjectValueBuiltin() // presumptuous; not explicitly stated in the spec?
}

type SubjectValueList []SubjectValue
