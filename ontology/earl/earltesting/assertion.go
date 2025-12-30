package earltesting

import (
	"testing"
	"time"

	"github.com/dpb587/rdfkit-go/ontology/earl/earliri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type Assertion struct {
	rs                ReportScope
	t                 *testing.T
	testIRI           rdf.IRI
	assertionNode     rdf.SubjectValue
	resultNode        rdf.SubjectValue
	startTime         time.Time
	testResultOutcome *rdf.IRI
}

func (a *Assertion) SetTestResultOutcome(outcome rdf.IRI) {
	a.testResultOutcome = &outcome
}

func (a *Assertion) AddTestResultDescription(description string) {
	a.AddTestResultStatement(rdfdescription.ObjectStatement{
		Predicate: rdf.IRI("http://purl.org/dc/terms/description"),
		Object: rdf.Literal{
			LexicalForm: description,
			Datatype:    xsdiri.String_Datatype,
		},
	})
}

func (a *Assertion) AddTestResultStatement(statements ...rdfdescription.Statement) {
	a.rs.report.builder.Add(rdfdescription.SubjectResource{
		Subject:    a.resultNode,
		Statements: statements,
	}.NewTriples()...)
}

func (a *Assertion) TSkip(t *testing.T, outcome rdf.IRI, description string) {
	a.SetTestResultOutcome(outcome)
	a.AddTestResultDescription(description)

	t.Helper()
	t.Skip(description)
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
			Object:    a.testIRI,
		},
		{
			Subject:   a.assertionNode,
			Predicate: earliri.Mode_ObjectProperty,
			Object:    earliri.Automatic_TestMode,
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
			Object:    a.getTestResultOutcome(),
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
}

func (a *Assertion) getTestResultOutcome() rdf.IRI {
	if a.testResultOutcome != nil {
		return *a.testResultOutcome
	} else if a.t.Failed() {
		return earliri.Failed_Fail
	} else if a.t.Skipped() {
		return earliri.CantTell_CannotTell
	}

	return earliri.Passed_Pass
}
