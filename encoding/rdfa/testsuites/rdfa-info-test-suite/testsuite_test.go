package testsuite

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/dpb587/rdfkit-go/dev/earltestingutil"
	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/encoding/html"
	"github.com/dpb587/rdfkit-go/encoding/rdfa"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/ontology/earl/earliri"
	"github.com/dpb587/rdfkit-go/ontology/earl/earltesting"
	"github.com/dpb587/rdfkit-go/ontology/foaf/foafiri"
	"github.com/dpb587/rdfkit-go/ontology/rdf/rdfiri"
	"github.com/dpb587/rdfkit-go/ontology/xsd/xsdobject"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/rdfutil"
	"github.com/dpb587/rdfkit-go/rdf/triples"
	"github.com/dpb587/rdfkit-go/rdfdescription"
	"github.com/dpb587/rdfkit-go/testing/testingarchive"
	"github.com/dpb587/rdfkit-go/testing/testingassert"
	"github.com/dpb587/rdfkit-go/testing/testingutil"
	"github.com/dpb587/rdfkit-go/x/rdfdescriptionstruct"
)

const manifestPrefix = "http://rdfa.info/"

type expectedOverride struct {
	Expected   []byte
	OutcomeIRI rdf.IRI
	Message    string
}

func Test(t *testing.T) {
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
					rdf.IRI("https://github.com/dpb587/rdfkit-go/tree/main/encoding/rdfa/testsuites/rdfa-info-test-suite"),
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
					xsdobject.String("rdfkit-go/encoding/rdfa"),
				},
				foafiri.Homepage_Property: rdf.ObjectValueList{
					rdf.IRI("https://pkg.go.dev/github.com/dpb587/rdfkit-go/encoding/rdfa"),
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
	rdfioDebug := testingutil.NewDebugRdfioBuilderFromEnv(t)

	for _, profile := range []struct {
		rdfaVersion  string
		hostLanguage string
		overrides    map[rdf.IRI]expectedOverride
	}{
		{"rdfa1.1", "html4", map[rdf.IRI]expectedOverride{
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/html4/manifest#0180": {
				Message:    "expectation seems incorrect? empty @prefix definition may be ignored, but reference still resolves against default vocabulary",
				OutcomeIRI: earliri.Untested_NotTested,
			},
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/html4/manifest#0258": {
				Message:    "expectation seems incorrect? a warning may be emitted for the prefix, but the _:* syntax is still a valid blank node",
				OutcomeIRI: earliri.Untested_NotTested,
			},
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/html4/manifest#0295": {
				Message:    "expectation assumes an HTML DOM implementation which immediately closes an incorrectly self-closed span tag (with go, Line 468 creates its following statements to have a subject of #b)",
				OutcomeIRI: earliri.Untested_NotTested,
			},
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/html4/manifest#0319": {
				Message: "leaky assumption about querying of results which expects iri resolution",
				Expected: []byte(`<http://example.com/> <relative/iri#prop> "value" .
<http://example.com/> <relative/uri#prop> "value" .`),
			},
		}},
		{"rdfa1.1", "html5", map[rdf.IRI]expectedOverride{
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/html5/manifest#0283": {
				Message:    "expectation seems incorrect? html-rdfa says innerText is used as fallback for @datetime attr, then xsd:date uses whiteSpace=collapse, then it becomes a valid xsd:date",
				OutcomeIRI: earliri.Untested_NotTested,
			},
		}},
		{"rdfa1.1", "xhtml1", map[rdf.IRI]expectedOverride{
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/xhtml1/manifest#0113": {
				// also, case description in manifest.ttl seems incorrect?
				Message:    "expectation assumes an HTML DOM implementation (go tokenizer includes spaces until the outer closing tag since span should not be self-closing)",
				OutcomeIRI: earliri.Untested_NotTested,
			},
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/xhtml1/manifest#0180": {
				Message:    "expectation seems incorrect? empty @prefix definition may be ignored, but reference still resolves against default vocabulary",
				OutcomeIRI: earliri.Untested_NotTested,
			},
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/xhtml1/manifest#0198": {
				Message: "verified equivalent aside from differing physical representations of xml (we use Exclusive XML Canonicalization for outer nodes)",
				Expected: []byte(`<http://www.example.org/me#mark> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> .
<http://www.example.org/me#mark> <http://xmlns.com/foaf/0.1/firstName> "Mark" .
<http://www.example.org/me#mark> <http://xmlns.com/foaf/0.1/name> "<span xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:foaf=\"http://xmlns.com/foaf/0.1/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" property=\"foaf:firstName\">Mark</span> <span xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:foaf=\"http://xmlns.com/foaf/0.1/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" property=\"foaf:surname\">Birbeck</span>"^^<http://www.w3.org/1999/02/22-rdf-syntax-ns#XMLLiteral> .
<http://www.example.org/me#mark> <http://xmlns.com/foaf/0.1/surname> "Birbeck" .`),
			},
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/xhtml1/manifest#0258": {
				Message:    "expectation seems incorrect? a warning may be emitted for the prefix, but the _:* syntax is still a valid blank node",
				OutcomeIRI: earliri.Untested_NotTested,
			},
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/xhtml1/manifest#0295": {
				Message:    "expectation assumes an HTML DOM implementation which immediately closes an incorrectly self-closed span tag (with go, Line 468 creates its following statements to have a subject of #b)",
				OutcomeIRI: earliri.Untested_NotTested,
			},
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/xhtml1/manifest#0319": {
				Message: "leaky assumption about querying of results which expects iri resolution",
				Expected: []byte(`<http://example.com/> <relative/iri#prop> "value" .
<http://example.com/> <relative/uri#prop> "value" .`),
			},
		}},
		{"rdfa1.1", "xhtml5", map[rdf.IRI]expectedOverride{
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/xhtml5/manifest#0198": {
				Message: "verified equivalent aside from differing physical representations of xml (we use Exclusive XML Canonicalization for outer nodes)",
				Expected: []byte(`<http://www.example.org/me#mark> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> .
<http://www.example.org/me#mark> <http://xmlns.com/foaf/0.1/firstName> "Mark" .
<http://www.example.org/me#mark> <http://xmlns.com/foaf/0.1/name> "<span xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:foaf=\"http://xmlns.com/foaf/0.1/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" property=\"foaf:firstName\">Mark</span> <span xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:foaf=\"http://xmlns.com/foaf/0.1/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" property=\"foaf:surname\">Birbeck</span>"^^<http://www.w3.org/1999/02/22-rdf-syntax-ns#XMLLiteral> .
<http://www.example.org/me#mark> <http://xmlns.com/foaf/0.1/surname> "Birbeck" .`),
			},
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/xhtml5/manifest#0283": {
				Message:    "expectation seems incorrect? html-rdfa says innerText is used as fallback for @datetime attr, then xsd:date uses whiteSpace=collapse, then it becomes a valid xsd:date",
				OutcomeIRI: earliri.Untested_NotTested,
			},
			"http://rdfa.info/test-suite/test-cases/rdfa1.1/xhtml5/manifest#0319": {
				Message: "leaky assumption about querying of results which expects iri resolution",
				Expected: []byte(`<http://example.com/> <relative/iri#prop> "value" .
<http://example.com/> <relative/uri#prop> "value" .`),
			},
		}},
	} {
		testdata, testdataManifest := requireTestdata(t, profile.rdfaVersion, profile.hostLanguage)

		for _, entry := range testdataManifest.Entries {
			t.Run(string(entry.ID), func(t *testing.T) {
				tAssertion := earlReport.NewAssertion(t, entry.ID)

				decodeAction := func() (encodingtest.TripleStatementList, error) {
					htmlDocument, err := html.ParseDocument(
						testdata.NewFileByteReader(t, string(entry.Action)),
						html.DocumentConfig{}.
							SetLocation(string(entry.Action)).
							SetCaptureTextOffsets(true),
					)
					if err != nil {
						return nil, fmt.Errorf("parse document: %v", err)
					}

					opts := rdfa.DecoderConfig{}

					if profile.rdfaVersion == "rdfa1.1" && profile.hostLanguage == "html4" {
						opts = opts.SetHtmlProcessingProfile(rdfa.HTML4_RDFa11_HtmlProcessProfile)
					} else if profile.rdfaVersion == "rdfa1.1" && profile.hostLanguage == "html5" {
						opts = opts.SetHtmlProcessingProfile(rdfa.HTML5_RDFa11_HtmlProcessProfile)
					} else if profile.rdfaVersion == "rdfa1.1" && profile.hostLanguage == "xhtml1" {
						opts = opts.SetHtmlProcessingProfile(rdfa.XHTML1_RDFa11_HtmlProcessProfile)
					} else if profile.rdfaVersion == "rdfa1.1" && profile.hostLanguage == "xhtml5" {
						opts = opts.SetHtmlProcessingProfile(rdfa.XHTML5_RDFa11_HtmlProcessProfile)
					} else {
						t.Fatalf("unsupported profile combination: %s + %s", profile.rdfaVersion, profile.hostLanguage)
					}

					return encodingtest.CollectTripleStatementsErr(rdfa.NewDecoder(
						htmlDocument,
						opts,
					))
				}

				decodeResult := func() (rdf.TripleList, error) {
					return triples.CollectErr(turtle.NewDecoder(
						testdata.NewFileByteReader(t, string(entry.Result)),
					))
				}

				if override, ok := profile.overrides[entry.ID]; ok {
					tAssertion.SetMode(earliri.Manual_TestMode)
					tAssertion.Logf("manual review: %s", override.Message)

					if len(override.OutcomeIRI) > 0 {
						tAssertion.SetResultOutcome(override.OutcomeIRI)

						return
					}

					decodeResult = func() (rdf.TripleList, error) {
						return triples.CollectErr(turtle.NewDecoder(
							bytes.NewReader(override.Expected),
						))
					}
				}

				switch entry.Type {
				case "http://rdfa.info/vocabs/rdfa-test#PositiveEvaluationTest":
					expectedStatements, err := decodeResult()
					if err != nil {
						tAssertion.Fatalf("setup error: decode result: %v", err)
					}

					actualStatements, err := decodeAction()
					if err != nil {
						tAssertion.Fatalf("error: %v", err)
					}

					testingassert.IsomorphicGraphs(t.Context(), tAssertion, expectedStatements, actualStatements.AsTriples())

					rdfioDebug.PutTriplesBundle(t.Name(), actualStatements)
				case "http://rdfa.info/vocabs/rdfa-test#NegativeEvaluationTest":
					actualStatements, err := decodeAction()
					if err != nil {
						tAssertion.Fatalf("error: %v", err)
					} else if len(actualStatements) > 0 {
						testingassert.IsomorphicGraphs(t.Context(), tAssertion, nil, actualStatements.AsTriples())
					}
				default:
					t.Fatalf("unsupported test type: %s", entry.Type)
				}
			})
		}
	}
}

func requireTestdata(t *testing.T, rdfaVersion, hostLanguage string) (testingarchive.Archive, *Manifest) {
	testdata := testingarchive.OpenTarGz(
		t,
		"testdata.tar.gz",
		func(v string) string {
			return manifestPrefix + strings.TrimPrefix(v, "./")
		},
	)

	manifestBase := manifestPrefix + "test-suite/test-cases/" + rdfaVersion + "/" + hostLanguage + "/manifest"
	manifestResources := rdfdescription.NewResourceListBuilder()

	{
		manifestDecoder, err := turtle.NewDecoder(
			testdata.NewFileByteReader(t, manifestBase+".ttl"),
			turtle.DecoderConfig{}.
				SetDefaultBase(manifestBase+".ttl"),
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

	manifestResource := manifestResources.ExportResource(rdf.IRI(manifestBase), rdfdescription.DefaultExportResourceOptions)
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
	Result rdf.IRI `rdf:"o,p=http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#result"`
}
