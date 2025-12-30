package testsuite

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/encoding/ntriples"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/testing/testingarchive"
	"github.com/dpb587/rdfkit-go/testing/testingutil"
	"github.com/dpb587/rdfkit-go/x/rdfdescriptionstruct"
)

const manifestPrefix = "http://www.w3.org/2013/N-TriplesTests/"

func Test(t *testing.T) {
	testdata, manifest := requireTestdata(t)

	rdfioDebug := testingutil.NewDebugRdfioBuilderFromEnv(t)

	for _, entry := range manifest.Entries {
		decodeAction := func() (encodingtest.TripleStatementList, error) {
			return encodingtest.CollectTripleStatementsErr(ntriples.NewDecoder(
				testdata.NewFileByteReader(t, string(entry.Action)),
				ntriples.DecoderConfig{}.
					SetCaptureTextOffsets(true),
			))
		}

		switch entry.Type {
		case "http://www.w3.org/ns/rdftest#TestNTriplesNegativeSyntax":
			t.Run("NegativeSyntax/"+entry.Name, func(t *testing.T) {
				_, err := decodeAction()
				if err != nil {
					t.Logf("error: %v", err)
				} else {
					t.Fatal("expected error, but got none")
				}
			})
		case "http://www.w3.org/ns/rdftest#TestNTriplesPositiveSyntax":
			t.Run("PositiveSyntax/"+entry.Name, func(t *testing.T) {
				actualStatements, err := decodeAction()
				if err != nil {
					t.Fatalf("error: %v", err)
				}

				rdfioDebug.PutTriplesBundle(t.Name(), actualStatements)
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
}
