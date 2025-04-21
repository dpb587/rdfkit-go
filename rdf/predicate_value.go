package rdf

// PredicateValue represents any value that can be used for a predicate property.
//
// Normative values are an [IRI] type.
type PredicateValue interface {
	Term

	isPredicateValueBuiltin()

	isObjectValueBuiltin() // presumptuous; not explicitly stated in the spec?
}

type PredicateValueList []PredicateValue
