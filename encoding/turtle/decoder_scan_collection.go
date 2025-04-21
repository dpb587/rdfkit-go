package turtle

import (
	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/turtle/internal/grammar"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
)

func reader_scan_collection(r *Decoder, ectx evaluationContext, r0 rune, openSubject rdf.SubjectValue, openSubjectRange *cursorio.TextOffsetRange) (readerStack, error) {
	if r0 == ')' {
		return r.emit(&statement{
			triple: rdf.Triple{
				Subject:   ectx.CurSubject,
				Predicate: ectx.CurPredicate,
				Object:    rdfiri.Nil_List,
			},
			offsets: r.buildTextOffsets(
				encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
				encoding.PredicateStatementOffsets, ectx.CurPredicateLocation,
				encoding.ObjectStatementOffsets, r.commitForTextOffsetRange([]rune{r0}),
			),
		})
	}

	r.buf.BacktrackRunes(r0)

	nectx := ectx
	nectx.CurSubject = openSubject
	nectx.CurSubjectLocation = openSubjectRange
	nectx.CurPredicate = rdfiri.First_Property
	nectx.CurPredicateLocation = nil

	r.pushState(nectx, reader_scan_collection_Continue)

	if ectx.CurSubject == nil {
		// collection as a subject
		return readerStack{nectx, reader_scan_Object}, nil
	}

	// TODO should emit immediately; but reader_scan_Object doesn't currently pop itself off the stack
	r.emit(&statement{
		triple: rdf.Triple{
			Subject:   ectx.CurSubject,
			Predicate: ectx.CurPredicate,
			Object:    openSubject,
		},
		offsets: r.buildTextOffsets(
			encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
			encoding.PredicateStatementOffsets, ectx.CurPredicateLocation,
			encoding.ObjectStatementOffsets, openSubjectRange,
		),
	})

	return readerStack{nectx, reader_scan_Object}, nil
}

func reader_scan_collection_Continue(r *Decoder, ectx evaluationContext, r0 rune, err error) (readerStack, error) {
	if err != nil {
		return readerStack{}, grammar.R_collection.Err(r.newOffsetError(err, nil, nil))
	}

	if r0 == ')' {
		return r.emit(&statement{
			triple: rdf.Triple{
				Subject:   ectx.CurSubject,
				Predicate: rdfiri.Rest_Property,
				Object:    rdfiri.Nil_List,
			},
			offsets: r.buildTextOffsets(
				encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
				encoding.ObjectStatementOffsets, r.commitForTextOffsetRange([]rune{r0}),
			),
		})
	}

	r.buf.BacktrackRunes(r0)

	nectx := ectx
	nectx.CurSubject = ectx.Global.BlankNodeFactory.NewBlankNode()
	nectx.CurSubjectLocation = nil

	r.pushState(nectx, reader_scan_collection_Continue)

	// TODO should emit immediately; but reader_scan_Object doesn't currently pop itself off the stack
	r.emit(&statement{
		triple: rdf.Triple{
			Subject:   ectx.CurSubject,
			Predicate: rdfiri.Rest_Property,
			Object:    nectx.CurSubject,
		},
		offsets: r.buildTextOffsets(
			encoding.SubjectStatementOffsets, ectx.CurSubjectLocation,
		),
	})

	return readerStack{nectx, reader_scan_Object}, nil
}
