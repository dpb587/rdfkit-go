package jsonld

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
)

type globalEvaluationContext struct {
	BlankNodeStringMapper blanknodeutil.StringMapper
	BlankNodeFactory      rdf.BlankNodeFactory
}

type evaluationContext struct {
	global *globalEvaluationContext

	CurrentContainer encoding.ContainerResource

	ActiveGraph      rdf.GraphNameValue
	ActiveGraphRange *cursorio.TextOffsetRange

	ActiveSubject      rdf.SubjectValue
	ActiveSubjectRange *cursorio.TextOffsetRange

	ActiveProperty      rdf.PredicateValue
	ActivePropertyRange *cursorio.TextOffsetRange

	Reverse bool
}
