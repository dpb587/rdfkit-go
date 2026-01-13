package testsuite

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/ontology/earl/earliri"
	"github.com/dpb587/rdfkit-go/ontology/earl/earltesting"
	"github.com/dpb587/rdfkit-go/ontology/foaf/foafiri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdobject"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/rdfutil"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/testing/testingarchive"
	"github.com/dpb587/rdfkit-go/testing/testingutil"
	"github.com/dpb587/rdfkit-go/x/rdfdescriptionstruct"
)

const manifestPrefix = "http://www.w3.org/2013/N-QuadsTests/"

func Test(t *testing.T) {
	testdata, testdataManifest := requireTestdata(t)

	earlReport := earltesting.NewReportFromEnv(t).
		WithAssertor(
			rdf.IRI("#assertor"),
			rdfdescription.NewStatementsFromObjectsByPredicate(rdfutil.ObjectsByPredicate{
				rdfiri.Type_Property: rdf.ObjectValueList{
					earliri.Software_Class,
				},
				foafiri.Name_Property: rdf.ObjectValueList{
					xsdobject.String("rdfkit-go test suite"),
				},
				foafiri.Homepage_Property: rdf.ObjectValueList{
					rdf.IRI("https://github.com/dpb587/rdfkit-go/tree/main/encoding/nquads/testsuites/w3-2013-N-QuadsTests"),
				},
			})...,
		).
		WithSubject(
			rdf.IRI("#subject"),
			rdfdescription.NewStatementsFromObjectsByPredicate(rdfutil.ObjectsByPredicate{
				rdfiri.Type_Property: rdf.ObjectValueList{
					earliri.Software_Class,
					rdf.IRI("http://usefulinc.com/ns/doap#Project"),
				},
				foafiri.Name_Property: rdf.ObjectValueList{
					xsdobject.String("rdfkit-go/encoding/nquads"),
				},
				foafiri.Homepage_Property: rdf.ObjectValueList{
					rdf.IRI("https://pkg.go.dev/github.com/dpb587/rdfkit-go/encoding/nquads"),
				},
				rdf.IRI("http://usefulinc.com/ns/doap#programming-language"): rdf.ObjectValueList{
					xsdobject.String("Go"),
				},
				rdf.IRI("http://usefulinc.com/ns/doap#repository"): rdf.ObjectValueList{
					rdf.IRI("https://github.com/dpb587/rdfkit-go"),
				},
			})...,
		)

	rdfioDebug := testingutil.NewDebugRdfioBuilderFromEnv(t)

	for _, entry := range testdataManifest.Entries {
		t.Run(string(entry.ID), func(t *testing.T) {
			tAssertion := earlReport.NewAssertion(t, entry.ID)

			decodeAction := func() (encodingtest.QuadStatementList, error) {
				return encodingtest.CollectQuadStatementsErr(
					nquads.NewDecoder(
						testdata.NewFileByteReader(t, string(entry.Action)),
						nquads.DecoderConfig{}.
							SetCaptureTextOffsets(true),
					),
				)
			}

			switch entry.Type {
			case "http://www.w3.org/ns/rdftest#TestNQuadsNegativeSyntax":
				_, err := decodeAction()
				if err != nil {
					tAssertion.Logf("error (expected): %v", err)
				} else {
					t.Fatal("expected error, but got none")
				}
			case "http://www.w3.org/ns/rdftest#TestNQuadsPositiveSyntax":
				actualStatements, err := decodeAction()
				if err != nil {
					tAssertion.Fatalf("error: %v", err)
				}

				rdfioDebug.PutQuadsBundle(t.Name(), actualStatements)
			default:
				t.Fatalf("unsupported test type: %s", entry.Type)
			}
		})
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

	manifestResource := manifestResources.ExportResource(rdf.IRI(manifestPrefix+"manifest.ttl"), rdfdescription.DefaultExportResourceOptions)
	if len(manifestResource.GetResourceStatements()) == 0 {
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
