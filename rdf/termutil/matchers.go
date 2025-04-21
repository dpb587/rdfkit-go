package termutil

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

type Equals struct {
	Expected rdf.Term
}

func (m Equals) MatchTerm(t rdf.Term) bool {
	return m.Expected.TermEquals(t)
}

//

func EqualsOneOf(expected ...rdf.Term) Matcher {
	mappedIRIs := map[rdf.IRI]struct{}{}
	mappedBlankNodes := map[rdf.BlankNodeIdentifier]struct{}{}
	allLiterals := map[rdf.IRI][]rdf.Literal{}

	for _, t := range expected {
		switch t := t.(type) {
		case rdf.IRI:
			mappedIRIs[t] = struct{}{}
		case rdf.BlankNode:
			if identifier := t.GetBlankNodeIdentifier(); identifier != nil {
				mappedBlankNodes[identifier] = struct{}{}
			}
		case rdf.Literal:
			allLiterals[t.Datatype] = append(allLiterals[t.Datatype], t)
		}
	}

	return MatcherFunc(func(t rdf.Term) bool {
		switch t := t.(type) {
		case rdf.IRI:
			_, ok := mappedIRIs[t]
			return ok
		case rdf.BlankNode:
			identifier := t.GetBlankNodeIdentifier()
			if identifier == nil {
				return false
			}

			_, ok := mappedBlankNodes[identifier]

			return ok
		case rdf.Literal:
			for _, l := range allLiterals[t.Datatype] {
				if l.TermEquals(t) {
					return true
				}
			}

			return false
		}

		return false
	})
}
