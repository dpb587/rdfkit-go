package rdf

type Triple struct {
	Subject   SubjectValue
	Predicate PredicateValue
	Object    ObjectValue
}

// TODO needs more research; subject+object? wait for next version of the spec?
// in the meantime, standalone struct is meaningful elsewhere and unlikely to change with new spec

// var _ Term = Triple{}
// var _ ObjectValue = Triple{}

// func (Triple) isTermBuiltin()        {}
// func (Triple) isObjectValueBuiltin() {}

// func (Triple) TermKind() TermKind {
// 	return TermKindTriple
// }

// func (t Triple) TermEquals(d Term) bool {
// 	dTriple, ok := d.(Triple)
// 	if !ok {
// 		return false
// 	} else if !t.Subject.TermEquals(dTriple.Subject) {
// 		return false
// 	} else if !t.Predicate.TermEquals(dTriple.Predicate) {
// 		return false
// 	}

// 	return t.Object.TermEquals(dTriple.Object)
// }

//

type TripleList []Triple
