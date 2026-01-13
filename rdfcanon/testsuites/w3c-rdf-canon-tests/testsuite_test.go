package testsuite

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/dev/earltestingutil"
	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/ontology/earl/earliri"
	"github.com/dpb587/rdfkit-go/ontology/earl/earltesting"
	"github.com/dpb587/rdfkit-go/ontology/foaf/foafiri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdobject"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodes"
	"github.com/dpb587/rdfkit-go/rdf/quads"
	"github.com/dpb587/rdfkit-go/rdf/rdfutil"
	"github.com/dpb587/rdfkit-go/rdfcanon"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/testing/testingarchive"
	"github.com/dpb587/rdfkit-go/testing/testingassert"
	"github.com/dpb587/rdfkit-go/x/rdfdescriptionstruct"
	"github.com/dpb587/rdfkit-go/x/storage/inmemory"
)

const manifestPrefix = "https://w3c.github.io/rdf-canon/tests/"

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
					rdf.IRI("https://github.com/dpb587/rdfkit-go/tree/main/rdfcanon/testsuites/w3c-rdf-canon-tests"),
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
					xsdobject.String("rdfkit-go/rdfcanon"),
				},
				foafiri.Homepage_Property: rdf.ObjectValueList{
					rdf.IRI("https://pkg.go.dev/github.com/dpb587/rdfkit-go/rdfcanon"),
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

	for _, entry := range testdataManifest.Entries {
		t.Run(string(entry.ID), func(t *testing.T) {
			tAssertion := earlReport.NewAssertion(t, entry.ID)

			canonicalizerOpts := []rdfcanon.CanonicalizeOption{}

			if entry.HashAlgorithm != nil {
				switch *entry.HashAlgorithm {
				case "SHA384":
					canonicalizerOpts = append(canonicalizerOpts, rdfcanon.CanonicalizeConfig{}.SetHashFunc(sha512.New384))
				default:
					t.Fatalf("unsupported hash algorithm: %s", *entry.HashAlgorithm)
				}
			}

			runCanonicalizer := func(ctx context.Context) (*rdfcanon.Canonicalization, error) {
				decoder, err := nquads.NewDecoder(
					bytes.NewReader(testdata.GetFileBytes(t, string(entry.Action))),
					nquads.DecoderConfig{},
				)
				if err != nil {
					return nil, fmt.Errorf("decode action: %v", err)
				}

				defer decoder.Close()

				// canonicalizer does not inherently dedup data
				dataset := inmemory.NewDataset()

				for decoder.Next() {
					err := dataset.AddQuad(ctx, decoder.Quad())
					if err != nil {
						return nil, fmt.Errorf("decode action: %v", err)
					}
				}

				if err := decoder.Err(); err != nil {
					return nil, fmt.Errorf("decode action: %v", err)
				}

				datasetIter, err := dataset.NewQuadIterator(ctx)
				if err != nil {
					return nil, fmt.Errorf("create dataset iterator: %v", err)
				}

				return rdfcanon.Canonicalize(t.Context(), datasetIter, canonicalizerOpts...)
			}

			switch entry.Type {
			case "https://w3c.github.io/rdf-canon/tests/vocab#RDFC10EvalTest":
				expectedStatements, err := quads.CollectErr(nquads.NewDecoder(
					bytes.NewReader(testdata.GetFileBytes(t, string(entry.Result))),
					nquads.DecoderConfig{},
				))
				if err != nil {
					tAssertion.Fatalf("decode expected result: %v", err)
				}

				canonicalization, err := runCanonicalizer(t.Context())
				if err != nil {
					tAssertion.Fatalf("canonicalize: %v", err)
				}

				actualBuffer := &bytes.Buffer{}
				_, err = canonicalization.WriteTo(actualBuffer)
				if err != nil {
					tAssertion.Fatalf("write canonicalized: %v", err)
				}

				actualStatements, err := quads.CollectErr(nquads.NewDecoder(
					bytes.NewReader(actualBuffer.Bytes()),
					nquads.DecoderConfig{},
				))
				if err != nil {
					tAssertion.Fatalf("collect: %v", err)
				}

				testingassert.IsomorphicDatasets(t.Context(), tAssertion, expectedStatements, actualStatements)
			case "https://w3c.github.io/rdf-canon/tests/vocab#RDFC10MapTest":
				if entry.Result == "" {
					t.Fatal("missing test result")
				}

				var expectedMap map[string]string
				err := json.Unmarshal(testdata.GetFileBytes(t, string(entry.Result)), &expectedMap)
				if err != nil {
					tAssertion.Fatalf("decode expected map: %v", err)
				}

				// Create a custom StringMapper that tracks the original string identifiers
				mapper := NewTrackingStringMapper()
				inputDecoder, err := nquads.NewDecoder(
					bytes.NewReader(testdata.GetFileBytes(t, string(entry.Action))),
					nquads.DecoderConfig{}.SetBlankNodeStringFactory(mapper),
				)
				if err != nil {
					tAssertion.Fatalf("decode action: %v", err)
				}

				canonicalization, err := rdfcanon.Canonicalize(t.Context(), inputDecoder, canonicalizerOpts...)
				if err != nil {
					tAssertion.Fatalf("canonicalize: %v", err)
				}

				// Build the actual map from original identifiers to canonical identifiers
				actualMap := make(map[string]string)
				for origID, bn := range mapper.mappings {
					actualMap[origID] = canonicalization.GetBlankNodeIdentifier(bn)
				}

				if len(expectedMap) != len(actualMap) {
					tAssertion.Fatalf("map size mismatch: expected %d, got %d", len(expectedMap), len(actualMap))
				}

				for expectedBN, expectedCanonical := range expectedMap {
					actualCanonical, ok := actualMap[expectedBN]
					if !ok {
						tAssertion.Fatalf("missing blank node in actual map: %s", expectedBN)
					}
					if expectedCanonical != actualCanonical {
						tAssertion.Fatalf("blank node %s: expected %s, got %s", expectedBN, expectedCanonical, actualCanonical)
					}
				}
			case "https://w3c.github.io/rdf-canon/tests/vocab#RDFC10NegativeEvalTest":
				_, err := runCanonicalizer(t.Context())
				if err != nil {
					tAssertion.Logf("error (expected): %v", err)
				} else {
					tAssertion.Fatalf("expected error, but got none")
				}
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
			return manifestPrefix + strings.TrimPrefix(v, "./tests/")
		},
	)

	manifestResources := rdfdescription.NewResourceListBuilder()

	{
		manifestDecoder, err := turtle.NewDecoder(
			bytes.NewReader(testdata.GetFileBytes(t, manifestPrefix+"manifest.ttl")),
			turtle.DecoderConfig{}.
				SetDefaultBase(manifestPrefix+"manifest.ttl"),
		)
		if err != nil {
			t.Fatal(fmt.Errorf("decode manifest: %v", err))
		}

		defer manifestDecoder.Close()

		for manifestDecoder.Next() {
			manifestResources.Add(manifestDecoder.Triple())
		}

		if err := manifestDecoder.Err(); err != nil {
			t.Fatal(fmt.Errorf("decode manifest: %v", err))
		}
	}

	manifest := &Manifest{}

	manifestResource := manifestResources.ExportResource(rdf.IRI(manifestPrefix+"manifest"), rdfdescription.DefaultExportResourceOptions)
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
	ID            rdf.IRI `rdf:"s"`
	Type          rdf.IRI `rdf:"o,p=http://www.w3.org/1999/02/22-rdf-syntax-ns#type"`
	Name          string  `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#name"`
	Action        rdf.IRI `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#action"`
	Result        rdf.IRI `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#result"`
	HashAlgorithm *string `rdf:"o,p=https://w3c.github.io/rdf-canon/tests/vocab#hashAlgorithm"`
}

// TrackingStringMapper is a custom StringMapper that tracks the original string identifiers
// so we can verify the blank node mappings after canonicalization
type TrackingStringMapper struct {
	factory  blanknodes.StringFactory
	mappings map[string]rdf.BlankNode
}

func NewTrackingStringMapper() *TrackingStringMapper {
	return &TrackingStringMapper{
		factory:  blanknodes.NewStringFactory(),
		mappings: make(map[string]rdf.BlankNode),
	}
}

func (t *TrackingStringMapper) NewBlankNode() rdf.BlankNode {
	return t.factory.NewBlankNode()
}

func (t *TrackingStringMapper) NewStringBlankNode(v string) rdf.BlankNode {
	if bn, ok := t.mappings[v]; ok {
		return bn
	}
	bn := t.factory.NewStringBlankNode(v)
	t.mappings[v] = bn
	return bn
}
