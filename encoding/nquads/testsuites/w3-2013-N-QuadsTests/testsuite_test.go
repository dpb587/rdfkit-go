package testsuite

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/internal/devencoding/rdfioutil"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/dpb587/rdfkit-go/testing/testingarchive"
)

const manifestPrefix = "http://www.w3.org/2013/N-QuadsTests/"

func Test(t *testing.T) {
	testdata, manifestResources := requireTestdata(t)

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
		var isNegativeSyntax, isPositiveSyntax bool
		var testName string
		var testAction rdf.IRI

		for _, triple := range manifestResource.AsTriples() {
			switch triple.Predicate.(rdf.IRI) {
			case rdfiri.Type_Property:
				if oIRI, ok := triple.Object.(rdf.IRI); ok {
					switch oIRI {
					case "http://www.w3.org/ns/rdftest#TestNQuadsNegativeSyntax":
						isNegativeSyntax = true
					case "http://www.w3.org/ns/rdftest#TestNQuadsPositiveSyntax":
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
			}
		}

		decodeAction := func() (rdfio.StatementList, error) {
			return rdfio.CollectStatementsErr(
				nquads.NewDecoder(
					testdata.NewFileByteReader(t, string(testAction)),
					nquads.DecoderConfig{}.
						SetCaptureTextOffsets(true),
				),
			)
		}

		if isNegativeSyntax {
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
				actualStatements, err := decodeAction()
				if err != nil {
					t.Fatalf("error: %v", err)
				}

				debugBundle.PutBundle(t.Name(), actualStatements)
			})
		}
	}
}

func requireTestdata(t *testing.T) (testingarchive.Archive, *rdfdescription.ResourceListBuilder) {
	testdata := testingarchive.OpenTarGz(
		t,
		"testdata.tar.gz",
		func(v string) string {
			return manifestPrefix + strings.TrimPrefix(v, "./")
		},
	)

	manifestResources := rdfdescription.NewResourceListBuilder()

	{
		manifestDecoder, err := turtle.NewDecoder(
			testdata.NewFileByteReader(t, manifestPrefix+"manifest.ttl"),
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

	return testdata, manifestResources
}
