package testsuite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/dev/earltestingutil"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/internal/jsonldinternal"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/ontology/earl/earliri"
	"github.com/dpb587/rdfkit-go/ontology/earl/earltesting"
	"github.com/dpb587/rdfkit-go/ontology/foaf/foafiri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdobject"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/rdfutil"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/testing/testingarchive"
)

const manifestPrefix = "https://w3c.github.io/json-ld-api/tests/"

var manifestPrefixURL, _ = iri.ParseIRI(manifestPrefix)

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
					xsdobject.String("rdfkit-go w3c-github-json-ld-api-expand"),
				},
				foafiri.Homepage_Property: rdf.ObjectValueList{
					rdf.IRI("https://github.com/dpb587/rdfkit-go/tree/main/encoding/jsonld/internal/jsonldinternal/testsuites/w3c-github-json-ld-api-expand"),
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
					xsdobject.String("rdfkit-go/encoding/jsonld/internal/jsonldinternal"),
				},
				foafiri.Homepage_Property: rdf.ObjectValueList{
					rdf.IRI("https://pkg.go.dev/github.com/dpb587/rdfkit-go/encoding/jsonld/internal/jsonldinternal"),
				},
				rdf.IRI("http://usefulinc.com/ns/doap#programming-language"): rdf.ObjectValueList{
					xsdobject.String("Go"),
				},
				rdf.IRI("http://usefulinc.com/ns/doap#repository"): rdf.ObjectValueList{
					rdf.IRI("https://github.com/dpb587/rdfkit-go"),
				},
			})...,
		)

	earltestingutil.ReportSummaryFromEnv(t, earlReport, earltestingutil.DefaultReportSummaryOptions)

	for _, sequence := range testdataManifest.Sequences {
		t.Run(string(sequence.ID), func(t *testing.T) {
			tAssertion := earlReport.NewAssertion(t, sequence.ID)

			decodeAction := func() (jsonldinternal.ExpandedValue, error) {
				dopt := jsonldtype.ProcessorOptions{
					BaseURL: manifestPrefix + sequence.Input,
					DocumentLoader: jsonldtype.DocumentLoaderFunc(func(ctx context.Context, u string, opts jsonldtype.DocumentLoaderOptions) (jsonldtype.RemoteDocument, error) {
						if !testdata.HasFile(u) {
							return jsonldtype.RemoteDocument{}, fmt.Errorf("unknown url: %s", u)
						}

						doc, err := inspectjson.Parse(testdata.NewFileByteReader(t, u))
						if err != nil {
							return jsonldtype.RemoteDocument{}, fmt.Errorf("parse: %v", err)
						}

						docURL, err := url.Parse(u)
						if err != nil {
							return jsonldtype.RemoteDocument{}, fmt.Errorf("parse url: %v", err)
						}

						return jsonldtype.RemoteDocument{
							ContentType: "application/ld+json",
							Document:    doc,
							DocumentURL: docURL,
						}, nil
					}),
				}

				if len(sequence.Option.Base) > 0 {
					dopt.BaseURL = sequence.Option.Base
				}

				if len(sequence.Option.ProcessingMode) > 0 {
					dopt.ProcessingMode = sequence.Option.ProcessingMode
				} else if len(sequence.Option.SpecVersion) > 0 {
					dopt.ProcessingMode = sequence.Option.SpecVersion
				}

				if len(sequence.Option.ExpandContext) > 0 {
					expandContextURL, err := manifestPrefixURL.Parse(sequence.Option.ExpandContext)
					if err != nil {
						return nil, fmt.Errorf("parse: %v", err)
					}

					dopt.ExpandContext = inspectjson.StringValue{
						Value: expandContextURL.String(),
					}
				}

				parsedInput, err := inspectjson.Parse(testdata.NewFileByteReader(t, manifestPrefix+sequence.Input))
				if err != nil {
					return nil, fmt.Errorf("parse: %v", err)
				}

				return jsonldinternal.Expand(parsedInput, dopt)
			}

			if slices.Contains(sequence.Type, "jld:NegativeEvaluationTest") {
				_, err := decodeAction()
				if err != nil {
					tAssertion.Logf("error (expected): %v", err)
				} else {
					t.Fatal("expected error, but got none")
				}
			} else if slices.Contains(sequence.Type, "jld:PositiveEvaluationTest") {
				var expectedBuiltin any

				err := json.NewDecoder(testdata.NewFileByteReader(t, manifestPrefix+sequence.Expect)).Decode(&expectedBuiltin)
				if err != nil {
					tAssertion.Fatalf("unmarshal: %v", err)
				}

				actual, err := decodeAction()
				if err != nil {
					tAssertion.Fatalf("error: %v", err)
				}

				actualBuiltin := actual.AsBuiltin()

				if !reflect.DeepEqual(expectedBuiltin, actualBuiltin) {
					actualBuffer := &bytes.Buffer{}

					{
						actualEncoder := json.NewEncoder(actualBuffer)
						actualEncoder.SetIndent("", "  ")

						err = actualEncoder.Encode(actualBuiltin)
						if err != nil {
							t.Fatal(fmt.Errorf("marshal: %v", err))
						}
					}

					expectedBuffer := &bytes.Buffer{}

					{
						expectedEncoder := json.NewEncoder(expectedBuffer)
						expectedEncoder.SetIndent("", "  ")

						err := expectedEncoder.Encode(expectedBuiltin)
						if err != nil {
							t.Fatal(fmt.Errorf("marshal: %v", err))
						}
					}

					t.Fatal(strings.Join([]string{
						"error: unexpected result",
						"=== ACTUAL",
						strings.Join(strings.Split(actualBuffer.String(), "\n"), "\n"),
						"=== EXPECTED",
						strings.Join(strings.Split(expectedBuffer.String(), "\n"), "\n"),
					}, "\n\n"))
				}
			} else {
				t.Fatalf("unsupported test type: %v", sequence.Type)
			}
		})
	}
}

func requireTestdata(t *testing.T) (testingarchive.Archive, manifestSchema) {
	testdata := testingarchive.OpenTarGz(
		t,
		"testdata.tar.gz",
		func(v string) string {
			return manifestPrefix + strings.TrimPrefix(v, "./")
		},
	)

	var loadedManifest manifestSchema

	if err := json.Unmarshal(testdata.GetFileBytes(t, manifestPrefix+"expand-manifest.jsonld"), &loadedManifest); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	for sequenceIdx, sequence := range loadedManifest.Sequences {
		loadedManifest.Sequences[sequenceIdx].ID = rdf.IRI(manifestPrefix + "expand-manifest" + sequence.ID)
	}

	return testdata, loadedManifest
}
