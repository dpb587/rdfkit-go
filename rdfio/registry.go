package rdfio

import (
	"net/http"

	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
)

var Registry = NewRegistry(RegistryOptions{
	HttpClient: http.DefaultClient,
	DocumentLoaderJSONLD: jsonldtype.NewCachingDocumentLoader(
		jsonldtype.NewDefaultDocumentLoader(
			http.DefaultClient,
		),
	),
})
