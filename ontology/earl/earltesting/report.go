package earltesting

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/dpb587/rdfkit-go/encoding/ntriples"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/ontology/earl/earliri"
	"github.com/dpb587/rdfkit-go/ontology/foaf/foafiri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/rdfs/rdfsiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/rdfdescription/rdfdescriptionutil"
)

type Report struct {
	t       *testing.T
	fromEnv bool

	mu         sync.Mutex
	builder    *rdfdescription.ResourceListBuilder
	assertions []AssertionProfile
}

func NewReport(t *testing.T) *Report {
	r := &Report{
		t:       t,
		builder: rdfdescription.NewResourceListBuilder(),
	}

	return r
}

func NewReportFromEnv(t *testing.T) *Report {
	r := NewReport(t)
	r.fromEnv = true

	t.Cleanup(func() {
		filePath := os.Getenv("TESTING_EARL_OUTPUT")
		if len(filePath) == 0 {
			return
		}

		file, err := os.Create(filePath)
		if err != nil {
			t.Errorf("earltesting: failed to create file: %v", err)
		}

		defer file.Close()

		var encoder rdfdescriptionutil.ResourceEncoder

		if filepath.Ext(filepath.Base(filePath)) == ".ttl" {
			encoder, err = turtle.NewEncoder(file, turtle.EncoderConfig{}.
				SetBuffered(true).
				SetBufferedSort(false).
				SetPrefixes(iri.PrefixMappingList{
					// based on usage in conventional reports
					iri.PrefixMapping{Prefix: "dc", Expanded: "http://purl.org/dc/terms/"},
					iri.PrefixMapping{Prefix: "dc11", Expanded: "http://purl.org/dc/elements/1.1/"},
					iri.PrefixMapping{Prefix: "doap", Expanded: "http://usefulinc.com/ns/doap#"},
					iri.PrefixMapping{Prefix: "rdf", Expanded: string(rdfiri.Base)},
					iri.PrefixMapping{Prefix: "rdfs", Expanded: string(rdfsiri.Base)},
					iri.PrefixMapping{Prefix: "earl", Expanded: string(earliri.Base)},
					iri.PrefixMapping{Prefix: "foaf", Expanded: string(foafiri.Base)},
					iri.PrefixMapping{Prefix: "xsd", Expanded: string(xsdiri.Base)},
				}),
			)
			if err != nil {
				t.Errorf("earltesting: failed to create encoder[turtle]: %v", err)
			}
		} else {
			ntriplesEncoder, err := ntriples.NewEncoder(file)
			if err != nil {
				t.Errorf("earltesting: failed to create encoder[ntriples]: %v", err)
			}

			encoder = rdfdescriptionutil.NewTriplesResourceEncoder(ntriplesEncoder)
		}

		defer encoder.Close()

		if err := r.ToResourceWriter(t.Context(), encoder); err != nil {
			t.Errorf("earltesting: failed to export resources: %v", err)
		}
	})

	return r
}

func (r *Report) GetReport() *Report {
	return r
}

func (r *Report) GetResourceListBuilder() *rdfdescription.ResourceListBuilder {
	return r.builder
}

func (r *Report) GetAssertionProfiles() []AssertionProfile {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]AssertionProfile{}, r.assertions...)
}

func (r *Report) WithAssertor(subjectValue rdf.SubjectValue, statements ...rdfdescription.Statement) ReportScope {
	return ReportScope{report: r}.WithAssertor(subjectValue, statements...)
}

func (r *Report) WithSubject(subjectValue rdf.SubjectValue, statements ...rdfdescription.Statement) ReportScope {
	return ReportScope{report: r}.WithSubject(subjectValue, statements...)
}

func (r *Report) ToResourceWriter(ctx context.Context, w rdfdescription.ResourceWriter) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// prefer deterministic ordering based on assertion ordering (assuming tests are ordered)

	resourcesByAssertionSubject := map[rdf.SubjectValue]struct{}{}

	for _, profile := range r.assertions {
		if err := w.AddResource(ctx, r.builder.ExportResource(profile.Node, rdfdescription.DefaultExportResourceOptions)); err != nil {
			return fmt.Errorf("failed to add assertion resource: %w", err)
		}

		resourcesByAssertionSubject[profile.Node] = struct{}{}
	}

	// TODO deterministic order?

	for s := range r.builder.Subjects() {
		if _, ok := resourcesByAssertionSubject[s]; ok {
			continue
		} else if sBlankNode, ok := s.(rdf.BlankNode); ok {
			if r.builder.GetBlankNodeReferences(sBlankNode) == 1 {
				continue
			}
		}

		if err := w.AddResource(ctx, r.builder.ExportResource(s, rdfdescription.DefaultExportResourceOptions)); err != nil {
			return fmt.Errorf("failed to add resource: %w", err)
		}
	}

	return nil
}
