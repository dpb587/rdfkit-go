package testsuite

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/encoding/trig"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/internal/devtest"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/rdfio"
)

const manifestPrefix = "http://www.w3.org/2013/TriGTests/"

func Test(t *testing.T) {
	archiveEntries, manifestResources := requireTestdata(t)
	oxigraphExec := os.Getenv("TESTING_OXIGRAPH_EXEC")

	for _, manifestResource := range manifestResources.GetResources() {
		var isEval, isNegativeSyntax, isPositiveSyntax bool
		var testName string
		var testAction, testResult rdf.IRI

		for _, triple := range manifestResource.AsTriples() {
			switch triple.Predicate.(rdf.IRI) {
			case rdfiri.Type_Property:
				if oIRI, ok := triple.Object.(rdf.IRI); ok {
					switch oIRI {
					case "http://www.w3.org/ns/rdftest#TestTrigEval":
						isEval = true
					case "http://www.w3.org/ns/rdftest#TestTrigNegativeSyntax":
						isNegativeSyntax = true
					case "http://www.w3.org/ns/rdftest#TestTrigPositiveSyntax":
						isPositiveSyntax = true
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

		if testName == "trig-syntax-bad-num-05" {
			// duplicate naming; disambiguate

			h := sha256.New()
			h.Write(archiveEntries[string(testAction)])

			testName += fmt.Sprintf("/%s", base64.RawStdEncoding.EncodeToString(h.Sum(nil))[0:8])
		}

		decodeAction := func() (rdfio.StatementList, error) {
			r, err := trig.NewDecoder(
				bytes.NewReader(archiveEntries[string(testAction)]),
				trig.DecoderConfig{}.
					SetDefaultBase(string(testAction)),
			)
			if err != nil {
				return nil, fmt.Errorf("decode action: %v", err)
			}

			defer r.Close()

			return rdfio.CollectStatements(r)
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
					return
				} else if len(oxigraphExec) == 0 {
					t.Log("eval: processor required, but TESTING_OXIGRAPH_EXEC is empty")
					t.Log(err.Error())
					t.SkipNow()
				}

				oxigraphErr := devtest.AssertOxigraphAsk(t.Context(), oxigraphExec, testAction, bytes.NewReader(archiveEntries[string(testResult)]), actualStatements)
				if oxigraphErr != nil {
					t.Logf("eval: %v", oxigraphErr)
					t.Log(err.Error())
					t.FailNow()
				}
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
		} else if isPositiveSyntax {
			t.Run("PositiveSyntax/"+testName, func(t *testing.T) {
				_, err := decodeAction()
				if err != nil {
					t.Fatalf("error: %v", err)
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
