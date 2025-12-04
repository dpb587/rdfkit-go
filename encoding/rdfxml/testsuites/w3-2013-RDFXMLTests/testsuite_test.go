package testsuite

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/encoding/rdfxml"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/internal/devencoding/rdfioutil"
	"github.com/dpb587/rdfkit-go/internal/devtest"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/rdfio"
)

const manifestPrefix = "http://www.w3.org/2013/RDFXMLTests/"

func Test(t *testing.T) {
	archiveEntries, manifestResources := requireTestdata(t)
	oxigraphExec := os.Getenv("TESTING_OXIGRAPH_EXEC")

	var debugWriter = io.Discard
	var debugBundle *rdfioutil.BundleEncoder

	if fhPath := os.Getenv("TESTING_DEBUG_DUMPFILE"); len(fhPath) > 0 {
		fh, err := os.OpenFile(fhPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			t.Fatalf("open debug file: %v", err)
		}

		defer fh.Close()

		debugWriter = fh
	}

	debugBundle = rdfioutil.NewBundleEncoder(debugWriter)
	defer debugBundle.Close()

	for _, manifestResource := range manifestResources.GetResources() {
		var isEval, isNegativeSyntax bool
		var testName string
		var testAction, testResult rdf.IRI

		for _, triple := range manifestResource.AsTriples() {
			switch triple.Predicate.(rdf.IRI) {
			case rdfiri.Type_Property:
				if oIRI, ok := triple.Object.(rdf.IRI); ok {
					switch oIRI {
					case "http://www.w3.org/ns/rdftest#TestXMLEval":
						isEval = true
					case "http://www.w3.org/ns/rdftest#TestXMLNegativeSyntax":
						isNegativeSyntax = true
					}
				}
			case "http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#name":
				if oLiteral, ok := triple.Object.(rdf.Literal); ok {
					testName = oLiteral.LexicalForm
				}
			case "http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#action":
				if oIRI, ok := triple.Object.(rdf.IRI); ok {
					testAction = oIRI
				}
			case "http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#result":
				if oIRI, ok := triple.Object.(rdf.IRI); ok {
					testResult = oIRI
				}
			}
		}

		decodeAction := func() (rdfio.StatementList, error) {
			return rdfio.CollectStatementsErr(rdfxml.NewDecoder(
				bytes.NewReader(archiveEntries[string(testAction)]),
				rdfxml.DecoderConfig{}.
					SetBaseURL(string(testAction)).
					SetWarningListener(func(err error) {
						t.Logf("warn: %s", err.Error())
					}).
					SetCaptureTextOffsets(true),
			))
		}

		if isEval {
			t.Run("Eval/"+testName, func(t *testing.T) {
				expectedStatements, err := rdfio.CollectStatementsErr(nquads.NewDecoder(
					bytes.NewReader(archiveEntries[string(testResult)]),
				))
				if err != nil {
					t.Fatalf("setup error: decode result: %v", err)
				}

				actualStatements, err := decodeAction()
				if err != nil {
					t.Fatalf("error: %v", err)
				}

				err = devtest.AssertStatementEquals(expectedStatements, actualStatements)
				if err == nil {
					// good
				} else if len(oxigraphExec) == 0 {
					t.Log("eval: processor required, but TESTING_OXIGRAPH_EXEC is empty")
					t.Log(err.Error())
					t.SkipNow()
				} else {
					oxigraphErr := devtest.AssertOxigraphAsk(t.Context(), oxigraphExec, testAction, bytes.NewReader(archiveEntries[string(testResult)]), actualStatements)
					if oxigraphErr != nil {
						t.Logf("eval: %v", oxigraphErr)
						t.Log(err.Error())
						t.FailNow()
					}
				}

				debugBundle.PutBundle(t.Name(), actualStatements)
			})
		} else if isNegativeSyntax {
			t.Run("NegativeSyntax/"+testName, func(t *testing.T) {
				_, err := decodeAction()
				if err != nil {
					t.Logf("error: %v", err)
				} else {
					t.Fatal("expected error, but got none")
				}
			})
		}
	}
}

func requireTestdata(t *testing.T) (map[string][]byte, *rdfdescription.ResourceListBuilder) {
	archiveEntries, err := devtest.OpenArchiveTarGz(
		"testdata.tar.gz",
		func(v string) string {
			return manifestPrefix + strings.TrimPrefix(v, "./")
		},
	)
	if err != nil {
		t.Fatal(fmt.Errorf("testdata: %v", err))
	}

	manifestResources := rdfdescription.NewResourceListBuilder()

	{
		manifestDecoder, err := turtle.NewDecoder(
			bytes.NewReader(archiveEntries[manifestPrefix+"manifest.ttl"]),
			turtle.DecoderConfig{}.
				SetDefaultBase(manifestPrefix),
		)
		if err != nil {
			t.Fatal(fmt.Errorf("decode: %v", err))
		}

		defer manifestDecoder.Close()

		for manifestDecoder.Next() {
			manifestResources.AddTriple(manifestDecoder.GetTriple())
		}

		if err := manifestDecoder.Err(); err != nil {
			t.Fatal(fmt.Errorf("decode: %v", err))
		}
	}

	return archiveEntries, manifestResources
}
