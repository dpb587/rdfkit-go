package sparql

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dpb587/rdfkit-go/rdf"
)

type QueryResponse struct {
	Head    QueryResponseHead
	Boolean *bool
	Results *QueryResponseResultList
}

type QueryResponseHead struct {
	Variables QueryResponseHeadVariableList
	Links     QueryResponseHeadLinkList
}

type QueryResponseHeadVariable struct {
	Name string
}

type QueryResponseHeadVariableList []QueryResponseHeadVariable

type QueryResponseHeadLink struct {
	Href string
}

type QueryResponseHeadLinkList []QueryResponseHeadLink

type QueryResponseResult struct {
	Bindings QueryResponseResultBindingMap
}

type QueryResponseResultList []QueryResponseResult

type QueryResponseResultBinding struct {
	Name string
	Term rdf.Term
}

type QueryResponseResultBindingMap map[string]QueryResponseResultBinding

func DecodeQueryResponse(r *http.Response) (*QueryResponse, error) {
	switch strings.SplitN(r.Header.Get("Content-Type"), ";", 2)[0] {
	case "application/sparql-results+json":
		return DecodeQueryResponseJSON(r.Body)
	default:
		io.Copy(os.Stderr, r.Body)
	}

	return nil, fmt.Errorf("unsupported content type: %s", r.Header.Get("Content-Type"))
}
