package testsuite

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/encoding/ntriples"
	"github.com/dpb587/rdfkit-go/encoding/rdfxml"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/internal/devencoding/rdfioutil"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/triples"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/testing/testingarchive"
	"github.com/dpb587/rdfkit-go/testing/testingassert"
	"github.com/dpb587/rdfkit-go/x/rdfdescriptionstruct"
)

const manifestPrefix = "http://www.w3.org/2013/RDFXMLTests/"

func Test(t *testing.T) {
	testdata, manifest := requireTestdata(t)

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

	for _, entry := range manifest.Entries {
		decodeAction := func() (encodingtest.TripleStatementList, error) {
			return encodingtest.CollectTripleStatementsErr(rdfxml.NewDecoder(
				testdata.NewFileByteReader(t, string(entry.Action)),
				rdfxml.DecoderConfig{}.
					SetBaseURL(string(entry.Action)).
					SetWarningListener(func(err error) {
						t.Logf("warn: %s", err.Error())
					}).
					SetCaptureTextOffsets(true),
			))
		}

		switch entry.Type {
		case "http://www.w3.org/ns/rdftest#TestXMLEval":
			t.Run("Eval/"+entry.Name, func(t *testing.T) {
				expectedStatements, err := triples.CollectErr(ntriples.NewDecoder(
					testdata.NewFileByteReader(t, string(entry.Result)),
					ntriples.DecoderConfig{},
				))
				if err != nil {
					t.Fatalf("setup error: decode result: %v", err)
				}

				actualStatements, err := decodeAction()
				if err != nil {
					t.Fatalf("error: %v", err)
				}

				testingassert.IsomorphicGraphs(t, expectedStatements, actualStatements.AsTriples())

				debugBundle.PutTriplesBundle(t.Name(), actualStatements)
			})
		case "http://www.w3.org/ns/rdftest#TestXMLNegativeSyntax":
			t.Run("NegativeSyntax/"+entry.Name, func(t *testing.T) {
				_, err := decodeAction()
				if err != nil {
					t.Logf("error: %v", err)
				} else {
					t.Fatal("expected error, but got none")
				}
			})
		default:
			t.Fatalf("unsupported test type: %s", entry.Type)
		}
	}
}

func requireTestdata(t *testing.T) (testingarchive.Archive, *Manifest) {
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
				SetDefaultBase(manifestPrefix+"manifest.ttl"),
		)
		if err != nil {
			t.Fatal(fmt.Errorf("decode: %v", err))
		}

		defer manifestDecoder.Close()

		for manifestDecoder.Next() {
			manifestResources.Add(manifestDecoder.Triple())
		}

		if err := manifestDecoder.Err(); err != nil {
			t.Fatal(fmt.Errorf("decode: %v", err))
		}
	}

	manifest := &Manifest{}

	manifestResource, ok := manifestResources.GetResource(rdf.IRI(manifestPrefix + "manifest.ttl"))
	if !ok {
		t.Fatalf("manifest resource not found")
	} else if err := rdfdescriptionstruct.Unmarshal(manifestResources, manifestResource, manifest); err != nil {
		t.Fatalf("unmarshal manifest: %v", err)
	}

	return testdata, manifest
}

type Manifest struct {
	Entries rdfdescriptionstruct.Collection[ManifestEntry] `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#entries"`
}

type ManifestEntry struct {
	ID     rdf.IRI `rdf:"s"`
	Type   rdf.IRI `rdf:"o,p=http://www.w3.org/1999/02/22-rdf-syntax-ns#type"`
	Name   string  `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#name"`
	Action rdf.IRI `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#action"`
	Result rdf.IRI `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#result"`
}
