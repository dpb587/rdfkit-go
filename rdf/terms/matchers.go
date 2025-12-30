package terms

import "github.com/dpb587/rdfkit-go/rdf"

type isBlankNode struct{}

func (m isBlankNode) MatchTerm(t rdf.Term) bool {
	_, ok := t.(rdf.BlankNode)

	return ok
}

var IsBlankNode = isBlankNode{}

//

type isIRI struct{}

func (m isIRI) MatchTerm(t rdf.Term) bool {
	_, ok := t.(rdf.IRI)

	return ok
}

var IsIRI = isIRI{}

//

type isLiteral struct{}

func (m isLiteral) MatchTerm(t rdf.Term) bool {
	_, ok := t.(rdf.Literal)

	return ok
}

var IsLiteral = isLiteral{}

//

type IsLiteralDatatype struct {
	Datatype rdf.TermMatcher
}

func (m IsLiteralDatatype) MatchTerm(t rdf.Term) bool {
	literal, ok := t.(rdf.Literal)
	if !ok {
		return false
	}

	return m.Datatype.MatchTerm(literal.Datatype)
}

//

type Equals struct {
	Expected rdf.Term
}

func (m Equals) MatchTerm(t rdf.Term) bool {
	return m.Expected.TermEquals(t)
}

//

type EqualsOneOfCompiled struct {
	mappedIRIs         map[rdf.IRI]struct{}
	mappedBlankNodes   map[rdf.BlankNodeIdentifier]struct{}
	literalsByDatatype map[rdf.IRI][]rdf.Literal
}

func (m EqualsOneOfCompiled) AsLogicalOrMatcher() LogicalOrMatcher {
	var as LogicalOrMatcher

	for iri := range m.mappedIRIs {
		as = append(as, Equals{
			Expected: iri,
		})
	}

	for identifier := range m.mappedBlankNodes {
		as = append(as, Equals{
			Expected: rdf.NewBlankNodeWithIdentifier(identifier),
		})
	}

	for _, literals := range m.literalsByDatatype {
		for _, literal := range literals {
			as = append(as, Equals{
				Expected: literal,
			})
		}
	}

	return as
}

func (m EqualsOneOfCompiled) MatchTerm(t rdf.Term) bool {
	switch t := t.(type) {
	case rdf.IRI:
		_, ok := m.mappedIRIs[t]
		return ok
	case rdf.BlankNode:
		identifier := t.GetBlankNodeIdentifier()
		if identifier == nil {
			return false
		}

		_, ok := m.mappedBlankNodes[identifier]

		return ok
	case rdf.Literal:
		for _, l := range m.literalsByDatatype[t.Datatype] {
			if l.TermEquals(t) {
				return true
			}
		}

		return false
	}

	return false
}

func EqualsOneOf(expected ...rdf.Term) rdf.TermMatcher {
	compiled := EqualsOneOfCompiled{
		mappedIRIs:         map[rdf.IRI]struct{}{},
		mappedBlankNodes:   map[rdf.BlankNodeIdentifier]struct{}{},
		literalsByDatatype: map[rdf.IRI][]rdf.Literal{},
	}

	for _, t := range expected {
		switch t := t.(type) {
		case rdf.IRI:
			compiled.mappedIRIs[t] = struct{}{}
		case rdf.BlankNode:
			if identifier := t.GetBlankNodeIdentifier(); identifier != nil {
				compiled.mappedBlankNodes[identifier] = struct{}{}
			}
		case rdf.Literal:
			compiled.literalsByDatatype[t.Datatype] = append(compiled.literalsByDatatype[t.Datatype], t)
		}
	}

	// simplification shortcut
	if len(compiled.mappedIRIs) == 1 && len(compiled.mappedBlankNodes) == 0 && len(compiled.literalsByDatatype) == 0 {
		for iri := range compiled.mappedIRIs {
			return Equals{
				Expected: iri,
			}
		}
	}

	return compiled
}
