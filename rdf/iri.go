package rdf

// An IRI (Internationalized Resource Identifier) is a string that conforms to the syntax defined in [RFC 3987] which
// complements URIs with support for Unicode characters. Values should be fully resolved to their absolute form.
//
// [RFC 3987]: https://www.rfc-editor.org/rfc/rfc3987
type IRI string

var _ SubjectValue = IRI("")
var _ PredicateValue = IRI("")
var _ ObjectValue = IRI("")
var _ GraphNameValue = IRI("")

func (IRI) isTermBuiltin()           {}
func (IRI) isSubjectValueBuiltin()   {}
func (IRI) isPredicateValueBuiltin() {}
func (IRI) isObjectValueBuiltin()    {}
func (IRI) isGraphNameValueBuiltin() {}

func (IRI) TermKind() TermKind {
	return TermKindIRI
}

func (t IRI) TermEquals(a Term) bool {
	a, ok := a.(IRI)
	if !ok {
		return false
	}

	return t == a
}
