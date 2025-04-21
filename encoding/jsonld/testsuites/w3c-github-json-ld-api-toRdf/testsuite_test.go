package testsuite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/jsonld"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/internal/devtest"
	"github.com/dpb587/rdfkit-go/rdfio"
)

const manifestPrefix = "https://w3c.github.io/json-ld-api/tests/"

func Test(t *testing.T) {
	archiveEntries, manifestResources := requireTestdata(t)
	oxigraphExec := os.Getenv("TESTING_OXIGRAPH_EXEC")

	for _, sequence := range manifestResources.Sequence {
		decodeAction := func() (rdfio.StatementList, error) {
			// var testInputURL, _ = urlutil.Parse()

			dopt := jsonld.DecoderConfig{}.
				SetDefaultBase(manifestPrefix + sequence.Input).
				SetDocumentLoader(jsonldtype.DocumentLoaderFunc(func(ctx context.Context, u string, opts jsonldtype.DocumentLoaderOptions) (jsonldtype.RemoteDocument, error) {
					buf, ok := archiveEntries[u]
					if !ok {
						return jsonldtype.RemoteDocument{}, fmt.Errorf("unknown url: %s", u)
					}

					doc, err := inspectjson.Parse(bytes.NewReader(buf))
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
				}))

			if len(sequence.Option.Base) > 0 {
				dopt = dopt.SetDefaultBase(sequence.Option.Base)
			}

			if len(sequence.Option.ProcessingMode) > 0 {
				dopt = dopt.SetProcessingMode(sequence.Option.ProcessingMode)
			}

			if len(sequence.Option.RDFDirection) > 0 {
				dopt = dopt.SetRDFDirection(sequence.Option.RDFDirection)
			}

			r, err := jsonld.NewDecoder(
				bytes.NewReader(archiveEntries[manifestPrefix+sequence.Input]),
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
				}

				expectedStatements, err := rdfio.CollectStatementsErr(nquads.NewDecoder(
					bytes.NewReader(archiveEntries[manifestPrefix+sequence.Expect]),
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

				oxigraphErr := devtest.AssertOxigraphAsk(t.Context(), oxigraphExec, manifestPrefix, bytes.NewReader(archiveEntries[manifestPrefix+sequence.Expect]), actualStatements)
				if oxigraphErr != nil {
					t.Logf("eval: %v", oxigraphErr)
					t.Log(err.Error())
					t.FailNow()
				}
			})
		}
	}
}

func requireTestdata(t *testing.T) (map[string][]byte, manifestSchema) {
	archiveEntries, err := devtest.OpenArchiveTarGz(
		"testdata.tar.gz",
		func(v string) string {
			return manifestPrefix + strings.TrimPrefix(v, "./")
		},
	)
	if err != nil {
		t.Fatal(fmt.Errorf("testdata: %v", err))
	}

	var loadedManifest manifestSchema

	err = json.Unmarshal(archiveEntries[manifestPrefix+"manifest.jsonld"], &loadedManifest)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	return archiveEntries, loadedManifest
}
