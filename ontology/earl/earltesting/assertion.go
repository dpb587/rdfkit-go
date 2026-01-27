package earltesting

import (
	"bytes"
	"testing"
	"time"

	"github.com/dpb587/rdfkit-go/ontology/earl/earliri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type Assertion struct {
	rs             ReportScope
	t              *testing.T
	test           rdf.SubjectValue
	assertionNode  rdf.SubjectValue
	assertionMode  *rdf.IRI
	resultNode     rdf.SubjectValue
	resultOutcome  *rdf.IRI
	startTime      time.Time
	descriptionLog *bytes.Buffer
}

func (a *Assertion) SetMode(mode rdf.IRI) {
	a.assertionMode = &mode
}

func (a *Assertion) SetResultOutcome(outcome rdf.IRI) {
	a.resultOutcome = &outcome
}

func (a *Assertion) AddResultStatement(statements ...rdfdescription.Statement) {
	a.rs.report.builder.Add(rdfdescription.SubjectResource{
		Subject:    a.resultNode,
		Statements: statements,
	}.NewTriples()...)
}

func (a *Assertion) Skip(outcome rdf.IRI, args ...any) {
	a.t.Helper()
	a.SetResultOutcome(outcome)
	a.Log(args...)
	a.t.SkipNow()
}

func (a *Assertion) Skipf(outcome rdf.IRI, format string, args ...any) {
	a.t.Helper()
	a.SetResultOutcome(outcome)
	a.Logf(format, args...)
	a.t.SkipNow()
}

func (a *Assertion) finalize() {
	a.rs.report.mu.Lock()
	defer a.rs.report.mu.Unlock()

	a.rs.report.builder.Add(rdf.TripleList{
		{
			Subject:   a.assertionNode,
			Predicate: rdfiri.Type_Property,
			Object:    earliri.Assertion_Class,
		},
		{
			Subject:   a.assertionNode,
			Predicate: earliri.Test_ObjectProperty,
			Object:    a.test,
		},
		{
			Subject:   a.assertionNode,
			Predicate: earliri.Mode_ObjectProperty,
			Object:    a.getMode(),
		},
		{
			Subject:   a.assertionNode,
			Predicate: earliri.Result_ObjectProperty,
			Object:    a.resultNode,
		},
		{
			Subject:   a.resultNode,
			Predicate: rdfiri.Type_Property,
			Object:    earliri.TestResult_Class,
		},
		{
			Subject:   a.resultNode,
			Predicate: earliri.Outcome_ObjectProperty,
			Object:    a.getResultOutcome(),
		},
		{
			Subject:   a.resultNode,
			Predicate: rdf.IRI("http://purl.org/dc/terms/date"),
			Object: rdf.Literal{
				LexicalForm: a.startTime.In(time.UTC).Format(time.RFC3339),
				Datatype:    rdf.IRI("http://www.w3.org/2001/XMLSchema#dateTime"),
			},
		},
	}...)

	if a.rs.assertor != nil {
		a.rs.report.builder.Add(rdf.Triple{
			Subject:   a.assertionNode,
			Predicate: earliri.AssertedBy_ObjectProperty,
			Object:    a.rs.assertor,
		})
	}

	if a.rs.subject != nil {
		a.rs.report.builder.Add(rdf.Triple{
			Subject:   a.assertionNode,
			Predicate: earliri.Subject_ObjectProperty,
			Object:    a.rs.subject,
		})
	}

	if a.descriptionLog.Len() > 0 {
		a.rs.report.builder.Add(rdf.Triple{
			Subject:   a.resultNode,
			Predicate: rdf.IRI("http://purl.org/dc/terms/description"),
			Object: rdf.Literal{
				LexicalForm: a.descriptionLog.String(),
				Datatype:    xsdiri.String_Datatype,
			},
		})
	}
}

func (a *Assertion) getMode() rdf.IRI {
	if a.assertionMode != nil {
		return *a.assertionMode
	}

	return earliri.Automatic_TestMode
}

func (a *Assertion) getResultOutcome() rdf.IRI {
	if a.resultOutcome != nil {
		return *a.resultOutcome
	} else if a.t.Failed() {
		return earliri.Failed_Fail
	} else if a.t.Skipped() {
		return earliri.CantTell_CannotTell
	}

	return earliri.Passed_Pass
}
