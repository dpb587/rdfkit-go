package earltesting

import (
	"context"
	"fmt"
	"os"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/ontology/earl/earliri"
	"github.com/dpb587/rdfkit-go/ontology/foaf/foafiri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdliteral"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

type Report struct {
	t      *testing.T
	output string

	mu                sync.Mutex
	builder           *rdfdescription.ResourceListBuilder
	assertionSubjects rdf.SubjectValueList
}

type ReportScope struct {
	report   *Report
	assertor rdf.SubjectValue
	subject  rdf.SubjectValue
}

func NewReport(t *testing.T) ReportScope {
	report := &Report{
		t:       t,
		builder: rdfdescription.NewResourceListBuilder(),
	}

	t.Cleanup(func() {
		if report.output == "" {
			return
		}

		if err := report.writeOutput(); err != nil {
			t.Errorf("failed to write EARL report: %v", err)
		}
	})

	return ReportScope{
		report: report,
	}
}

func NewReportFromEnv(t *testing.T) ReportScope {
	r := NewReport(t)

	if v := os.Getenv("TESTING_EARL_OUTPUT"); len(v) > 0 {
		r.SetOutput(v)
	}

	if v := os.Getenv("TESTING_EARL_SUBJECT_RELEASE_REVISION"); len(v) > 0 {
		releaseStatements := rdfdescription.StatementList{
			rdfdescription.ObjectStatement{
				Predicate: rdf.IRI("http://usefulinc.com/ns/doap#name"),
				Object:    xsdliteral.NewString("jelly-rdf-go"),
			},
			rdfdescription.ObjectStatement{
				Predicate: rdf.IRI("http://usefulinc.com/ns/doap#revision"),
				Object:    xsdliteral.NewString(v),
			},
		}

		if v := os.Getenv("TESTING_EARL_SUBJECT_RELEASE_DATE"); len(v) > 0 {
			if value, err := xsdliteral.MapDateTime(v); err == nil {
				releaseStatements = append(releaseStatements, rdfdescription.ObjectStatement{
					Predicate: rdf.IRI("http://purl.org/dc/terms/created"),
					Object:    value.AsLiteralTerm(),
				})
			} else if value, err := xsdliteral.MapDate(v); err == nil {
				releaseStatements = append(releaseStatements, rdfdescription.ObjectStatement{
					Predicate: rdf.IRI("http://purl.org/dc/terms/created"),
					Object:    value.AsLiteralTerm(),
				})
			} else {
				t.Fatalf("configure: TESTING_EARL_SUBJECT_RELEASE_DATE: %v", err)
			}
		}

		r.AddSubjectStatements(rdf.IRI("#subject"), rdfdescription.AnonResourceStatement{
			Predicate: rdf.IRI("http://usefulinc.com/ns/doap#release"),
			AnonResource: rdfdescription.AnonResource{
				Statements: releaseStatements,
			},
		})
	}

	return r
}

func (rs ReportScope) SetOutput(path string) {
	rs.report.output = path
}

func (rs ReportScope) AddSubjectStatements(s rdf.SubjectValue, statements ...rdfdescription.Statement) {
	rs.report.mu.Lock()
	defer rs.report.mu.Unlock()

	for _, t := range (rdfdescription.SubjectResource{
		Subject:    s,
		Statements: statements,
	}.AsTriples()) {
		rs.report.builder.AddTriple(t)
	}
}

func (rs ReportScope) WithAssertor(subjectValue rdf.SubjectValue, statements ...rdfdescription.Statement) ReportScope {
	rs.report.mu.Lock()
	defer rs.report.mu.Unlock()

	rs.report.builder.AddTriple(rdf.Triple{
		Subject:   subjectValue,
		Predicate: rdfiri.Type_Property,
		Object:    earliri.Assertor_Class,
	})

	if len(statements) > 0 {
		for _, t := range rdfdescription.StatementList(statements).NewTriples(subjectValue) {
			rs.report.builder.AddTriple(t)
		}
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

	rs.report.builder.AddTriple(rdf.Triple{
		Subject:   subjectValue,
		Predicate: rdfiri.Type_Property,
		Object:    earliri.TestSubject_Class,
	})

	if len(statements) > 0 {
		for _, t := range rdfdescription.StatementList(statements).NewTriples(subjectValue) {
			rs.report.builder.AddTriple(t)
		}
	}

	return ReportScope{
		report:   rs.report,
		assertor: rs.assertor,
		subject:  subjectValue,
	}
}

func (rs ReportScope) NewAssertion(t *testing.T, testIRI rdf.IRI) *Assertion {
	rs.report.mu.Lock()
	defer rs.report.mu.Unlock()

	assertion := &Assertion{
		rs:            rs,
		t:             t,
		testIRI:       testIRI,
		assertionNode: rdf.NewBlankNode(),
		resultNode:    rdf.NewBlankNode(),
		startTime:     time.Now(),
	}

	rs.report.assertionSubjects = append(rs.report.assertionSubjects, assertion.assertionNode)

	t.Cleanup(func() {
		assertion.finalize()
	})

	return assertion
}

func (r *Report) writeOutput() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.Create(r.output)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder, err := turtle.NewEncoder(file, turtle.EncoderConfig{}.
		SetBuffered(true).
		SetPrefixes(iriutil.NewPrefixMap(
			iriutil.PrefixMapping{Prefix: "dc", Expanded: "http://purl.org/dc/terms/"},
			iriutil.PrefixMapping{Prefix: "doap", Expanded: "http://usefulinc.com/ns/doap#"},
			iriutil.PrefixMapping{Prefix: "earl", Expanded: earliri.Base},
			iriutil.PrefixMapping{Prefix: "foaf", Expanded: foafiri.Base},
			iriutil.PrefixMapping{Prefix: "xsd", Expanded: xsdiri.Base},
		)),
	)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}

	ctx := context.Background()
	resources := r.builder.GetResources()

	for _, assertionSubject := range r.assertionSubjects {
		for _, resource := range resources {
			if resource.GetResourceSubject() == assertionSubject {
				if err := encoder.PutResource(ctx, resource); err != nil {
					return fmt.Errorf("failed to add resource: %w", err)
				}

				break
			}
		}
	}

	for _, resource := range resources {
		if slices.Contains(r.assertionSubjects, resource.GetResourceSubject()) {
			continue
		}

		if err := encoder.PutResource(ctx, resource); err != nil {
			return fmt.Errorf("failed to add resource: %w", err)
		}
	}

	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close encoder: %w", err)
	}

	return nil
}
