package jsonldinternal

import (
	"context"

	"github.com/dpb587/inspectjson-go/inspectjson"
	"github.com/dpb587/rdfkit-go/encoding/jsonld/jsonldtype"
	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

var ProcessorRemoteContextsLimit = 128

var profileContextIRI = "http://www.w3.org/ns/json-ld#context"

const (
	ProcessingMode_JSON_LD_1_0 = "json-ld-1.0"
	ProcessingMode_JSON_LD_1_1 = "json-ld-1.1"
)

type contextProcessor struct {
	ctx            context.Context
	documentLoader jsonldtype.DocumentLoader

	processingMode string

	dereferencedDocumentByIRI map[string]dereferencedDocument
}

type dereferencedDocument struct {
	documentURL          *iriutil.ParsedIRI
	documentContextValue inspectjson.Value
}
