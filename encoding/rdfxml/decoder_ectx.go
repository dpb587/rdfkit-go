package rdfxml

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

type uniqRefID struct {
	Base string
	ID   string
}

type globalEvaluationContext struct {
	BlankNodeStringMapper blanknodeutil.StringMapper
	BlankNodeFactory      rdf.BlankNodeFactory

	uniqRefID map[uniqRefID]struct{}
}

type evaluationContext struct {
	Base     *iriutil.ParsedIRI
	Language *string

	ParentSubject           rdf.SubjectValue
	ParentSubjectLocation   *cursorio.TextOffsetRange
	ParentPredicate         rdf.PredicateValue
	ParentPredicateLocation *cursorio.TextOffsetRange
	ParentContainerIndex    *int

	CurrentContainer encoding.ContainerResource // unrelated to ParentContainerIndex

	Global  *globalEvaluationContext
	UsedIDs map[string]struct{}
}

func (ectx evaluationContext) ResolveIRI(v string) rdf.IRI {
	var vURL *iriutil.ParsedIRI

	if ectx.Base == nil {
		return rdf.IRI(v)
	} else if v == "" {
		// [spec 5.3] An empty same document reference "" resolves against the URI part of the base URI; any fragment part is ignored.

		vURL, _ = ectx.Base.Parse("")
		vURL.DropFragment()
	} else {
		var err error

		vURL, err = ectx.Base.Parse(v)
		if err != nil {
			// TODO warn
			return rdf.IRI(v)
		}
	}

	return rdf.IRI(vURL.String())
}
