package jsonldinternal

import (
	"errors"
	"net/url"
	"strings"

	"github.com/dpb587/rdfkit-go/iri"
)

var errEmptyURL = errors.New("empty URL")

// TODO refactor; does this require a resolved, absolute url?
func resolveURL(base *iri.ParsedIRI, r string) (*iri.ParsedIRI, error) {
	if r == "" && base == nil {
		return nil, errEmptyURL
	} else if r == "" {
		return base, nil
	}

	rURL, err := iri.ParseIRI(r)
	if err != nil {
		return nil, err
	} else if rURL.IsAbs() || base == nil {
		return rURL, nil
	}

	return base.ResolveReference(rURL), nil
}

// TODO attempting to solve spec passing baseURL from term definitions where it is empty
func coalesceBaseURL(bb ...*iri.ParsedIRI) *iri.ParsedIRI {
	for _, b := range bb {
		if b != nil {
			return b
		}
	}

	return nil
}

func isIRI(processingMode string, v string) bool {
	vSplit := strings.SplitN(v, ":", 2)
	if len(vSplit) != 2 {
		return false
	} else if strings.HasPrefix(vSplit[0], "_") {
		return false
	}

	if processingMode == ProcessingMode_JSON_LD_1_1 {
		// test case t0123
		// url.Parse(RequestURI)? does not actually error on an unencoded space
		// there is probably a better way to validate?

		if strings.Contains(v, " ") {
			return false
		} else if _, err := url.ParseRequestURI(v); err != nil {
			return false
		}
	}

	return true
}
