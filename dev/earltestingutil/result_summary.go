package earltestingutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/ontology/earl/earliri"
	"github.com/dpb587/rdfkit-go/ontology/earl/earltesting"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
	"github.com/dpb587/rdfkit-go/rdfdescription"
)

var DefaultReportSummaryOptions = ReportSummaryOptions{
	ExtraPredicates: rdf.PredicateValueList{
		rdf.IRI("http://purl.org/dc/terms/description"),
	},
}

type ReportSummaryOptions struct {
	HeaderComment   string
	ExtraPredicates rdf.PredicateValueList
}

type ReportSummary struct {
	report *earltesting.Report
	output string
	opts   ReportSummaryOptions
}

func ReportSummaryFromEnv(t *testing.T, rp earltesting.ReportProvider, opts ReportSummaryOptions) {
	fhPath := os.Getenv("TESTING_DEV_EARL_SUMMARY_OUTPUT")
	if len(fhPath) == 0 {
		return
	}

	t.Cleanup(func() {
		fh, err := os.OpenFile(fhPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			t.Fatalf("open debug file: %v", err)
		}

		defer fh.Close()

		err = WriteResultSummary(fh, rp.GetReport(), opts)
		if err != nil {
			t.Errorf("close debug rdfio writer: %v", err)
		}
	})
}

func WriteResultSummary(w io.Writer, report *earltesting.Report, opts ReportSummaryOptions) error {
	buf := &bytes.Buffer{}

	formatter := turtle.NewTermFormatter(turtle.TermFormatterOptions{
		Prefixes: iriutil.NewPrefixMap(
			iriutil.PrefixMapping{
				Prefix:   "earl",
				Expanded: earliri.Base,
			},
		),
	})

	builder := report.GetResourceListBuilder()

	if len(opts.HeaderComment) > 0 {
		fmt.Fprintf(buf, "# %s\n", strings.ReplaceAll(opts.HeaderComment, "\n", "\n# "))
		fmt.Fprintln(buf, "")
	}

	for assertionIdx, assertion := range report.GetAssertionProfiles() {
		if assertionIdx > 0 {
			fmt.Fprintln(buf, "")
		}

		fmt.Fprintln(buf, formatter.FormatTerm(assertion.Test))

		resultDescription := builder.ExportResourceStatements(
			assertion.ResultNode,
			rdfdescription.ExportResourceOptions{},
		).GroupByPredicate()

		// only supporting IRIs

		var outcomeIRI rdf.IRI

		if outcomeStatements, ok := resultDescription[earliri.Outcome_ObjectProperty]; ok {
			if statement, ok := outcomeStatements[0].(rdfdescription.ObjectStatement); ok {
				if iri, ok := statement.Object.(rdf.IRI); ok {
					outcomeIRI = iri
				}
			}
		}

		if len(outcomeIRI) > 0 {
			fmt.Fprintf(buf, "  %s %s\n", formatter.FormatTerm(earliri.Outcome_ObjectProperty), formatter.FormatTerm(outcomeIRI))
		} else {
			fmt.Fprintf(buf, "  %s %s\n", formatter.FormatTerm(earliri.Outcome_ObjectProperty), earliri.CantTell_CannotTell)
		}

		for _, predicate := range opts.ExtraPredicates {
			statements, ok := resultDescription[predicate]
			if !ok {
				continue
			}

			for _, statement := range statements {
				switch s := statement.(type) {
				case rdfdescription.ObjectStatement:
					fmt.Fprintf(buf, "  %s %s\n", formatter.FormatTerm(s.Predicate), formatter.FormatTerm(s.Object))
				}
			}
		}
	}

	_, err := w.Write(buf.Bytes())

	return err
}
