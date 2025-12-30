package rdfdescriptionstruct

// Collection represents an RDF list that should be traversed and unmarshaled.
// When unmarshaling, if the object value is an IRI or Blank Node representing
// a resource with type rdf:List, the unmarshaler will follow the chain of values
// and decode each list value according to the element type T.
type Collection[T any] []T
