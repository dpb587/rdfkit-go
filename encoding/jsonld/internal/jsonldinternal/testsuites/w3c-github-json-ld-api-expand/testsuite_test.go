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
	"github.com/dpb587/rdfkit-go/encoding/jsonld/internal/jsonldinternal"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/internal/devtest"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

const manifestPrefix = "https://w3c.github.io/json-ld-api/tests/"

var manifestPrefixURL, _ = iriutil.ParseIRI(manifestPrefix)

func Test(t *testing.T) {
	archiveEntries, manifestResources := requireTestdata(t)

	for _, sequence := range manifestResources.Sequence {
		decodeAction := func() (inspectjson.Value, error) {
			dopt := jsonldtype.ProcessorOptions{
				BaseURL: manifestPrefix + sequence.Input,
				DocumentLoader: jsonldtype.DocumentLoaderFunc(func(ctx context.Context, u string, opts jsonldtype.DocumentLoaderOptions) (jsonldtype.RemoteDocument, error) {
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

			parsedInput, err := inspectjson.Parse(bytes.NewReader(archiveEntries[manifestPrefix+sequence.Input]))
			if err != nil {
				return nil, fmt.Errorf("parse: %v", err)
			}

			return jsonldinternal.Expand(parsedInput, dopt)
		}

		if slices.Contains(sequence.Type, "jld:NegativeEvaluationTest") {
			t.Run("NegativeSyntax/"+strings.TrimPrefix(sequence.ID, "#"), func(t *testing.T) {
				_, err := decodeAction()
				if err != nil {
					t.Logf("error: %v", err)
				} else {
					t.Fatal("expected error, but got none")
				}
			})
		} else if slices.Contains(sequence.Type, "jld:PositiveEvaluationTest") {
			t.Run("Eval/"+strings.TrimPrefix(sequence.ID, "#"), func(t *testing.T) {
				var expectedBuiltin any

				err := json.NewDecoder(bytes.NewReader(archiveEntries[manifestPrefix+sequence.Expect])).Decode(&expectedBuiltin)
				if err != nil {
					t.Fatalf("unmarshal: %v", err)
				}

				actual, err := decodeAction()
				if err != nil {
					t.Fatalf("error: %v", err)
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
