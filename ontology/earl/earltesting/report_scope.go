package earltesting

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/dpb587/rdfkit-go/ontology/earl/earliri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdobject"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/rdfutil"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type ReportScope struct {
	report   *Report
	assertor rdf.SubjectValue
	subject  rdf.SubjectValue
}

func (rs ReportScope) GetReport() *Report {
	return rs.report
}

func (rs ReportScope) AddSubjectStatements(s rdf.SubjectValue, statements ...rdfdescription.Statement) {
	rs.report.mu.Lock()
	defer rs.report.mu.Unlock()

	rs.report.builder.Add(rdfdescription.SubjectResource{
		Subject:    s,
		Statements: statements,
	}.NewTriples()...)
}

func (rs ReportScope) WithAssertor(subjectValue rdf.SubjectValue, statements ...rdfdescription.Statement) ReportScope {
	rs.report.mu.Lock()
	defer rs.report.mu.Unlock()

	rs.report.builder.Add(rdf.Triple{
		Subject:   subjectValue,
		Predicate: rdfiri.Type_Property,
		Object:    earliri.Assertor_Class,
	})

	if len(statements) > 0 {
		rs.report.builder.Add(rdfdescription.StatementList(statements).NewTriples(subjectValue)...)
	}

	return ReportScope{
		report:   rs.report,
		assertor: subjectValue,
		subject:  rs.subject,
	}
}

func (rs ReportScope) WithSubject(subjectValue rdf.SubjectValue, statements ...rdfdescription.Statement) ReportScope {
	rs.report.mu.Lock()
	defer rs.report.mu.Unlock()

	rs.report.builder.Add(rdf.Triple{
		Subject:   subjectValue,
		Predicate: rdfiri.Type_Property,
		Object:    earliri.TestSubject_Class,
	})

	if len(statements) > 0 {
		rs.report.builder.Add(rdfdescription.StatementList(statements).NewTriples(subjectValue)...)
	}

	if rs.report.fromEnv {
		releaseStatements := rdfdescription.StatementList{}

		if v := os.Getenv("TESTING_EARL_SUBJECT_RELEASE_NAME"); len(v) > 0 {
			releaseStatements = append(releaseStatements, rdfdescription.StatementList{
				rdfdescription.ObjectStatement{
					Predicate: rdf.IRI("http://usefulinc.com/ns/doap#name"),
					Object:    xsdobject.String(v),
				},
			})
		}

		if v := os.Getenv("TESTING_EARL_SUBJECT_RELEASE_REVISION"); len(v) > 0 {
			releaseStatements = append(releaseStatements, rdfdescription.StatementList{
				rdfdescription.ObjectStatement{
					Predicate: rdf.IRI("http://usefulinc.com/ns/doap#revision"),
					Object:    xsdobject.String(v),
				},
			})
		}

		if v := os.Getenv("TESTING_EARL_SUBJECT_RELEASE_DATE"); len(v) > 0 {
			literalValue, err := rdfutil.CoalesceObjectValue(v, xsdobject.MapDateTime, xsdobject.MapDate)
			if err != nil {
				rs.report.t.Fatalf("configure: TESTING_EARL_SUBJECT_RELEASE_DATE: %v", err)
			}

			releaseStatements = append(releaseStatements, rdfdescription.ObjectStatement{
				Predicate: rdf.IRI("http://purl.org/dc/terms/created"),
				Object:    literalValue,
			})
		}

		rs.report.builder.Add(rdfdescription.AnonResourceStatement{
			Predicate: rdf.IRI("http://usefulinc.com/ns/doap#release"),
			AnonResource: rdfdescription.AnonResource{
				Statements: releaseStatements,
			},
		}.NewTriples(subjectValue)...)
	}

	return ReportScope{
		report:   rs.report,
		assertor: rs.assertor,
		subject:  subjectValue,
	}
}

func (rs ReportScope) NewAssertion(t *testing.T, test rdf.SubjectValue) *Assertion {
	rs.report.mu.Lock()
	defer rs.report.mu.Unlock()

	assertionNode := rdf.NewBlankNode()
	resultNode := rdf.NewBlankNode()

	assertion := &Assertion{
		rs:             rs,
		t:              t,
		test:           test,
		assertionNode:  assertionNode,
		resultNode:     resultNode,
		startTime:      time.Now(),
		descriptionLog: &bytes.Buffer{},
	}

	rs.report.assertions = append(rs.report.assertions, AssertionProfile{
		Test:        test,
		Node:        assertionNode,
		ResultNode:  resultNode,
		TestingName: t.Name(),
	})

	t.Cleanup(assertion.finalize)

	return assertion
}
