package rdf

// PredicateValue represents any value that can be used for a predicate property.
//
// This is a closed interface. See [IRI].
type PredicateValue interface {
	Term

	isPredicateValueBuiltin()

	isObjectValueBuiltin() // presumptuous; not explicitly stated in the spec?
}

type PredicateValueList []PredicateValue
