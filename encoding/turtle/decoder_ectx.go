package turtle

import (
	"fmt"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
)

type evaluationContext struct {
	CurSubject         rdf.SubjectValue
	CurSubjectLocation *cursorio.TextOffsetRange

	CurPredicate         rdf.PredicateValue
	CurPredicateLocation *cursorio.TextOffsetRange

	Global *globalEvaluationContext
}

func (ectx evaluationContext) ResolveURL(v string) (*iri.ParsedIRI, error) {
	if ectx.Global.Base == nil {
		return iri.ParseIRI(v)
	}

	return ectx.Global.Base.Parse(v)
}

func (ectx evaluationContext) ResolveIRI(v string) (rdf.IRI, error) {
	if ectx.Global.Base == nil {
		return rdf.IRI(v), nil
	}

	u, err := ectx.Global.Base.Parse(v)
	if err != nil {
		return "", fmt.Errorf("resolve iri: %v", err)
	}

	return rdf.IRI(u.String()), nil
}

type globalEvaluationContext struct {
	Base                   *iri.ParsedIRI
	Prefixes               *iri.PrefixManager
	BlankNodeStringFactory blanknodes.StringFactory
}
