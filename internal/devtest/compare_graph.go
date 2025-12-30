package devtest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/internal/oxigraph"
	"github.com/dpb587/rdfkit-go/internal/oxigraph/nquadsutil"
	"github.com/dpb587/rdfkit-go/rdf"
)

func AssertOxigraphAsk(ctx context.Context, exec string, expectedBase rdf.IRI, expectedNquads io.Reader, actual rdf.QuadList) error {
	actualEncoded, err := func() ([]byte, error) {
		buf := bytes.NewBuffer(nil)
		encoder, err := nquads.NewEncoder(buf)
		if err != nil {
			return nil, fmt.Errorf("unexpected error: %v", err)
		}

		for _, statement := range actual {
			if err := encoder.AddQuad(ctx, statement); err != nil {
				return nil, fmt.Errorf("unexpected error: %v", err)
			}
		}

		if err := encoder.Close(); err != nil {
			return nil, fmt.Errorf("unexpected error: %v", err)
		}

		return buf.Bytes(), nil
	}()
	if err != nil {
		return fmt.Errorf("encode: %v", err)
	}

	sparqlEval := oxigraph.NewService(oxigraph.ServiceOptions{
		Exec: exec,
	})

	defer sparqlEval.Close()

	err = sparqlEval.ImportReader(bytes.NewReader(actualEncoded), oxigraph.NQuadsFormat, oxigraph.ImportOptions{
		// Base: expectedBase,
	})
	if err != nil {
		return fmt.Errorf("import: %v", err)
	}

	//

	sparqlEvalClient, err := sparqlEval.NewClient()
	if err != nil {
		return fmt.Errorf("client: %v", err)
	}

	askBuffer, err := io.ReadAll(expectedNquads)
	if err != nil {
		return fmt.Errorf("read: %v", err)
	}

	askTransformBytes, err := nquadsutil.NewSpaqlAsk(bytes.NewReader(askBuffer))
	if err != nil {
		return fmt.Errorf("transform: %v", err)
	}

	res, err := sparqlEvalClient.Query(
		ctx,
		"ASK { "+string(askTransformBytes)+" }",
	)
	if err != nil {
		return fmt.Errorf("query: %v", err)
	} else if !*res.Boolean {
		return fmt.Errorf("query did not match")
	}

	return nil
}

// eventually should be a proper, isomorphic comparison
func AssertStatementEquals(expected, actual rdf.QuadList) error {
	var encodedStatements = [2]*bytes.Buffer{
		bytes.NewBuffer(nil),
		bytes.NewBuffer(nil),
	}

	for i, statements := range [2]rdf.QuadList{expected, actual} {
		ctx := context.Background()
		encoder, err := nquads.NewEncoder(encodedStatements[i])
		if err != nil {
			return fmt.Errorf("unexpected error: %v", err)
		}

		for _, statement := range statements {
			if err := encoder.AddQuad(ctx, statement); err != nil {
				return fmt.Errorf("unexpected error: %v", err)
			}
		}

		if err := encoder.Close(); err != nil {
			return fmt.Errorf("unexpected error: %v", err)
		}
	}

	var expectedStatementsList = strings.Split(encodedStatements[0].String(), "\n")
	slices.SortFunc(expectedStatementsList, strings.Compare)
	expectedStatements := strings.Join(expectedStatementsList, "\n")

	var actualStatementsList = strings.Split(encodedStatements[1].String(), "\n")
	slices.SortFunc(actualStatementsList, strings.Compare)
	actualStatements := strings.Join(actualStatementsList, "\n")

	if expectedStatements != actualStatements {
		return fmt.Errorf("expected does not match actual\n\n=== EXPECTED\n%s\n\n=== ACTUAL\n%s", expectedStatements, actualStatements)
	}

	return nil
}

func AssertTripleEquals(expected, actual rdf.TripleList) error {
	return AssertStatementEquals(expected.AsQuads(nil), actual.AsQuads(nil))
}
