package htmlmicrodata

import (
	"errors"

	"github.com/dpb587/cursorio-go/cursorio"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/iri"
	"github.com/dpb587/rdfkit-go/rdf"
	"golang.org/x/net/html"
)

type globalEvaluationContext struct {
	ResolvedItemscopes map[*html.Node]rdf.SubjectValue
	DocumentContainer  *DocumentResource
	BlankNodeFactory   rdf.BlankNodeFactory
}

type evaluationContext struct {
	BaseURL *iri.ParsedIRI

	CurrentContainer encoding.ContainerResource

	CurrentSubject      rdf.SubjectValue
	CurrentSubjectRange *cursorio.TextOffsetRange

	CurrentItemtypes []string

	RecursedItemrefs map[string]struct{}

	Global *globalEvaluationContext
}

func (w evaluationContext) ResolveURL(u string) (string, error) {
	if w.BaseURL == nil {
		if valid, err := iri.ParseIRI(u); err == nil && valid.IsAbs() {
			return u, nil
		}

		return "", errors.New("no base url")
	}

	parsed, err := w.BaseURL.Parse(u)
	if err != nil {
		return "", err
	}

	return parsed.String(), nil
}
