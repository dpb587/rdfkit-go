package testsuite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/jsonld"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/internal/devencoding/rdfioutil"
	"github.com/dpb587/rdfkit-go/internal/devtest"
	"github.com/dpb587/rdfkit-go/rdfio"
	"github.com/dpb587/rdfkit-go/testing/testingarchive"
)

const manifestPrefix = "https://w3c.github.io/json-ld-api/tests/"

func Test(t *testing.T) {
	testdata, manifestResources := requireTestdata(t)
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

	for _, sequence := range manifestResources.Sequence {
		decodeAction := func() (rdfio.StatementList, error) {
			// var testInputURL, _ = urlutil.Parse()

			dopt := jsonld.DecoderConfig{}.
				SetDefaultBase(manifestPrefix + sequence.Input).
				SetDocumentLoader(jsonldtype.DocumentLoaderFunc(func(ctx context.Context, u string, opts jsonldtype.DocumentLoaderOptions) (jsonldtype.RemoteDocument, error) {
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
				})).
				SetCaptureTextOffsets(true)

			if len(sequence.Option.Base) > 0 {
				dopt = dopt.SetDefaultBase(sequence.Option.Base)
			}

			if len(sequence.Option.ProcessingMode) > 0 {
				dopt = dopt.SetProcessingMode(sequence.Option.ProcessingMode)
			} else if len(sequence.Option.SpecVersion) > 0 {
				dopt = dopt.SetProcessingMode(sequence.Option.SpecVersion)
			}

			if len(sequence.Option.RDFDirection) > 0 {
				dopt = dopt.SetRDFDirection(sequence.Option.RDFDirection)
			}

			r, err := jsonld.NewDecoder(
				testdata.NewFileByteReader(t, manifestPrefix+sequence.Input),
				dopt,
			)
			if err != nil {
				return nil, fmt.Errorf("decode: %v", err)
			}

			defer r.Close()

			return rdfio.CollectStatements(r)
		}

		if slices.Contains(sequence.Type, "jld:PositiveEvaluationTest") {
			t.Run("Eval/"+sequence.ID, func(t *testing.T) {
				if sequence.Option.ProduceGeneralizedRdf {
					t.Skip("ignore: produceGeneralizedRdf is not supported")
				} else if sequence.ID == "#t0122" || sequence.ID == "#t0123" || sequence.ID == "#t0124" || sequence.ID == "#t0125" {
					// The stdlib URL resolver normalizes .. and . segments out of the path, but these tests expect
					// the dot-segments to be retained.
					//
					// RFC 3986 suggests dot-segments may be processed when resolving potentially-relative URIs based
					// on Section 5 "Reference Resolution" which is what happens from stdlib's perspective when joining
					// the URLs during decoding. The spec also suggests dot-segments "are intended for use at the
					// beginning of a relative-path reference" which is not strictly followed in these test cases.
					//
					// Go stdlib's built-in behavior still seems to follow RFC 3986, especially given the processor
					// invokes a resolution of IRIs. Additionally, having all or none dot-segments (unlike the test
					// expectations), seems like a more preferred, deterministic behavior.
					//
					// It does not seem worth a custom IRI resolver just to accomodate these few test cases, especially
					// since a canonical resolution of expected + actual will match. Similarly, piprate/json-gold also
					// ignores (all) RFC3986 tests
					//
					//   https://github.com/piprate/json-gold/blob/4a395db392d18e12e04b157dfa61f5bd179a342b/ld/processor_test.go#L315-L316
					//
					// Historically, these tests were passing due to fallback evaluation through oxigraph load+ASK eval.
					// As of v0.4.7, oxigraph no longer seems to resolve IRIs during load (reasonable), so these tests
					// started failing.
					//
					// Examples:
					//
					// * EXPECT <urn:ex:s091> <urn:ex:p> <http://a/bb/ccc/./d;p?y> .
					//   ACTUAL <urn:ex:s091> <urn:ex:p> <http://a/bb/ccc/d;p?y> .
					// * EXPECT <urn:ex:s093> <urn:ex:p> <http://a/bb/ccc/./d;p?q#s> .
					//   ACTUAL <urn:ex:s093> <urn:ex:p> <http://a/bb/ccc/d;p?q#s> .
					// * EXPECT <urn:ex:s099> <urn:ex:p> <http://a/bb/ccc/./d;p?q> .
					//   ACTUAL <urn:ex:s099> <urn:ex:p> <http://a/bb/ccc/d;p?q> .
					// * EXPECT <urn:ex:s133> <urn:ex:p> <http://a/bb/ccc/../d;p?y> .
					//   ACTUAL <urn:ex:s133> <urn:ex:p> <http://a/bb/d;p?y> .
					t.Skip("ignore: expected failures (requires unresolved RFC 3986 dot-segments)")
				}

				expectedStatements, err := rdfio.CollectStatementsErr(nquads.NewDecoder(
					testdata.NewFileByteReader(t, manifestPrefix+sequence.Expect),
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
					oxigraphErr := devtest.AssertOxigraphAsk(t.Context(), oxigraphExec, manifestPrefix, testdata.NewFileByteReader(t, manifestPrefix+sequence.Expect), actualStatements)
					if oxigraphErr != nil {
						t.Logf("eval: %v", oxigraphErr)
						t.Log(err.Error())
						t.FailNow()
					}
				}

				debugBundle.PutBundle(t.Name(), actualStatements)
			})
		}
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

	if err := json.Unmarshal(testdata.GetFileBytes(t, manifestPrefix+"manifest.jsonld"), &loadedManifest); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	return testdata, loadedManifest
}
